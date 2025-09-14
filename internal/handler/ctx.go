package handler

import (
	"context"
	"strconv"

	"google.golang.org/grpc/metadata"
)

// GetUserIDFromCtx извлекает user_id из gRPC metadata.
func GetUserIDFromCtx(ctx context.Context) int64 {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0
	}

	ids := md.Get("user_id")
	if len(ids) == 0 {
		return 0
	}

	userID, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		return 0
	}

	return userID
}
