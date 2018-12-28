package dockertest

import (
	"context"
	"database/sql"
	"testing"
)

func TestContainerPool(t *testing.T) {
	pool, err := NewContainerPool(3, "postgres:alpine", "5432", func(addr string) error {
		db, err := sql.Open("postgres", "postgres://postgres:postgres@"+addr+"?sslmode=disable")
		if err != nil {
			return err
		}
		return db.Ping()
	})
	if err != nil {
		t.Fatalf("could not start postgres, %s", err)
	}

	c, err := pool.GetContainer(context.Background())
	if err != nil {
		t.Fatalf("could not get container, %s", err)
	}

	pool.ReleaseContainer(c)

	if count := checkDocker(t, c); count != 3 {
		t.Fatal("container did not start or died early", count)
	}

	pool.Shutdown()

	if count := checkDocker(t, c); count != 0 {
		t.Fatal("container is still running after shutdown", count)
	}
}
