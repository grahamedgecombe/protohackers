package protohackers

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
)

type ServeUDPFunc func(ctx context.Context, logger *zap.Logger, conn net.PacketConn, addr net.Addr, packet []byte) error

func ListenAndServeUDP(ctx context.Context, logger *zap.Logger, addr string, maxLen int, serve ServeUDPFunc) error {
	conn, err := new(net.ListenConfig).ListenPacket(ctx, "udp", addr)
	if err != nil {
		return fmt.Errorf("protohackers: listen failed: %w", err)
	}
	defer conn.Close()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	for {
		buf := make([]byte, maxLen)

		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			return fmt.Errorf("protohackers: read failed: %w", err)
		}

		logger := logger.With(
			zap.Stringer("local", conn.LocalAddr()),
			zap.Stringer("remote", addr),
		)

		if err := serve(ctx, logger, conn, addr, buf[:n]); err != nil {
			return fmt.Errorf("protohackers: serve failed: %w", err)
		}
	}
}
