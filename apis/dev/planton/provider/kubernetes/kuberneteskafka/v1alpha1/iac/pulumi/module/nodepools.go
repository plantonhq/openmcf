package module

import (
	kubernetesprovider "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	kuberneteskafkav1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kuberneteskafka/v1alpha1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createNodePools renders one kafka.strimzi.io/v1 KafkaNodePool per
// spec.node_pools entry — untyped CustomResources for the same reason as
// the Kafka CR (see kafka.go). Each pool binds to the cluster through the
// strimzi.io/cluster label; the resource name is the pool name (the same
// key the Terraform module's for_each uses, so import IDs derive blind).
//
// The pools are created BEFORE the Kafka CR in the dependency graph:
// Strimzi tolerates either order, but a Kafka CR with no matching pools
// reports a transient warning state the lanes would otherwise race.
func createNodePools(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) ([]pulumi.Resource, error) {
	created := make([]pulumi.Resource, 0, len(locals.Spec.GetNodePools()))
	for _, pool := range locals.Spec.GetNodePools() {
		// The binding label rides ON TOP of the identity labels — twin
		// of TF's merge(local.labels, {cluster label}).
		labels := make(map[string]string, len(locals.Labels)+1)
		for k, v := range locals.Labels {
			labels[k] = v
		}
		labels[vars.ClusterLabelKey] = locals.ClusterName

		specBody := map[string]interface{}{
			"replicas": int(pool.GetReplicas()),
			"roles":    stringSliceToInterface(pool.GetRoles()),
			"storage":  storageBody(pool.GetStorage()),
		}
		if resources := resourcesMap(pool.GetResources()); resources != nil {
			specBody["resources"] = resources
		}
		if template := poolTemplateBody(pool); template != nil {
			specBody["template"] = template
		}

		createdPool, err := apiextensions.NewCustomResource(ctx, pool.GetName(),
			&apiextensions.CustomResourceArgs{
				ApiVersion: pulumi.String(vars.ApiVersion),
				Kind:       pulumi.String("KafkaNodePool"),
				Metadata: &kubernetesmeta.ObjectMetaArgs{
					Name:      pulumi.String(pool.GetName()),
					Namespace: pulumi.String(locals.Namespace),
					Labels:    pulumi.ToStringMap(labels),
				},
				OtherFields: kubernetes.UntypedArgs{
					"spec": specBody,
				},
			}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
		if err != nil {
			return nil, err
		}
		created = append(created, createdPool)
	}
	return created, nil
}

// storageBody is the twin of TF's pool storage rendering. The three
// shapes: persistent-claim (single volume — size required, class/
// deleteClaim optional), ephemeral (emptyDir), jbod (per-volume bodies,
// each persistent-claim with its own size/class/deleteClaim and at most
// one kraftMetadata: shared marker).
func storageBody(storage *kuberneteskafkav1alpha1.KubernetesKafkaStorage) map[string]interface{} {
	storageType := storage.GetType()
	if storageType == "" {
		storageType = "persistent-claim"
	}

	switch storageType {
	case "ephemeral":
		return map[string]interface{}{"type": "ephemeral"}
	case "jbod":
		volumes := make([]interface{}, 0, len(storage.GetVolumes()))
		for _, volume := range storage.GetVolumes() {
			volumes = append(volumes, jbodVolumeBody(volume))
		}
		return map[string]interface{}{
			"type":    "jbod",
			"volumes": volumes,
		}
	default: // persistent-claim
		body := map[string]interface{}{
			"type": "persistent-claim",
			"size": storage.GetSize(),
		}
		if class := storage.GetStorageClass().GetValue(); class != "" {
			body["class"] = class
		}
		if storage.GetDeleteClaim() {
			body["deleteClaim"] = true
		}
		return body
	}
}

func jbodVolumeBody(volume *kuberneteskafkav1alpha1.KubernetesKafkaStorageVolume) map[string]interface{} {
	body := map[string]interface{}{
		"id":   int(volume.GetId()),
		"type": "persistent-claim",
		"size": volume.GetSize(),
	}
	if class := volume.GetStorageClass().GetValue(); class != "" {
		body["class"] = class
	}
	if volume.GetDeleteClaim() {
		body["deleteClaim"] = true
	}
	if volume.GetKraftMetadata() {
		body["kraftMetadata"] = "shared"
	}
	return body
}

// poolTemplateBody renders the pool's scheduling knobs into Strimzi's pod
// template. The Strimzi pod template carries affinity and tolerations but
// NO nodeSelector — a node_selector map therefore translates to a
// requiredDuringSchedulingIgnoredDuringExecution nodeAffinity with one
// matchExpressions entry per label (semantically identical for exact-match
// selection; the Terraform module renders the same translation).
func poolTemplateBody(pool *kuberneteskafkav1alpha1.KubernetesKafkaNodePool) map[string]interface{} {
	podBody := map[string]interface{}{}

	if len(pool.GetTolerations()) > 0 {
		podBody["tolerations"] = tolerationsSlice(pool.GetTolerations())
	}

	if len(pool.GetNodeSelector()) > 0 {
		matchExpressions := make([]interface{}, 0, len(pool.GetNodeSelector()))
		// Sorted iteration keeps the rendered CR deterministic across
		// runs (Go map order is random; TF for-expressions sort keys).
		for _, key := range sortedKeys(pool.GetNodeSelector()) {
			matchExpressions = append(matchExpressions, map[string]interface{}{
				"key":      key,
				"operator": "In",
				"values":   []interface{}{pool.GetNodeSelector()[key]},
			})
		}
		podBody["affinity"] = map[string]interface{}{
			"nodeAffinity": map[string]interface{}{
				"requiredDuringSchedulingIgnoredDuringExecution": map[string]interface{}{
					"nodeSelectorTerms": []interface{}{
						map[string]interface{}{
							"matchExpressions": matchExpressions,
						},
					},
				},
			},
		}
	}

	if len(podBody) == 0 {
		return nil
	}
	return map[string]interface{}{
		"pod": podBody,
	}
}

func tolerationsSlice(tolerations []*kubernetesprovider.WorkloadToleration) []interface{} {
	out := make([]interface{}, 0, len(tolerations))
	for _, t := range tolerations {
		tol := map[string]interface{}{}
		if t.GetKey() != "" {
			tol["key"] = t.GetKey()
		}
		if t.GetOperator() != "" {
			tol["operator"] = t.GetOperator()
		}
		if t.GetValue() != "" {
			tol["value"] = t.GetValue()
		}
		if t.GetEffect() != "" {
			tol["effect"] = t.GetEffect()
		}
		if t.TolerationSeconds != nil {
			tol["tolerationSeconds"] = t.GetTolerationSeconds()
		}
		out = append(out, tol)
	}
	return out
}

func sortedKeys(in map[string]string) []string {
	keys := make([]string, 0, len(in))
	for k := range in {
		keys = append(keys, k)
	}
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[j] < keys[i] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	return keys
}
