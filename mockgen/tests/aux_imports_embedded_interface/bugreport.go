//go:generate mockgen -aux_files aux=aux/aux.go -destination bugreport_mock.go -package bugreport -source=bugreport.go Example

package bugreport

import (
	"log"

	"github.com/golang/mock/mockgen/tests/aux_imports_embedded_interface/aux"
)

// Source is an interface w/ an embedded foreign interface
type Source interface {
	aux.Foreign
}

func CallForeignMethod(s Source) {
	log.Println(s.Method())
}
