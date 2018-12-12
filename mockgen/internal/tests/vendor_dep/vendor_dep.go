package vendor_dep

import "test.src/a"

type VendorsDep interface {
	Foo() a.Ifc
}
