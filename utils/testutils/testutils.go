package testutils

import (
	"fmt"
	"math/rand"
	"time"
)

func SocketFileName() string {
	rand.Seed(time.Now().Unix())
	return "/tmp/madnet-test-" + fmt.Sprint(rand.Int())
}

func WaitUntil(f func() bool) {
	end := time.Now().Add(5 * time.Second)
	for time.Now().Before(end) {
		time.Sleep(time.Millisecond)
		if f() {
			time.Sleep(10 * time.Millisecond)
			return
		}
	}
	panic("timeout")
}

func TestAsync(f func() error) <-chan error {
	errchan := make(chan error)
	go func() {
		defer close(errchan)
		errchan <- f()
	}()

	return errchan
}
