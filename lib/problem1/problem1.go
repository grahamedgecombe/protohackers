package problem1

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"strings"

	protohackers "github.com/grahamedgecombe/protohackers/lib"
	"go.uber.org/zap"
)

type request struct {
	Method string           `json:"method"`
	Number *json.RawMessage `json:"number"`
}

type response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func isPrime(line []byte) (bool, error) {
	var r request
	if err := json.Unmarshal(line, &r); err != nil {
		return false, fmt.Errorf("problem1: unmarshal failed: %w", err)
	}

	if r.Method != "isPrime" {
		return false, errors.New("problem1: invalid method name")
	} else if r.Number == nil {
		return false, errors.New("problem1: number missing")
	}

	s := strings.SplitN(string(*r.Number), ".", 2)
	if len(s) == 2 {
		return false, nil
	}

	n, ok := new(big.Int).SetString(s[0], 10)
	if !ok {
		return false, errors.New("problem1: number invalid")
	}
	return n.ProbablyPrime(80), nil
}

func Solve(ctx context.Context, logger *zap.Logger, addr string) error {
	return protohackers.ListenAndServeTCP(ctx, logger, addr, func(ctx context.Context, logger *zap.Logger, conn net.Conn) error {
		scanner := bufio.NewScanner(conn)

		for scanner.Scan() {
			prime, err := isPrime(scanner.Bytes())
			if err != nil {
				logger.Warn("request invalid", zap.Error(err))

				if _, err := io.WriteString(conn, "!\n"); err != nil {
					return fmt.Errorf("problem1: write failed: %w", err)
				}

				return nil
			}

			b, err := json.Marshal(&response{
				Method: "isPrime",
				Prime:  prime,
			})
			if err != nil {
				return fmt.Errorf("problem1: marshal failed: %w", err)
			}

			b = append(b, '\n')

			if _, err := conn.Write(b); err != nil {
				return fmt.Errorf("problem1: write failed: %w", err)
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("problem1: scan failed: %w", err)
		}

		return nil
	})
}
