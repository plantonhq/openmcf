package module

import (
	"github.com/pkg/errors"
	kubernetesrbacv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kubernetesrbac/v1alpha1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	rbacv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/rbac/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createBinding creates the RoleBinding (namespace scope) or ClusterRoleBinding
// (cluster scope) pointing every subject at the grant's role. No-op when the spec
// has no subjects — the grant then only publishes a role definition.
//
// createdRole is the Role/ClusterRole created by this module, or nil when binding
// to an existing role. When non-nil it is added as an explicit dependency so the
// binding is never applied before the role it references exists.
func createBinding(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource, createdRole pulumi.Resource) error {
	if len(locals.Spec.GetSubjects()) == 0 {
		return nil
	}

	metadata := &metav1.ObjectMetaArgs{
		Name:        pulumi.String(locals.BindingName),
		Labels:      pulumi.ToStringMap(locals.Labels),
		Annotations: pulumi.ToStringMap(locals.Annotations),
	}

	// roleRef is immutable in Kubernetes and always addresses the role through the
	// rbac.authorization.k8s.io API group. Its kind is the grant's resolved role
	// kind: the created role's kind, or for existing roles, ClusterRole in cluster
	// scope and Role-or-ClusterRole in namespace scope per is_cluster_role.
	roleRef := rbacv1.RoleRefArgs{
		ApiGroup: pulumi.String(rbacApiGroup),
		Kind:     pulumi.String(locals.RoleKind),
		Name:     pulumi.String(locals.RoleName),
	}

	subjects := buildSubjects(locals)

	resourceOptions := []pulumi.ResourceOption{pulumi.Provider(provider)}
	if createdRole != nil {
		resourceOptions = append(resourceOptions, pulumi.DependsOn([]pulumi.Resource{createdRole}))
	}

	if locals.IsNamespaceScoped {
		metadata.Namespace = pulumi.String(locals.Namespace)
		if _, err := rbacv1.NewRoleBinding(
			ctx,
			locals.BindingName,
			&rbacv1.RoleBindingArgs{
				Metadata: metadata,
				RoleRef:  roleRef,
				Subjects: subjects,
			},
			resourceOptions...,
		); err != nil {
			return errors.Wrapf(err, "failed to create role binding %s", locals.BindingName)
		}
		return nil
	}

	if _, err := rbacv1.NewClusterRoleBinding(
		ctx,
		locals.BindingName,
		&rbacv1.ClusterRoleBindingArgs{
			Metadata: metadata,
			RoleRef:  roleRef,
			Subjects: subjects,
		},
		resourceOptions...,
	); err != nil {
		return errors.Wrapf(err, "failed to create cluster role binding %s", locals.BindingName)
	}
	return nil
}

// buildSubjects maps spec subjects onto Kubernetes rbac/v1 Subjects:
//
//   - service_account → kind ServiceAccount with a namespace: the subject's own
//     when set, otherwise the grant's namespace (namespace scope only — spec
//     validation guarantees cluster-scoped grants set it explicitly, because a
//     ServiceAccount always lives in some namespace).
//   - user / group → kind User/Group under the rbac.authorization.k8s.io API
//     group. These are plain strings matched against what the cluster's
//     authenticator asserts; Kubernetes has no User or Group objects.
//
// PARITY-EXCEPTION: the Terraform kubernetes provider's subject schema defaults
// namespace to "default" for every subject kind, so its User/Group subjects carry
// namespace "default" in the stored object while this module omits it. The RBAC
// authorizer ignores namespace on User/Group subjects, so authorization behavior
// is identical across both engines.
func buildSubjects(locals *Locals) rbacv1.SubjectArray {
	specSubjects := locals.Spec.GetSubjects()
	subjects := make(rbacv1.SubjectArray, 0, len(specSubjects))

	for _, specSubject := range specSubjects {
		switch subject := specSubject.GetSubject().(type) {
		case *kubernetesrbacv1alpha1.KubernetesRbacSubject_ServiceAccount:
			subjectNamespace := subject.ServiceAccount.GetNamespace().GetValue()
			if subjectNamespace == "" {
				subjectNamespace = locals.Namespace
			}
			subjects = append(subjects, rbacv1.SubjectArgs{
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String(subject.ServiceAccount.GetName().GetValue()),
				Namespace: pulumi.String(subjectNamespace),
			})

		case *kubernetesrbacv1alpha1.KubernetesRbacSubject_User:
			subjects = append(subjects, rbacv1.SubjectArgs{
				Kind:     pulumi.String("User"),
				Name:     pulumi.String(subject.User),
				ApiGroup: pulumi.String(rbacApiGroup),
			})

		case *kubernetesrbacv1alpha1.KubernetesRbacSubject_Group:
			subjects = append(subjects, rbacv1.SubjectArgs{
				Kind:     pulumi.String("Group"),
				Name:     pulumi.String(subject.Group),
				ApiGroup: pulumi.String(rbacApiGroup),
			})
		}
	}

	return subjects
}
