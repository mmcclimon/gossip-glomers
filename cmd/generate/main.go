package main

import (
	"fmt"
	"log"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	pid := os.Getpid()
	i := 1

	n.Handle("generate", func(msg maelstrom.Message) error {
		body := map[string]any{
			"type": "generate_ok",
			"id":   fmt.Sprintf("%d.%d", pid, i),
		}

		i++

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
