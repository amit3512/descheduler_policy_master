/*
Copyright 2022 The Kubernetes Authors.
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

package custompolicy

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"github.com/amit3512/descheduler_policy_master/pkg/api"
)

func TestSetDefaults_CustomPolicyOneArgs(t *testing.T) {
	tests := []struct {
		name string
		in   runtime.Object
		want runtime.Object
	}{
		{
			name: "CustomPolicyOneArgs empty",
			in:   &CustomPolicyOneArgs{},
			want: &CustomPolicyOneArgs{
				UseDeviationThresholds: false,
				Thresholds:             nil,
				TargetThresholds:       nil,
				NumberOfNodes:          0,
			},
		},
		{
			name: "CustomPolicyOneArgs with value",
			in: &CustomPolicyOneArgs{
				UseDeviationThresholds: true,
				Thresholds: api.ResourceThresholds{
					v1.ResourceCPU:    20,
					v1.ResourceMemory: 120,
				},
				TargetThresholds: api.ResourceThresholds{
					v1.ResourceCPU:    80,
					v1.ResourceMemory: 80,
				},
				NumberOfNodes: 10,
			},
			want: &CustomPolicyOneArgs{
				UseDeviationThresholds: true,
				Thresholds: api.ResourceThresholds{
					v1.ResourceCPU:    20,
					v1.ResourceMemory: 120,
				},
				TargetThresholds: api.ResourceThresholds{
					v1.ResourceCPU:    80,
					v1.ResourceMemory: 80,
				},
				NumberOfNodes: 10,
			},
		},
	}
	for _, tc := range tests {
		scheme := runtime.NewScheme()
		utilruntime.Must(AddToScheme(scheme))
		t.Run(tc.name, func(t *testing.T) {
			scheme.Default(tc.in)
			if diff := cmp.Diff(tc.in, tc.want); diff != "" {
				t.Errorf("Got unexpected defaults (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestSetDefaults_HighNodeUtilizationArgs(t *testing.T) {
	tests := []struct {
		name string
		in   runtime.Object
		want runtime.Object
	}{
		{
			name: "HighNodeUtilizationArgs empty",
			in:   &HighNodeUtilizationArgs{},
			want: &HighNodeUtilizationArgs{
				Thresholds:    nil,
				NumberOfNodes: 0,
			},
		},
		{
			name: "HighNodeUtilizationArgs with value",
			in: &HighNodeUtilizationArgs{
				Thresholds: api.ResourceThresholds{
					v1.ResourceCPU:    20,
					v1.ResourceMemory: 120,
				},
				NumberOfNodes: 10,
			},
			want: &HighNodeUtilizationArgs{
				Thresholds: api.ResourceThresholds{
					v1.ResourceCPU:    20,
					v1.ResourceMemory: 120,
				},
				NumberOfNodes: 10,
			},
		},
	}
	for _, tc := range tests {
		scheme := runtime.NewScheme()
		utilruntime.Must(AddToScheme(scheme))
		t.Run(tc.name, func(t *testing.T) {
			scheme.Default(tc.in)
			if diff := cmp.Diff(tc.in, tc.want); diff != "" {
				t.Errorf("Got unexpected defaults (-want, +got):\n%s", diff)
			}
		})
	}
}
