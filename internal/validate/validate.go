package validate

import (
	"encoding/json"
	"fmt"

	"github.com/colinbruner/argo-workflows-webhook/internal/logger"
	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Validate(ar v1.AdmissionReview) *v1.AdmissionResponse {
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
		return validateCronWorkflows()
	/*
		//TODO: implement Workflow and WorkflowTemplate validation
		case "Workflow":
			logger.Debug("Handling workflow resource")
			return nil
		case "WorkflowTemplate":
			logger.Debug("Handling workflowtemplate resource")
			return nil
	*/
	default:
		logger.Error(fmt.Sprintf("Unsupported resource: %s", ar.Request.Kind.Kind))
		return nil
	}
}

func validateCronWorkflows() *v1.AdmissionResponse {
	logger.Debug("Validating cronworkflow")

	// TODO: Add logic currently just allow them all
	return &v1.AdmissionResponse{
		Allowed: true,
	}
}
