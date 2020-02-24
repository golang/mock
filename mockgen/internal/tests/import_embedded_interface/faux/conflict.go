package faux

import "github.com/golang/mock/mockgen/internal/tests/import_embedded_interface/other/log"

func Conflict1() {
	log.Foo()
}
