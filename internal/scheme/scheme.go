package scheme

// https://github.com/kubernetes/kubernetes/blob/master/test/images/agnhost/webhook/scheme.go

import (
	argoalphav1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	v1 "k8s.io/api/admission/v1"
	registrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var scheme = runtime.NewScheme()

func init() {
	addToScheme(scheme)
}

func addToScheme(scheme *runtime.Scheme) {
	utilruntime.Must(v1.AddToScheme(scheme))
	utilruntime.Must(registrationv1.AddToScheme(scheme))
	utilruntime.Must(argoalphav1.AddToScheme(scheme))
}

var Codecs = serializer.NewCodecFactory(scheme)
