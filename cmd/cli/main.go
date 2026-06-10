package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/YetAnotherMechanicusEnjoyer/shardgo/internal/network"
)

func main() {
	addr := "localhost:6379"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	client := network.NewClient(addr)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("ShardGo CLI (connected to %s)\n", addr)
	fmt.Println("Command: SET key value [ttl], GET key, DEL key, EXIT")

	for {
		fmt.Print("> ")
		if !scanner.Scan() || scanner.Err() != nil {
			break
		}
		input := scanner.Text()
		if strings.ToLower(input) == "exit" {
			break
		}

		resp, err := client.SendRequest(input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Println(resp)
	}
}
