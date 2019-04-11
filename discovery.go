package main

import (
	"encoding/json"
	"fmt"

	rabbithole "github.com/michaelklishin/rabbit-hole"
)

func discovery(
	rmqc *rabbithole.Client,
	allQueues map[string][]rabbithole.QueueInfo,
	aggGroup string,
	aggregate bool,
) error {
	discoveryData := make(map[string][]map[string]string)

	var discoveredItems []map[string]string

	for vhost, queues := range allQueues {

		for _, queue := range queues {
			discoveredItem := make(map[string]string)
			if aggregate {

				discoveredItem["{#GROUPNAME}"] = aggGroup
				discoveredItem["{#AGGQUEUENAME}"] = queue.Name
				discoveredItem["{#AGGVHOST}"] = vhost

				discoveredItems = append(discoveredItems, discoveredItem)

				continue
			}

			discoveredItem["{#QUEUENAME}"] = queue.Name
			discoveredItem["{#VHOST}"] = vhost

			discoveredItems = append(discoveredItems, discoveredItem)
		}
	}

	discoveryData["data"] = discoveredItems

	out, err := json.Marshal(discoveryData)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", out)
	return nil
}
