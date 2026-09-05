package resources

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Listener hostname admission follows the Gateway API's rule: empty admits
// all, a wildcard admits exactly one more label, otherwise exact.
func TestListenerAdmitsHostname(t *testing.T) {
	cases := []struct {
		listener, hostname string
		want               bool
	}{
		{"", "planton.example.com", true},
		{"planton.example.com", "planton.example.com", true},
		{"other.example.com", "planton.example.com", false},
		{"*.example.com", "planton.example.com", true},
		{"*.example.com", "a.b.example.com", false}, // one label only
		{"*.example.com", "example.com", false},     // the suffix alone is not a match
		{"*.example.com", "planton.example.org", false},
		{"*.example.com", ".example.com", false},
	}
	for _, c := range cases {
		if got := ListenerAdmitsHostname(c.listener, c.hostname); got != c.want {
			t.Errorf("ListenerAdmitsHostname(%q, %q) = %v, want %v", c.listener, c.hostname, got, c.want)
		}
	}
}

// The HTTPRoute is the route table, rule for rule, with numeric backend
// ports and a disabled request timeout only where server streams flow.
func TestHTTPRouteRendersTheRouteTable(t *testing.T) {
	route := HTTPRoute(HTTPRouteConfig{
		CRName: "planton", Namespace: "planton", Hostname: "planton.example.com",
		HostnameDerived: true, GatewayName: "main", GatewayNamespace: "gw", SectionName: "https",
	})

	if route.GetAnnotations()[DerivedHostnameAnnotation] != "planton.example.com" {
		t.Error("a derived hostname must be recorded on the route")
	}
	parents, _, _ := unstructured.NestedSlice(route.Object, "spec", "parentRefs")
	parent := parents[0].(map[string]any)
	if parent["name"] != "main" || parent["namespace"] != "gw" || parent["sectionName"] != "https" {
		t.Errorf("parentRef = %v", parent)
	}
	rules, _, _ := unstructured.NestedSlice(route.Object, "spec", "rules")
	table := FrontDoorRoutes()
	if len(rules) != len(table) {
		t.Fatalf("%d rules for %d routes", len(rules), len(table))
	}
	for idx, raw := range rules {
		rule := raw.(map[string]any)
		matches, _, _ := unstructured.NestedSlice(rule, "matches")
		path, _, _ := unstructured.NestedMap(matches[0].(map[string]any), "path")
		if path["type"] != "PathPrefix" || path["value"] != table[idx].PathPrefix {
			t.Errorf("rule %d path = %v, want PathPrefix %s", idx, path, table[idx].PathPrefix)
		}
		backends, _, _ := unstructured.NestedSlice(rule, "backendRefs")
		backend := backends[0].(map[string]any)
		if backend["name"] != table[idx].ServiceName("planton") || backend["port"] != int64(table[idx].ServicePort()) {
			t.Errorf("rule %d backend = %v", idx, backend)
		}
		_, hasTimeout, _ := unstructured.NestedMap(rule, "timeouts")
		if hasTimeout != (table[idx].Backend == BackendControlPlane) {
			t.Errorf("rule %d timeouts present = %v; only control-plane routes carry the streaming timeout", idx, hasTimeout)
		}
	}
}

// The grant lives in the SECRET's namespace and names the Gateway's.
func TestTLSReferenceGrantPermitsTheGatewayNamespace(t *testing.T) {
	grant := TLSReferenceGrant("planton", "planton", "gw", nil)
	if grant.GetNamespace() != "planton" {
		t.Errorf("grant namespace = %s, want the Secret's namespace", grant.GetNamespace())
	}
	from, _, _ := unstructured.NestedSlice(grant.Object, "spec", "from")
	if from[0].(map[string]any)["namespace"] != "gw" || from[0].(map[string]any)["kind"] != "Gateway" {
		t.Errorf("from = %v", from[0])
	}
	to, _, _ := unstructured.NestedSlice(grant.Object, "spec", "to")
	if to[0].(map[string]any)["name"] != IngressTLSSecretName("planton") {
		t.Errorf("to = %v", to[0])
	}
}

// Listener parsing lifts the facts the edge reasons about, defaulting the
// certificate reference namespace to the Gateway's own.
func TestParseGatewayListeners(t *testing.T) {
	gw := &unstructured.Unstructured{Object: map[string]any{
		"metadata": map[string]any{"name": "main", "namespace": "gw"},
		"spec": map[string]any{"listeners": []any{
			map[string]any{"name": "http", "protocol": "HTTP", "port": int64(80)},
			map[string]any{
				"name": "https", "protocol": "HTTPS", "port": int64(443), "hostname": "*.example.com",
				"allowedRoutes": map[string]any{"namespaces": map[string]any{
					"from": "Selector", "selector": map[string]any{"matchLabels": map[string]any{"planton": "yes"}},
				}},
				"tls": map[string]any{"certificateRefs": []any{
					map[string]any{"name": "wild"},
					map[string]any{"name": "planton-ingress-tls", "namespace": "planton"},
				}},
			},
		}},
	}}
	listeners := ParseGatewayListeners(gw)
	if len(listeners) != 2 {
		t.Fatalf("%d listeners", len(listeners))
	}
	if listeners[0].AllowedNamespaces != "Same" {
		t.Errorf("allowedRoutes default = %s, want Same", listeners[0].AllowedNamespaces)
	}
	https := listeners[1]
	if https.AllowedNamespaces != "Selector" || https.NamespaceSelector == nil || https.NamespaceSelector.MatchLabels["planton"] != "yes" {
		t.Errorf("selector not lifted: %+v", https)
	}
	if len(https.CertificateRefs) != 2 || https.CertificateRefs[0] != "gw/wild" || https.CertificateRefs[1] != "planton/planton-ingress-tls" {
		t.Errorf("certificateRefs = %v", https.CertificateRefs)
	}
}
