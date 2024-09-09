package servers

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrInternal = status.New(codes.Internal, "internal server error").Err()
