package module

import (
	"github.com/pkg/errors"
	kubernetesharborv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesharbor/v1"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *kubernetesharborv1.KubernetesHarborStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	//create kubernetes-provider from the credential in the stack-input
	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesProviderConfig(ctx,
		stackInput.GetProviderConfig(), "kubernetes")
	if err != nil {
		return errors.Wrap(err, "failed to create kubernetes provider")
	}

	// Conditionally create namespace based on create_namespace flag
	_, err = namespace(ctx, locals, stackInput.Target.Spec, kubernetesProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create namespace")
	}

	//deploy Harbor using helm-chart
	if err := harbor(ctx, locals, kubernetesProvider); err != nil {
		return errors.Wrap(err, "failed to create harbor helm-chart resources")
	}

	//create Harbor Core/Portal ingress resources using Gateway API
	if err := createCoreIngress(ctx, locals, kubernetesProvider); err != nil {
		return errors.Wrap(err, "failed to create harbor core ingress resources")
	}

	//create Notary ingress resources using Gateway API
	if err := createNotaryIngress(ctx, locals, kubernetesProvider); err != nil {
		return errors.Wrap(err, "failed to create notary ingress resources")
	}

	return nil
}
