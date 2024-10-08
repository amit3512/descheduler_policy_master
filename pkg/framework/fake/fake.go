package fake

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"

	"github.com/amit3512/descheduler_policy_master/pkg/descheduler/evictions"
	podutil "github.com/amit3512/descheduler_policy_master/pkg/descheduler/pod"
	frameworktypes "github.com/amit3512/descheduler_policy_master/pkg/framework/types"
)

type HandleImpl struct {
	ClientsetImpl                 clientset.Interface
	GetPodsAssignedToNodeFuncImpl podutil.GetPodsAssignedToNodeFunc
	SharedInformerFactoryImpl     informers.SharedInformerFactory
	EvictorFilterImpl             frameworktypes.EvictorPlugin
	PodEvictorImpl                *evictions.PodEvictor
}

var _ frameworktypes.Handle = &HandleImpl{}

func (hi *HandleImpl) ClientSet() clientset.Interface {
	return hi.ClientsetImpl
}

func (hi *HandleImpl) GetPodsAssignedToNodeFunc() podutil.GetPodsAssignedToNodeFunc {
	return hi.GetPodsAssignedToNodeFuncImpl
}

func (hi *HandleImpl) SharedInformerFactory() informers.SharedInformerFactory {
	return hi.SharedInformerFactoryImpl
}

func (hi *HandleImpl) Evictor() frameworktypes.Evictor {
	return hi
}

func (hi *HandleImpl) Filter(pod *v1.Pod) bool {
	return hi.EvictorFilterImpl.Filter(pod)
}

func (hi *HandleImpl) PreEvictionFilter(pod *v1.Pod) bool {
	return hi.EvictorFilterImpl.PreEvictionFilter(pod)
}

func (hi *HandleImpl) Evict(ctx context.Context, pod *v1.Pod, opts evictions.EvictOptions) error {
	return hi.PodEvictorImpl.EvictPod(ctx, pod, opts)
}
