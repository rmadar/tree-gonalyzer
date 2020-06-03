package ana

// The returned type of TreeFunc must be a boolean.
type Selection struct {
	Name     string
	TreeFunc TreeFunc
}

// EmptySelection returns an empty selection type,
// i.e. without name (empty string) and returning
// always true.
func EmptySelection() *Selection {
	return &Selection{
		Name: "",
		TreeFunc: TreeFunc{
			VarsName: []string{},
			Fct:      func() bool { return true },
		},
	}
}

// NewSelection returns a selection with the specified name
// and TreeFunc fct. Fct must return a boolean.
func NewSelection(name string, fct TreeFunc) *Selection {
	return &Selection{
		Name:     name,
		TreeFunc: fct,
	}
}
