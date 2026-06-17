package v1alpha2

import (
	"context"
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type fakeUsers struct {
	user *User
	err  error
}

func (f *fakeUsers) Get(_ context.Context, _ string, _ metav1.GetOptions) (*User, error) {
	return f.user, f.err
}

func TestGetUserZone(t *testing.T) {
	user := &User{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "olaresid",
			Annotations: map[string]string{UserAnnotationZoneKey: "olares.cn"},
		},
	}
	got, err := GetUserZone(context.Background(), &fakeUsers{user: user}, "olaresid")
	if err != nil {
		t.Fatalf("GetUserZone() error = %v", err)
	}
	if got != "olares.cn" {
		t.Fatalf("GetUserZone() = %q, want %q", got, "olares.cn")
	}
}

func TestUserAnnotationMissing(t *testing.T) {
	user := &User{ObjectMeta: metav1.ObjectMeta{Name: "alice"}}
	_, err := UserAnnotation(user, UserAnnotationZoneKey)
	if err == nil {
		t.Fatal("UserAnnotation() expected error for missing annotation")
	}
}

func TestGetUserAnnotation(t *testing.T) {
	user := &User{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{"custom.example/key": "value"},
		},
	}
	got, err := GetUserAnnotation(context.Background(), &fakeUsers{user: user}, "alice", "custom.example/key")
	if err != nil {
		t.Fatalf("GetUserAnnotation() error = %v", err)
	}
	if got != "value" {
		t.Fatalf("GetUserAnnotation() = %q, want %q", got, "value")
	}
}

func TestGetUserZoneUsesFakeClient(t *testing.T) {
	users := &fakeUsers{
		user: &User{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{UserAnnotationZoneKey: "demo.zone"},
			},
		},
	}
	got, err := GetUserZone(context.Background(), users, "bob")
	if err != nil {
		t.Fatalf("GetUserZone() error = %v", err)
	}
	if !reflect.DeepEqual(got, "demo.zone") {
		t.Fatalf("GetUserZone() = %q, want %q", got, "demo.zone")
	}
}
