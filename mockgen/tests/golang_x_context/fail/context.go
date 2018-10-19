package fail

import "golang.org/x/net/context"

// Context generates a package collision when this file is parsed
// from another file
type Context interface {
	Get() context.Context
}
