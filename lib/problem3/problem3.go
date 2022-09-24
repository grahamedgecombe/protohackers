package problem3

import (
	"bufio"
	"container/list"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
	"sync"

	protohackers "github.com/grahamedgecombe/protohackers/lib"
	"go.uber.org/zap"
)

var nameRegex = regexp.MustCompile(`(?i)^[a-z0-9]+$`)

type client struct {
	conn    net.Conn
	name    string
	element *list.Element
}

type room struct {
	clients list.List
	mutex   sync.Mutex
}

func (r *room) Join(c *client) {
	message := fmt.Sprintf("* %s has entered the room\n", c.name)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	front := r.clients.Front()

	var b strings.Builder
	if front != nil {
		b.WriteString("* The room contains: ")
	} else {
		b.WriteString("* The room is empty")
	}

	for e := front; e != nil; e = e.Next() {
		otherC := e.Value.(*client)
		_, _ = io.WriteString(otherC.conn, message)

		b.WriteString(otherC.name)
		if e.Next() != nil {
			b.WriteString(", ")
		}
	}

	b.WriteRune('\n')

	_, _ = io.WriteString(c.conn, b.String())

	c.element = r.clients.PushBack(c)
}

func (r *room) Leave(c *client) {
	message := fmt.Sprintf("* %s has left the room\n", c.name)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.clients.Remove(c.element)

	for e := r.clients.Front(); e != nil; e = e.Next() {
		otherC := e.Value.(*client)
		_, _ = io.WriteString(otherC.conn, message)
	}
}

func (r *room) Broadcast(c *client, message string) {
	message = fmt.Sprintf("[%s] %s\n", c.name, message)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	for e := r.clients.Front(); e != nil; e = e.Next() {
		otherC := e.Value.(*client)
		if c != otherC {
			_, _ = io.WriteString(otherC.conn, message)
		}
	}
}

func Solve(ctx context.Context, logger *zap.Logger, addr string) error {
	var r room

	return protohackers.ListenAndServeTCP(ctx, logger, addr, func(ctx context.Context, logger *zap.Logger, conn net.Conn) error {
		if _, err := io.WriteString(conn, "What is your name?\n"); err != nil {
			return fmt.Errorf("problem3: write failed: %w", err)
		}

		scanner := bufio.NewScanner(conn)
		if !scanner.Scan() {
			return nil
		} else if err := scanner.Err(); err != nil {
			return fmt.Errorf("problem3: scan failed: %w", err)
		}

		name := scanner.Text()
		if !nameRegex.MatchString(name) {
			return errors.New("problem3: invalid name")
		}

		c := client{
			conn: conn,
			name: name,
		}

		r.Join(&c)
		defer r.Leave(&c)

		for scanner.Scan() {
			r.Broadcast(&c, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("problem3: scan failed: %w", err)
		}

		return nil
	})
}
