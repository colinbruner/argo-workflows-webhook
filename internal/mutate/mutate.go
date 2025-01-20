package mutate

import (
	"encoding/json"
	"fmt"

	"github.com/colinbruner/argo-workflows-webhook/internal/logger"
	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// NOTE: Could dynamically look up aspects of this from ConfigMap, such as value?
	patchStartingDeadlineSeconds string = `[
         { "op": "add", "path": "/spec/startingDeadlineSeconds", "value": 300 }
     ]`
)

// ServeMutate handles the /mutate/ endpoint
func Mutate(ar v1.AdmissionReview) *v1.AdmissionResponse {
	customResource := struct {
		metav1.ObjectMeta
		Data map[string]string
	}{}

	raw := ar.Request.Object.Raw
	err := json.Unmarshal(raw, &customResource)
	if err != nil {
		logger.Error(fmt.Sprintf("Error unmarshalling request: %s", err))
	}

	switch ar.Request.Kind.Kind {
	case "CronWorkflow":
		logger.Debug("Handling cronworkflow resource")
		return mutateCronWorkflows()
	case "Workflow":
		//TODO
		logger.Debug("Handling workflow resource")
		return nil
	case "WorkflowTemplate":
		//TODO
		logger.Debug("Handling workflowtemplate resource")
		return nil
	default:
		logger.Error(fmt.Sprintf("Unsupported resource: %s", ar.Request.Kind.Kind))
		return nil
	}

}

func mutateCronWorkflows() *v1.AdmissionResponse {
	logger.Debug("Mutating cronworkflows")

	patchType := v1.PatchTypeJSONPatch

	return &v1.AdmissionResponse{
		Allowed:   true,
		PatchType: &patchType,
		Patch:     []byte(patchStartingDeadlineSeconds),
	}
}
