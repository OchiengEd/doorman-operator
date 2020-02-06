package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DoormanSpec defines the desired state of Doorman
// +k8s:openapi-gen=true
type DoormanSpec struct {
	// Define number of instances to be deployed
	Replicas int32 `json:"replicas,omitempty"`
	// Define database config;
	// Database size should only be provided for operator managed database
	Storage DoormanDatabaseSpec `json:"database,omitempty"`
}

// DoormanDatabaseSpec defines the storage of Doorman persistent data
type DoormanDatabaseSpec struct {
	Username       string `json:"username,omitempty"`
	PasswordLength int    `json:"password_length,omitempty"`
	DatabaseName   string `json:"name,omitempty"`
	Size           int32  `json:"size,omitempty"`
}

// DoormanStatus defines the observed state of Doorman
// +k8s:openapi-gen=true
type DoormanStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Doorman is the Schema for the doormen API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=doormen,scope=Namespaced
type Doorman struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DoormanSpec   `json:"spec,omitempty"`
	Status DoormanStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DoormanList contains a list of Doorman
type DoormanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Doorman `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Doorman{}, &DoormanList{})
}
