package main

import (
	"fmt"
	"os"
	"strconv"

	zsend "github.com/blacked/go-zabbix"
	docopt "github.com/docopt/docopt-go"
	rabbithole "github.com/michaelklishin/rabbit-hole"
)

var version = "[manual build]"

func main() {
	usage := `zabbix-extension-rabbitmq

Usage:
    zabbix-extension-rabbitmq [-r <address>] [-u <name>] [-s <password>] [-z <host>] [-p <number>] [-d [-g <name>]]

RabbitMQ options:
    -r --rabbitmq <address>          Listen address of rabbitmq server [default: 127.0.0.1:15672]
    -u --rabbitmq-user <name>        Rabbitmq management username [default: guest]
    -s --rabbitmq-secret <password>  Rabbitmq management password [default: guest]

Zabbix options:
    -z --zabbix <host>         Hostname or IP address of zabbix server [default: 127.0.0.1]
    -p --zabbix-port <number>  Port of zabbix server [default: 10051]
	-d --discovery             Run low-level discovery for determine queues, exchanges, etc.
	-g --group <name>          Group name which will be use for aggregate item values.[default: None]

Misc options:
    -h --help                  Show this screen.
	-v --version               Show version.
`

	args, _ := docopt.Parse(usage, nil, true, version, false)

	aggGroup := args["--group"].(string)

	zabbixHost := args["--zabbix"].(string)
	zabbixPort, err := strconv.Atoi(args["--zabbix-port"].(string))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var metrics []*zsend.Metric

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	rmqc, err := rabbithole.NewClient(
		fmt.Sprintf("http://%s", args["--rabbitmq"].(string)),
		args["--rabbitmq-user"].(string),
		args["--rabbitmq-secret"].(string),
	)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	queues, err := rmqc.ListQueues()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	overview, err := rmqc.Overview()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	nodeInfo, err := rmqc.GetNode(
		overview.Node,
	)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if args["--discovery"].(bool) {

		err = discovery(rmqc, queues, aggGroup)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	metrics = getQueuesMetrics(
		hostname,
		queues,
		metrics,
	)

	metrics = getOverview(
		hostname,
		overview,
		metrics,
	)

	metrics = getNodeMetrics(
		hostname,
		nodeInfo,
		metrics,
	)

	packet := zsend.NewPacket(metrics)
	sender := zsend.NewSender(
		zabbixHost,
		zabbixPort,
	)
	sender.Send(packet)

	fmt.Println("OK")
}
