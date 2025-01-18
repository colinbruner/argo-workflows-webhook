package scheme

import (
	"testing"

	argoalphav1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	v1 "k8s.io/api/admission/v1"
	registrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestAddToScheme(t *testing.T) {
	scheme := runtime.NewScheme()
	addToScheme(scheme)

	if !scheme.IsGroupRegistered(v1.SchemeGroupVersion.Group) {
		t.Errorf("admissionv1 group not registered")
	}

	if !scheme.IsGroupRegistered(registrationv1.SchemeGroupVersion.Group) {
		t.Errorf("admissionregistrationv1 group not registered")
	}

	if !scheme.IsGroupRegistered(argoalphav1.SchemeGroupVersion.Group) {
		t.Errorf("argoalphav1 group not registered")
	}
}

func TestCodecs(t *testing.T) {
	if Codecs.UniversalDeserializer() == nil {
		t.Errorf("UniversalDeserializer is nil")
	}

	if Codecs.LegacyCodec(v1.SchemeGroupVersion) == nil {
		t.Errorf("LegacyCodec for admissionv1 is nil")
	}

	if Codecs.LegacyCodec(registrationv1.SchemeGroupVersion) == nil {
		t.Errorf("LegacyCodec for admissionregistrationv1 is nil")
	}

	if Codecs.LegacyCodec(argoalphav1.SchemeGroupVersion) == nil {
		t.Errorf("LegacyCodec for argoalphav1 is nil")
	}
}
