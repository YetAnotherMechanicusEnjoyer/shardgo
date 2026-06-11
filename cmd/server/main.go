package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/YetAnotherMechanicusEnjoyer/shardgo/internal/cache"
	"github.com/YetAnotherMechanicusEnjoyer/shardgo/internal/cluster"
	"github.com/YetAnotherMechanicusEnjoyer/shardgo/internal/network"
)

func main() {
	addr := flag.String("addr", ":6379", "Server address")
	isMaster := flag.Bool("master", false, "Is the node Master")
	nodes := flag.String("nodes", "127.0.0.1:6379,127.0.0.1:6380,127.0.0.1:6381", "List cluster nodes (separator: ',')")
	replicas := flag.Int("replicas", 100, "Number of virtual nodes per node")
	flag.Parse()

	if strings.HasPrefix(*addr, ":") {
		*addr = "127.0.0.1" + *addr
	}

	nodeList := strings.Split(*nodes, ",")

	c := cache.New()

	clusterManager := cluster.NewManager(*replicas)
	for _, node := range nodeList {
		clusterManager.AddNode(node)
	}

	server := network.NewServer(*addr, c, clusterManager, *isMaster)
	log.Printf("ShardGo server starting on %s (Master: %v)...\n", *addr, *isMaster)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down...")
		os.Exit(0)
	}()

	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
