package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/colinbruner/argo-workflows-webhook/internal/config"
	"k8s.io/klog/v2"
)

var (
	certFile string
	keyFile  string
	port     int
)

func main() {
	flag.StringVar(&certFile, "tls-cert-file", os.Getenv("TLS_CERT_FILE"), "TLS Certificate File")
	flag.StringVar(&keyFile, "tls-key-file", os.Getenv("TLS_KEY_FILE"), "TLS Private Key File")
	flag.IntVar(&port, "port", 8443, "HTTP Server Address") // TODO: envVar with default fallback
	flag.Parse()

	cfg := config.Config{
		CertFile: certFile,
		KeyFile:  keyFile,
	}
	err := cfg.Validate()
	if err != nil {
		panic(err)
	}

	configureHandlers()

	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		TLSConfig: cfg.SetupTLS(),
	}
	klog.Info("Starting server on", server.Addr)
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		panic(err)
	}
}
