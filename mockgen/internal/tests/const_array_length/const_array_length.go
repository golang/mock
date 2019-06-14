//go:generate mockgen -aux_files consts=consts/consts.go -destination const_array_length_mock.go -package const_array_length -source const_array_length.go -self_package github.com/golang/mock/mockgen/internal/tests/const_array_length
package const_array_length

import "github.com/golang/mock/mockgen/internal/tests/const_array_length/consts"

const C = 2

type ConstArrayLength interface {
	consts.AuxInterface
	PackagePrefixConstLength() [consts.C]int
	ConstLength() [C]int
	LiteralLength() [3]int
	SliceLength() []int
}
