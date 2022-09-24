package problem2

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	protohackers "github.com/grahamedgecombe/protohackers/lib"
	"go.uber.org/zap"
)

type request struct {
	Type       uint8
	Arg1, Arg2 int32
}

type price struct {
	Time  int32
	Value int32
}

func Solve(ctx context.Context, logger *zap.Logger, addr string) error {
	return protohackers.ListenAndServeTCP(ctx, logger, addr, func(ctx context.Context, logger *zap.Logger, conn net.Conn) error {
		var prices []price

		for {
			var r request
			if err := binary.Read(conn, binary.BigEndian, &r); err != nil {
				return fmt.Errorf("problem2: read failed: %w", err)
			}

			switch r.Type {
			case 'I':
				prices = append(prices, price{
					Time:  r.Arg1,
					Value: r.Arg2,
				})

			case 'Q':
				var sum, n int64

				for _, p := range prices {
					if p.Time >= r.Arg1 && p.Time <= r.Arg2 {
						sum += int64(p.Value)
						n++
					}
				}

				var mean int32
				if n != 0 {
					mean = int32(sum / n)
				}

				if err := binary.Write(conn, binary.BigEndian, mean); err != nil {
					return fmt.Errorf("problem2: write failed: %w", err)
				}

			default:
				return errors.New("problem2: invalid type")
			}
		}
	})
}
