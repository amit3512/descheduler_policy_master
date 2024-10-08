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

package scheme

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"github.com/amit3512/descheduler_policy_master/pkg/api"
	"github.com/amit3512/descheduler_policy_master/pkg/api/v1alpha1"
	"github.com/amit3512/descheduler_policy_master/pkg/api/v1alpha2"
	"github.com/amit3512/descheduler_policy_master/pkg/apis/componentconfig"
	componentconfigv1alpha1 "github.com/amit3512/descheduler_policy_master/pkg/apis/componentconfig/v1alpha1"
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
	"github.com/amit3512/descheduler_policy_master/pkg/framework/plugins/custompolicy"
)

var (
	Scheme = runtime.NewScheme()
	Codecs = serializer.NewCodecFactory(Scheme)
)

func init() {
	utilruntime.Must(api.AddToScheme(Scheme))
	utilruntime.Must(defaultevictor.AddToScheme(Scheme))
	utilruntime.Must(nodeutilization.AddToScheme(Scheme))
	utilruntime.Must(podlifetime.AddToScheme(Scheme))
	utilruntime.Must(removeduplicates.AddToScheme(Scheme))
	utilruntime.Must(removefailedpods.AddToScheme(Scheme))
	utilruntime.Must(removepodshavingtoomanyrestarts.AddToScheme(Scheme))
	utilruntime.Must(removepodsviolatinginterpodantiaffinity.AddToScheme(Scheme))
	utilruntime.Must(removepodsviolatingnodeaffinity.AddToScheme(Scheme))
	utilruntime.Must(removepodsviolatingnodetaints.AddToScheme(Scheme))
	utilruntime.Must(removepodsviolatingtopologyspreadconstraint.AddToScheme(Scheme))
	utilruntime.Must(custompolicy.AddToScheme(Scheme))


	utilruntime.Must(componentconfig.AddToScheme(Scheme))
	utilruntime.Must(componentconfigv1alpha1.AddToScheme(Scheme))
	utilruntime.Must(v1alpha1.AddToScheme(Scheme))
	utilruntime.Must(v1alpha2.AddToScheme(Scheme))
	utilruntime.Must(Scheme.SetVersionPriority(
		v1alpha2.SchemeGroupVersion,
		v1alpha1.SchemeGroupVersion,
	))
}
