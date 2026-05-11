package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

//+genclient
//+genclient:nonNamespaced
//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster,shortName={aag}, categories={all}
//+kubebuilder:printcolumn:JSONPath=.spec.appName, name=application name, type=string
//+kubebuilder:printcolumn:JSONPath=.spec.user, name=user, type=string
//+kubebuilder:printcolumn:JSONPath=.spec.grantedBy, name=grantedBy, type=string
//+kubebuilder:printcolumn:JSONPath=.status.state, name=state, type=string
//+kubebuilder:printcolumn:JSONPath=.metadata.creationTimestamp, name=age, type=date
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppAccessGrant grants a non-admin user permission to open a running
// application that lives in their own user-space. Owner/admin users are
// implicitly granted by their role and do not need this CR.
type AppAccessGrant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:",inline"`
	Spec              AppAccessGrantSpec   `json:"spec"`
	Status            AppAccessGrantStatus `json:"status,omitempty"`
}

type AppAccessGrantSpec struct {
	AppName      string `json:"appName"`      // matches Application.Spec.Name
	AppNamespace string `json:"appNamespace"` // matches Application.Spec.Namespace
	User         string `json:"user"`         // the normal user being granted
	GrantedBy    string `json:"grantedBy"`    // owner/admin username
}

type AppAccessGrantStatus struct {
	// pending, active
	State     string      `json:"state"`
	GrantedAt metav1.Time `json:"grantedAt,omitempty"`
}

//+kubebuilder:object:root=true
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppAccessGrantList containers a list of AppAccessGrant
type AppAccessGrantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppAccessGrant `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AppAccessGrant{}, &AppAccessGrantList{})
}
