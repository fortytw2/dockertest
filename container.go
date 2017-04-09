package dockertest

import (
	"fmt"
	"net"
	"os/exec"
	"syscall"
	"time"
)

// A Container is a container inside docker
type Container struct {
	Name string
	Addr string

	cmd *exec.Cmd
}

// Shutdown ends the container
func (c *Container) Shutdown() {
	c.cmd.Process.Signal(syscall.SIGINT)
	time.Sleep(250 * time.Millisecond)
}

// RunContainer runs a given docker container and returns a port on which the
// container can be reached
func RunContainer(container string, port string, waitFunc func(addr string) error) (*Container, error) {
	free := freePort()
	addr := fmt.Sprintf("localhost:%d", free)
	cmd := exec.Command("docker", "run", "-p", fmt.Sprintf("%d:%s", free, port), container)
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
		cmd:  cmd,
	}, nil
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
