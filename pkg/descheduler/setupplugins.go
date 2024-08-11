/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package descheduler

import (
	"github.com/amit3512/descheduler_policy_master/pkg/framework/pluginregistry"
	"github.com/amit3512/descheduler_policy_master/pkg/framework/plugins/defaultevictor"
	"github.com/amit3512/descheduler_policy_master/pkg/framework/plugins/nodeutilization"
	"github.com/amit3512/descheduler_policy_master/pkg/framework/plugins/podlifetime"
	"github.com/amit3512/descheduler_policy_master/pkg/framework/plugins/removeduplicates"
	"github.com/amit3512/descheduler_policy_master/pkg/framework/plugins/removefailedpods"
	"github.com/amit3512/descheduler_policy_master/pkg/framework/plugins/removepodshavingtoomanyrestarts"
	"github.com/amit3512/descheduler_policy_master/pkg/framework/plugins/removepodsviolatinginterpodantiaffinity"
	"github.com/amit3512/descheduler_policy_master/pkg/framework/plugins/removepodsviolatingnodeaffinity"
	"github.com/amit3512/descheduler_policy_master/pkg/framework/plugins/removepodsviolatingnodetaints"
	"github.com/amit3512/descheduler_policy_master/pkg/framework/plugins/removepodsviolatingtopologyspreadconstraint"
)

func SetupPlugins() {
	pluginregistry.PluginRegistry = pluginregistry.NewRegistry()
	RegisterDefaultPlugins(pluginregistry.PluginRegistry)
}

func RegisterDefaultPlugins(registry pluginregistry.Registry) {
	pluginregistry.Register(defaultevictor.PluginName, defaultevictor.New, &defaultevictor.DefaultEvictor{}, &defaultevictor.DefaultEvictorArgs{}, defaultevictor.ValidateDefaultEvictorArgs, defaultevictor.SetDefaults_DefaultEvictorArgs, registry)
	pluginregistry.Register(nodeutilization.LowNodeUtilizationPluginName, nodeutilization.NewLowNodeUtilization, &nodeutilization.LowNodeUtilization{}, &nodeutilization.LowNodeUtilizationArgs{}, nodeutilization.ValidateLowNodeUtilizationArgs, nodeutilization.SetDefaults_LowNodeUtilizationArgs, registry)
	pluginregistry.Register(nodeutilization.HighNodeUtilizationPluginName, nodeutilization.NewHighNodeUtilization, &nodeutilization.HighNodeUtilization{}, &nodeutilization.HighNodeUtilizationArgs{}, nodeutilization.ValidateHighNodeUtilizationArgs, nodeutilization.SetDefaults_HighNodeUtilizationArgs, registry)
	pluginregistry.Register(podlifetime.PluginName, podlifetime.New, &podlifetime.PodLifeTime{}, &podlifetime.PodLifeTimeArgs{}, podlifetime.ValidatePodLifeTimeArgs, podlifetime.SetDefaults_PodLifeTimeArgs, registry)
	pluginregistry.Register(removeduplicates.PluginName, removeduplicates.New, &removeduplicates.RemoveDuplicates{}, &removeduplicates.RemoveDuplicatesArgs{}, removeduplicates.ValidateRemoveDuplicatesArgs, removeduplicates.SetDefaults_RemoveDuplicatesArgs, registry)
	pluginregistry.Register(removefailedpods.PluginName, removefailedpods.New, &removefailedpods.RemoveFailedPods{}, &removefailedpods.RemoveFailedPodsArgs{}, removefailedpods.ValidateRemoveFailedPodsArgs, removefailedpods.SetDefaults_RemoveFailedPodsArgs, registry)
	pluginregistry.Register(removepodshavingtoomanyrestarts.PluginName, removepodshavingtoomanyrestarts.New, &removepodshavingtoomanyrestarts.RemovePodsHavingTooManyRestarts{}, &removepodshavingtoomanyrestarts.RemovePodsHavingTooManyRestartsArgs{}, removepodshavingtoomanyrestarts.ValidateRemovePodsHavingTooManyRestartsArgs, removepodshavingtoomanyrestarts.SetDefaults_RemovePodsHavingTooManyRestartsArgs, registry)
	pluginregistry.Register(removepodsviolatinginterpodantiaffinity.PluginName, removepodsviolatinginterpodantiaffinity.New, &removepodsviolatinginterpodantiaffinity.RemovePodsViolatingInterPodAntiAffinity{}, &removepodsviolatinginterpodantiaffinity.RemovePodsViolatingInterPodAntiAffinityArgs{}, removepodsviolatinginterpodantiaffinity.ValidateRemovePodsViolatingInterPodAntiAffinityArgs, removepodsviolatinginterpodantiaffinity.SetDefaults_RemovePodsViolatingInterPodAntiAffinityArgs, registry)
	pluginregistry.Register(removepodsviolatingnodeaffinity.PluginName, removepodsviolatingnodeaffinity.New, &removepodsviolatingnodeaffinity.RemovePodsViolatingNodeAffinity{}, &removepodsviolatingnodeaffinity.RemovePodsViolatingNodeAffinityArgs{}, removepodsviolatingnodeaffinity.ValidateRemovePodsViolatingNodeAffinityArgs, removepodsviolatingnodeaffinity.SetDefaults_RemovePodsViolatingNodeAffinityArgs, registry)
	pluginregistry.Register(removepodsviolatingnodetaints.PluginName, removepodsviolatingnodetaints.New, &removepodsviolatingnodetaints.RemovePodsViolatingNodeTaints{}, &removepodsviolatingnodetaints.RemovePodsViolatingNodeTaintsArgs{}, removepodsviolatingnodetaints.ValidateRemovePodsViolatingNodeTaintsArgs, removepodsviolatingnodetaints.SetDefaults_RemovePodsViolatingNodeTaintsArgs, registry)
	pluginregistry.Register(removepodsviolatingtopologyspreadconstraint.PluginName, removepodsviolatingtopologyspreadconstraint.New, &removepodsviolatingtopologyspreadconstraint.RemovePodsViolatingTopologySpreadConstraint{}, &removepodsviolatingtopologyspreadconstraint.RemovePodsViolatingTopologySpreadConstraintArgs{}, removepodsviolatingtopologyspreadconstraint.ValidateRemovePodsViolatingTopologySpreadConstraintArgs, removepodsviolatingtopologyspreadconstraint.SetDefaults_RemovePodsViolatingTopologySpreadConstraintArgs, registry)
}
