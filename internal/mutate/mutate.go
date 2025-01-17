package mutate

import (
	"encoding/json"
	"net/http"

	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

const (
	customResourcePatch1 string = `[
         { "op": "add", "path": "/spec/startingDeadlineSeconds", "value": 300 }
     ]`
)

// ServeMutate handles the /mutate/{resource} endpoint
func Mutate(ar v1.AdmissionReview) *v1.AdmissionResponse {
	cr := struct {
		metav1.ObjectMeta
		Data map[string]string
	}{}

	raw := ar.Request.Object.Raw
	err := json.Unmarshal(raw, &cr)
	if err != nil {
		klog.Error("Error unmarshalling request", err)
	}

	switch ar.Request.Kind.Kind {
	case "CronWorkflow":
		klog.Info("Handling cronworkflows resource")
		return mutateCronWorkflows()
	case "Workflow":
		//TODO
		klog.Info("Handling workflows resource")
		return nil
	case "WorkflowTemplate":
		//TODO
		klog.Info("Handling workflowtemplate resource")
		return nil
	default:
		// TODO: klog error?
		klog.Error("Unsupported resource", http.StatusBadRequest)
		return nil
	}

}

// func mutateCronWorkflows(ar v1.AdmissionReview) *v1.AdmissionResponse {
func mutateCronWorkflows() *v1.AdmissionResponse {
	klog.Info("Mutating cronworkflows")

	patchType := v1.PatchTypeJSONPatch

	return &v1.AdmissionResponse{
		Allowed:   true,
		PatchType: &patchType,
		Patch:     []byte(customResourcePatch1),
	}
}