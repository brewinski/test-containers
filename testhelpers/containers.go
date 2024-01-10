package testhelpers

import (
	"context"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/k6"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

type ContainerTerminator interface {
	Terminate(ctx context.Context) error
	GetContainerID() string
}

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func CreatePostgresContainer(ctx context.Context, initScripts ...string) (*PostgresContainer, error) {
	if len(initScripts) == 0 {
		initScripts = []string{filepath.Join("..", "testdata", "init-db.sql")}
	}

	mountTarget, err := filepath.Abs("./testdata/db")
	if err != nil {
		return nil, err
	}

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts(initScripts...),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Mounts: testcontainers.ContainerMounts{
					{
						Source: testcontainers.GenericVolumeMountSource{
							Name: "test",
						},
						Target: testcontainers.ContainerMountTarget(mountTarget),
					},
				},
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}, nil
}

type RedisContainer struct {
	*redis.RedisContainer
	ConnectionString string
}

func CreateRedisContainer(ctx context.Context) (*RedisContainer, error) {
	redisContainer, err := redis.RunContainer(
		ctx,
		testcontainers.WithImage("redis:7"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("* Ready to accept connections tcp"),
		),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}

	return &RedisContainer{
		RedisContainer:   redisContainer,
		ConnectionString: connStr,
	}, nil
}

func CreateK6Container(ctx context.Context) (*k6.K6Container, error) {
	k6Container, err := k6.RunContainer(ctx, testcontainers.WithImage("szkiba/k6x:v0.3.1"), k6.WithTestScript(filepath.Join("./", "testdata", "k6", "load-test.js")))
	if err != nil {
		return nil, err
	}

	return k6Container, nil
}
