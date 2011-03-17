// An example package with an interface.
package user

// Random bunch of imports to test mockgen.
import "io"
import (
	btz "bytes"
	"hash"
	"log"
	"os"
)

// Dependencies outside the standard library.
import (
	"gomock.googlecode.com/hg/sample/imp1"
	renamed2                               "gomock.googlecode.com/hg/sample/imp2"
	.                                      "gomock.googlecode.com/hg/sample/imp3"
	"gomock.googlecode.com/hg/sample/imp4" // calls itself "imp_four"
)

// A bizarre interface to test corner cases in mockgen.
// This would normally be in its own file or package,
// separate from the user of it (e.g. io.Reader).
type Index interface {
	Get(key string) interface{}
	GetTwo(key1, key2 string) (v1, v2 interface{})
	Put(key string, value interface{})

	// Check that imports are handled correctly.
	Summary(buf *btz.Buffer, w io.Writer)
	Other() hash.Hash

	// A method with an anonymous argument.
	Anon(string)

	// Methods using foreign types outside the standard library.
	ForeignOne(imp1.Imp1)
	ForeignTwo(renamed2.Imp2)
	ForeignThree(Imp3)
	ForeignFour(imp_four.Imp4)

	// A method that returns a nillable type.
	NillableRet() os.Error
}

// some random use of another package that isn't needed by the interface.
var _ os.Errno

// A function that we will test that uses the above interface.
// It takes a list of keys and values, and puts them in the index.
func Remember(index Index, keys []string, values []interface{}) {
	for i, k := range keys {
		index.Put(k, values[i])
	}
	err := index.NillableRet()
	if err != nil {
		log.Fatalf("Woah! %v", err)
	}
}
