package config

import (
	"crypto/tls"
	"fmt"

	"k8s.io/klog/v2"
)

type Config struct {
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

func (c *Config) Validate() error {
	if c.CertFile == "" {
		return fmt.Errorf("must specify config key tls: cert_file")
	}
	if c.KeyFile == "" {
		return fmt.Errorf("must specify config key tls: key_file")
	}

	return nil
}

func (c Config) SetupTLS() *tls.Config {
	sCert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if err != nil {
		klog.Fatal(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{sCert},
		// TODO: uses mutual tls after we agree on what cert the apiserver should use.
		// ClientAuth:   tls.RequireAndVerifyClientCert,
	}
}
