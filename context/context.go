package contextext

import "context"

type detachedContext struct {
	parent context.Context
}
