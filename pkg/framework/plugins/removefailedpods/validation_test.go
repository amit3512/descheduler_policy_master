package removefailedpods

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/amit3512/descheduler_policy_master/pkg/api"
)

func TestValidateRemoveFailedPodsArgs(t *testing.T) {
	var oneHourPodLifetimeSeconds uint = 3600
	testCases := []struct {
		description string
		args        *RemoveFailedPodsArgs
		expectError bool
	}{
		{
			description: "valid namespace args, no errors",
			args: &RemoveFailedPodsArgs{
				Namespaces: &api.Namespaces{
					Include: []string{"default"},
				},
				ExcludeOwnerKinds:     []string{"Job"},
				Reasons:               []string{"ReasonDoesNotMatch"},
				MinPodLifetimeSeconds: &oneHourPodLifetimeSeconds,
			},
			expectError: false,
		},
		{
			description: "invalid namespaces args, expects error",
			args: &RemoveFailedPodsArgs{
				Namespaces: &api.Namespaces{
					Include: []string{"default"},
					Exclude: []string{"kube-system"},
				},
			},
			expectError: true,
		},
		{
			description: "valid label selector args, no errors",
			args: &RemoveFailedPodsArgs{
				LabelSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"role.kubernetes.io/node": ""},
				},
			},
			expectError: false,
		},
		{
			description: "invalid label selector args, expects errors",
			args: &RemoveFailedPodsArgs{
				LabelSelector: &metav1.LabelSelector{
					MatchExpressions: []metav1.LabelSelectorRequirement{
						{
							Operator: metav1.LabelSelectorOpIn,
						},
					},
				},
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := ValidateRemoveFailedPodsArgs(tc.args)
			hasError := err != nil
			if tc.expectError != hasError {
				t.Error("unexpected arg validation behavior")
			}
		})
	}
}
