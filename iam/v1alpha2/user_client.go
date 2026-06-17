package v1alpha2

import (
	"context"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/gentype"
	"k8s.io/client-go/rest"
)

// UserInterface has methods to work with User resources. It mirrors the shape
// of client-gen output so this hand-rolled client can be swapped for a
// generated one in the future without breaking callers.
type UserInterface interface {
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*User, error)
}

// Client holds typed clients for the iam.kubesphere.io/v1alpha2 API group.
//
// We hand-roll this instead of using k8s code-generator because the iam
// package lives at iam/v1alpha2/ rather than under api/<group>/<version>/.
// kube_codegen.sh discovers groups by walking <in_dir>/<group>/<version>/
// and we deliberately keep this import path stable for downstream consumers.
type Client struct {
	Users UserInterface
}

// NewClient creates a typed client for iam.kubesphere.io/v1alpha2.
func NewClient(cfg *rest.Config) (*Client, error) {
	config := *cfg
	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}
	gv := SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = localCodecs.WithoutConversion()

	httpClient, err := rest.HTTPClientFor(&config)
	if err != nil {
		return nil, err
	}
	restClient, err := rest.RESTClientForConfigAndClient(&config, httpClient)
	if err != nil {
		return nil, err
	}
	return &Client{Users: newUsers(restClient)}, nil
}

// NewClientForConfigAndClient is like NewClient but accepts a pre-built http client.
func NewClientForConfigAndClient(cfg *rest.Config, httpClient *http.Client) (*Client, error) {
	config := *cfg
	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}
	gv := SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = localCodecs.WithoutConversion()

	restClient, err := rest.RESTClientForConfigAndClient(&config, httpClient)
	if err != nil {
		return nil, err
	}
	return &Client{Users: newUsers(restClient)}, nil
}

type users struct {
	*gentype.Client[*User]
}

func newUsers(c rest.Interface) *users {
	return &users{
		gentype.NewClient[*User](
			"users",
			c,
			localParameterCodec,
			"",
			func() *User { return &User{} },
		),
	}
}
