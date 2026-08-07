package verify

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// RbacVerifier verifies a KubernetesRbac grant: the role object (when the
// grant defines one) and the binding object (when the grant has subjects).
// The created object names follow the module's naming rules — role name is
// the explicit createRole.name or the component's metadata.name; the binding
// is always named after the component — so the verifier re-derives them from
// the manifest instead of guessing from ManifestInfo's generic fields.
type RbacVerifier struct {
	ManifestPath string
}

// rbacManifestShape captures just the fields of a KubernetesRbac manifest the
// verifier needs. YAML manifests carry protojson field names (camelCase).
type rbacManifestShape struct {
	Metadata struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		NamespaceScope *struct {
			Namespace *struct {
				Value string `yaml:"value"`
			} `yaml:"namespace"`
		} `yaml:"namespaceScope"`
		ClusterScope *struct{} `yaml:"clusterScope"`
		CreateRole   *struct {
			Name string `yaml:"name"`
		} `yaml:"createRole"`
		ExistingRole *struct {
			Name string `yaml:"name"`
		} `yaml:"existingRole"`
		Subjects []map[string]interface{} `yaml:"subjects"`
	} `yaml:"spec"`
}

// rbacExpectations resolves which objects the grant deploys and their names.
type rbacExpectations struct {
	namespace   string // empty for cluster scope
	roleKind    string // "role" or "clusterrole"; empty when no role is created
	roleName    string
	bindingKind string // "rolebinding" or "clusterrolebinding"; empty when no binding
	bindingName string
}

func (v *RbacVerifier) expectations() (*rbacExpectations, error) {
	data, err := os.ReadFile(v.ManifestPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read manifest %s", v.ManifestPath)
	}
	var m rbacManifestShape
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, errors.Wrapf(err, "failed to parse manifest YAML %s", v.ManifestPath)
	}

	e := &rbacExpectations{}
	clusterScoped := m.Spec.ClusterScope != nil
	if !clusterScoped {
		e.namespace = "default"
		if m.Spec.NamespaceScope != nil && m.Spec.NamespaceScope.Namespace != nil && m.Spec.NamespaceScope.Namespace.Value != "" {
			e.namespace = m.Spec.NamespaceScope.Namespace.Value
		}
	}

	if m.Spec.CreateRole != nil {
		e.roleKind = "role"
		if clusterScoped {
			e.roleKind = "clusterrole"
		}
		e.roleName = m.Spec.CreateRole.Name
		if e.roleName == "" {
			e.roleName = m.Metadata.Name
		}
	}

	if len(m.Spec.Subjects) > 0 {
		e.bindingKind = "rolebinding"
		if clusterScoped {
			e.bindingKind = "clusterrolebinding"
		}
		e.bindingName = m.Metadata.Name
	}

	return e, nil
}

func (v *RbacVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	e, err := v.expectations()
	if err != nil {
		return err
	}
	if e.roleKind != "" {
		if err := KubectlResourceExists(ctx, kubeconfig, e.roleKind, e.roleName, e.namespace); err != nil {
			return err
		}
	}
	if e.bindingKind != "" {
		if err := KubectlResourceExists(ctx, kubeconfig, e.bindingKind, e.bindingName, e.namespace); err != nil {
			return err
		}
	}
	return nil
}

func (v *RbacVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	e, err := v.expectations()
	if err != nil {
		return err
	}
	if e.roleKind != "" {
		if err := KubectlResourceAbsent(ctx, kubeconfig, e.roleKind, e.roleName, e.namespace); err != nil {
			return err
		}
	}
	if e.bindingKind != "" {
		if err := KubectlResourceAbsent(ctx, kubeconfig, e.bindingKind, e.bindingName, e.namespace); err != nil {
			return err
		}
	}
	return nil
}
