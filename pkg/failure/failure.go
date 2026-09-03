// Package failure holds the one shape every explained refusal takes, in the
// engines, the CLI, and the repository guards alike: what was observed (the
// fact, with its value), what it most likely means (one root cause, never a
// list of maybes), and the exact next step (a field, a flag, a command, an
// upstream fact to check). The three labelled lines are stable so an agent can
// match them and act; the values inside them vary.
//
// A message that names only the mechanism ("count mismatch", "not found",
// "connection refused") is a defect; wrap it in a Failure at the first place
// that knows what it means.
package failure

// Failure is an error whose text is the three-part explanation.
type Failure struct {
	Observed string
	Meaning  string
	NextStep string
}

func (f *Failure) Error() string {
	return "observed: " + f.Observed + "\nmeaning: " + f.Meaning + "\nnext step: " + f.NextStep
}
