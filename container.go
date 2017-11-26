package dockertest

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// A Container is a container inside docker
type Container struct {
	Name string
	Args []string
	Addr string

	cmd *exec.Cmd
}

// Shutdown ends the container
func (c *Container) Shutdown() {
	c.cmd.Process.Signal(syscall.SIGINT)
	// Check if the process is excited.  If it doesn't exit in a reasonable period of time, just continue
	for i := 0; i < 200 && !(c.cmd == nil || c.cmd.ProcessState == nil || c.cmd.ProcessState.Exited()); i++ {
		time.Sleep(5 * time.Millisecond)
	}
}

// RunContainer runs a given docker container and returns a port on which the
// container can be reached
func RunContainer(container string, port string, waitFunc func(addr string) error, args ...string) (*Container, error) {
	free := freePort()
	host := getHost()
	addr := fmt.Sprintf("%s:%d", host, free)
	argsFull := append(args, "run", "-p", fmt.Sprintf("%d:%s", free, port), container)
	cmd := exec.Command("docker", argsFull...)
	// run this in the background

	start := make(chan struct{})
	go func() {
		err := cmd.Start()
		if err != nil {
			fmt.Printf("could not run container, %s\n", err)
		}
		start <- struct{}{}
		cmd.Wait()
	}()

	<-start
	for {
		err := waitFunc(addr)
		if err == nil {
			break
		}

		time.Sleep(time.Millisecond * 150)
	}

	return &Container{
		Name: container,
		Addr: addr,
		Args: args,
		cmd:  cmd,
	}, nil
}

func getHost() string {
	out, err := exec.Command("docker-machine", "ip", os.Getenv("DOCKER_MACHINE_NAME")).Output()
	if err == nil {
		return strings.TrimSpace(string(out[:]))
	}
	return "localhost"
}

func freePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port
}
