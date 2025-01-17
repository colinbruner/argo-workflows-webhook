package main

import (
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  config
		wantErr bool
	}{
		{
			name: "valid config",
			config: config{
				CertFile: "server.crt",
				KeyFile:  "server.key",
			},
			wantErr: false,
		},
		{
			name: "missing cert file",
			config: config{
				CertFile: "",
				KeyFile:  "server.key",
			},
			wantErr: true,
		},
		{
			name: "missing key file",
			config: config{
				CertFile: "server.crt",
				KeyFile:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetupTLS(t *testing.T) {
	// Relative file reference
	cfg := config{
		CertFile: "../../testdata/server.crt",
		KeyFile:  "../../testdata/server.key",
	}

	tlsConfig := cfg.SetupTLS()
	if tlsConfig == nil {
		t.Fatal("Expected non-nil tls.Config")
	}

	if len(tlsConfig.Certificates) != 1 {
		t.Fatalf("Expected 1 certificate, got %d", len(tlsConfig.Certificates))
	}
}
