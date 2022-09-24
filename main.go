package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	protohackers "github.com/grahamedgecombe/protohackers/lib"
	"github.com/grahamedgecombe/protohackers/lib/problem0"
	"go.uber.org/zap"
)

var solvers = []protohackers.SolveFunc{
	problem0.Solve,
}

func main() {
	var (
		addr    = flag.String("addr", ":10000", "")
		problem = flag.Int("problem", len(solvers)-1, "")
	)
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer stop()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("protohackers: failed to create zap logger: %v\n", err)
	}
	defer logger.Sync()

	if err := solvers[*problem](ctx, logger, *addr); err != nil {
		logger.Fatal("solve failed", zap.Error(err))
	}
}
