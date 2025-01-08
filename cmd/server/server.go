package server

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	certFile     string
	keyFile      string
	port         int
	sidecarImage string
)

// CmdWebhook is used by agnhost Cobra.
var CmdWebhook = &cobra.Command{
	Use:   "webhook",
	Short: "Starts a HTTP server, useful for testing MutatingAdmissionWebhook and ValidatingAdmissionWebhook",
	Long: `Starts a HTTP server, useful for testing MutatingAdmissionWebhook and ValidatingAdmissionWebhook.
After deploying it to Kubernetes cluster, the Administrator needs to create a ValidatingWebhookConfiguration
in the Kubernetes cluster to register remote webhook admission controllers.`,
	Args: cobra.MaximumNArgs(0),
	Run:  main,
}

func init() {
	CmdWebhook.Flags().StringVar(&certFile, "tls-cert-file", "",
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated after server cert).")
	CmdWebhook.Flags().StringVar(&keyFile, "tls-private-key-file", "",
		"File containing the default x509 private key matching --tls-cert-file.")
	CmdWebhook.Flags().IntVar(&port, "port", 443,
		"Secure port that the webhook listens on")
	CmdWebhook.Flags().StringVar(&sidecarImage, "sidecar-image", "",
		"Image to be used as the injected sidecar")
}

func main(cmd *cobra.Command, args []string) {
	config := Config{
		CertFile: certFile,
		KeyFile:  keyFile,
	}

	http.HandleFunc("/readyz", func(w http.ResponseWriter, req *http.Request) { w.Write([]byte("ok")) })
	//http.HandleFunc("/version", serveVersion)
	//http.HandleFunc("/mutate", serveMutate)
	//http.HandleFunc("/validate", serveValidate)

	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		TLSConfig: configTLS(config),
	}
	err := server.ListenAndServeTLS("", "")
	if err != nil {
		panic(err)
	}
}
