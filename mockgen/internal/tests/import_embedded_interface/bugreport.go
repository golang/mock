//go:generate mockgen -destination bugreport_mock.go -package bugreport -source=bugreport.go

package bugreport

import (
	"log"

	"github.com/golang/mock/mockgen/internal/tests/import_embedded_interface/ersatz"
	"github.com/golang/mock/mockgen/internal/tests/import_embedded_interface/faux"
)

// Source is an interface w/ an embedded foreign interface
type Source interface {
	ersatz.Embedded
	faux.Foreign
}

func CallForeignMethod(s Source) {
	log.Println(s.Ersatz())
	log.Println(s.OtherErsatz())
}
