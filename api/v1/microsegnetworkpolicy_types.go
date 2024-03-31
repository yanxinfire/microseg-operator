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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type MicrosegNetworkPolicyIngress struct {
	ResourceSelector  map[string]string `json:"resourceSelector,omitempty"`
	NamespaceSelector map[string]string `json:"namespaceSelector,omitempty"`
	IPBlocks          []string          `json:"ipBlock,omitempty"`
	Action            string            `json:"action"`
	Protocol          string            `json:"protocol,omitempty"`
	Ports             string            `json:"ports,omitempty"`
}
type MicrosegNetworkPolicyEgress struct {
	ResourceSelector  map[string]string `json:"resourceSelector,omitempty"`
	NamespaceSelector map[string]string `json:"namespaceSelector,omitempty"`
	IPBlocks          []string          `json:"ipBlock,omitempty"`
	Action            string            `json:"action"`
	Protocol          string            `json:"protocol,omitempty"`
	Ports             string            `json:"ports,omitempty"`
}

type MicrosegNetworkPolicyRule struct {
	Ingress []MicrosegNetworkPolicyIngress `json:"ingress,omitempty"`
	Egress  []MicrosegNetworkPolicyEgress  `json:"egress,omitempty"`
}

// MicrosegNetworkPolicySpec defines the desired state of MicrosegNetworkPolicy
type MicrosegNetworkPolicySpec struct {
	ResourceSelector  map[string]string         `json:"resourceSelector,omitempty"`
	NamespaceSelector map[string]string         `json:"namespaceSelector,omitempty"`
	Order             *int                      `json:"order,omitempty"`
	Rules             MicrosegNetworkPolicyRule `json:"rules,omitempty"`
	PolicyTypes       []string                  `json:"policyTypes,omitempty"`
}

// MicrosegNetworkPolicyStatus defines the observed state of MicrosegNetworkPolicy
type MicrosegNetworkPolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MicrosegNetworkPolicy is the Schema for the microsegnetworkpolicies API
type MicrosegNetworkPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MicrosegNetworkPolicySpec   `json:"spec,omitempty"`
	Status MicrosegNetworkPolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MicrosegNetworkPolicyList contains a list of MicrosegNetworkPolicy
type MicrosegNetworkPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MicrosegNetworkPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MicrosegNetworkPolicy{}, &MicrosegNetworkPolicyList{})
}
