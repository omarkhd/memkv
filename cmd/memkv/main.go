package main

import (
	"log"

	"omarkhd/memkv/metrics"
	"omarkhd/memkv/server"
	"omarkhd/memkv/store"
)

func main() {
	go metrics.Expose()

	// Creating data store
	ds, err := store.New()
	if err != nil {
		panic(err.Error())
	}
	// Creating http server
	s, err := server.New(ds)
	if err != nil {
		panic(err.Error())
	}
	log.Print("Starting omarkhd/memkv")
	go s.Start()
	select {}
}
