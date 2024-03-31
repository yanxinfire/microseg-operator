package policy

import (
	"context"

	microsegv1 "github.com/yanxinfire/microseg-operator/api/v1"
)

type Engine interface {
	CreateNetworkPolicy(ctx context.Context, policyName, namespace string, spec microsegv1.MicrosegNetworkPolicySpec) error
	DeleteNetworkPolicy(ctx context.Context, policyName, namespace string) error
	GetNetworkPolicy(ctx context.Context, policyName, namespace string) error
}
