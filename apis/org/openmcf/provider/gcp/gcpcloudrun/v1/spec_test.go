package gcpcloudrunv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
	foreignkeyv1 "github.com/plantonhq/openmcf/apis/org/openmcf/shared/foreignkey/v1"
)

func TestGcpCloudRunSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "GcpCloudRunSpec Custom Validation Tests")
}

var _ = ginkgo.Describe("GcpCloudRunSpec Custom Validation Tests", func() {

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.Context("gcp_cloud_run", func() {

			ginkgo.It("should not return a validation error for minimal valid fields", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region: "us-central1",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Port:   8080,
							Cpu:    1,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
						MaxConcurrency: 80,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return validation error with all optional fields set", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run-full",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region:         "us-central1",
						ServiceName:    "custom-service-name",
						ServiceAccount: "sa-name@test-project-123.iam.gserviceaccount.com",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "v1.0.0",
							},
							Env: &GcpCloudRunContainerEnv{
								Variables: map[string]string{
									"ENV":  "production",
									"PORT": "8080",
								},
								Secrets: map[string]string{
									"API_KEY": "projects/123/secrets/api-key:latest",
								},
							},
							Port:   8080,
							Cpu:    2,
							Memory: 1024,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 1,
								Max: 100,
							},
						},
						MaxConcurrency:       100,
						TimeoutSeconds:       600,
						Ingress:              GcpCloudRunIngress_INGRESS_TRAFFIC_ALL,
						AllowUnauthenticated: true,
						VpcAccess: &GcpCloudRunVpcAccess{
							Network: &foreignkeyv1.StringValueOrRef{
								LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "projects/test-project-123/global/networks/default"},
							},
							Subnet: &foreignkeyv1.StringValueOrRef{
								LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "projects/test-project-123/regions/us-central1/subnetworks/default"},
							},
							Egress: "ALL_TRAFFIC",
						},
						ExecutionEnvironment: GcpCloudRunExecutionEnvironment_EXECUTION_ENVIRONMENT_GEN2,
						Dns: &GcpCloudRunDns{
							Enabled:     true,
							Hostnames:   []string{"api.example.com", "api2.example.com"},
							ManagedZone: "example-zone",
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {
		ginkgo.Context("project_id validation", func() {

			ginkgo.It("should return validation error for missing project_id", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						Region: "us-central1",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Cpu:    1,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should not return validation error for project_id with valueRef", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
								ValueFrom: &foreignkeyv1.ValueFromRef{
									Name: "my-gcp-project",
								},
							},
						},
						Region: "us-central1",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Port:   8080,
							Cpu:    1,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("region validation", func() {

			ginkgo.It("should return validation error for missing region", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Cpu:    1,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return validation error for invalid region pattern", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region: "invalid-region",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Cpu:    1,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("service_name validation", func() {

			ginkgo.It("should return validation error for invalid service_name pattern", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region:      "us-central1",
						ServiceName: "Invalid_Service_Name",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Cpu:    1,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return validation error for service_name that is too long", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region:      "us-central1",
						ServiceName: "this-is-a-very-long-service-name-that-exceeds-the-maximum-allowed-length-of-63-characters",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Cpu:    1,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("service_account validation", func() {

			ginkgo.It("should return validation error for invalid service_account pattern", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region:         "us-central1",
						ServiceAccount: "invalid-service-account-email",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Cpu:    1,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("container validation", func() {

			ginkgo.It("should return validation error for missing container", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region: "us-central1",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return validation error for missing image repo", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region: "us-central1",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Tag: "1.0.0",
							},
							Cpu:    1,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return validation error for missing image tag", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region: "us-central1",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
							},
							Cpu:    1,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return validation error for invalid cpu value", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region: "us-central1",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Cpu:    3,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return validation error for memory below minimum", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region: "us-central1",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Cpu:    1,
							Memory: 100,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return validation error for memory above maximum", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region: "us-central1",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Cpu:    4,
							Memory: 40000,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return validation error for invalid port range", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region: "us-central1",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Port:   70000,
							Cpu:    1,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return validation error for missing replicas", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region: "us-central1",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Cpu:    1,
							Memory: 512,
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("max_concurrency validation", func() {

			ginkgo.It("should return validation error for max_concurrency below minimum", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region: "us-central1",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Cpu:    1,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
						MaxConcurrency: 0,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return validation error for max_concurrency above maximum", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region: "us-central1",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Cpu:    1,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
						MaxConcurrency: 1500,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("timeout_seconds validation", func() {

			ginkgo.It("should return validation error for timeout_seconds below minimum", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region: "us-central1",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Cpu:    1,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
						TimeoutSeconds: 0,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return validation error for timeout_seconds above maximum", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region: "us-central1",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Cpu:    1,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
						TimeoutSeconds: 4000,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("vpc_access validation", func() {

			ginkgo.It("should return validation error for invalid egress value", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region: "us-central1",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Cpu:    1,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
						VpcAccess: &GcpCloudRunVpcAccess{
							Network: &foreignkeyv1.StringValueOrRef{
								LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "projects/test-project-123/global/networks/default"},
							},
							Subnet: &foreignkeyv1.StringValueOrRef{
								LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "projects/test-project-123/regions/us-central1/subnetworks/default"},
							},
							Egress: "INVALID_EGRESS_VALUE",
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("dns validation", func() {

			ginkgo.It("should return validation error when dns.enabled is true but hostnames is empty", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region: "us-central1",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Cpu:    1,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
						Dns: &GcpCloudRunDns{
							Enabled:     true,
							ManagedZone: "example-zone",
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return validation error when dns.enabled is true but managed_zone is empty", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region: "us-central1",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Cpu:    1,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
						Dns: &GcpCloudRunDns{
							Enabled:   true,
							Hostnames: []string{"api.example.com"},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return validation error for invalid hostname pattern", func() {
				input := &GcpCloudRun{
					ApiVersion: "gcp.openmcf.org/v1",
					Kind:       "GcpCloudRun",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cloud-run",
					},
					Spec: &GcpCloudRunSpec{
						ProjectId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-project-123"},
						},
						Region: "us-central1",
						Container: &GcpCloudRunContainer{
							Image: &GcpCloudRunContainerImage{
								Repo: "us-docker.pkg.dev/prj/registry/app",
								Tag:  "1.0.0",
							},
							Cpu:    1,
							Memory: 512,
							Replicas: &GcpCloudRunContainerReplicas{
								Min: 0,
								Max: 10,
							},
						},
						Dns: &GcpCloudRunDns{
							Enabled:     true,
							Hostnames:   []string{"Invalid_Hostname"},
							ManagedZone: "example-zone",
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})
})
