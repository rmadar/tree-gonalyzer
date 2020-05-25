package ana

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
