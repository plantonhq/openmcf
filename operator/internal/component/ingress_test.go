package component

import (
	"context"
	"strings"
	"testing"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/plantonhq/planton/operator/internal/resources"
)

// A hostname derived by a prior reconcile is REUSED, never re-derived: the
// identity realm bakes the advertised URL in at first import, so a derivation
// that drifted with a re-published controller address would break sign-in
// silently. The recorded annotation on the Ingress is the whole memory.
func TestResolveAutoHostnameReusesTheRecordedDerivation(t *testing.T) {
	sticky := "planton-planton.203-0-113-7.sslip.io"
	existing := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        resources.IngressName("planton"),
			Namespace:   "planton",
			Annotations: map[string]string{resources.DerivedHostnameAnnotation: sticky},
		},
	}
	c := fake.NewClientBuilder().WithScheme(clientgoscheme.Scheme).WithObjects(existing).Build()

	ingress := &Ingress{}
	hostname, waitMsg, err := ingress.resolveAutoHostname(context.Background(), c,
		resources.IngressConfig{CRName: "planton", Namespace: "planton"})
	if err != nil {
		t.Fatal(err)
	}
	if waitMsg != "" {
		t.Fatalf("a recorded derivation must resolve immediately, got wait: %s", waitMsg)
	}
	if hostname != sticky {
		t.Fatalf("hostname = %s, want the recorded %s (derivation is once-only)", hostname, sticky)
	}
}

// The certificate-wait message builders are pure so every arm's wording is
// pinned here -- the message IS the product on this surface (a person acts
// on it directly, with no engineer to translate).

func TestCertificateWaitingMessageNamesTheDNSRecordWhenAddressKnown(t *testing.T) {
	msg := certificateWaitingMessage("planton.example.com",
		`Issuing certificate as Secret does not exist`, "203.0.113.7")
	for _, want := range []string{
		"waiting for the certificate for planton.example.com",
		"Issuing certificate as Secret does not exist",
		"point planton.example.com at 203.0.113.7",
		"validate over the internet",
	} {
		if !strings.Contains(msg, want) {
			t.Fatalf("waiting message missing %q:\n%s", want, msg)
		}
	}
}

func TestCertificateWaitingMessageDegradesWithoutAnAddress(t *testing.T) {
	// Datacenter controllers (RKE2's node-served nginx) legitimately publish
	// no address; the guidance must still name the way out, never require
	// what the cluster does not offer.
	msg := certificateWaitingMessage("planton.example.com", "", "")
	if !strings.Contains(msg, "the address your ingress controller serves on") {
		t.Fatalf("the no-address arm must still guide the DNS task:\n%s", msg)
	}
	if strings.Contains(msg, "at  ") || strings.Contains(msg, ": .") {
		t.Fatalf("empty inputs must not leave grammatical holes:\n%s", msg)
	}
}

func TestCertificateReadyConditionLiftsStatusAndMessage(t *testing.T) {
	cert := &unstructured.Unstructured{Object: map[string]any{
		"status": map[string]any{
			"conditions": []any{
				map[string]any{"type": "Issuing", "status": "True", "message": "in progress"},
				map[string]any{"type": "Ready", "status": "False", "message": "waiting on the order"},
			},
		},
	}}
	ready, message := certificateReadyCondition(cert)
	if ready || message != "waiting on the order" {
		t.Fatalf("expected the Ready condition's own message, got ready=%v message=%q", ready, message)
	}

	cert.Object["status"] = map[string]any{
		"conditions": []any{map[string]any{"type": "Ready", "status": "True", "message": "up to date"}},
	}
	if ready, _ := certificateReadyCondition(cert); !ready {
		t.Fatal("a True Ready condition must report issued")
	}

	if ready, message := certificateReadyCondition(&unstructured.Unstructured{Object: map[string]any{}}); ready || message != "" {
		t.Fatal("no conditions means not ready, with nothing to relay")
	}
}
