package policy

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	microsegv1 "github.com/yanxinfire/microseg-operator/api/v1"
	"github.com/yanxinfire/microseg-operator/internal/services/types"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var _ Engine = (*k8sEngine)(nil)

type k8sEngine struct {
	client kubernetes.Interface
}

func NewK8sEngine(client kubernetes.Interface) Engine {
	return &k8sEngine{
		client: client,
	}
}

func (e *k8sEngine) CreateNetworkPolicy(ctx context.Context, policyName, namespace string, spec microsegv1.MicrosegNetworkPolicySpec) error {
	np, err := e.generatePolicy(policyName, namespace, spec)
	if err != nil {
		return errors.Wrap(err, "failed to build network policy")
	}
	npResponse, err := e.client.NetworkingV1().NetworkPolicies(namespace).Create(ctx, np, metav1.CreateOptions{})
	nprBytes, _ := json.Marshal(npResponse)
	logrus.Infoln(string(nprBytes))
	return err
}

func (e *k8sEngine) DeleteNetworkPolicy(ctx context.Context, policyName, namespace string) error {
	err := e.client.NetworkingV1().NetworkPolicies(namespace).Delete(ctx, policyName, metav1.DeleteOptions{})
	return err
}

func (e *k8sEngine) buildK8sIngressRule(ingressRules []microsegv1.MicrosegNetworkPolicyIngress) ([]netv1.NetworkPolicyIngressRule, error) {
	ingresses := make([]netv1.NetworkPolicyIngressRule, 0, len(ingressRules))
	for _, inRule := range ingressRules {
		if inRule.Action == types.Deny {
			continue
		}
		var ruleProtocol *v1.Protocol
		if strings.TrimSpace(inRule.Protocol) != "" {
			p := v1.Protocol(inRule.Protocol)
			ruleProtocol = &p
		}
		npIngressRule := netv1.NetworkPolicyIngressRule{
			Ports: []netv1.NetworkPolicyPort{},
			From:  []netv1.NetworkPolicyPeer{},
		}
		if strings.TrimSpace(inRule.Ports) != "" {
			if err := validatePorts(inRule.Ports); err != nil {
				return nil, err
			}
			ports := k8sRulePorts(inRule.Ports, ruleProtocol)
			npIngressRule.Ports = append(npIngressRule.Ports, ports...)

		} else {
			npIngressRule.Ports = append(npIngressRule.Ports,
				netv1.NetworkPolicyPort{Protocol: ruleProtocol})
		}
		for _, ipBlock := range inRule.IPBlocks {
			npIngressRule.From = append(npIngressRule.From,
				e.fillEntityRuleField(nil, nil, ipBlock))
		}
		npIngressRule.From = append(npIngressRule.From,
			e.fillEntityRuleField(inRule.NamespaceSelector, inRule.ResourceSelector, ""))
		ingresses = append(ingresses, npIngressRule)
	}
	return ingresses, nil
}

func (e *k8sEngine) buildK8sEgressRule(egressRules []microsegv1.MicrosegNetworkPolicyEgress) ([]netv1.NetworkPolicyEgressRule, error) {
	egresses := make([]netv1.NetworkPolicyEgressRule, 0, len(egressRules))
	for _, eRule := range egressRules {
		if eRule.Action == types.Deny {
			continue
		}
		var ruleProtocol *v1.Protocol
		if strings.TrimSpace(eRule.Protocol) != "" {
			p := v1.Protocol(eRule.Protocol)
			ruleProtocol = &p
		}
		npEgressRule := netv1.NetworkPolicyEgressRule{
			Ports: []netv1.NetworkPolicyPort{},
			To:    []netv1.NetworkPolicyPeer{},
		}
		if len(eRule.Ports) > 0 {
			if err := validatePorts(eRule.Ports); err != nil {
				return nil, err
			}
			ports := k8sRulePorts(eRule.Ports, ruleProtocol)
			npEgressRule.Ports = append(npEgressRule.Ports, ports...)
		} else {
			npEgressRule.Ports = append(npEgressRule.Ports,
				netv1.NetworkPolicyPort{Protocol: ruleProtocol})
		}
		for _, ipBlock := range eRule.IPBlocks {
			npEgressRule.To = append(npEgressRule.To,
				e.fillEntityRuleField(nil, nil, ipBlock))
		}
		npEgressRule.To = append(npEgressRule.To,
			e.fillEntityRuleField(eRule.NamespaceSelector, eRule.ResourceSelector, ""))
		egresses = append(egresses, npEgressRule)
	}
	return egresses, nil
}

func (e *k8sEngine) generatePolicy(policyName, namespace string, spec microsegv1.MicrosegNetworkPolicySpec) (*netv1.NetworkPolicy, error) {
	isController := true
	np := &netv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      policyName,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "microseg.security.cn",
					Kind:       "MicrosegNetworkPolicy",
					Name:       policyName,
					Controller: &isController,
				},
			},
		},
		Spec: netv1.NetworkPolicySpec{
			PolicyTypes: []netv1.PolicyType{},
		},
	}
	if len(spec.ResourceSelector) > 0 {
		np.Spec.PodSelector = metav1.LabelSelector{
			MatchLabels: spec.ResourceSelector,
		}
	}

	np.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.k8s.io",
		Version: "v1",
		Kind:    "NetworkPolicy",
	})

	for _, v := range spec.PolicyTypes {
		if v == types.Ingress {
			np.Spec.PolicyTypes = append(np.Spec.PolicyTypes, netv1.PolicyTypeIngress)
		} else if v == types.Egress {
			np.Spec.PolicyTypes = append(np.Spec.PolicyTypes, netv1.PolicyTypeEgress)
		}
	}

	ingressRules, err := e.buildK8sIngressRule(spec.Rules.Ingress)
	if err != nil {
		return nil, err
	}
	egressRules, err := e.buildK8sEgressRule(spec.Rules.Egress)
	if err != nil {
		return nil, err
	}
	np.Spec.Ingress = ingressRules
	np.Spec.Egress = egressRules
	return np, nil
}

func (e *k8sEngine) fillEntityRuleField(namespaceSelector, resourceSelector map[string]string, ipBlock string) netv1.NetworkPolicyPeer {
	entRule := netv1.NetworkPolicyPeer{
		NamespaceSelector: &metav1.LabelSelector{},
	}
	if len(resourceSelector) > 0 {
		entRule.PodSelector = &metav1.LabelSelector{MatchLabels: resourceSelector}
	}
	if len(namespaceSelector) > 0 {
		entRule.NamespaceSelector = &metav1.LabelSelector{MatchLabels: namespaceSelector}
	}
	if len(ipBlock) > 0 {
		entRule.IPBlock = &netv1.IPBlock{CIDR: ipBlock}
	}
	return entRule
}

func (e *k8sEngine) GetNetworkPolicy(ctx context.Context, policyName, namespace string) error {
	_, err := e.client.NetworkingV1().NetworkPolicies(namespace).Get(ctx, policyName, metav1.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "failed to get k8s Network Policies %s", policyName)
	}
	return nil
}
