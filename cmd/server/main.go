package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/colinbruner/argo-workflows-webhook/internal/logger"
	"github.com/colinbruner/argo-workflows-webhook/internal/router"
)

var (
	certFile string
	keyFile  string
	port     int
)

func configureHandlers() {
	http.HandleFunc("/", router.ServeIndex)
	http.HandleFunc("/version", router.ServeVersion)
	http.HandleFunc("/mutate", router.ServeMutate)
	//http.HandleFunc("/validate", router.ServeValidate) // TODO
}

func main() {
	flag.StringVar(&certFile, "tls-cert-file", os.Getenv("TLS_CERT_FILE"), "TLS Certificate File")
	flag.StringVar(&keyFile, "tls-key-file", os.Getenv("TLS_KEY_FILE"), "TLS Private Key File")
	flag.IntVar(&port, "port", 8443, "HTTP Server Address") // TODO: envVar with default fallback
	flag.Parse()

	cfg := config{
		CertFile: certFile,
		KeyFile:  keyFile,
	}
	logger.Info("Loading configuration")

	err := cfg.Validate()
	if err != nil {
		panic(err)
	}
	logger.Info("Configuration validated")

	configureHandlers()

	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		TLSConfig: cfg.SetupTLS(),
	}
	logger.Info(fmt.Sprintf("Starting server on %s", server.Addr))
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		panic(err)
	}
}
