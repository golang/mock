package pass

import "golang.org/x/net/context"

// Context parses correctly when using go context and golang.org/x/net/context
// included in this package and called from another file in another package
type Context interface {
	Get() context.Context
}
