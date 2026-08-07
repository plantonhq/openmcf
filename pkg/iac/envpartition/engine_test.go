package envpartition

import (
	"reflect"
	"testing"

	rulev1 "github.com/plantonhq/planton/iac/environmentpartitionrule/v1"
)

// testRule compiles a small vocabulary shaped like the default rule.
func testRule(t *testing.T, mutate func(spec *rulev1.EnvironmentPartitionRuleSpec)) *Rule {
	t.Helper()
	doc := &rulev1.EnvironmentPartitionRule{
		ApiVersion: "iac.planton.dev/v1",
		Kind:       "EnvironmentPartitionRule",
		Spec: &rulev1.EnvironmentPartitionRuleSpec{
			AuthoritativeTagKeys: []string{"planton.ai/environment", "env", "environment", "stage"},
			Environments: []*rulev1.EnvironmentDefinition{
				{Name: "prod", Tokens: []string{"prod", "prd", "production"}},
				{Name: "staging", Tokens: []string{"staging", "stage", "stg"}},
			},
		},
	}
	if mutate != nil {
		mutate(doc.Spec)
	}
	rule, err := CompileRule(doc)
	if err != nil {
		t.Fatalf("compiling test rule: %v", err)
	}
	return rule
}

func assignmentOf(t *testing.T, result *Result, identifier string) Assignment {
	t.Helper()
	for _, a := range result.Assignments {
		if a.Identifier == identifier {
			return a
		}
	}
	t.Fatalf("no assignment for %q", identifier)
	return Assignment{}
}

func TestDefaultRuleCompiles(t *testing.T) {
	rule := DefaultRule()
	if !rule.containmentInheritance {
		t.Fatal("default rule must keep containment inheritance on")
	}
	if rule.fallbackEnvironment != "" {
		t.Fatal("default rule must not declare a fallback environment -- unassigned stays honest")
	}
	if got := rule.tokenToEnv["prd"]; got != "prod" {
		t.Fatalf("default rule: token prd -> %q, want prod", got)
	}
	if got := rule.normalizeEnvValue("Stage"); got != "staging" {
		t.Fatalf("default rule: tag value Stage -> %q, want staging", got)
	}
}

// DefaultRuleDocument is the transport/display form of the same embedded
// document DefaultRule compiles; the two must never diverge.
func TestDefaultRuleDocumentMatchesDefaultRule(t *testing.T) {
	doc := DefaultRuleDocument()
	if doc.GetKind() != "EnvironmentPartitionRule" {
		t.Fatalf("document kind = %q, want EnvironmentPartitionRule", doc.GetKind())
	}
	recompiled, err := CompileRule(doc)
	if err != nil {
		t.Fatalf("compiling DefaultRuleDocument: %v", err)
	}
	if !reflect.DeepEqual(recompiled, DefaultRule()) {
		t.Fatal("DefaultRuleDocument compiles to a different rule than DefaultRule")
	}
	// A fresh copy every call: mutating one caller's document must not
	// poison the next caller's.
	doc.Spec.AuthoritativeTagKeys = nil
	if len(DefaultRuleDocument().GetSpec().GetAuthoritativeTagKeys()) == 0 {
		t.Fatal("DefaultRuleDocument must return a fresh copy per call")
	}
}

// The ladder itself: each tier outranks everything below it.
func TestPrecedenceLadder(t *testing.T) {
	rule := testRule(t, func(spec *rulev1.EnvironmentPartitionRuleSpec) {
		spec.Overrides = []*rulev1.PartitionOverride{
			{TypeName: "AWS::S3::Bucket", Identifier: "orders-stg-archive", Environment: "prod"},
		}
	})
	result := Partition(rule, []Resource{
		// Tier 1 beats tier 2 AND tier 3: tagged staging, named stg, taught prod.
		{TypeName: "AWS::S3::Bucket", Identifier: "orders-stg-archive", Name: "orders-stg-archive", Tags: map[string]string{"env": "staging"}},
		// Tier 2 beats tier 3: named stg, tagged prod.
		{TypeName: "AWS::SQS::Queue", Identifier: "q-1", Name: "orders-stg-events", Tags: map[string]string{"env": "prod"}},
		// Tier 3 beats tier 4: named prod, contained in a staging VPC.
		{TypeName: "AWS::EC2::VPC", Identifier: "vpc-stg", Name: "orders-stg-vpc"},
		{TypeName: "AWS::EC2::Subnet", Identifier: "subnet-1", Name: "orders-prod-subnet", Containers: []string{"vpc-stg"}},
	})

	if a := assignmentOf(t, result, "orders-stg-archive"); a.Environment != "prod" || a.Tier != TierOverride {
		t.Fatalf("override must beat tag and name: got %+v", a)
	}
	if a := assignmentOf(t, result, "q-1"); a.Environment != "prod" || a.Tier != TierTag {
		t.Fatalf("tag must beat name token: got %+v", a)
	}
	subnet := assignmentOf(t, result, "subnet-1")
	if subnet.Environment != "prod" || subnet.Tier != TierNameToken {
		t.Fatalf("name token must beat containment: got %+v", subnet)
	}
	// The tension with the container is flagged, never silently resolved.
	if len(subnet.Conflicts) != 1 {
		t.Fatalf("own-signal-vs-container tension must be flagged: got %+v", subnet.Conflicts)
	}
}

