package mappingeval

import (
	"fmt"
	"strings"
)

// Report is the scored result of grading one proposal against one ground
// truth. It is JSON-serializable so runs leave durable artifacts, and its
// Summary renders the human-readable verdict.
type Report struct {
	Grouping         GroupingScore             `json:"grouping"`
	Spec             SpecScore                 `json:"spec"`
	Refs             RefScore                  `json:"refs"`
	Partition        PartitionScore            `json:"partition"`
	Coverage         Coverage                  `json:"coverage"`
	NameDerivability []NameDerivabilityFinding `json:"nameDerivability,omitempty"`
	Instances        []InstanceReport          `json:"instances"`
}

// GroupingScore grades resource-to-instance assignment over the scored
// universe.
type GroupingScore struct {
	// UniverseSize is the number of scan-visible ground-truth resources --
	// the recall denominator.
	UniverseSize int `json:"universeSize"`
	// InUniverseClaims counts every claim a proposed instance makes on a
	// universe resource (duplicates included) -- the precision denominator.
	InUniverseClaims int `json:"inUniverseClaims"`
	// Correct counts universe resources claimed by exactly the matched
	// counterpart of their ground-truth owner.
	Correct int `json:"correct"`

	Unclaimed       []AccountResourceRef  `json:"unclaimed,omitempty"`
	DuplicateClaims []AccountResourceRef  `json:"duplicateClaims,omitempty"`
	Misassigned     []MisassignedResource `json:"misassigned,omitempty"`
}

// Precision is Correct over every in-universe claim made.
func (g GroupingScore) Precision() float64 { return ratio(g.Correct, g.InUniverseClaims) }

// Recall is Correct over the universe.
func (g GroupingScore) Recall() float64 { return ratio(g.Correct, g.UniverseSize) }

// MisassignedResource is a universe resource claimed by the wrong instance.
type MisassignedResource struct {
	Resource        AccountResourceRef `json:"resource"`
	GroundTruthName string             `json:"groundTruthName"`
	ClaimedBy       string             `json:"claimedBy"`
}

// SpecScore grades manifest reconstruction: recall over the leaves the
// ground-truth specs set.
type SpecScore struct {
	// GroundTruthLeaves is the number of scored leaves the ground truth
	// sets (config-only-excluded and ref-arm leaves not included).
	GroundTruthLeaves int `json:"groundTruthLeaves"`
	Matched           int `json:"matched"`
	// ExcludedConfigOnly counts leaves excluded because no scan can see
	// them (declared config-only attributes) -- reported so the exclusion
	// is visible, never silent.
	ExcludedConfigOnly int `json:"excludedConfigOnly"`

	Mismatched []FieldFinding `json:"mismatched,omitempty"`
	Missing    []FieldFinding `json:"missing,omitempty"`
	// ProposerExtra lists leaves the proposal sets beyond the ground truth
	// -- informational: materializing a cloud-observed default explicitly
	// is honest, not wrong.
	ProposerExtra []FieldFinding `json:"proposerExtra,omitempty"`
}

// Recall is Matched over GroundTruthLeaves.
func (s SpecScore) Recall() float64 { return ratio(s.Matched, s.GroundTruthLeaves) }

// FieldFinding locates one spec-leaf observation.
type FieldFinding struct {
	Instance  string `json:"instance"`
	FieldPath string `json:"fieldPath"`
	Expected  string `json:"expected,omitempty"`
	Actual    string `json:"actual,omitempty"`
}

// RefScore grades value_from wiring.
type RefScore struct {
	GroundTruthEdges int `json:"groundTruthEdges"`
	CorrectEdges     int `json:"correctEdges"`

	MissingEdges     []EdgeFinding `json:"missingEdges,omitempty"`
	WrongTargetEdges []EdgeFinding `json:"wrongTargetEdges,omitempty"`
	// UnexpectedEdges are proposed references with no ground-truth
	// counterpart -- a frozen dependency the ground truth never declared.
	UnexpectedEdges []EdgeFinding `json:"unexpectedEdges,omitempty"`
}

// Recall is CorrectEdges over GroundTruthEdges.
func (r RefScore) Recall() float64 { return ratio(r.CorrectEdges, r.GroundTruthEdges) }

