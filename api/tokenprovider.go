package prismasdk2

import (
	"context"
)

// TokenProvider is an interface for getting a new PrismaToken.
// This is done differently based on the cloud provider, etc.
type TokenProvider interface {
	Token(context.Context) (string, error)
	AccountID(ctx context.Context) (string, error)
}
