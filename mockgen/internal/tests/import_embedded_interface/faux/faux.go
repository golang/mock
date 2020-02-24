package faux

import (
	"log"

	"github.com/golang/mock/mockgen/internal/tests/import_embedded_interface/other/ersatz"
)

type Foreign interface {
	ersatz.Embedded
}

func Conflict0() {
	log.Println()
}