// Precision is CorrectEdges over every edge the proposal drew on matched
// instances.
func (r RefScore) Precision() float64 {
	return ratio(r.CorrectEdges, r.CorrectEdges+len(r.WrongTargetEdges)+len(r.UnexpectedEdges))
}

// EdgeFinding locates one reference observation.
type EdgeFinding struct {
	Instance  string `json:"instance"`
	FieldPath string `json:"fieldPath"`
	Detail    string `json:"detail"`
}

// PartitionScore grades environment assignment: each proposed manifest's
// metadata.env against its ground-truth counterpart's. Graded only on
// suites that DECLARE it (grade_environment_partition) -- an exam-fairness
// rule: fixture manifests also carry operational environments that leave no
// scan-visible trace, and owing those would be debt no proposer could pay.
type PartitionScore struct {
	// Graded reports whether the suite declared this axis. When false the
	// score is vacuously zero and Perfect() ignores it.
	Graded bool `json:"graded"`
	// GroundTruthInstances counts the ground-truth instances that set
	// metadata.env -- the denominator. Unproposed instances owe theirs
	// (mirroring the spec and refs axes: skipping a resource never
	// shrinks a denominator).
	GroundTruthInstances int `json:"groundTruthInstances"`
	// Correct counts matched instances whose proposed env equals the
	// ground truth's.
	Correct int `json:"correct"`

	// WrongEnv: matched instances assigned to a DIFFERENT environment --
	// the worst class (a wrong split poisons env-scoped references).
	WrongEnv []EnvFinding `json:"wrongEnv,omitempty"`
	// MissingEnv: ground-truth env owed but not delivered (the proposal
	// left env empty, or never proposed the instance). Honest
	// unassignment lands here -- visible, never penalized as wrong.
	MissingEnv []EnvFinding `json:"missingEnv,omitempty"`
	// ExtraEnv: proposal sets env where the ground truth has none --
	// informational, mirroring spec proposer-extra.
	ExtraEnv []EnvFinding `json:"extraEnv,omitempty"`
}

// Recall is Correct over GroundTruthInstances.
func (p PartitionScore) Recall() float64 { return ratio(p.Correct, p.GroundTruthInstances) }

// EnvFinding locates one environment-assignment observation.
type EnvFinding struct {
	Instance string `json:"instance"`
	Expected string `json:"expected,omitempty"`
	Actual   string `json:"actual,omitempty"`
}

// Coverage is the honesty ledger over the scored universe.
type Coverage struct {
	UniverseSize      int `json:"universeSize"`
	ClaimedInUniverse int `json:"claimedInUniverse"`

	// UnmappedInUniverse: universe resources the proposal explicitly
	// declared unmappable -- honest, visible, and better than silence.
	UnmappedInUniverse []AccountResourceRef `json:"unmappedInUniverse,omitempty"`
	// Unaccounted: universe resources the proposal neither claimed nor
	// declared unmapped -- the silent gap, the worst class.
	Unaccounted []AccountResourceRef `json:"unaccounted,omitempty"`
	// OutOfUniverseClaims: claims on resources outside the ground truth
	// (pre-existing account contents, AWS-implicit resources). Reported for
	// information; no answer key exists for them.
	OutOfUniverseClaims []AccountResourceRef `json:"outOfUniverseClaims,omitempty"`
}

// NameDerivabilityFinding flags a proposed name that breaks a
// from_metadata_name import recipe.
type NameDerivabilityFinding struct {
	Instance   string `json:"instance"`
	Kind       string `json:"kind"`
	Identifier string `json:"identifier"`
}

// InstanceReport is the per-ground-truth-instance ledger line.
type InstanceReport struct {
	Kind            string `json:"kind"`
	GroundTruthName string `json:"groundTruthName"`
	// ProposedName is empty when no proposed instance matched.
	ProposedName string `json:"proposedName,omitempty"`
	ClaimOverlap int    `json:"claimOverlap,omitempty"`
	// InvisibleResourceTypes: IaC resource types of this instance no scan
	// can discover (no cloud_control_type_name declared).
	InvisibleResourceTypes []string `json:"invisibleResourceTypes,omitempty"`
}

