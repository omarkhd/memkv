package main

import (
	"log"

	"omarkhd/memkv/server"
)

func main() {
	s, err := server.New()
	if err != nil {
		panic(err.Error())
	}

	log.Print("Starting omarkhd/memkv")
	go s.Start()
	select {}
}