func TestAuthoritativeTagValues(t *testing.T) {
	rule := testRule(t, nil)
	result := Partition(rule, []Resource{
		// A declared alias normalizes to the canonical name.
		{Identifier: "r-1", Tags: map[string]string{"env": "prd"}},
		// An undeclared value is taken literally: the tag is the user's word.
		{Identifier: "r-2", Tags: map[string]string{"env": "ops"}},
		// Agreeing keys are fine (prod twice = one environment).
		{Identifier: "r-3", Tags: map[string]string{"env": "prod", "stage": "production"}},
		// Disagreeing authoritative tags: conflict, never a guess.
		{Identifier: "r-4", Tags: map[string]string{"env": "prod", "stage": "staging"}},
		// Non-authoritative tags never assign.
		{Identifier: "r-5", Tags: map[string]string{"team": "prod"}},
	})
	if a := assignmentOf(t, result, "r-1"); a.Environment != "prod" {
		t.Fatalf("prd must normalize to prod: got %+v", a)
	}
	if a := assignmentOf(t, result, "r-2"); a.Environment != "ops" || a.Tier != TierTag {
		t.Fatalf("literal tag value must name its environment: got %+v", a)
	}
	if a := assignmentOf(t, result, "r-3"); a.Environment != "prod" || len(a.Conflicts) != 0 {
		t.Fatalf("agreeing keys must assign cleanly: got %+v", a)
	}
	if a := assignmentOf(t, result, "r-4"); a.Environment != "" || len(a.Conflicts) != 1 {
		t.Fatalf("disagreeing tags must conflict and stay unassigned: got %+v", a)
	}
	if a := assignmentOf(t, result, "r-5"); a.Environment != "" {
		t.Fatalf("undeclared tag keys must never assign: got %+v", a)
	}
}

func TestNameTokensMatchWholeTokensOnly(t *testing.T) {
	rule := testRule(t, nil)
	result := Partition(rule, []Resource{
		{Identifier: "n-1", Name: "orders-prod-vpc"},
		{Identifier: "n-2", Name: "orders_stg.queue"},
		// Substrings must never match: neither embedded ("reproduction")
		// nor as part of a longer token ("prodfix").
		{Identifier: "n-3", Name: "reproduction-pipeline"},
		{Identifier: "n-4", Name: "sg-prodfix"},
		// Tokens of different environments in one name: conflict.
		{Identifier: "n-5", Name: "prod-to-staging-sync"},
		// Repeated tokens of the SAME environment agree.
		{Identifier: "n-6", Name: "prod-orders-prod"},
	})
	if a := assignmentOf(t, result, "n-1"); a.Environment != "prod" || a.Tier != TierNameToken {
		t.Fatalf("whole token must match: got %+v", a)
	}
	if a := assignmentOf(t, result, "n-2"); a.Environment != "staging" {
		t.Fatalf("underscore/dot separators must split tokens: got %+v", a)
	}
	if a := assignmentOf(t, result, "n-3"); a.Environment != "" || len(a.Conflicts) != 0 {
		t.Fatalf("embedded substring must not match: got %+v", a)
	}
	if a := assignmentOf(t, result, "n-4"); a.Environment != "" {
		t.Fatalf("token-embedded substring must not match: got %+v", a)
	}
	if a := assignmentOf(t, result, "n-5"); a.Environment != "" || len(a.Conflicts) != 1 {
		t.Fatalf("cross-environment tokens must conflict: got %+v", a)
	}
	if a := assignmentOf(t, result, "n-6"); a.Environment != "prod" || len(a.Conflicts) != 0 {
		t.Fatalf("same-environment tokens must agree: got %+v", a)
	}
}

