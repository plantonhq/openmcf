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
//
// Where that first place is, is layered. Code that can see the fact refuses
// in place (a Pulumi program, an HCL precondition on a data source's result,
// a CLI flag check). A failure that kills a read itself (a repository host
// that does not resolve, an API server answering Forbidden) reaches an HCL
// module only as the provider's raw text, so the layer that runs the engine
// reads the output afterwards with Explain and adds what it means and what to
// do, from the constructors here, never repeating what the module already
// said. Both engines therefore end in the same words for the same failure.
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
