package main

import (
	"fmt"
	"strconv"

	zsend "github.com/blacked/go-zabbix"
	rabbithole "github.com/michaelklishin/rabbit-hole"
)

func makePrefix(key string) string {
	return fmt.Sprintf(
		"ngs.srv.rabbitmq.%s", key,
	)
}

func getOverview(
	hostname string,
	overview *rabbithole.Overview,
	metrics []*zsend.Metric,
) []*zsend.Metric {
	metrics = append(
		metrics,
		zsend.NewMetric(
			hostname,
			makePrefix(
				"overview.consumers",
			),
			strconv.Itoa(overview.ObjectTotals.Consumers),
		),
	)
	metrics = append(
		metrics,
		zsend.NewMetric(
			hostname,
			makePrefix(
				"overview.connections",
			),
			strconv.Itoa(overview.ObjectTotals.Connections),
		),
	)
	metrics = append(
		metrics,
		zsend.NewMetric(
			hostname,
			makePrefix(
				"overview.exchanges",
			),
			strconv.Itoa(overview.ObjectTotals.Exchanges),
		),
	)
	metrics = append(
		metrics,
		zsend.NewMetric(
			hostname,
			makePrefix(
				"overview.queues",
			),
			strconv.Itoa(overview.ObjectTotals.Queues),
		),
	)
	metrics = append(
		metrics,
		zsend.NewMetric(
			hostname,
			makePrefix(
				"overview.channels",
			),
			strconv.Itoa(overview.ObjectTotals.Channels),
		),
	)

	return metrics
}

func getQueuesMetrics(
	hostname string,
	queues []rabbithole.QueueInfo,
	metrics []*zsend.Metric,
) []*zsend.Metric {

	for _, queue := range queues {

		metrics = append(
			metrics,
			zsend.NewMetric(
				hostname,
				makePrefix(
					fmt.Sprintf("queue.node[%s]", queue.Name),
				),
				queue.Node,
			),
		)

		metrics = append(
			metrics,
			zsend.NewMetric(
				hostname,
				makePrefix(
					fmt.Sprintf("queue.consumers[%s]", queue.Name),
				),
				strconv.Itoa(queue.Consumers),
			),
		)

		metrics = append(
			metrics,
			zsend.NewMetric(
				hostname,
				makePrefix(
					fmt.Sprintf("queue.memory_usage[%s]", queue.Name),
				),
				strconv.FormatInt(queue.Memory, 10),
			),
		)

		metrics = append(
			metrics,
			zsend.NewMetric(
				hostname,
				makePrefix(
					fmt.Sprintf("queue.messages[%s]", queue.Name),
				),
				strconv.Itoa(queue.Messages),
			),
		)

		metrics = append(
			metrics,
			zsend.NewMetric(
				hostname,
				makePrefix(
					fmt.Sprintf("queue.messages_ready[%s]", queue.Name),
				),
				strconv.Itoa(queue.MessagesReady),
			),
		)

		metrics = append(
			metrics,
			zsend.NewMetric(
				hostname,
				makePrefix(
					fmt.Sprintf(
						"queue.messages_unacknowledged[%s]",
						queue.Name,
					),
				),
				strconv.Itoa(queue.MessagesUnacknowledged),
			),
		)

		metrics = append(
			metrics,
			zsend.NewMetric(
				hostname,
				makePrefix(
					fmt.Sprintf(
						"queue.messages_stats.publish[%s]", queue.Name,
					),
				),
				strconv.FormatInt(queue.MessageStats.Publish, 10),
			),
		)

		metrics = append(
			metrics,
			zsend.NewMetric(
				hostname,
				makePrefix(
					fmt.Sprintf(
						"queue.messages_stats.deliver[%s]", queue.Name,
					),
				),
				strconv.FormatInt(queue.MessageStats.Deliver, 10),
			),
		)

		metrics = append(
			metrics,
			zsend.NewMetric(
				hostname,
				makePrefix(
					fmt.Sprintf(
						"queue.messages_stats.deliver_no_ack[%s]",
						queue.Name,
					),
				),
				strconv.FormatInt(queue.MessageStats.DeliverNoAck, 10),
			),
		)

		metrics = append(
			metrics,
			zsend.NewMetric(
				hostname,
				makePrefix(
					fmt.Sprintf(
						"queue.messages_stats.deliver_get[%s]",
						queue.Name,
					),
				),
				strconv.FormatInt(queue.MessageStats.DeliverGet, 10),
			),
		)

		metrics = append(
			metrics,
			zsend.NewMetric(
				hostname,
				makePrefix(
					fmt.Sprintf(
						"queue.messages_stats.redeliver[%s]",
						queue.Name,
					),
				),
				strconv.FormatInt(queue.MessageStats.Redeliver, 10),
			),
		)

		metrics = append(
			metrics,
			zsend.NewMetric(
				hostname,
				makePrefix(
					fmt.Sprintf(
						"queue.messages_stats.get[%s]",
						queue.Name,
					),
				),
				strconv.FormatInt(queue.MessageStats.Get, 10),
			),
		)

		metrics = append(
			metrics,
			zsend.NewMetric(
				hostname,
				makePrefix(
					fmt.Sprintf(
						"queue.messages_stats.get_no_ack[%s]",
						queue.Name,
					),
				),
				strconv.FormatInt(queue.MessageStats.GetNoAck, 10),
			),
		)
	}

	return metrics
}

func getNodeMetrics(
	hostname string,
	nodeInfo *rabbithole.NodeInfo,
	metrics []*zsend.Metric,
) []*zsend.Metric {

	metrics = append(
		metrics,
		zsend.NewMetric(
			hostname,
			makePrefix(
				"memory_limit",
			),
			strconv.Itoa(nodeInfo.MemLimit),
		),
	)
	metrics = append(
		metrics,
		zsend.NewMetric(
			hostname,
			makePrefix(
				"memory_used",
			),
			strconv.Itoa(nodeInfo.MemUsed),
		),
	)
	metrics = append(
		metrics,
		zsend.NewMetric(
			hostname,
			makePrefix(
				"memory_alarm",
			),
			strconv.FormatBool(nodeInfo.MemAlarm),
		),
	)

	return metrics
}
