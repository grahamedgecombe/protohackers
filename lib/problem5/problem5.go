package problem5

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"regexp"

	protohackers "github.com/grahamedgecombe/protohackers/lib"
	"go.uber.org/zap"
)

var (
	bogusCoinRegex   = regexp.MustCompile(`(?i)^7[a-z0-9]{25,34}$`)
	bogusCoinAddress = []byte("7YWHMfk9JZe0LM0g1ZauHuiSxhI")
	space            = []byte(" ")
)

func Solve(ctx context.Context, logger *zap.Logger, addr string) error {
	return protohackers.ListenAndServeTCP(ctx, logger, addr, func(ctx context.Context, logger *zap.Logger, conn net.Conn) error {
		remote, err := new(net.Dialer).DialContext(ctx, "tcp", "chat.protohackers.com:16963")
		if err != nil {
			return fmt.Errorf("problem5: dial failed: %w", err)
		}
		defer remote.Close()

		go func() {
			defer conn.Close()

			if err := proxy(remote, conn); err != nil {
				logger.Warn("downstream proxy failed", zap.Error(err))
			}
		}()

		if err := proxy(conn, remote); err != nil {
			return fmt.Errorf("problem5: upstream proxy failed: %w", err)
		}

		return nil
	})
}

func proxy(r io.Reader, w io.Writer) error {
	br := bufio.NewReader(r)
	for {
		var eof bool

		line, err := br.ReadBytes('\n')
		if err == io.EOF {
			eof = true
		} else if err != nil {
			return fmt.Errorf("problem5: read failed: %w", err)
		} else {
			line = line[:len(line)-1]
		}

		words := bytes.Split(line, space)

		for i, word := range words {
			if bogusCoinRegex.Match(word) {
				words[i] = bogusCoinAddress
			}
		}

		b := bytes.Join(words, space)
		if !eof {
			b = append(b, '\n')
		}

		if _, err := w.Write(b); err != nil {
			return fmt.Errorf("problem5: write failed: %w", err)
		}

		if eof {
			return nil
		}
	}
}
