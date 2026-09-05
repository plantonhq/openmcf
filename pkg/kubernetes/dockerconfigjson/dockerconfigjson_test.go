package dockerconfigjson

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
)

func decode(t *testing.T, doc string) map[string]map[string]string {
	t.Helper()
	var parsed struct {
		Auths map[string]map[string]string `json:"auths"`
	}
	if err := json.Unmarshal([]byte(doc), &parsed); err != nil {
		t.Fatalf("encoded document is not JSON: %v\n%s", err, doc)
	}
	return parsed.Auths
}

func TestEncode_OneRecordPerServerWithBase64AuthPair(t *testing.T) {
	doc, err := Encode([]Auth{
		{Server: "ghcr.io", Username: "acme-pull-bot", Password: "ghp_token", Email: "ops@acme.io"},
		{Server: "quay.io", Username: "acme+robot", Password: "robot-token"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	auths := decode(t, doc)
	if len(auths) != 2 {
		t.Fatalf("want 2 auths records, got %d: %s", len(auths), doc)
	}
	ghcr := auths["ghcr.io"]
	if ghcr["username"] != "acme-pull-bot" || ghcr["password"] != "ghp_token" || ghcr["email"] != "ops@acme.io" {
		t.Fatalf("ghcr record mismatch: %v", ghcr)
	}
	wantAuth := base64.StdEncoding.EncodeToString([]byte("acme-pull-bot:ghp_token"))
	if ghcr["auth"] != wantAuth {
		t.Fatalf("auth pair: want %q, got %q", wantAuth, ghcr["auth"])
	}
	if _, hasEmail := auths["quay.io"]["email"]; hasEmail {
		t.Fatalf("an empty email must be omitted, got %v", auths["quay.io"])
	}
}

func TestEncode_IsDeterministicAcrossInputOrder(t *testing.T) {
	a := []Auth{{Server: "ghcr.io", Username: "u", Password: "p"}, {Server: "quay.io", Username: "u", Password: "p"}}
	b := []Auth{a[1], a[0]}
	docA, errA := Encode(a)
	docB, errB := Encode(b)
	if errA != nil || errB != nil {
		t.Fatalf("unexpected errors: %v %v", errA, errB)
	}
	if docA != docB {
		t.Fatalf("documents differ by input order:\n%s\n%s", docA, docB)
	}
}

func TestEncode_EmptyInputIsAnEmptyAuthsMap(t *testing.T) {
	doc, err := Encode(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doc != `{"auths":{}}` {
		t.Fatalf("want an empty auths map, got %s", doc)
	}
}

func TestEncode_RefusesADuplicateServerWithTheSentence(t *testing.T) {
	_, err := Encode([]Auth{
		{Server: "ghcr.io", Username: "a", Password: "p"},
		{Server: "ghcr.io", Username: "b", Password: "p"},
	})
	if err == nil {
		t.Fatal("a duplicate server must be refused")
	}
	for _, want := range []string{`"ghcr.io"`, "one login per registry", "merge them into one entry"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("refusal must say %q, got: %s", want, err)
		}
	}
}

func TestServers_SortsForStableLogLines(t *testing.T) {
	got := Servers([]Auth{{Server: "quay.io"}, {Server: "ghcr.io"}})
	if strings.Join(got, ",") != "ghcr.io,quay.io" {
		t.Fatalf("want sorted servers, got %v", got)
	}
}
