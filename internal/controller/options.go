package controller

import (
	"context"
	"github.com/yanxinfire/microseg-operator/configs"
	"github.com/yanxinfire/microseg-operator/internal/services/policy"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MicrosegNetworkPolicyOption func(r *MicrosegNetworkPolicyReconciler)

func WithEngine(engineName string) MicrosegNetworkPolicyOption {
	switch engineName {
	case configs.CalicoEngine:
		return func(r *MicrosegNetworkPolicyReconciler) {
			r.engine = policy.NewCalicoEngine(configs.CalicoCli)
		}
	case configs.K8sEngine:
		return func(r *MicrosegNetworkPolicyReconciler) {
			r.engine = policy.NewK8sEngine(configs.MicrosegK8sCli)
		}
	default:
		return func(r *MicrosegNetworkPolicyReconciler) {
			r.engine = policy.NewK8sEngine(configs.MicrosegK8sCli)
		}
	}

}

func cacheNodeIPs(r *MicrosegNetworkPolicyReconciler) {
	nodeList, _ := configs.MicrosegK8sCli.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	for _, node := range nodeList.Items {
		for _, addr := range node.Status.Addresses {
			if addr.Type == v1.NodeInternalIP {
				r.nodeIPs = append(r.nodeIPs, addr.Address+"/32")
				break
			}
		}
	}
}
