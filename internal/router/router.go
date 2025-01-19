package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/colinbruner/argo-workflows-webhook/internal/logger"
	"github.com/colinbruner/argo-workflows-webhook/internal/mutate"
	"github.com/colinbruner/argo-workflows-webhook/internal/scheme"
	v1 "k8s.io/api/admission/v1"
)

// admit is a generic function signature for a webhook (validate or mutate)
type admit func(v1.AdmissionReview) *v1.AdmissionResponse

// admitHandler is a handler, that CAN support multiple admission review versions
type admitHandler struct {
	v1 admit
}

// newAdminHandler returns a new admithandler for the given function
func newAdmitHandler(f admit) admitHandler {
	return admitHandler{
		v1: f,
	}
}

func serve(w http.ResponseWriter, r *http.Request, admit admitHandler) {
	logger.Debug("Reading body in request")
	var body []byte
	if r.Body != nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	logger.Debug("Reading headers in request")
	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		logger.Error(fmt.Sprintf("contentType=%s, expect application/json", contentType))
		return
	}

	logger.Debug(fmt.Sprintf("handling request body: %s", body))
	deserializer := scheme.Codecs.UniversalDeserializer()
	obj, gvk, err := deserializer.Decode(body, nil, nil)
	if err != nil {
		msg := fmt.Sprintf("Request could not be decoded: %v", err)
		logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	requestedAdmissionReview, ok := obj.(*v1.AdmissionReview)
	if !ok {
		msg := fmt.Sprintf("Expected v1.AdmissionReview but got: %T", obj)
		logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	responseAdmissionReview := &v1.AdmissionReview{}
	responseAdmissionReview.SetGroupVersionKind(*gvk)
	responseAdmissionReview.Response = admit.v1(*requestedAdmissionReview)
	responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID

	responseObj := responseAdmissionReview

	logger.Debug(fmt.Sprintf("sending response: %v", responseObj)) // TODO: Probably dont send over wire
	respBytes, err := json.Marshal(responseObj)
	if err != nil {
		logger.Error(fmt.Sprintf("Error marshalling JSON: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(respBytes); err != nil {
		logger.Error(fmt.Sprintf("Error writing HTTP response: %s", err))
	}
}

func ServeVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "v0.0.1") // TODO: read version from build
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func ServeMutate(w http.ResponseWriter, r *http.Request) {
	logger.Info("Request received for /mutate")
	serve(w, r, newAdmitHandler(mutate.Mutate))
}

/*
func ServeValidate(w http.ResponseWriter, r *http.Request) {
	klog.Info("Request received for /validate")
	serve(w, r, newAdmitHandler(validate.Validate))
}
*/