func TestContainmentInheritsTransitively(t *testing.T) {
	rule := testRule(t, nil)
	// The chain arrives OUT of dependency order to prove the fixpoint owes
	// nothing to input order: NAT -> subnet -> VPC(named prod).
	resources := []Resource{
		{TypeName: "AWS::EC2::NatGateway", Identifier: "nat-1", Containers: []string{"subnet-1", "vpc-1"}},
		{TypeName: "AWS::EC2::Subnet", Identifier: "subnet-1", Containers: []string{"vpc-1"}},
		{TypeName: "AWS::EC2::VPC", Identifier: "vpc-1", Name: "orders-prod-vpc"},
		// A container absent from the set is ignored; no present container
		// means no inheritance.
		{TypeName: "AWS::EC2::SecurityGroup", Identifier: "sg-1", Containers: []string{"vpc-absent"}},
	}
	result := Partition(rule, resources)
	if a := assignmentOf(t, result, "subnet-1"); a.Environment != "prod" || a.Tier != TierContainment {
		t.Fatalf("subnet must inherit from its VPC: got %+v", a)
	}
	if a := assignmentOf(t, result, "nat-1"); a.Environment != "prod" || a.Tier != TierContainment {
		t.Fatalf("NAT must inherit through the chain: got %+v", a)
	}
	if a := assignmentOf(t, result, "sg-1"); a.Environment != "" || len(a.Conflicts) != 0 {
		t.Fatalf("absent containers must not assign or conflict: got %+v", a)
	}
}

func TestContainmentDisagreementAndDisable(t *testing.T) {
	rule := testRule(t, nil)
	resources := []Resource{
		{Identifier: "vpc-p", Name: "prod-vpc"},
		{Identifier: "vpc-s", Name: "stg-vpc"},
		// Sits in two containers with different environments: flagged,
		// never guessed.
		{Identifier: "torn", Containers: []string{"vpc-p", "vpc-s"}},
	}
	result := Partition(rule, resources)
	if a := assignmentOf(t, result, "torn"); a.Environment != "" || len(a.Conflicts) != 1 {
		t.Fatalf("disagreeing containers must conflict: got %+v", a)
	}

	disabled := testRule(t, func(spec *rulev1.EnvironmentPartitionRuleSpec) {
		spec.DisableContainmentInheritance = true
	})
	result = Partition(disabled, []Resource{
		{Identifier: "vpc-p", Name: "prod-vpc"},
		{Identifier: "subnet-1", Containers: []string{"vpc-p"}},
	})
	if a := assignmentOf(t, result, "subnet-1"); a.Environment != "" {
		t.Fatalf("disabled inheritance must leave the subnet unassigned: got %+v", a)
	}
}

func TestFallbackAppliesOnlyToCleanUnassigned(t *testing.T) {
	rule := testRule(t, func(spec *rulev1.EnvironmentPartitionRuleSpec) {
		spec.FallbackEnvironment = "default"
	})
	result := Partition(rule, []Resource{
		{Identifier: "clean", Name: "telemetry-bucket"},
		{Identifier: "conflicted", Tags: map[string]string{"env": "prod", "stage": "staging"}},
	})
	if a := assignmentOf(t, result, "clean"); a.Environment != "default" || a.Tier != TierFallback {
		t.Fatalf("clean unassigned must take the fallback: got %+v", a)
	}
	// A conflict is a question for the review gate; the fallback must not
	// bury it.
	if a := assignmentOf(t, result, "conflicted"); a.Environment != "" {
		t.Fatalf("conflicted resources must never take the fallback: got %+v", a)
	}
}

func TestPartitionIsDeterministic(t *testing.T) {
	rule := testRule(t, nil)
	resources := []Resource{
		{Identifier: "vpc-1", Name: "orders-prod-vpc"},
		{Identifier: "subnet-1", Containers: []string{"vpc-1"}},
		{Identifier: "q-1", Name: "orders-stg-events", Tags: map[string]string{"env": "prd", "environment": "production"}},
		{Identifier: "loose", Name: "telemetry"},
	}
	first := Partition(rule, resources)
	for i := 0; i < 10; i++ {
		if next := Partition(rule, resources); !reflect.DeepEqual(first, next) {
			t.Fatalf("run %d diverged:\nfirst: %+v\nnext:  %+v", i, first, next)
		}
	}
	if got := first.Environments(); !reflect.DeepEqual(got, []string{"prod"}) {
		t.Fatalf("Environments() = %v, want [prod]", got)
	}
}

func TestCompileRuleRejectsAmbiguousVocabulary(t *testing.T) {
	doc := &rulev1.EnvironmentPartitionRule{
		Kind: "EnvironmentPartitionRule",
		Spec: &rulev1.EnvironmentPartitionRuleSpec{
			Environments: []*rulev1.EnvironmentDefinition{
				{Name: "prod", Tokens: []string{"prod"}},
				{Name: "staging", Tokens: []string{"prod"}},
			},
		},
	}
	if _, err := CompileRule(doc); err == nil {
		t.Fatal("a token declared by two environments must be rejected")
	}

	doc = &rulev1.EnvironmentPartitionRule{Kind: "SomethingElse"}
	if _, err := CompileRule(doc); err == nil {
		t.Fatal("a wrong kind string must be rejected")
	}
}
