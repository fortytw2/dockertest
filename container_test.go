package dockertest

import (
	"database/sql"
	"os/exec"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

func TestRunContainer(t *testing.T) {
	container, err := RunContainer("postgres:alpine", "5432", func(addr string) error {
		db, err := sql.Open("postgres", "postgres://postgres:postgres@"+addr+"?sslmode=disable")
		if err != nil {
			return err
		}

		return db.Ping()
	})
	if err != nil {
		t.Fatalf("could not start postgres, %s", err)
	}

	container.Shutdown()

	buf, err := exec.Command("docker", "ps").CombinedOutput()
	if err != nil {
		t.Fatalf("could not docker ps, %s", err)
	}

	lines := strings.Split(strings.TrimRight(string(buf), "\n"), "\n")
	if len(lines) != 1 {
		t.Fatal("container is still running after shutdown", len(lines))
	}
}
