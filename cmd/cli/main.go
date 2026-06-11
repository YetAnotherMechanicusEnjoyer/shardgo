package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/YetAnotherMechanicusEnjoyer/shardgo/internal/cluster"
	"github.com/YetAnotherMechanicusEnjoyer/shardgo/internal/network"
)

func main() {
	nodes := flag.String("nodes", "127.0.0.1:6379,127.0.0.1:6380,127.0.0.1:6381", "List cluster nodes (separator: ',')")
	replicas := flag.Int("replicas", 100, "Number of virtual nodes per node")
	flag.Parse()

	nodeList := strings.Split(*nodes, ",")

	clusterManager := cluster.NewManager(*replicas)
	for _, node := range nodeList {
		clusterManager.AddNode(node)
	}

	firstNode := nodeList[0]

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("ShardGo CLI (cluster: %v)\n", nodeList)
	fmt.Println("Command: SET key value [ttl], GET key, DEL key, ADD_NODE addr, EXIT")

	for {
		fmt.Print("> ")
		if !scanner.Scan() || scanner.Err() != nil {
			break
		}

		input := scanner.Text()
		if strings.ToLower(input) == "exit" {
			break
		}

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		if parts[0] == "STATS" {
			totalKeys := 0
			for _, node := range nodeList {
				client := network.NewClient(node)
				resp, err := client.SendRequest("STATS")
				if err != nil {
					fmt.Printf("Error getting stats from %s: %v\n", node, err)
					continue
				}
				if len(resp) > 1 && resp[0] == ':' {
					numStr := resp[1:]
					keys, err := strconv.Atoi(numStr)
					if err != nil {
						fmt.Printf("Error parsing stats from %s: %v\n", node, err)
						continue
					}
					fmt.Printf("%s: %d keys\n", node, keys)
					totalKeys += keys
				} else {
					fmt.Printf("%s: invalid response '%s'\n", node, resp)
				}
			}
			fmt.Printf("Total keys in cluster: %d\n", totalKeys)
			continue
		}

		if parts[0] == "NODE" && len(parts) >= 2 {
			key := parts[1]
			node, ok := clusterManager.GetNodeForKey(key)
			if !ok {
				fmt.Printf("No node found for key %s\n", key)
			} else {
				fmt.Printf("Key '%s' is on node: %s\n", key, node)
			}
			continue
		}

		if parts[0] == "ADD_NODE" && len(parts) >= 2 {
			node := parts[1]
			client := network.NewClient(firstNode)
			resp, err := client.SendRequest(input)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Println(resp)

			clusterManager.AddNode(node)
			continue
		}

		key := ""
		if len(parts) >= 2 {
			key = parts[1]
		}

		node, ok := clusterManager.GetNodeForKey(key)
		if !ok || node == "" {
			node = firstNode
		}

		client := network.NewClient(node)
		resp, err := client.SendRequest(input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Println(resp)
	}
}
