package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	admission "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"github.com/cloudzero/cloudzero-insights-controller/pkg/build"
)

const (
	defaultPort            = ":8443"
	tlsKeyFilePath  string = "/etc/certs/tls.key"
	tlsCertFilePath string = "/etc/certs/tls.crt"
)

const (
	httpHeaderContentType = "Content-Type"
	httpContentTypeJSON   = "application/json"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecFactory  = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecFactory.UniversalDeserializer()
)

// add kind AdmissionReview in scheme
func init() {
	_ = corev1.AddToScheme(runtimeScheme)
	_ = admission.AddToScheme(runtimeScheme)
	_ = appsv1.AddToScheme(runtimeScheme)
}

type admitv1Func func(admission.AdmissionReview) *admission.AdmissionResponse

type admitHandler struct {
	v1 admitv1Func
}

func AdmitHandler(f admitv1Func) admitHandler {
	return admitHandler{
		v1: f,
	}
}

// serve handles the http portion of a request prior to handing to an admit function
func serve(w http.ResponseWriter, r *http.Request, admit admitHandler) {
	var body []byte
	if r.Body != nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// verify the content type is accurate
	contentType := r.Header.Get(httpHeaderContentType)
	if contentType != httpContentTypeJSON {
		log.Error().Msgf("contentType=%s, expect application/json", contentType)
		return
	}

	log.Info().Msgf("handling request: %s", body)
	var responseObj runtime.Object
	if obj, gvk, err := deserializer.Decode(body, nil, nil); err != nil {
		msg := fmt.Sprintf("Request could not be decoded: %v", err)
		log.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return

	} else {
		requestedAdmissionReview, ok := obj.(*admission.AdmissionReview)
		if !ok {
			log.Error().Msgf("Expected appsv1.AdmissionReview but got: %T", obj)
			return
		}
		responseAdmissionReview := &admission.AdmissionReview{}
		responseAdmissionReview.SetGroupVersionKind(*gvk)
		responseAdmissionReview.Response = admit.v1(*requestedAdmissionReview)
		responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID
		responseObj = responseAdmissionReview

	}
	log.Info().Msgf("sending response: %v", responseObj)
	respBytes, err := json.Marshal(responseObj)
	if err != nil {
		log.Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set(httpHeaderContentType, httpContentTypeJSON)
	if _, err := w.Write(respBytes); err != nil {
		log.Err(err)
	}
}

func serveValidate(w http.ResponseWriter, r *http.Request) {
	serve(w, r, AdmitHandler(validate))
}

// verify if a Deployment has the 'prod' prefix name
func validate(ar admission.AdmissionReview) *admission.AdmissionResponse {
	log.Info().Msgf("validating deployments")
	deploymentResource := metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	if ar.Request.Resource != deploymentResource {
		log.Error().Msgf("expect resource to be %s", deploymentResource)
		return nil
	}

	deployment := appsv1.Deployment{}
	if _, _, err := deserializer.Decode(ar.Request.Object.Raw, nil, &deployment); err != nil {
		log.Err(err)
		return &admission.AdmissionResponse{Result: &metav1.Status{Message: err.Error()}}
	}

	if b, err := json.MarshalIndent(deployment, "", "  "); err == nil {
		fmt.Println(string(b))
	}

	return &admission.AdmissionResponse{Allowed: true}
}

func main() {
	var tlsKey, tlsCert string
	flag.StringVar(&tlsKey, "tlsKey", tlsKeyFilePath, "Path to the TLS key")
	flag.StringVar(&tlsCert, "tlsCert", tlsCertFilePath, "Path to the TLS certificate")
	flag.Parse()

	log.Info().Msgf("Starting CloudZero Insights Controller %s", build.GetVersion())

	http.HandleFunc("/validate", serveValidate)
	log.Info().Msg("Server started ...")

	err := http.ListenAndServeTLS(defaultPort, tlsCert, tlsKey, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}

	log.Info().Msg("Server stopped")
}
