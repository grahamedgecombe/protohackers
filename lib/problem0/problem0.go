package problem0

import (
	"context"
	"fmt"
	"io"
	"net"

	protohackers "github.com/grahamedgecombe/protohackers/lib"
	"go.uber.org/zap"
)

func Solve(ctx context.Context, logger *zap.Logger, addr string) error {
	return protohackers.ListenAndServeTCP(ctx, logger, addr, func(ctx context.Context, logger *zap.Logger, conn net.Conn) error {
		if _, err := io.Copy(conn, conn); err != nil {
			return fmt.Errorf("problem0: copy failed: %w", err)
		}

		return nil
	})
}
