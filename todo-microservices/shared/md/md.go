package md

import (
	"context"
	"errors"
	"strconv"

	"google.golang.org/grpc/metadata"
)

const metadataKeyUserID string = "x-user-id"

var ErrNotFoundUserID = errors.New("not found user id")

func GetUserIDFromContext(ctx context.Context) uint64 {
	userID, err := safeGetUserIDFromContext(ctx)
	if err != nil {
		panic(err)
	}

	return userID
}

func safeGetUserIDFromContext(ctx context.Context) (uint64, error) {
	var userID uint64
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return userID, ErrNotFoundUserID
	}

	values := md.Get(metadataKeyUserID)
	if len(values) == 0 {
		return userID, ErrNotFoundUserID
	}

	userID, err := strconv.ParseUint(values[0], 10, 64)
	if err != nil {
		return userID, err
	}

	return userID, nil
}
