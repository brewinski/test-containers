package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/testcontainers/testcontainers-go-demo/testhelpers"
)

func StopContainer(ctx context.Context, container testhelpers.ContainerTerminator, containerType string) error {
	name := container.GetContainerID()[:12]

	slog.Info("Stopping container", "name", name, "type", containerType)

	err := container.Terminate(ctx)
	if err != nil {
		slog.Error("Error stopping container", "name", name, "error", err, "type", containerType)
		return err
	}

	slog.Info("Container stopped", "name", name, "note", containerType)
	return nil
}

func main() {
	// start postgres container
	ctx := context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(ctx, filepath.Join(".", "testdata", "init-db.sql"))
	if err != nil {
		log.Fatal(err)
	}

	defer StopContainer(ctx, pgContainer, "PostgreSQL")
	slog.Info("Postgres container started successfully")

	redisContainer, err := testhelpers.CreateRedisContainer(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer StopContainer(ctx, redisContainer, "Redis")
	slog.Info("Redis container started successfully")

	// k6Container, err := testhelpers.CreateK6Container(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer StopContainer(ctx, k6Container, "K6")
	// slog.Info("K6 container started successfully")

	slog.Info("Containers running, connection strings", "redis", redisContainer.ConnectionString, "postgres", pgContainer.ConnectionString)
	// run the service... (http, grpc, cron... etc)

	sigKill := make(chan os.Signal, 1)
	signal.Notify(sigKill, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigKill:
			slog.Info("Received termination signal, shutting down...")
			return
		default:
			// TODO: implement service logic here
			// e.g. start HTTP server, run cron jobs etc
			time.Sleep(5 * time.Second)
			slog.Info("Containers running, connection strings", "redis", redisContainer.ConnectionString, "postgres", pgContainer.ConnectionString)
		}
	}
}
