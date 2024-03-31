package policy

import (
	"context"
	"fmt"
	"strings"

	microsegv1 "github.com/yanxinfire/microseg-operator/api/v1"
	"github.com/yanxinfire/microseg-operator/internal/services/types"

	"github.com/pkg/errors"
	v3 "github.com/projectcalico/api/pkg/apis/projectcalico/v3"
	projectcalicov3 "github.com/projectcalico/api/pkg/client/clientset_generated/clientset/typed/projectcalico/v3"
	"github.com/projectcalico/api/pkg/lib/numorstring"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ Engine = (*calicoEngine)(nil)

type calicoEngine struct {
	client projectcalicov3.ProjectcalicoV3Interface
}

func NewCalicoEngine(calicoCli projectcalicov3.ProjectcalicoV3Interface) Engine {
	return &calicoEngine{
		client: calicoCli,
	}
}

func (e *calicoEngine) CreateNetworkPolicy(ctx context.Context, policyName, namespace string, spec microsegv1.MicrosegNetworkPolicySpec) error {
	policyReq, err := e.generatePolicy(policyName, namespace, spec)
	if err != nil {
		return err
	}
	_, err = e.client.GlobalNetworkPolicies().Create(ctx, policyReq, metav1.CreateOptions{})
	if err != nil {
		//return errors.Wrap(err, "failed to create calico GlobalNetworkPolicy")
		return err
	}
	return nil
}

func (e *calicoEngine) DeleteNetworkPolicy(ctx context.Context, policyName, namespace string) error {
	err := e.client.GlobalNetworkPolicies().Delete(ctx, policyName, metav1.DeleteOptions{})
	if err != nil {
		return errors.Wrapf(err, "failed to delete calico Network Policies")
	}
	return nil
}

func (e *calicoEngine) generatePolicy(policyName, namespace string, spec microsegv1.MicrosegNetworkPolicySpec) (*v3.GlobalNetworkPolicy, error) {
	isController := true
	np := &v3.GlobalNetworkPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GlobalNetworkPolicy",
			APIVersion: "projectcalico.org/v3",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: policyName,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "microseg.security.cn",
					Kind:       "MicrosegNetworkPolicy",
					Name:       policyName,
					Controller: &isController,
				},
			},
		},
		Spec: v3.GlobalNetworkPolicySpec{},
	}
	if spec.Order != nil {
		policyOrder := float64(*spec.Order)
		np.Spec.Order = &policyOrder
	}

	var resSelectors []string
	for k, v := range spec.ResourceSelector {
		resSelectors = append(resSelectors, fmt.Sprintf("%s == \"%s\"", k, v))
	}
	if len(resSelectors) > 0 {
		np.Spec.Selector = strings.Join(resSelectors, " && ")
	}

	var nsSelectors []string
	for k, v := range spec.NamespaceSelector {
		nsSelectors = append(nsSelectors, fmt.Sprintf("%s == \"%s\"", k, v))
	}
	if len(nsSelectors) > 0 {
		np.Spec.NamespaceSelector = strings.Join(nsSelectors, " && ")
	}

	if len(spec.PolicyTypes) > 0 {
		np.Spec.Types = []v3.PolicyType{}
		for _, v := range spec.PolicyTypes {
			if v == types.Ingress {
				np.Spec.Types = append(np.Spec.Types, v3.PolicyTypeIngress)
			} else if v == types.Egress {
				np.Spec.Types = append(np.Spec.Types, v3.PolicyTypeEgress)
			}
		}
	}
	ingresses, err := e.buildIngressRule(spec.Rules.Ingress)
	if err != nil {
		return nil, err
	}
	egresses, err := e.buildEgressRule(spec.Rules.Egress)
	if err != nil {
		return nil, err
	}
	np.Spec.Ingress = ingresses
	np.Spec.Egress = egresses
	return np, nil
}

func (e *calicoEngine) buildIngressRule(ingressRules []microsegv1.MicrosegNetworkPolicyIngress) ([]v3.Rule, error) {
	ingressRules = append(ingressRules, microsegv1.MicrosegNetworkPolicyIngress{
		ResourceSelector: map[string]string{"console-job": "apiscan"},
		Action:           "Allow",
		Protocol:         "TCP",
	})
	ingresses := make([]v3.Rule, 0, len(ingressRules))
	for _, inRule := range ingressRules {
		calicoIngressRule := v3.Rule{
			Source: e.fillEntityRuleField(inRule.NamespaceSelector, inRule.ResourceSelector, inRule.IPBlocks),
			Action: v3.Action(inRule.Action),
		}
		protocol := numorstring.ProtocolFromString(inRule.Protocol)
		calicoIngressRule.Protocol = &protocol
		if strings.TrimSpace(inRule.Ports) != "" {
			destRule := v3.EntityRule{}
			if err := validatePorts(inRule.Ports); err != nil {
				return nil, err
			}
			ports := calicoRulePorts(inRule.Ports)
			destRule.Ports = append(destRule.Ports, ports...)
			calicoIngressRule.Destination = destRule
		}
		ingresses = append(ingresses, calicoIngressRule)
	}
	return ingresses, nil
}

func (e *calicoEngine) buildEgressRule(egressRules []microsegv1.MicrosegNetworkPolicyEgress) ([]v3.Rule, error) {
	egressRules = append(egressRules, microsegv1.MicrosegNetworkPolicyEgress{
		ResourceSelector: map[string]string{"k8s-app": "kube-dns"},
		Action:           "Allow",
		Protocol:         "UDP",
	})
	egresses := make([]v3.Rule, 0, len(egressRules))
	for _, eRule := range egressRules {
		calicoEgressRule := v3.Rule{
			Destination: e.fillEntityRuleField(eRule.NamespaceSelector, eRule.ResourceSelector, eRule.IPBlocks),
			Action:      v3.Action(eRule.Action),
		}
		protocol := numorstring.ProtocolFromString(eRule.Protocol)
		calicoEgressRule.Protocol = &protocol
		if strings.TrimSpace(eRule.Ports) != "" {
			if err := validatePorts(eRule.Ports); err != nil {
				return nil, err
			}
			ports := calicoRulePorts(eRule.Ports)
			calicoEgressRule.Destination.Ports = append(calicoEgressRule.Destination.Ports, ports...)
		}
		egresses = append(egresses, calicoEgressRule)
	}
	return egresses, nil
}

func (e *calicoEngine) fillEntityRuleField(namespaceSelector, resourceSelector map[string]string, ipBlocks []string) v3.EntityRule {
	entRule := v3.EntityRule{}
	var selectors []string
	for k, v := range resourceSelector {
		selectors = append(selectors, fmt.Sprintf("%s == \"%s\"", k, v))
	}

	var nsSelectors []string
	for k, v := range namespaceSelector {
		nsSelectors = append(nsSelectors, fmt.Sprintf("%s == \"%s\"", k, v))
	}
	if len(nsSelectors) > 0 {
		entRule.NamespaceSelector = strings.Join(nsSelectors, " && ")
	}

	if len(selectors) > 0 {
		entRule.Selector = strings.Join(selectors, " && ")
	}
	if len(ipBlocks) > 0 {
		entRule.Nets = append(entRule.Nets, ipBlocks...)
	}
	return entRule
}

func (e *calicoEngine) GetNetworkPolicy(ctx context.Context, policyName, namespace string) error {
	_, err := e.client.GlobalNetworkPolicies().Get(ctx, policyName, metav1.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "failed to get calico Network Policies %s", policyName)
	}
	return nil
}
