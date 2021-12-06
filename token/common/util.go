package common

import (
	"time"

	prisma_types "github.com/aporeto-se/prisma-sdk-go-v2/types"
)

// TokenExpired returns error if token is expired otherwise nil
func TokenExpired(exp int64) error {
	if time.Now().Unix() > exp {
		return prisma_types.NewTokenExpiredError()
	}
	return nil
}
