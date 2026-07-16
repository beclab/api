package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProxyListenerSpec describes a single dynamic port to expose through l4-bfl-proxy.
type ProxyListenerSpec struct {
	// Appid is the DNS subdomain prefix of the owning application, i.e.
	// md5(appName)[:8] for user apps. Combined with Port and the owner's zone it
	// yields the hostname `<appid>-<port>.<userzone>`.
	Appid string `json:"appid"`

	// Port is the TCP port the backend process listens on inside the pod. It is
	// used both as the DNS label and as the upstream port.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port"`

	// PodIP is the pod IP the backend listens on. l4-bfl-proxy dials this address
	// directly (Envoy STATIC cluster), so no Service/endpoint is required.
	PodIP string `json:"podIP"`

	// Owner is the Olares username whose zone the hostname is built under. The
	// listener is served on that user's HTTPS filter chain.
	Owner string `json:"owner"`

	// Protocol is the application-layer protocol. Only "http" is supported today.
	// +kubebuilder:validation:Enum=http
	// +kubebuilder:default=http
	Protocol string `json:"protocol,omitempty"`
}

// ProxyListenerStatus reports the generated URL and readiness of the listener.
type ProxyListenerStatus struct {
	// URL is the externally reachable hostname `<appid>-<port>.<userzone>` (no scheme).
	// +optional
	URL string `json:"url,omitempty"`
	// Ready is true once l4-bfl-proxy has programmed the virtual host.
	// +optional
	Ready bool `json:"ready,omitempty"`
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Namespaced, shortName={pl}, categories={all}
//+kubebuilder:printcolumn:JSONPath=.spec.appid, name=appid, type=string
//+kubebuilder:printcolumn:JSONPath=.spec.port, name=port, type=integer
//+kubebuilder:printcolumn:JSONPath=.spec.podIP, name=podIP, type=string
//+kubebuilder:printcolumn:JSONPath=.status.url, name=url, type=string
//+kubebuilder:printcolumn:JSONPath=.status.ready, name=ready, type=boolean
//+kubebuilder:printcolumn:JSONPath=.metadata.creationTimestamp, name=age, type=date
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProxyListener declares one dynamic port that l4-bfl-proxy should expose as an
// HTTPS virtual host `<appid>-<port>.<userzone>` routed directly to a pod IP. It
// is created by an in-pod agent (e.g. the OpenCode port-watcher) whenever a
// process starts listening on a new TCP port, and consumed by l4-bfl-proxy.
type ProxyListener struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProxyListenerSpec   `json:"spec,omitempty"`
	Status ProxyListenerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProxyListenerList contains a list of ProxyListener
type ProxyListenerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProxyListener `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProxyListener{}, &ProxyListenerList{})
}
