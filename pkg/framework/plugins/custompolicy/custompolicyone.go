package custompolicy

import (
	"context"
	"fmt"
	"sort"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"github.com/amit3512/descheduler_policy_master/pkg/descheduler/evictions"
	nodeutil "github.com/amit3512/descheduler_policy_master/pkg/descheduler/node"
	podutil "github.com/amit3512/descheduler_policy_master/pkg/descheduler/pod"
	frameworktypes "github.com/amit3512/descheduler_policy_master/pkg/framework/types"
)

const CustomPolicyOnePluginName = "CustomPolicyOne"

type CustomPolicyOne struct {
	handle    frameworktypes.Handle
	args      *CustomPolicyOneArgs
	podFilter func(pod *v1.Pod) bool
}

var _ frameworktypes.BalancePlugin = &CustomPolicyOne{}

// NewCustomPolicyOne builds plugin from its arguments while passing a handle
func NewCustomPolicyOne(args runtime.Object, handle frameworktypes.Handle) (frameworktypes.Plugin, error) {
	CustomPolicyOneArgsArgs, ok := args.(*CustomPolicyOneArgs)
	if !ok {
		return nil, fmt.Errorf("want args to be of type CustomPolicyOneArgs, got %T", args)
	}

	podFilter, err := podutil.NewOptions().
		WithFilter(handle.Evictor().Filter).
		BuildFilterFunc()
	if err != nil {
		return nil, fmt.Errorf("error initializing pod filter function: %v", err)
	}

	return &CustomPolicyOne{
		handle:    handle,
		args:      CustomPolicyOneArgsArgs,
		podFilter: podFilter,
	}, nil
}

// Name retrieves the plugin name
func (l *CustomPolicyOne) Name() string {
	return CustomPolicyOnePluginName
}

// Balance extension point implementation for the plugin
func (l *CustomPolicyOne) Balance(ctx context.Context, nodes []*v1.Node) *frameworktypes.Status {
	useDeviationThresholds := l.args.UseDeviationThresholds
	thresholds := l.args.Thresholds
	targetThresholds := l.args.TargetThresholds
	resourceNames := getResourceNames(thresholds)

	// Correct function call with proper arguments
	nodeUsageFunc := l.handle.GetPodsAssignedToNodeFunc()
	lowNodes, sourceNodes := classifyNodes(
		getNodeUsage(nodes, resourceNames, nodeUsageFunc),
		getNodeThresholds(nodes, thresholds, targetThresholds, resourceNames, nodeUsageFunc, useDeviationThresholds),
		func(node *v1.Node, usage NodeUsage, threshold NodeThresholds) bool {
			if nodeutil.IsNodeUnschedulable(node) {
				klog.V(2).InfoS("Node is unschedulable, thus not considered as underutilized", "node", klog.KObj(node))
				return false
			}
			return isNodeWithLowUtilization(usage, threshold.lowResourceThreshold)
		},
		func(node *v1.Node, usage NodeUsage, threshold NodeThresholds) bool {
			return isNodeAboveTargetUtilization(usage, threshold.highResourceThreshold)
		},
	)

	// Calculate CPU utilization for each node in sourceNodes
	nodeCPUUtilization := make(map[*v1.Node]float64)
	for _, node := range sourceNodes {
		pods, _ := nodeUsageFunc(node.spec.nodeName, l.podFilter)
		nodeCPUUtilization[node] = calculateCPUUtilization(node, pods)
	}

	// Sort sourceNodes by CPU utilization in descending order
	sort.SliceStable(sourceNodes, func(i, j int) bool {
		return nodeCPUUtilization[sourceNodes[i]] > nodeCPUUtilization[sourceNodes[j]]
	})

	// Get the node with the highest CPU utilization
	highestUtilizedNode := sourceNodes[0]
	pods, _ := nodeUsageFunc(highestUtilizedNode.spec.nodeName, l.podFilter)

	// Find the pod with the least CPU usage on the highest utilized node
	var podToEvict *v1.Pod
	minCPUUsage := resource.NewQuantity(int64(1<<63-1), resource.DecimalSI) // Max int value
	for _, pod := range pods {
		podCPUUsage := resource.NewQuantity(0, resource.DecimalSI)
		for _, container := range pod.Spec.Containers {
			podCPUUsage.Add(container.Resources.Requests[v1.ResourceCPU])
		}

		if podCPUUsage.Cmp(*minCPUUsage) < 0 {
			minCPUUsage = podCPUUsage
			podToEvict = pod
		}
	}

	if podToEvict == nil {
		klog.V(1).InfoS("No pod found to evict from the most utilized node", "node", klog.KObj(highestUtilizedNode))
		return nil
	}

	klog.V(1).InfoS("Evicting pod with the least CPU usage", "pod", klog.KObj(podToEvict), "node", klog.KObj(highestUtilizedNode))

	// Find the target node for rescheduling the pod
	var targetNode *v1.Node
	for _, node := range lowNodes {
		if isRelatedToPod(podToEvict, node) {
			targetNode = node
			break
		}
	}

	if targetNode == nil {
		// If no related node found, pick the node with the most available resources
		targetNode = findNodeWithMostResources(lowNodes)
	}

	if targetNode == nil {
		klog.V(1).InfoS("No suitable target node found for rescheduling", "pod", klog.KObj(podToEvict))
		return nil
	}

	klog.V(1).InfoS("Rescheduling pod to target node", "pod", klog.KObj(podToEvict), "targetNode", klog.KObj(targetNode))

	// Evict the pod and reschedule it to the target node
	err := l.handle.Evictor().Evict(ctx, podToEvict, evictions.EvictOptions{StrategyName: CustomPolicyOnePluginName})
	if err != nil {
		klog.V(1).InfoS("Failed to evict pod", "pod", klog.KObj(podToEvict), "error", err)
		return frameworktypes.NewStatus(frameworktypes.Error, fmt.Sprintf("Failed to evict pod: %v", err))
	}

	return nil
}

