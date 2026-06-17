package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

// localScheme/codecs/parameterCodec back the hand-rolled typed REST client in
// this package. They are intentionally package-private to avoid colliding with
// callers that also import "k8s.io/client-go/kubernetes/scheme".
var (
	localScheme         = runtime.NewScheme()
	localCodecs         = serializer.NewCodecFactory(localScheme)
	localParameterCodec = runtime.NewParameterCodec(localScheme)
)

func init() {
	metav1.AddToGroupVersion(localScheme, SchemeGroupVersion)
	_ = AddToScheme(localScheme)
}
