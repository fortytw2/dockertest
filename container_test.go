package dockertest

import (
	"database/sql"
	"os/exec"
	"strings"
	"testing"
	"time"

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
	buf, err := exec.Command("docker", "ps").CombinedOutput()
	if err != nil {
		t.Fatalf("could not docker ps, %s", err)
	}

	lines := strings.Split(strings.TrimRight(string(buf), "\n"), "\n")
	count := 0
	for _, line := range lines {
		if strings.Contains(line, container.Name) {
			count++
		}
	}
	if count != 1 {
		t.Fatal("container did not start or died early", len(lines))
	}

	container.Shutdown()
	time.Sleep(1 * time.Second) // takes a bit for the docker daemon to realize that the process has ended

	buf, err = exec.Command("docker", "ps").CombinedOutput()
	if err != nil {
		t.Fatalf("could not docker ps, %s", err)
	}

	lines = strings.Split(strings.TrimRight(string(buf), "\n"), "\n")
	count = 0
	for _, line := range lines {
		if strings.Contains(line, container.Name) {
			count++
		}
	}
	if count != 0 {
		t.Fatal("container is still running after shutdown", count)
	}
}
