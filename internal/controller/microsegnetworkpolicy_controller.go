/*
Copyright 2024 Xin Yan.

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

package controller

import (
	"context"
	"encoding/json"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	"github.com/yanxinfire/microseg-operator/internal/services/policy"
	"github.com/yanxinfire/microseg-operator/internal/services/types"
	"github.com/yanxinfire/microseg-operator/pkg/event"
	"github.com/yanxinfire/microseg-operator/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"

	microsegv1 "github.com/yanxinfire/microseg-operator/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	MicrosegNetworkPolicyFinalizer string = "microsegnetworkpolicies.finalizers.security.cn"
)

// MicrosegNetworkPolicyReconciler reconciles a MicrosegNetworkPolicy object
type MicrosegNetworkPolicyReconciler struct {
	client.Client
	Log     logr.Logger
	Scheme  *runtime.Scheme
	record  event.Recorder
	engine  policy.Engine
	nodeIPs []string
}

func NewMicrosegNetworkPolicyReconciler(client client.Client, scheme *runtime.Scheme, opts ...MicrosegNetworkPolicyOption) *MicrosegNetworkPolicyReconciler {
	p := &MicrosegNetworkPolicyReconciler{
		Client: client,
		Scheme: scheme,
		Log:    ctrl.Log.WithName("MicrosegNetworkPolicy"),
	}
	for _, opt := range opts {
		opt(p)
	}

	cacheNodeIPs(p)
	p.CreateIsolationPolicy()
	return p
}

func (r *MicrosegNetworkPolicyReconciler) CreateIsolationPolicy() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	order := 10
	tnp := microsegv1.MicrosegNetworkPolicy{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: "security-isolation-gnp",
		},
		Spec: microsegv1.MicrosegNetworkPolicySpec{
			Order: &order,
			ResourceSelector: map[string]string{
				util.IsolationLabelKey: "true",
			},
			PolicyTypes: []string{types.Ingress, types.Egress},
		},
	}
	err := r.applyNetworkPolicy(ctx, "security-isolation-gnp", "", tnp.Spec)
	if err != nil {
		logrus.Error(err)
	}
}

//+kubebuilder:rbac:groups=microseg.xinyan.cn,resources=microsegnetworkpolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=microseg.xinyan.cn,resources=microsegnetworkpolicies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=microseg.xinyan.cn,resources=microsegnetworkpolicies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MicrosegNetworkPolicy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *MicrosegNetworkPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	mLog := r.Log.WithValues("microseg-network-policy", req.NamespacedName)
	mLog.Info("Reconcile MicrosegNetworkPolicy")
	// your logic here
	var tnp microsegv1.MicrosegNetworkPolicy
	if err := r.Get(ctx, req.NamespacedName, &tnp); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if tnp.ObjectMeta.DeletionTimestamp.IsZero() {
		if !util.ContainsString(tnp.ObjectMeta.Finalizers, MicrosegNetworkPolicyFinalizer) {
			tnp.ObjectMeta.Finalizers = append(tnp.ObjectMeta.Finalizers, MicrosegNetworkPolicyFinalizer)
			if err := r.Update(ctx, &tnp); err != nil {
				return ctrl.Result{}, err
			}
		}

		nowSpec, err := json.Marshal(tnp.Spec)
		if err != nil {
			return ctrl.Result{}, err
		}
		if tnp.ObjectMeta.Annotations == nil {
			tnp.ObjectMeta.Annotations = map[string]string{
				"oldSpec": string(nowSpec),
			}
			if err := r.Update(ctx, &tnp); err != nil {
				return ctrl.Result{}, err
			}
			if err := r.updateNetworkPolicy(ctx, req.Name, req.Namespace, tnp.Spec); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
		} else if tnp.ObjectMeta.Annotations["oldSpec"] != string(nowSpec) {
			tnp.ObjectMeta.Annotations["oldSpec"] = string(nowSpec)
			if err := r.Update(ctx, &tnp); err != nil {
				return ctrl.Result{}, err
			}
			if err := r.updateNetworkPolicy(ctx, req.Name, req.Namespace, tnp.Spec); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
		}

		// in case of misoperating k8s network policy by operators of customer
		err = r.applyNetworkPolicy(ctx, req.Name, req.Namespace, tnp.Spec)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else {
		if util.ContainsString(tnp.ObjectMeta.Finalizers, MicrosegNetworkPolicyFinalizer) {
			if err := r.cleanNetworkPolicy(ctx, req.Name, req.Namespace); err != nil {
				return ctrl.Result{}, err
			}
			tnp.ObjectMeta.Finalizers = util.RemoveString(tnp.ObjectMeta.Finalizers, MicrosegNetworkPolicyFinalizer)
			if err := r.Update(ctx, &tnp); err != nil {
				return ctrl.Result{}, err
			}
		}
	}
	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
}

func (r *MicrosegNetworkPolicyReconciler) applyNetworkPolicy(ctx context.Context, name, namespace string, spec microsegv1.MicrosegNetworkPolicySpec) error {
	err := r.engine.GetNetworkPolicy(ctx, name, namespace)
	if err != nil {
		err = r.engine.CreateNetworkPolicy(ctx, name, namespace, spec)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *MicrosegNetworkPolicyReconciler) cleanNetworkPolicy(ctx context.Context, name, namespace string) error {
	if err := r.engine.DeleteNetworkPolicy(ctx, name, namespace); err != nil {
		return err
	}
	return nil
}

func (r *MicrosegNetworkPolicyReconciler) updateNetworkPolicy(ctx context.Context, name, namespace string, spec microsegv1.MicrosegNetworkPolicySpec) error {
	err := r.engine.GetNetworkPolicy(ctx, name, namespace)
	if err == nil {
		if err := r.engine.DeleteNetworkPolicy(ctx, name, namespace); err != nil {
			return err
		}
	}
	if err := r.engine.CreateNetworkPolicy(ctx, name, namespace, spec); err != nil {
		return err
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MicrosegNetworkPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&microsegv1.MicrosegNetworkPolicy{}).
		Complete(r)
}
