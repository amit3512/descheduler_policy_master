package e2e

import (
	"context"
	"os"
	"testing"

	componentbaseconfig "k8s.io/component-base/config"
	"github.com/amit3512/descheduler_policy_master/cmd/descheduler/app/options"
	deschedulerapi "github.com/amit3512/descheduler_policy_master/pkg/api"
	"github.com/amit3512/descheduler_policy_master/pkg/descheduler"
	"github.com/amit3512/descheduler_policy_master/pkg/descheduler/client"
	eutils "github.com/amit3512/descheduler_policy_master/pkg/descheduler/evictions/utils"
)

func TestClientConnectionConfiguration(t *testing.T) {
	ctx := context.Background()
	clientConnection := componentbaseconfig.ClientConnectionConfiguration{
		Kubeconfig: os.Getenv("KUBECONFIG"),
		QPS:        50,
		Burst:      100,
	}
	clientSet, err := client.CreateClient(clientConnection, "")
	if err != nil {
		t.Errorf("Error during client creation with %v", err)
	}

	s, err := options.NewDeschedulerServer()
	if err != nil {
		t.Fatalf("Unable to initialize server: %v", err)
	}
	s.Client = clientSet
	evictionPolicyGroupVersion, err := eutils.SupportEviction(s.Client)
	if err != nil || len(evictionPolicyGroupVersion) == 0 {
		t.Errorf("Error when checking support for eviction: %v", err)
	}
	if err := descheduler.RunDeschedulerStrategies(ctx, s, &deschedulerapi.DeschedulerPolicy{}, evictionPolicyGroupVersion); err != nil {
		t.Errorf("Error running descheduler strategies: %+v", err)
	}
}
