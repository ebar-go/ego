package wait

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestUntil(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	go func() {
		time.Sleep(time.Second * 5)
		cancel()
	}()
	Until(func() {
		log.Printf("test")
	}, time.Second, ctx.Done())
}

func TestUntilWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	go func() {
		time.Sleep(time.Second * 5)
		cancel()
	}()
	UntilWithContext(func() {
		log.Printf("test")
	}, time.Second, ctx)
}
