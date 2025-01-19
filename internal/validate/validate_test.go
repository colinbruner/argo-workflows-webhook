package validate

import (
	"testing"

	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		kind        string
		raw         []byte
		wantAllowed bool
	}{
		{
			name:        "Valid CronWorkflow",
			kind:        "CronWorkflow",
			raw:         []byte(`{"metadata": {"name": "test-cronworkflow"}, "data": {"key": "value"}}`),
			wantAllowed: true,
		},
		{
			name:        "Unsupported Resource",
			kind:        "Unsupported",
			raw:         []byte(`{"metadata": {"name": "test-unsupported"}, "data": {"key": "value"}}`),
			wantAllowed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := v1.AdmissionReview{
				Request: &v1.AdmissionRequest{
					Kind: metav1.GroupVersionKind{
						Kind: tt.kind,
					},
					Object: runtime.RawExtension{
						Raw: tt.raw,
					},
				},
			}

			got := Validate(ar)
			if (got != nil && got.Allowed) != tt.wantAllowed {
				t.Errorf("Validate() = %v, want %v", got.Allowed, tt.wantAllowed)
			}
		})
	}
}