// Perfect reports whether the proposal reconstructed the ground truth
// completely: every universe resource correctly grouped, every scored spec
// leaf matched, every edge wired, every graded environment assigned,
// nothing unaccounted, no name-derivability breaks. The deterministic
// baseline is pinned to this bar on the seeded suites.
func (r *Report) Perfect() bool {
	return r.Grouping.Correct == r.Grouping.UniverseSize &&
		r.Grouping.InUniverseClaims == r.Grouping.UniverseSize &&
		r.Spec.Matched == r.Spec.GroundTruthLeaves &&
		len(r.Spec.Mismatched) == 0 && len(r.Spec.Missing) == 0 &&
		r.Refs.CorrectEdges == r.Refs.GroundTruthEdges &&
		len(r.Refs.WrongTargetEdges) == 0 && len(r.Refs.UnexpectedEdges) == 0 &&
		r.Partition.Correct == r.Partition.GroundTruthInstances &&
		len(r.Partition.WrongEnv) == 0 &&
		len(r.Coverage.Unaccounted) == 0 &&
		len(r.NameDerivability) == 0
}

// Summary renders the human-readable verdict.
func (r *Report) Summary() string {
	var b strings.Builder
	fmt.Fprintf(&b, "grouping: %d/%d correct (precision %.2f, recall %.2f)",
		r.Grouping.Correct, r.Grouping.UniverseSize, r.Grouping.Precision(), r.Grouping.Recall())
	if len(r.Grouping.Unclaimed) > 0 || len(r.Grouping.Misassigned) > 0 || len(r.Grouping.DuplicateClaims) > 0 {
		fmt.Fprintf(&b, " [unclaimed %d, misassigned %d, duplicate %d]",
			len(r.Grouping.Unclaimed), len(r.Grouping.Misassigned), len(r.Grouping.DuplicateClaims))
	}
	fmt.Fprintf(&b, "\nspec: %d/%d leaves matched (recall %.2f; %d mismatched, %d missing, %d excluded config-only, %d proposer-extra)",
		r.Spec.Matched, r.Spec.GroundTruthLeaves, r.Spec.Recall(),
		len(r.Spec.Mismatched), len(r.Spec.Missing), r.Spec.ExcludedConfigOnly, len(r.Spec.ProposerExtra))
	fmt.Fprintf(&b, "\nrefs: %d/%d edges wired (recall %.2f, precision %.2f; %d missing, %d wrong-target, %d unexpected)",
		r.Refs.CorrectEdges, r.Refs.GroundTruthEdges, r.Refs.Recall(), r.Refs.Precision(),
		len(r.Refs.MissingEdges), len(r.Refs.WrongTargetEdges), len(r.Refs.UnexpectedEdges))
	if r.Partition.Graded {
		fmt.Fprintf(&b, "\npartition: %d/%d environments assigned (recall %.2f; %d wrong, %d missing, %d extra informational)",
			r.Partition.Correct, r.Partition.GroundTruthInstances, r.Partition.Recall(),
			len(r.Partition.WrongEnv), len(r.Partition.MissingEnv), len(r.Partition.ExtraEnv))
	}
	fmt.Fprintf(&b, "\ncoverage: %d/%d universe resources claimed, %d declared unmapped, %d unaccounted, %d out-of-universe claims (informational)",
		r.Coverage.ClaimedInUniverse, r.Coverage.UniverseSize,
		len(r.Coverage.UnmappedInUniverse), len(r.Coverage.Unaccounted), len(r.Coverage.OutOfUniverseClaims))
	if len(r.NameDerivability) > 0 {
		fmt.Fprintf(&b, "\nname-derivability: %d proposed names break a from_metadata_name import recipe", len(r.NameDerivability))
	}
	for _, instance := range r.Instances {
		if instance.ProposedName == "" {
			fmt.Fprintf(&b, "\n  UNMATCHED %s %q", instance.Kind, instance.GroundTruthName)
		}
	}
	return b.String()
}

func ratio(numerator, denominator int) float64 {
	if denominator == 0 {
		return 1
	}
	return float64(numerator) / float64(denominator)
}
