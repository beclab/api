package v1alpha1

import "testing"

func TestApplicationDeepCopyCopiesAttachments(t *testing.T) {
	app := &Application{Spec: ApplicationSpec{Attachments: []Attachment{{
		Workload: "frigate",
		Type:     "folder",
		Ref:      "drive/Home/Cameras",
		Name:     "cameras",
	}}}}

	copy := app.DeepCopy()
	copy.Spec.Attachments[0].Ref = "drive/Home/Other"
	if app.Spec.Attachments[0].Ref != "drive/Home/Cameras" {
		t.Fatalf("DeepCopy shared attachment storage with the original")
	}
}
