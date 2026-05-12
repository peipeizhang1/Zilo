package userctx

import (
	"context"
	"errors"
	"fmt"
	"strconv"
)

func GetUserIDFromCtx(ctx context.Context) (int64, error) {
	value := ctx.Value("userId")
	if value == nil {
		return 0, errors.New("missing userId in token")
	}

	switch v := value.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case string:
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, errors.New("invalid userId in token")
		}
		return id, nil
	default:
		return 0, fmt.Errorf("unsupported userId type: %T", value)
	}
}
