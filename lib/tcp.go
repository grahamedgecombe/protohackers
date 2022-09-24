package protohackers

import (
	"context"
	"errors"
	"fmt"
	"net"

	"go.uber.org/zap"
)

type ServeTCPFunc func(ctx context.Context, logger *zap.Logger, conn net.Conn) error

func ListenAndServeTCP(ctx context.Context, logger *zap.Logger, addr string, serve ServeTCPFunc) error {
	listener, err := new(net.ListenConfig).Listen(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("protohackers: listen failed: %w", err)
	}
	defer listener.Close()

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()

		var nerr net.Error
		if errors.As(err, &nerr) && nerr.Timeout() {
			continue
		} else if err != nil {
			return fmt.Errorf("protohackers: accept failed: %w", err)
		}

		logger := logger.With(
			zap.Stringer("local", conn.LocalAddr()),
			zap.Stringer("remote", conn.RemoteAddr()),
		)

		go func() {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			go func() {
				<-ctx.Done()
				conn.Close()
			}()

			if err := serve(ctx, logger, conn); err != nil {
				logger.Warn("serve failed", zap.Error(err))
			}
		}()
	}
}
