package removeduplicates

import (
	"testing"

	"github.com/amit3512/descheduler_policy_master/pkg/api"
)

func TestValidateRemovePodsViolatingNodeTaintsArgs(t *testing.T) {
	testCases := []struct {
		description string
		args        *RemoveDuplicatesArgs
		expectError bool
	}{
		{
			description: "valid namespace args, no errors",
			args: &RemoveDuplicatesArgs{
				ExcludeOwnerKinds: []string{"Job"},
				Namespaces: &api.Namespaces{
					Include: []string{"default"},
				},
			},
			expectError: false,
		},
		{
			description: "invalid namespaces args, expects error",
			args: &RemoveDuplicatesArgs{
				ExcludeOwnerKinds: []string{"Job"},
				Namespaces: &api.Namespaces{
					Include: []string{"default"},
					Exclude: []string{"kube-system"},
				},
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := ValidateRemoveDuplicatesArgs(tc.args)

			hasError := err != nil
			if tc.expectError != hasError {
				t.Error("unexpected arg validation behavior")
			}
		})
	}
}
