package component

import (
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"

	v1 "github.com/plantonhq/planton/operator/api/v1"
)

// Deleting a PlantonPlatform must delete the platform. The operator has no
// finalizers, so the only thing standing between a deleted resource and a
// namespace full of orphaned workloads is the owner reference on every object
// the operator applies -- these tests pin that contract at the seam every
// chart-rendered object passes through.

func ownershipPlatform() *v1.PlantonPlatform {
	return &v1.PlantonPlatform{
		TypeMeta:   metav1.TypeMeta{APIVersion: "planton.ai/v1", Kind: "PlantonPlatform"},
		ObjectMeta: metav1.ObjectMeta{Name: "planton", Namespace: "planton", UID: types.UID("uid-platform")},
	}
}

func rendered(kind, namespace, name string) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetAPIVersion("v1")
	obj.SetKind(kind)
	obj.SetName(name)
	if namespace != "" {
		obj.SetNamespace(namespace)
	}
	return obj
}

func TestOwnedByPlatformStampsEveryNamespacedObject(t *testing.T) {
	platform := ownershipPlatform()
	base := &Base{}
	objs := []*unstructured.Unstructured{
		rendered("StatefulSet", "planton", "planton-redis-primary"),
		rendered("Service", "planton", "planton-redis-headless"),
		rendered("ConfigMap", "planton", "planton-temporal-config"),
		rendered("Job", "planton", "planton-temporal-schema-1"),
	}

	if err := ownedByPlatform(platform, base.OwnerReferenceFor(platform), objs); err != nil {
		t.Fatalf("expected the platform's own objects to be ownable, got: %v", err)
	}
	for _, obj := range objs {
		refs := obj.GetOwnerReferences()
		if len(refs) != 1 {
			t.Fatalf("%s %s: expected exactly one owner reference, got %d", obj.GetKind(), obj.GetName(), len(refs))
		}
		ref := refs[0]
		if ref.Kind != "PlantonPlatform" || ref.APIVersion != "planton.ai/v1" || ref.Name != "planton" || ref.UID != "uid-platform" {
			t.Errorf("%s %s: owner reference does not name the platform: %+v", obj.GetKind(), obj.GetName(), ref)
		}
		if ref.Controller != nil || ref.BlockOwnerDeletion != nil {
			t.Errorf("%s %s: the platform's owner reference is plain (no controller flag), matching every typed builder; got %+v", obj.GetKind(), obj.GetName(), ref)
		}
	}
}

func TestOwnedByPlatformKeepsAnExistingReferenceToTheSameOwner(t *testing.T) {
	// A CloudNativePG Cluster arrives with a controller reference already set
	// by its builder; stamping must not add a second entry for the same UID.
	platform := ownershipPlatform()
	base := &Base{}
	controller := true
	obj := rendered("Cluster", "planton", "planton-postgres")
	obj.SetOwnerReferences([]metav1.OwnerReference{{
		APIVersion: "planton.ai/v1", Kind: "PlantonPlatform", Name: "planton", UID: "uid-platform",
		Controller: &controller, BlockOwnerDeletion: &controller,
	}})

	if err := ownedByPlatform(platform, base.OwnerReferenceFor(platform), []*unstructured.Unstructured{obj}); err != nil {
		t.Fatalf("unexpected refusal: %v", err)
	}
	refs := obj.GetOwnerReferences()
	if len(refs) != 1 {
		t.Fatalf("expected the builder's reference to be kept as the only entry, got %d", len(refs))
	}
	if refs[0].Controller == nil || !*refs[0].Controller {
		t.Errorf("the builder's controller reference must survive stamping, got %+v", refs[0])
	}
}

func TestOwnedByPlatformIsIdempotent(t *testing.T) {
	platform := ownershipPlatform()
	base := &Base{}
	obj := rendered("Deployment", "planton", "planton-temporal-frontend")
	objs := []*unstructured.Unstructured{obj}
	for i := range 3 {
		if err := ownedByPlatform(platform, base.OwnerReferenceFor(platform), objs); err != nil {
			t.Fatalf("pass %d: unexpected refusal: %v", i, err)
		}
	}
	if got := len(obj.GetOwnerReferences()); got != 1 {
		t.Fatalf("expected one owner reference after repeated stamping, got %d", got)
	}
}

func TestOwnedByPlatformRefusesAClusterScopedObjectInWords(t *testing.T) {
	platform := ownershipPlatform()
	base := &Base{}
	err := ownedByPlatform(platform, base.OwnerReferenceFor(platform), []*unstructured.Unstructured{
		rendered("StatefulSet", "planton", "fine"),
		rendered("ClusterRoleBinding", "", "planton-auth-delegator"),
	})
	if err == nil {
		t.Fatal("expected a cluster-scoped rendered object to be refused")
	}
	for _, want := range []string{"ClusterRoleBinding", "planton-auth-delegator", "cluster-scoped", "never garbage-collected", "shared cluster infrastructure"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("refusal should say %q; got: %v", want, err)
		}
	}
}

func TestOwnedByPlatformRefusesAnObjectInAnotherNamespaceInWords(t *testing.T) {
	platform := ownershipPlatform()
	base := &Base{}
	err := ownedByPlatform(platform, base.OwnerReferenceFor(platform), []*unstructured.Unstructured{
		rendered("Secret", "istio-ingress", "planton-ingress-tls"),
	})
	if err == nil {
		t.Fatal("expected an object outside the platform's namespace to be refused")
	}
	for _, want := range []string{"istio-ingress/planton-ingress-tls", `platform's namespace "planton"`, "treats an owner in another namespace as absent", "garbage-collected immediately"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("refusal should say %q; got: %v", want, err)
		}
	}
}
