package policy

import (
	"strconv"
	"strings"

	"github.com/projectcalico/api/pkg/lib/numorstring"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/pkg/errors"
)

func validatePorts(ports string) error {
	splitPorts := []string{}
	if strings.Contains(ports, "-") {
		splitPorts = strings.Split(ports, "-")
		if len(splitPorts) != 2 {
			return errors.New("range ports should be specified like 100-200")
		}
	} else {
		splitPorts = strings.Split(ports, ",")
	}
	for _, p := range splitPorts {
		if _, err := strconv.Atoi(strings.TrimSpace(p)); err != nil {
			return errors.New("port must be a valid number")
		}
	}
	return nil
}

func calicoRulePorts(strPorts string) []numorstring.Port {
	var policyPorts []numorstring.Port
	ps := strings.TrimSpace(strPorts)
	if strings.Contains(ps, "-") {
		ports := strings.Split(ps, "-")
		if len(ports) != 2 {
			logrus.Warnf("Invalid ports: %s", strPorts)
			return policyPorts
		}

		mini, err := strconv.Atoi(ports[0])
		if err != nil {
			return policyPorts
		}
		max, err := strconv.Atoi(ports[1])
		p, err := numorstring.PortFromRange(uint16(mini), uint16(max))
		policyPorts = append(policyPorts, p)
		return policyPorts
	} else if strings.Contains(ps, ",") {
		ports := strings.Split(ps, ",")
		for _, p := range ports {
			p = strings.TrimSpace(p)
			_, err := strconv.Atoi(p)
			if err != nil {
				return policyPorts
			}
			cp, err := numorstring.PortFromString(p)
			if err != nil {
				logrus.Errorf("Port From String err: %v", err)
				return policyPorts
			}
			policyPorts = append(policyPorts, cp)
		}
		return policyPorts
	} else if _, err := strconv.Atoi(ps); err == nil {
		cp, err := numorstring.PortFromString(ps)
		if err != nil {
			return nil
		}
		policyPorts = append(policyPorts, cp)
		return policyPorts
	} else {
		logrus.Warnf("Invalid ports")
	}

	return policyPorts
}

func k8sRulePorts(strPorts string, ruleProtocol *v1.Protocol) []netv1.NetworkPolicyPort {
	ports := []netv1.NetworkPolicyPort{}
	splitPorts := []string{}
	if strings.Contains(strPorts, "-") {
		splitPorts = strings.Split(strPorts, "-")
		startPort, _ := strconv.Atoi(strings.TrimSpace(splitPorts[0]))
		endPort, _ := strconv.Atoi(strings.TrimSpace(splitPorts[1]))
		for i := startPort; i < endPort+1; i++ {
			ports = append(ports, netv1.NetworkPolicyPort{
				Protocol: ruleProtocol,
				Port: &intstr.IntOrString{
					Type:   0,
					IntVal: int32(i),
				},
			})
		}
	} else {
		splitPorts = strings.Split(strPorts, ",")
		for _, p := range splitPorts {
			rp, _ := strconv.Atoi(strings.TrimSpace(p))
			ports = append(ports, netv1.NetworkPolicyPort{
				Protocol: ruleProtocol,
				Port: &intstr.IntOrString{
					Type:   0,
					IntVal: int32(rp),
				},
			})
		}
	}
	return ports
}
