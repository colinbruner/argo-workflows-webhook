package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/colinbruner/argo-workflows-webhook/internal/argo"
	"github.com/colinbruner/argo-workflows-webhook/internal/scheme"
	v1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
)

// admit is a generic function signature for an admission webhook
type admit func(v1.AdmissionReview) *v1.AdmissionResponse

// admitHandler is a handler, for both validators and mutators, that CAN support multiple admission review versions
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
	var body []byte
	if r.Body != nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		klog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	klog.V(2).Info(fmt.Sprintf("handling request: %s", body))

	deserializer := scheme.Codecs.UniversalDeserializer()
	obj, _, err := deserializer.Decode(body, nil, nil)
	if err != nil {
		msg := fmt.Sprintf("Request could not be decoded: %v", err)
		klog.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	var responseObj runtime.Object
	requestedAdmissionReview, ok := obj.(*v1.AdmissionReview)
	if !ok {
		klog.Errorf("Expected v1.AdmissionReview but got: %T", obj)
		return
	}
	responseAdmissionReview := &v1.AdmissionReview{}

	responseAdmissionReview.Response = admit.v1(*requestedAdmissionReview)
	responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID

	responseObj = responseAdmissionReview

	klog.V(2).Info(fmt.Sprintf("sending response: %v", responseObj))
	respBytes, err := json.Marshal(responseObj)
	if err != nil {
		klog.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(respBytes); err != nil {
		klog.Error(err)
	}
}

func serveVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "v0.0.1") // TODO: read version from build
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK") // TODO: usage message
}

func serveMutate(w http.ResponseWriter, r *http.Request) {
	serve(w, r, newAdmitHandler(argo.Mutate))
}

func configureHandlers() {
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/version", serveVersion)
	http.HandleFunc("/mutate", serveMutate)
}
