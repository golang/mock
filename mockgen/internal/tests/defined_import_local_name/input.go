package defined_import_local_name

import (
	"bytes"
	"context"
)

//go:generate mockgen -package defined_import_local_name -destination mock.go -source input.go -imports b_mock=bytes,c_mock=context

type WithImports interface {
	Method1() bytes.Buffer
	Method2() context.Context
}
