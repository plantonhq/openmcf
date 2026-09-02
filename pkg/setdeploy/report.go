package setdeploy

// Severity classifies one preflight report entry. The wall never prints, so
// these are semantic classes for the caller's renderer, not styling hints.
type Severity string

const (
	// SeverityRefusal blocks the deploy: the wall found something that would
	// fail or lie later, and the entry's message names what, where, why, and
	// what fixes it.
	SeverityRefusal Severity = "refusal"

	// SeverityWarning is surfaced but does not block: the deploy can proceed
	// and the user should still know (a module artifact that will fall back
	// to a slower source checkout, a sensitive value about to enter IaC
	// diff output).
	SeverityWarning Severity = "warning"

	// SeverityAssumption is a fact the wall cannot verify and deploys on
	// trust: an external relationship target that carries no value need, a
	// provider whose credentials have no probe. Stated, never silent.
	SeverityAssumption Severity = "assumption"
)

// Entry is one line of a preflight check's outcome. Source and FieldPath are
// empty when the entry concerns the whole set rather than one document.
type Entry struct {
	Severity  Severity
	Source    string
	FieldPath string
	Message   string
}

// Check is one wall check's outcome: the facts it verified (rendered as the
// pass lines) and the entries it raised. A check with no entries and no
// verified facts still renders — an empty check line is itself information
// ("no references to verify").
type Check struct {
	// Name is the check's stable identifier (load-and-schema, identity,
	// references, backend-resolved-values, cycles, engine-and-modules,
	// state-backend, provider-credentials).
	Name string
	// Title is the human heading the renderer prints.
	Title string
	// Verified are the facts the check positively confirmed, one line each.
	Verified []string
	// Entries are the refusals, warnings, and assumptions the check raised.
	Entries []Entry
}

func (c *Check) refusals() int {
	n := 0
	for _, e := range c.Entries {
		if e.Severity == SeverityRefusal {
			n++
		}
	}
	return n
}

// Report is the preflight wall's whole outcome: every check, in wall order,
// with all failures collected — the wall never stops at the first problem,
// because a user in CI fixes the full list in one commit, not one problem
// per run.
type Report struct {
	Checks []Check
}

// Refused reports whether any check raised a refusal — the deploy must not
// start.
func (r *Report) Refused() bool {
	return r.RefusalCount() > 0
}

// RefusalCount counts refusal entries across all checks.
func (r *Report) RefusalCount() int {
	n := 0
	for i := range r.Checks {
		n += r.Checks[i].refusals()
	}
	return n
}

// add appends a check outcome in wall order.
func (r *Report) add(c Check) {
	r.Checks = append(r.Checks, c)
}
