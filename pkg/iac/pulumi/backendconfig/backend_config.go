package backendconfig

import (
	"fmt"
	"strings"

	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumiannotationkeys"
	"github.com/plantonhq/planton/pkg/reflection/metadatareflect"
	"google.golang.org/protobuf/proto"
)

// PulumiBackendConfig represents the Pulumi backend configuration
type PulumiBackendConfig struct {
	// StackFqdn is the fully qualified stack name (org/project/stack)
	StackFqdn string
	// Organization is the Pulumi organization name
	Organization string
	// Project is the Pulumi project name
	Project string
	// StackName is the Pulumi stack name
	StackName string
}

// ExtractFromManifest extracts Pulumi backend configuration from manifest annotations
// Priority: stack.fqdn > (organization + project + stack.name)
func ExtractFromManifest(manifest proto.Message) (*PulumiBackendConfig, error) {
	annotations := metadatareflect.ExtractAnnotations(manifest)
	if annotations == nil {
		return nil, fmt.Errorf("no annotations found in manifest")
	}

	config := &PulumiBackendConfig{}

	// First priority: Check for stack.fqdn
	if stackFqdn, ok := annotations[pulumiannotationkeys.StackFqdnAnnotationKey]; ok && stackFqdn != "" {
		config.StackFqdn = stackFqdn

		// Parse the FQDN to extract components
		org, project, stack, err := parseStackFqdn(stackFqdn)
		if err != nil {
			return nil, fmt.Errorf("invalid stack.fqdn format: %w", err)
		}

		config.Organization = org
		config.Project = project
		config.StackName = stack

		return config, nil
	}

	// Second priority: Check for individual components
	org, hasOrg := annotations[pulumiannotationkeys.OrganizationAnnotationKey]
	project, hasProject := annotations[pulumiannotationkeys.ProjectAnnotationKey]
	stack, hasStack := annotations[pulumiannotationkeys.StackNameAnnotationKey]

	if !hasOrg || !hasProject || !hasStack {
		return nil, fmt.Errorf("missing required Pulumi backend annotations: need either %s or all of (%s, %s, %s)",
			pulumiannotationkeys.StackFqdnAnnotationKey,
			pulumiannotationkeys.OrganizationAnnotationKey,
			pulumiannotationkeys.ProjectAnnotationKey,
			pulumiannotationkeys.StackNameAnnotationKey)
	}

	if org == "" || project == "" || stack == "" {
		return nil, fmt.Errorf("Pulumi backend annotations cannot be empty")
	}

	config.Organization = org
	config.Project = project
	config.StackName = stack
	config.StackFqdn = fmt.Sprintf("%s/%s/%s", org, project, stack)

	return config, nil
}

// parseStackFqdn splits "org/project/stack" into components
func parseStackFqdn(fqdn string) (org, project, stack string, err error) {
	parts := strings.Split(fqdn, "/")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("stack FQDN must be in format 'organization/project/stack', got: %s", fqdn)
	}

	org = strings.TrimSpace(parts[0])
	project = strings.TrimSpace(parts[1])
	stack = strings.TrimSpace(parts[2])

	if org == "" || project == "" || stack == "" {
		return "", "", "", fmt.Errorf("stack FQDN components cannot be empty")
	}

	return org, project, stack, nil
}
