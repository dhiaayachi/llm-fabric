package discoverer

import (
	"context"
)

type Discoverer interface {
	Join(ctx context.Context, addresses []string) error
}
