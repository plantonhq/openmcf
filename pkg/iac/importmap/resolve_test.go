//go:build !codegen
// +build !codegen

package importmap

import (
	"reflect"
	"testing"

	componentv1 "github.com/plantonhq/planton/apis/dev/planton/iac/componentimportmap/v1"
	awsecrrepov1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsecrrepo/v1"
)

func TestPlaceholders(t *testing.T) {
	cases := []struct {
		idFormat string
		want     []string
	}{
		{"{bucket}", []string{"bucket"}},
		{"{bucket}:{intelligent_tiering_name}", []string{"bucket", "intelligent_tiering_name"}},
		{"{vpc_id}_{association_id}", []string{"vpc_id", "association_id"}},
		// Optional placeholders are still placeholders -- they need a
		// declaration and a resolution attempt.
		{"name:{table_name}/index:{index_name?}/{account_id}",
			[]string{"table_name", "index_name", "account_id"}},
		{"literal-id", nil},
	}
	for _, tc := range cases {
		got := Placeholders(tc.idFormat)
		if len(got) == 0 && len(tc.want) == 0 {
			continue
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Placeholders(%q) = %v, want %v", tc.idFormat, got, tc.want)
		}
	}
}

func TestRequiredPlaceholders(t *testing.T) {
	got := RequiredPlaceholders("name:{table_name}/index:{index_name?}/{account_id}")
	want := []string{"table_name", "account_id"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("RequiredPlaceholders = %v, want %v", got, want)
	}
}

func TestRenderID(t *testing.T) {
	id, err := RenderID("{bucket}:{name}", map[string]string{"bucket": "my-bucket", "name": "archive"})
	if err != nil {
		t.Fatalf("RenderID: %v", err)
	}
	if id != "my-bucket:archive" {
		t.Fatalf("RenderID = %q, want my-bucket:archive", id)
	}

	// A missing value must be an ERROR, never an empty substitution -- a
	// partially-filled ID would import the wrong resource.
	if _, err := RenderID("{bucket}:{name}", map[string]string{"bucket": "my-bucket"}); err == nil {
		t.Fatal("RenderID with a missing value did not error")
	}

	// An OPTIONAL placeholder renders as the empty string when unresolved --
	// the provider-documented shape for variants like table-level DynamoDB
	// contributor insights ("name:tbl/index:/123456789012").
	id, err = RenderID("name:{table_name}/index:{index_name?}/{account_id}",
		map[string]string{"table_name": "tbl", "account_id": "123456789012"})
	if err != nil {
		t.Fatalf("RenderID with optional placeholder: %v", err)
	}
	if id != "name:tbl/index:/123456789012" {
		t.Fatalf("RenderID = %q, want name:tbl/index:/123456789012", id)
	}
}

func TestParseTofuAddress(t *testing.T) {
	cases := []struct {
		address string
		resType string
		name    string
		key     string
		ok      bool
	}{
		{"aws_s3_bucket.this", "aws_s3_bucket", "this", "", true},
		{"aws_s3_bucket_intelligent_tiering_configuration.this[\"archive\"]",
			"aws_s3_bucket_intelligent_tiering_configuration", "this", "archive", true},
		{"aws_vpc_ipv4_cidr_block_association.secondary[\"10.1.0.0/16\"]",
			"aws_vpc_ipv4_cidr_block_association", "secondary", "10.1.0.0/16", true},
		// A count index is positional, never identity -- it must not leak
		// into an import ID as if it were a for_each discriminator.
		{"aws_instance.web[0]", "aws_instance", "web", "", true},
		{"module.nested.aws_s3_bucket.this", "", "", "", false},
		{"data.aws_ami.ubuntu", "", "", "", false},
	}
	for _, tc := range cases {
		resType, name, key, ok := ParseTofuAddress(tc.address)
		if ok != tc.ok || resType != tc.resType || name != tc.name || key != tc.key {
			t.Errorf("ParseTofuAddress(%q) = (%q,%q,%q,%v), want (%q,%q,%q,%v)",
				tc.address, resType, name, key, ok, tc.resType, tc.name, tc.key, tc.ok)
		}
	}
}

func TestResolveValues(t *testing.T) {
	m := &componentv1.ComponentImportMap{
		Spec: &componentv1.ComponentImportMapSpec{
			Values: []*componentv1.ImportValue{
				{
					Name: "bucket",
					Derivations: []*componentv1.ImportValueDerivation{
						{Source: &componentv1.ImportValueDerivation_FromMetadataName{FromMetadataName: true}},
					},
				},
				{
					Name: "repository_name",
					Derivations: []*componentv1.ImportValueDerivation{
						{Source: &componentv1.ImportValueDerivation_FromSpecField{FromSpecField: "repository_name"}},
					},
				},
				{
					Name: "vpc_id",
					Derivations: []*componentv1.ImportValueDerivation{
						{Source: &componentv1.ImportValueDerivation_FromStackOutput{FromStackOutput: "vpc_id"}},
						{Source: &componentv1.ImportValueDerivation_FromArnPart{FromArnPart: "resource_id"}},
					},
				},
				{
					Name: "tiering_name",
					Derivations: []*componentv1.ImportValueDerivation{
						{Source: &componentv1.ImportValueDerivation_FromAddressKey{FromAddressKey: true}},
					},
				},
				{Name: "association_id", WhereToFind: "describe-vpcs"},
			},
		},
	}

	rctx := ResolveContext{
		MetadataName: "my-bucket",
		Spec:         &awsecrrepov1.AwsEcrRepoSpec{RepositoryName: "team/app"},
		StackOutputs: map[string]string{},
		AddressKey:   "archive",
		ArnParts:     map[string]string{"resource_id": "vpc-0abc"},
	}

	resolved, unresolved := ResolveValues(
		m,
		[]string{"bucket", "repository_name", "vpc_id", "tiering_name", "association_id"},
		rctx,
	)

	want := map[string]string{
		"bucket":          "my-bucket",
		"repository_name": "team/app",
		// stack output empty -> falls through to the ARN part.
		"vpc_id":       "vpc-0abc",
		"tiering_name": "archive",
	}
	if !reflect.DeepEqual(resolved, want) {
		t.Errorf("resolved = %v, want %v", resolved, want)
	}
	if !reflect.DeepEqual(unresolved, []string{"association_id"}) {
		t.Errorf("unresolved = %v, want [association_id]", unresolved)
	}
}
