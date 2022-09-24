package protohackers

import (
	"context"

	"go.uber.org/zap"
)

type SolveFunc func(ctx context.Context, logger *zap.Logger, addr string) error
