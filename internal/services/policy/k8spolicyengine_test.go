package policy

import (
	"context"
	"testing"
	"time"

	microsegv1 "github.com/yanxinfire/microseg-operator/api/v1"

	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestCreatePolicy(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientset := testclient.NewSimpleClientset()
	e := NewK8sEngine(clientset)

	rules := microsegv1.MicrosegNetworkPolicySpec{
		ResourceSelector: map[string]string{
			"xinyansec-resource": "123456",
		},
		Rules: microsegv1.MicrosegNetworkPolicyRule{
			Ingress: []microsegv1.MicrosegNetworkPolicyIngress{
				{
					ResourceSelector: map[string]string{
						"xinyansec-segment": "11111",
					},
					Action:   "Allow",
					Protocol: "TCP",
					Ports:    "80-84",
				},
				{
					ResourceSelector: map[string]string{
						"xinyansec-segment": "44444",
					},
					Action:   "Allow",
					Protocol: "TCP",
				},
				{
					ResourceSelector: map[string]string{
						"xinyansec-segment": "44444",
					},
					Action: "Allow",
					Ports:  "80",
				},
			},
			Egress: []microsegv1.MicrosegNetworkPolicyEgress{
				{
					ResourceSelector: map[string]string{
						"xinyansec-resource": "22222",
					},
					Action:   "Allow",
					Protocol: "UDP",
					Ports:    "80,84",
				},
				{
					ResourceSelector: map[string]string{
						"xinyansec-resource": "33333",
					},
					Action:   "Deny",
					Protocol: "TCP",
					Ports:    "80-84",
				},
			},
		},
		PolicyTypes: nil,
	}

	err := e.CreateNetworkPolicy(ctx, "resource-50", "testns", rules)
	if err != nil {
		t.Error(err)
	}
}
