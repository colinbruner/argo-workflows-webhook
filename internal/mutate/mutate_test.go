package mutate

import (
	"encoding/json"
	"testing"

	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestMutateCronWorkflow(t *testing.T) {
	customResource := struct {
		metav1.ObjectMeta
		Data map[string]string
	}{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-cronworkflow",
			UID:  "123",
		},
		Data: map[string]string{},
	}

	raw, err := json.Marshal(customResource)
	if err != nil {
		t.Fatalf("Error marshalling cronworkflow: %v", err)
	}

	ar := v1.AdmissionReview{
		Request: &v1.AdmissionRequest{
			Kind: metav1.GroupVersionKind{
				Kind: "CronWorkflow",
			},
			Object: runtime.RawExtension{
				Raw: raw,
			},
		},
	}

	response := Mutate(ar)
	if !response.Allowed {
		t.Errorf("Expected response to be allowed")
	}

	expectedPatch := patchStartingDeadlineSeconds
	if string(response.Patch) != expectedPatch {
		t.Errorf("Expected patch %s, but got %s", expectedPatch, string(response.Patch))
	}
}

func TestMutateUnsupportedResource(t *testing.T) {
	cr := struct {
		metav1.ObjectMeta
		Data map[string]string
	}{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-unsupported",
		},
		Data: map[string]string{},
	}

	raw, err := json.Marshal(cr)
	if err != nil {
		t.Fatalf("Error marshalling unsupported resource: %v", err)
	}

	ar := v1.AdmissionReview{
		Request: &v1.AdmissionRequest{
			Kind: metav1.GroupVersionKind{
				Kind: "UnsupportedResource",
			},
			Object: runtime.RawExtension{
				Raw: raw,
			},
		},
	}

	response := Mutate(ar)
	if response != nil {
		t.Errorf("Expected response to nil for unsupported resource")
	}
}
