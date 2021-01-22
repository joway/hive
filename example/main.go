package main

import (
	"context"
	"github.com/joway/hive"
	"log"
	"time"
)

func main() {
	h, err := hive.New(
		hive.WithSize(10),
		hive.WithNonblocking(false),
	)
	if err != nil {
		log.Fatal(err)
	}
	err = h.Submit(context.Background(), func() {
		//do something
		time.Sleep(time.Millisecond * 100)
	})
	if err != nil {
		log.Fatal(err)
	}
}
