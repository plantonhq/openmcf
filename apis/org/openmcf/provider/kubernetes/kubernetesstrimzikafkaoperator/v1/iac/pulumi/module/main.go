package module

import (
	"github.com/pkg/errors"
	kubernetesstrimzikafkaoperatorv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesstrimzikafkaoperator/v1"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the entry‑point called by the OpenMCF engine.
func Resources(
	ctx *pulumi.Context,
	stackInput *kubernetesstrimzikafkaoperatorv1.KubernetesStrimziKafkaOperatorStackInput,
) error {
	// ------------------------------------------------------------------
	// Provider set‑up
	// ------------------------------------------------------------------
	k8sProvider, err := pulumikubernetesprovider.GetWithKubernetesProviderConfig(
		ctx,
		stackInput.ProviderConfig,
		"kubernetes",
	)
	if err != nil {
		return errors.Wrap(err, "failed to set up Kubernetes provider")
	}

	// ------------------------------------------------------------------
	// Helm install
	// ------------------------------------------------------------------
	if err := kafkaOperator(ctx, stackInput.Target, k8sProvider); err != nil {
		return errors.Wrap(err, "failed to install Kafka operator resources")
	}

	return nil
}
