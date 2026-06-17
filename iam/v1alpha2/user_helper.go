package v1alpha2

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	AnnotationGroup       = "bytetrade.io"
	UserAnnotationZoneKey = AnnotationGroup + "/zone"
)

// UserAnnotation returns a single annotation from a User object.
func UserAnnotation(user *User, annotation string) (string, error) {
	if user == nil {
		return "", fmt.Errorf("user is nil")
	}
	a, ok := user.GetAnnotations()[annotation]
	if !ok {
		return "", fmt.Errorf("user annotation %s not found", annotation)
	}
	return a, nil
}

// GetUserAnnotation returns user annotation, an error if there is any.
func GetUserAnnotation(ctx context.Context, users UserInterface, username, annotation string) (string, error) {
	user, err := users.Get(ctx, username, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return UserAnnotation(user, annotation)
}

// GetUserZone returns user zone, an error if there is any.
func GetUserZone(ctx context.Context, users UserInterface, username string) (string, error) {
	return GetUserAnnotation(ctx, users, username, UserAnnotationZoneKey)
}
