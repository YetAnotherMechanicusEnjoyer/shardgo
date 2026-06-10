package main

import (
	"log"

	"github.com/YetAnotherMechanicusEnjoyer/shardgo/internal/cache"
	"github.com/YetAnotherMechanicusEnjoyer/shardgo/internal/network"
)

func main() {
	c := cache.New()
	server := network.NewServer(":6379", c)
	log.Println("ShardGo server starting on :6379...")
	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
