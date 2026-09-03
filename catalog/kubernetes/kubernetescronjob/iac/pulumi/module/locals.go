package module

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	kubernetescronjobv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kubernetescronjob/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/kuberneteslabelkeys"
	"github.com/plantonhq/planton/pkg/kubernetes/kubernetesannotationkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	KubernetesCronJob *kubernetescronjobv1alpha1.KubernetesCronJob
	Namespace         string

	// Labels are the full governance label set stamped on every created object
	// and on the pod template of every scheduled run's Job (workload identity +
	// Planton resource tracking).
	Labels map[string]string

	// Computed satellite resource names, prefixed with metadata.name so multiple
	// instances sharing a namespace never collide.
	EnvSecretName       string
	ImagePullSecretName string
	ImagePullSecretData map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *kubernetescronjobv1alpha1.KubernetesCronJobStackInput) (*Locals, error) {
	locals := &Locals{}

	locals.KubernetesCronJob = stackInput.Target
	target := stackInput.Target

	if target.Spec.JobTemplate == nil || target.Spec.JobTemplate.Container == nil ||
		target.Spec.JobTemplate.Container.App == nil || target.Spec.JobTemplate.Container.App.Image == nil {
		return nil, errors.New("spec.job_template.container.app.image is required")
	}

	locals.Labels = map[string]string{
		"app":                            target.Metadata.Name,
		kuberneteslabelkeys.Resource:     strconv.FormatBool(true),
		kuberneteslabelkeys.ResourceName: target.Metadata.Name,
		kuberneteslabelkeys.ResourceKind: cloudresourcekind.CloudResourceKind_KubernetesCronJob.String(),
	}

	if target.Metadata.Id != "" {
		locals.Labels[kuberneteslabelkeys.ResourceId] = target.Metadata.Id
	}
	if target.Metadata.Org != "" {
		locals.Labels[kuberneteslabelkeys.Organization] = target.Metadata.Org
	}
	if target.Metadata.Env != "" {
		locals.Labels[kuberneteslabelkeys.Environment] = target.Metadata.Env
	}

	locals.Namespace = target.Spec.Namespace.GetValue()
	ctx.Export(OpNamespace, pulumi.String(locals.Namespace))
	ctx.Export(OpCronJobName, pulumi.String(target.Metadata.Name))

	// The schedule is exported so dependents and audits read the deployed truth
	// from outputs rather than from the spec.
	ctx.Export(OpSchedule, pulumi.String(target.Spec.Schedule))

	locals.EnvSecretName = fmt.Sprintf("%s-env-secrets", target.Metadata.Name)
	locals.ImagePullSecretName = fmt.Sprintf("%s-image-pull", target.Metadata.Name)

	// Docker registry credential resolution, in priority order:
	// 1. stackInput.DockerConfigJson — injected by the platform at deploy time.
	// 2. The docker-config-json-file annotation — open-source/local workflows.
	// 3. Neither — no pull secret is created (public images or SA-attached secrets).
	if stackInput.DockerConfigJson != "" {
		locals.ImagePullSecretData = map[string]string{".dockerconfigjson": stackInput.DockerConfigJson}
	} else if dockerConfigFilePath := target.Metadata.Annotations[kubernetesannotationkeys.DockerConfigJsonFileAnnotationKey]; dockerConfigFilePath != "" {
		dockerConfigJson, err := loadDockerConfigFromFile(dockerConfigFilePath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to load docker config from file specified in annotation: %s", dockerConfigFilePath)
		}
		locals.ImagePullSecretData = map[string]string{".dockerconfigjson": dockerConfigJson}
	}

	return locals, nil
}

// loadDockerConfigFromFile reads docker-config JSON from a path the OPERATOR
// supplied through the docker-config-json-file annotation. Break-glass only:
// an expert applying this module from their own laptop points it at a file
// on their own disk. The read is of the operator's machine, never of the
// module's directory, so it sits deliberately outside the module
// self-containment invariant (a module reads nothing beside itself).
func loadDockerConfigFromFile(filePath string) (string, error) {
	if strings.HasPrefix(filePath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", errors.Wrap(err, "failed to get user home directory")
		}
		filePath = filepath.Join(homeDir, filePath[2:])
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", errors.Errorf("docker config file does not exist: %s", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read docker config file: %s", filePath)
	}

	if len(content) == 0 {
		return "", errors.Errorf("docker config file is empty: %s", filePath)
	}

	var js json.RawMessage
	if err := json.Unmarshal(content, &js); err != nil {
		return "", errors.Wrapf(err, "docker config file contains invalid JSON: %s", filePath)
	}

	return string(content), nil
}
