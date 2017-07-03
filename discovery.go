package main

import (
	"encoding/json"
	"fmt"

	rabbithole "github.com/michaelklishin/rabbit-hole"
)

func discovery(
	rmqc *rabbithole.Client,
	queues []rabbithole.QueueInfo,
	aggGroup string,
) error {
	discoveryData := make(map[string][]map[string]string)

	var discoveredItems []map[string]string

	for _, queue := range queues {
		discoveredItem := make(map[string]string)
		discoveredItem["{#QUEUENAME}"] = queue.Name

		if aggGroup != "None" {

			aggregateItem := make(map[string]string)
			aggregateItem["{#GROUPNAME}"] = aggGroup
			aggregateItem["{#AGGQUEUENAME}"] = queue.Name

			discoveredItems = append(discoveredItems, aggregateItem)
		}

		discoveredItems = append(discoveredItems, discoveredItem)
	}

	discoveryData["data"] = discoveredItems

	out, err := json.Marshal(discoveryData)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", out)
	return nil
}
