package embed

//go:generate mockgen -embed -package embed -destination mock.go . Hoge
//go:generate mockgen -embed -destination mock/mock.go . Hoge

type Hoge interface {
	Fuga() error
	mustImplementedFunction()
}

type HogeImpl struct {
	s string
}

func (h *HogeImpl) Fuga() error {
	return nil
}