// Dummy function to calculate CPU utilization of a node based on its pods
func calculateCPUUtilization(node *v1.Node, pods []*v1.Pod) float64 {
	totalCPU := resource.NewMilliQuantity(0, resource.DecimalSI)
	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			totalCPU.Add(container.Resources.Requests[v1.ResourceCPU])
		}
	}
	return float64(totalCPU.MilliValue()) / float64(node.Status.Capacity[v1.ResourceCPU].MilliValue()) * 100
}

// Dummy function to find the node with the most available resources
func findNodeWithMostResources(nodes []*v1.Node) *v1.Node {
	var targetNode *v1.Node
	maxAvailableResources := resource.NewQuantity(0, resource.DecimalSI)

	for _, node := range nodes {
		availableResources := calculateAvailableResources(node)
		if availableResources.Cmp(*maxAvailableResources) > 0 {
			maxAvailableResources = availableResources
			targetNode = node
		}
	}

	return targetNode
}

// Dummy function to calculate available resources on a node
func calculateAvailableResources(node *v1.Node) *resource.Quantity {
	// Implement logic to calculate the available resources (CPU, memory, etc.) on the node
	// This function should return a quantity representing the aggregate available resources
	return resource.NewQuantity(1000, resource.DecimalSI) // Example value
}

// Implement logic to determine if a node is related to the pod
func isRelatedToPod(pod *v1.Pod, node *v1.Node) bool {
	podLabels := pod.Labels
	nodeLabels := node.Labels
	relatedLabels := []string{"app", "service", "database"}

	for _, key := range relatedLabels {
		if podValue, podHasLabel := podLabels[key]; podHasLabel {
			if nodeValue, nodeHasLabel := nodeLabels[key]; nodeHasLabel {
				if podValue == nodeValue {
					return true
				}
			}
		}
	}

	return false
}
