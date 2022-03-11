package typed

//go:generate mockgen -typed -aux_files faux=faux/faux.go -destination bugreport_mock.go -package typed -source=bugreport.go Example

import (
	"log"

	"github.com/golang/mock/mockgen/internal/tests/typed/faux"
)

// Source is an interface w/ an embedded foreign interface
type Source interface {
	faux.Foreign
}

func CallForeignMethod(s Source) {
	log.Println(s.Method())
}
