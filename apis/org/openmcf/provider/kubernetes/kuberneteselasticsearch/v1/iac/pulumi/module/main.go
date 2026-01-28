package module

import (
	"github.com/pkg/errors"
	kuberneteselasticsearchv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kuberneteselasticsearch/v1"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *kuberneteselasticsearchv1.KubernetesElasticsearchStackInput) error {
	locals := initializeLocals(ctx, stackInput)
	//create kubernetes-provider from the credential in the stack-input
	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesProviderConfig(ctx,
		stackInput.ProviderConfig, "kubernetes")
	if err != nil {
		return errors.Wrap(err, "failed to setup gcp provider")
	}

	// Conditionally create namespace based on create_namespace flag
	if stackInput.Target.Spec.CreateNamespace {
		createdNamespace, err := kubernetescorev1.NewNamespace(ctx, locals.Namespace,
			&kubernetescorev1.NamespaceArgs{
				Metadata: metav1.ObjectMetaPtrInput(
					&metav1.ObjectMetaArgs{
						Name:   pulumi.String(locals.Namespace),
						Labels: pulumi.ToStringMap(locals.Labels),
					}),
			}, pulumi.Provider(kubernetesProvider))
		if err != nil {
			return errors.Wrapf(err, "failed to create namespace")
		}
		//export name of the namespace
		ctx.Export(OpNamespace, createdNamespace.Metadata.Name().Elem())
	} else {
		// Use existing namespace - just reference it by name
		//export name of the namespace
		ctx.Export(OpNamespace, pulumi.String(locals.Namespace))
	}

	if err := elasticsearch(ctx, locals, kubernetesProvider); err != nil {
		return errors.Wrap(err, "failed to create elastic search resources")
	}

	if (locals.KubernetesElasticsearch.Spec.Elasticsearch.Ingress != nil &&
		locals.KubernetesElasticsearch.Spec.Elasticsearch.Ingress.Enabled) ||
		(locals.KubernetesElasticsearch.Spec.Kibana != nil &&
			locals.KubernetesElasticsearch.Spec.Kibana.Enabled &&
			locals.KubernetesElasticsearch.Spec.Kibana.Ingress != nil &&
			locals.KubernetesElasticsearch.Spec.Kibana.Ingress.Enabled) {
		if err := ingress(ctx, locals, kubernetesProvider); err != nil {
			return errors.Wrap(err, "failed to create ingress resources")
		}
	}

	return nil
}
