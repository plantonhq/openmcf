package backendconfig

import (
	"testing"

	awsvpcv1 "github.com/plantonhq/project-planton/apis/org/project_planton/provider/aws/awsvpc/v1"
	"github.com/plantonhq/project-planton/apis/org/project_planton/shared"
	"github.com/plantonhq/project-planton/pkg/iac/tofu/tofulabels"
	"github.com/stretchr/testify/assert"
)

func TestExtractFromManifest_TerraformProvisioner(t *testing.T) {
	tests := []struct {
		name      string
		manifest  *awsvpcv1.AwsVpc
		want      *TofuBackendConfig
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid s3 backend with terraform labels",
			manifest: &awsvpcv1.AwsVpc{
				Metadata: &shared.CloudResourceMetadata{
					Labels: map[string]string{
						tofulabels.BackendTypeLabelKey("terraform"):   "s3",
						tofulabels.BackendObjectLabelKey("terraform"): "my-terraform-state/aws-vpc/dev",
					},
				},
			},
			want: &TofuBackendConfig{
				BackendType:   "s3",
				BackendObject: "my-terraform-state/aws-vpc/dev",
			},
			wantError: false,
		},
		{
			name: "valid gcs backend with terraform labels",
			manifest: &awsvpcv1.AwsVpc{
				Metadata: &shared.CloudResourceMetadata{
					Labels: map[string]string{
						tofulabels.BackendTypeLabelKey("terraform"):   "gcs",
						tofulabels.BackendObjectLabelKey("terraform"): "my-gcs-bucket/terraform/state",
					},
				},
			},
			want: &TofuBackendConfig{
				BackendType:   "gcs",
				BackendObject: "my-gcs-bucket/terraform/state",
			},
			wantError: false,
		},
		{
			name: "valid azurerm backend with terraform labels",
			manifest: &awsvpcv1.AwsVpc{
				Metadata: &shared.CloudResourceMetadata{
					Labels: map[string]string{
						tofulabels.BackendTypeLabelKey("terraform"):   "azurerm",
						tofulabels.BackendObjectLabelKey("terraform"): "my-container/terraform/state",
					},
				},
			},
			want: &TofuBackendConfig{
				BackendType:   "azurerm",
				BackendObject: "my-container/terraform/state",
			},
			wantError: false,
		},
		{
			name: "valid local backend with terraform labels",
			manifest: &awsvpcv1.AwsVpc{
				Metadata: &shared.CloudResourceMetadata{
					Labels: map[string]string{
						tofulabels.BackendTypeLabelKey("terraform"):   "local",
						tofulabels.BackendObjectLabelKey("terraform"): "/tmp/terraform.tfstate",
					},
				},
			},
			want: &TofuBackendConfig{
				BackendType:   "local",
				BackendObject: "/tmp/terraform.tfstate",
			},
			wantError: false,
		},
		{
			name: "no backend labels - returns nil without error",
			manifest: &awsvpcv1.AwsVpc{
				Metadata: &shared.CloudResourceMetadata{
					Labels: map[string]string{
						"other.label": "value",
					},
				},
			},
			want:      nil,
			wantError: false,
		},
		{
			name: "missing backend object",
			manifest: &awsvpcv1.AwsVpc{
				Metadata: &shared.CloudResourceMetadata{
					Labels: map[string]string{
						tofulabels.BackendTypeLabelKey("terraform"): "s3",
						// Missing backend object
					},
				},
			},
			want:      nil,
			wantError: true,
			errorMsg:  "both",
		},
		{
			name: "missing backend type",
			manifest: &awsvpcv1.AwsVpc{
				Metadata: &shared.CloudResourceMetadata{
					Labels: map[string]string{
						tofulabels.BackendObjectLabelKey("terraform"): "my-bucket/state",
						// Missing backend type
					},
				},
			},
			want:      nil,
			wantError: true,
			errorMsg:  "both",
		},
		{
			name: "empty backend type",
			manifest: &awsvpcv1.AwsVpc{
				Metadata: &shared.CloudResourceMetadata{
					Labels: map[string]string{
						tofulabels.BackendTypeLabelKey("terraform"):   "",
						tofulabels.BackendObjectLabelKey("terraform"): "my-bucket/state",
					},
				},
			},
			want:      nil,
			wantError: true,
			errorMsg:  "cannot be empty",
		},
		{
			name: "empty backend object",
			manifest: &awsvpcv1.AwsVpc{
				Metadata: &shared.CloudResourceMetadata{
					Labels: map[string]string{
						tofulabels.BackendTypeLabelKey("terraform"):   "s3",
						tofulabels.BackendObjectLabelKey("terraform"): "",
					},
				},
			},
			want:      nil,
			wantError: true,
			errorMsg:  "cannot be empty",
		},
		{
			name: "unsupported backend type",
			manifest: &awsvpcv1.AwsVpc{
				Metadata: &shared.CloudResourceMetadata{
					Labels: map[string]string{
						tofulabels.BackendTypeLabelKey("terraform"):   "unsupported",
						tofulabels.BackendObjectLabelKey("terraform"): "some/path",
					},
				},
			},
			want:      nil,
			wantError: true,
			errorMsg:  "unsupported backend type",
		},
		{
			name: "no labels",
			manifest: &awsvpcv1.AwsVpc{
				Metadata: &shared.CloudResourceMetadata{},
			},
			want:      nil,
			wantError: true,
			errorMsg:  "no labels found in manifest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractFromManifest(tt.manifest, "terraform")

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestExtractFromManifest_TofuProvisioner(t *testing.T) {
	tests := []struct {
		name      string
		manifest  *awsvpcv1.AwsVpc
		want      *TofuBackendConfig
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid s3 backend with tofu labels",
			manifest: &awsvpcv1.AwsVpc{
				Metadata: &shared.CloudResourceMetadata{
					Labels: map[string]string{
						tofulabels.BackendTypeLabelKey("tofu"):   "s3",
						tofulabels.BackendObjectLabelKey("tofu"): "my-tofu-state/aws-vpc/dev",
					},
				},
			},
			want: &TofuBackendConfig{
				BackendType:   "s3",
				BackendObject: "my-tofu-state/aws-vpc/dev",
			},
			wantError: false,
		},
		{
			name: "valid gcs backend with tofu labels",
			manifest: &awsvpcv1.AwsVpc{
				Metadata: &shared.CloudResourceMetadata{
					Labels: map[string]string{
						tofulabels.BackendTypeLabelKey("tofu"):   "gcs",
						tofulabels.BackendObjectLabelKey("tofu"): "my-gcs-bucket/tofu/state",
					},
				},
			},
			want: &TofuBackendConfig{
				BackendType:   "gcs",
				BackendObject: "my-gcs-bucket/tofu/state",
			},
			wantError: false,
		},
		{
			name: "missing backend object with tofu labels",
			manifest: &awsvpcv1.AwsVpc{
				Metadata: &shared.CloudResourceMetadata{
					Labels: map[string]string{
						tofulabels.BackendTypeLabelKey("tofu"): "s3",
						// Missing backend object
					},
				},
			},
			want:      nil,
			wantError: true,
			errorMsg:  "both",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractFromManifest(tt.manifest, "tofu")

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestExtractFromManifest_LegacyFallback(t *testing.T) {
	tests := []struct {
		name            string
		manifest        *awsvpcv1.AwsVpc
		provisionerType string
		want            *TofuBackendConfig
		wantError       bool
		errorMsg        string
	}{
		{
			name: "tofu provisioner falls back to legacy terraform labels",
			manifest: &awsvpcv1.AwsVpc{
				Metadata: &shared.CloudResourceMetadata{
					Labels: map[string]string{
						// Using legacy terraform.* labels
						tofulabels.LegacyBackendTypeLabelKey:   "s3",
						tofulabels.LegacyBackendObjectLabelKey: "legacy-state/aws-vpc/prod",
					},
				},
			},
			provisionerType: "tofu",
			want: &TofuBackendConfig{
				BackendType:   "s3",
				BackendObject: "legacy-state/aws-vpc/prod",
			},
			wantError: false,
		},
		{
			name: "terraform provisioner uses terraform labels directly (same as legacy)",
			manifest: &awsvpcv1.AwsVpc{
				Metadata: &shared.CloudResourceMetadata{
					Labels: map[string]string{
						// terraform.* labels are both provisioner-specific AND legacy
						tofulabels.LegacyBackendTypeLabelKey:   "gcs",
						tofulabels.LegacyBackendObjectLabelKey: "terraform-state/aws-vpc/staging",
					},
				},
			},
			provisionerType: "terraform",
			want: &TofuBackendConfig{
				BackendType:   "gcs",
				BackendObject: "terraform-state/aws-vpc/staging",
			},
			wantError: false,
		},
		{
			name: "provisioner-specific labels take precedence over legacy",
			manifest: &awsvpcv1.AwsVpc{
				Metadata: &shared.CloudResourceMetadata{
					Labels: map[string]string{
						// Both tofu.* and terraform.* labels present
						tofulabels.BackendTypeLabelKey("tofu"):   "s3",
						tofulabels.BackendObjectLabelKey("tofu"): "tofu-specific-bucket/state",
						tofulabels.LegacyBackendTypeLabelKey:     "gcs",
						tofulabels.LegacyBackendObjectLabelKey:   "legacy-bucket/state",
					},
				},
			},
			provisionerType: "tofu",
			want: &TofuBackendConfig{
				BackendType:   "s3",
				BackendObject: "tofu-specific-bucket/state",
			},
			wantError: false,
		},
		{
			name: "legacy fallback with partial labels returns error",
			manifest: &awsvpcv1.AwsVpc{
				Metadata: &shared.CloudResourceMetadata{
					Labels: map[string]string{
						// Only type, missing object
						tofulabels.LegacyBackendTypeLabelKey: "s3",
					},
				},
			},
			provisionerType: "tofu",
			want:            nil,
			wantError:       true,
			errorMsg:        "both",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractFromManifest(tt.manifest, tt.provisionerType)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
