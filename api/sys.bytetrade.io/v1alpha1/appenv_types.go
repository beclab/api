package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//+genclient
//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Namespaced, shortName={appenv}, categories={all}
//+kubebuilder:printcolumn:JSONPath=.appName, name=app name, type=string
//+kubebuilder:printcolumn:JSONPath=.appOwner, name=owner, type=string
//+kubebuilder:printcolumn:JSONPath=.metadata.creationTimestamp, name=age, type=date

// AppEnv is the Schema for the application environment variables API
type AppEnv struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	AppName  string      `json:"appName" yaml:"appName" validate:"required"`
	AppOwner string      `json:"appOwner" yaml:"appOwner" validate:"required"`
	Envs     []AppEnvVar `json:"envs,omitempty" yaml:"envs,omitempty"`
	// Secrets records the app's secrets[] declarations. They are kept here,
	// alongside Envs, so the SystemEnv/UserEnv controllers can discover which
	// apps reference a changed variable and drive the same sync/apply flow.
	// Only the declaration lives here — never the resolved value.
	Secrets   []AppSecretVar `json:"secrets,omitempty" yaml:"secrets,omitempty"`
	NeedApply bool           `json:"needApply,omitempty" yaml:"needApply,omitempty"`
}

type AppEnvVar struct {
	EnvVarSpec    `json:",inline" yaml:",inline"`
	ApplyOnChange bool       `json:"applyOnChange,omitempty" yaml:"applyOnChange,omitempty"`
	ValueFrom     *ValueFrom `json:"valueFrom,omitempty" yaml:"valueFrom,omitempty"`
}

// AppSecretValueKey is the data key every generated app Secret stores its value
// under. It is a constant so chart authors can write a plain secretKeyRef
// without templating anything:
//
//	valueFrom:
//	  secretKeyRef:
//	    name: <AppSecretVar.Name>
//	    key: value
const AppSecretValueKey = "value"

// AppSecretVar declares one Olares-provided env var (a SystemEnv or UserEnv,
// referenced exactly like AppEnvVar.ValueFrom) to be materialized into a
// dedicated Kubernetes Secret in the app's namespace. The app's chart consumes
// it through a standard secretKeyRef instead of receiving the value as a
// plaintext env var, which keeps it out of the pod spec.
//
// Unlike AppEnvVar this type deliberately has NO value field, and none may be
// added: the resolved value belongs only in the Kubernetes Secret. Giving this
// struct somewhere to hold it would put the plaintext back into the AppEnv CR
// and defeat the entire point of declaring the variable as a secret.
//
// TODO(structured-secrets): one declaration produces exactly one Secret holding
// a single value under AppSecretValueKey. Once the Olares UX can author and
// manage multi-key secrets, allow several keys to share one Secret.
type AppSecretVar struct {
	// Name is the Secret's name in the app namespace, used verbatim. Chart
	// authors hardcode this exact string in their secretKeyRef.
	Name string `json:"name" yaml:"name"`

	// ValueFrom references the Olares-provided env var to pull the value from.
	// Reuses the AppEnvVar reference type so envs[] and secrets[] pull from the
	// same value space with identical syntax.
	ValueFrom *ValueFrom `json:"valueFrom,omitempty" yaml:"valueFrom,omitempty"`

	// ApplyOnChange mirrors AppEnvVar.ApplyOnChange: when true, rotating the
	// referenced variable redeploys the app so pods pick the new value up.
	// Secret values are injected at pod start, so without this the Secret is
	// refreshed but running pods keep serving the old value until they restart.
	ApplyOnChange bool `json:"applyOnChange,omitempty" yaml:"applyOnChange,omitempty"`
}

// ValueFrom defines a reference to an environment variable (UserEnv or SystemEnv)
type ValueFrom struct {
	EnvName string `json:"envName" yaml:"envName" validate:"required"`
	Status  string `json:"status,omitempty" yaml:"status,omitempty"`
}

type EnvValueOptionItem struct {
	Title string `json:"title" yaml:"title"`
	Value string `json:"value" yaml:"value"`
}

//+kubebuilder:object:root=true
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppEnvList contains a list of AppEnv
type AppEnvList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppEnv `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AppEnv{}, &AppEnvList{})
}
