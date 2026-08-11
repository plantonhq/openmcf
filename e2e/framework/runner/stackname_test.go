package runner

import (
	"strings"
	"testing"
)

// The two dependency labels that motivated boundStackName: a two-database
// Firestore chain whose labels agree on the first 50 characters. A plain
// head-truncate collapsed them onto one stack (the second up REPLACED the
// first fixture; teardown then failed "no stack named" on the survivor).
func TestBoundStackName_PreservesUniquenessUnderTruncation(t *testing.T) {
	a := boundStackName(GenerateStackName("dep-gcpfirestoredatabase-planton-oss-e2e-gcpfsdb-prereq", "run1234abcd"))
	b := boundStackName(GenerateStackName("dep-gcpfirestoredatabase-planton-oss-e2e-gcpfsdb-ent-idx", "run1234abcd"))

	if len(a) > maxStackNameLen || len(b) > maxStackNameLen {
		t.Fatalf("bounded names exceed %d chars: %q (%d), %q (%d)", maxStackNameLen, a, len(a), b, len(b))
	}
	if a == b {
		t.Fatalf("distinct labels collapsed onto one stack name: %q", a)
	}
	// Deterministic: up and destroy must compute the same name.
	if a != boundStackName(GenerateStackName("dep-gcpfirestoredatabase-planton-oss-e2e-gcpfsdb-prereq", "run1234abcd")) {
		t.Fatal("boundStackName is not deterministic")
	}
	// Distinct runs must still produce distinct names for the same label.
	if a == boundStackName(GenerateStackName("dep-gcpfirestoredatabase-planton-oss-e2e-gcpfsdb-prereq", "run5678wxyz")) {
		t.Fatal("run id no longer disambiguates after truncation")
	}
}

func TestBoundStackName_ShortNamesPassThrough(t *testing.T) {
	name := GenerateStackName("dep-gcpvpcnetwork-vpc", "run1234abcd")
	if got := boundStackName(name); got != name {
		t.Fatalf("short name modified: %q -> %q", name, got)
	}
	if strings.Contains(boundStackName(name), "--") {
		t.Fatalf("unexpected separator artifact in %q", boundStackName(name))
	}
}
