package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/YetAnotherMechanicusEnjoyer/shardgo/internal/cluster"
)

func main() {
	n := flag.Int("n", 1000, "Number of keys to generate")
	flag.Parse()

	nodes := []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"}
	ch := cluster.NewConsistentHash(100)
	for _, node := range nodes {
		ch.AddNode(node)
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	nodeCounts := make(map[string]int)
	for range *n {
		key := fmt.Sprintf("key%d", rand.Intn(1000000))
		node, _ := ch.GetNode(key)
		nodeCounts[node]++
	}

	for node, count := range nodeCounts {
		bar := strings.Repeat("█", count*30 / *n)
		fmt.Printf("%s: %d keys (%.1f%%) %s\n", node, count, float64(count)/10, bar)
	}
}
