package ana

// The returned type of TreeFunc must be a boolean.
// It will panic otherwise.
type Selection struct {
	Name     string
	TreeFunc TreeFunc
}

func NewSelection() *Selection {
	return &Selection{
		Name: "",
		TreeFunc: TreeFunc{
			VarsName: []string{},
			Fct:      func() bool { return true },
		},
	}
}
