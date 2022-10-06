package problem4

import (
	"bytes"
	"context"
	"fmt"
	"net"

	protohackers "github.com/grahamedgecombe/protohackers/lib"
	"go.uber.org/zap"
)

var version = []byte("Ken's Key-Value Store 1.0")

func Solve(ctx context.Context, logger *zap.Logger, addr string) error {
	m := make(map[string][]byte)

	return protohackers.ListenAndServeUDP(ctx, logger, addr, 1000, func(ctx context.Context, logger *zap.Logger, conn net.PacketConn, addr net.Addr, buf []byte) error {
		sep := bytes.IndexByte(buf, '=')
		if sep != -1 {
			// insert
			k := string(buf[:sep])
			m[k] = buf[sep+1:]
		} else {
			// retrieve
			k := string(buf)

			var (
				v  []byte
				ok bool
			)
			if k == "version" {
				v = version
				ok = true
			} else {
				v, ok = m[k]
			}

			resp := append(buf, '=')
			if ok {
				resp = append(resp, v...)
			}

			if _, err := conn.WriteTo(resp, addr); err != nil {
				return fmt.Errorf("protohackers: write failed: %w", err)
			}
		}

		return nil
	})
}
