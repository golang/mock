package multi_pkg

//go:generate mockgen -destination ./multi_mock.go -package=multi_pkg github.com/golang/mock/mockgen/internal/tests/multi_pkg/pkg_x InterfaceA,InterfaceB github.com/golang/mock/mockgen/internal/tests/multi_pkg/pkg_y InterfaceD github.com/golang/mock/mockgen/internal/tests/multi_pkg InterfaceF

type InterfaceE interface {
	GetE(string) string
}

type InterfaceF interface {
	GetF(string) string
}
