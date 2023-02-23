package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"golang.org/x/exp/maps"
)

var node *maelstrom.Node
var messages map[int]bool

type Message = maelstrom.Message

func main() {
	node = maelstrom.NewNode()
	messages = make(map[int]bool)

	node.Handle("broadcast", broadcast)
	node.Handle("read", read)
	node.Handle("topology", topology)

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}

func broadcast(msg Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	value := int(body["message"].(float64))

	have := messages[value]

	if !have {
		// send it everywhere else
		for _, id := range node.NodeIDs() {
			if id == node.ID() {
				continue
			}

			node.RPC(id, body, nil)
		}

		messages[value] = true
	}

	return node.Reply(msg, map[string]string{"type": "broadcast_ok"})
}

func read(msg Message) error {
	resp := map[string]any{
		"type":     "read_ok",
		"messages": maps.Keys(messages),
	}

	return node.Reply(msg, resp)
}

func topology(msg Message) error {
	resp := map[string]string{
		"type": "topology_ok",
	}

	return node.Reply(msg, resp)
}
