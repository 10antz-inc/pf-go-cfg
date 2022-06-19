package errors

import (
	"github.com/tys-muta/go-ers"
	"google.golang.org/grpc/codes"
)

var (
	ErrNotFoundSubscription = ers.New(codes.NotFound, "NotFoundSubscription", "subscription is not found")
)
