package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// databaseFirewall provisions the cluster's inbound trusted-sources rule
// set and exports its outputs.
//
// The spec's five typed lists fan out to the provider's polymorphic
// {type, value} rows, exactly mirroring the Terraform module's locals. The
// type tokens are the PROVIDER's values (ip_addr, droplet, k8s, app, tag)
// -- the SDK's doc comment renders "ipAddr", but that is bridged-docs
// camelization; values pass through to the provider's own validator
// untranslated, and "ipAddr" would be rejected.
//
// Destroy semantics: deleting this resource PUTs an EMPTY rule list -- the
// cluster then accepts connections from anywhere again (there is no
// object to 404 on).
func databaseFirewall(
	ctx *pulumi.Context,
	locals *Locals,
	digitalOceanProvider *digitalocean.Provider,
) (*digitalocean.DatabaseFirewall, error) {
	spec := locals.DigitalOceanDatabaseFirewall.Spec

	var rules digitalocean.DatabaseFirewallRuleArray
	appendRules := func(ruleType string, values []string) {
		for _, v := range values {
			rules = append(rules, digitalocean.DatabaseFirewallRuleArgs{
				Type:  pulumi.String(ruleType),
				Value: pulumi.String(v),
			})
		}
	}

	appendRules("ip_addr", spec.IpRules)

	// References are resolved to literal ids before the module runs.
	dropletIds := make([]string, 0, len(spec.DropletIds))
	for _, ref := range spec.DropletIds {
		dropletIds = append(dropletIds, ref.GetValue())
	}
	appendRules("droplet", dropletIds)

	clusterIds := make([]string, 0, len(spec.KubernetesClusterIds))
	for _, ref := range spec.KubernetesClusterIds {
		clusterIds = append(clusterIds, ref.GetValue())
	}
	appendRules("k8s", clusterIds)

	appIds := make([]string, 0, len(spec.AppIds))
	for _, ref := range spec.AppIds {
		appIds = append(appIds, ref.GetValue())
	}
	appendRules("app", appIds)

	appendRules("tag", spec.Tags)

	createdFirewall, err := digitalocean.NewDatabaseFirewall(
		ctx,
		"firewall",
		&digitalocean.DatabaseFirewallArgs{
			// References are resolved to the literal cluster UUID before
			// the module runs.
			ClusterId: pulumi.String(spec.Cluster.GetValue()),
			Rules:     rules,
		},
		pulumi.Provider(digitalOceanProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create digitalocean database firewall")
	}

	ctx.Export(OpClusterId, createdFirewall.ClusterId)

	return createdFirewall, nil
}
