package main

import (
	"fmt"
	"os"
	"strconv"

	zsend "github.com/blacked/go-zabbix"
	docopt "github.com/docopt/docopt-go"
	rabbithole "github.com/michaelklishin/rabbit-hole"
)

var (
	version   = "[manual build]"
	noneValue = "None"
)

func main() {
	usage := `zabbix-agent-extension-rabbitmq

Usage:
    zabbix-agent-extension-rabbitmq [-r <address>] [-u <name>] [-s <password>] [-c <path>] [-z <host>] [-p <number>] [-o <name>] [-t <timeout>] [-d [-g <name>] [-a]]

RabbitMQ options:
    -r --rabbitmq <address>          Listen address of RabbitMQ server [default: http://127.0.0.1:15672]
    -u --rabbitmq-user <name>        RabbitMQ management username [default: guest]
	-s --rabbitmq-secret <password>  RabbitMQ management password [default: guest]
	-t --rabbitmq-timeout <timeout>  RabbitMQ request timeout in ms [default: 5000]
    -c --ca <path>                   Path to CA file. [default: ` + noneValue + `]

Zabbix options:
    -z --zabbix <host>               Hostname or IP address of Zabbix server [default: 127.0.0.1]
    -p --zabbix-port <number>        Port of Zabbix server [default: 10051]
    -d --discovery                   Run low-level discovery for determine queues, exchanges, etc.
    -a --aggregate                   Discovery aggregate items.
    -g --group <name>                Group name which will be use for aggregate item values.[default: None]
    -o --hostname <name>             Hostname which will be used in zabbix-sender protocol data. [default: ` + obtainHostname() + `

Misc options:
    -h --help                        Show this screen.
    -v --version                     Show version.
`

	var (
		rmqc *rabbithole.Client
		err  error
	)

	args, _ := docopt.Parse(usage, nil, true, version, false)

	aggGroup := args["--group"].(string)

	zabbixHost := args["--zabbix"].(string)
	zabbixPort, err := strconv.Atoi(args["--zabbix-port"].(string))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	rmqTimeout, err := strconv.Atoi(args["--rabbitmq-timeout"].(string))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var metrics []*zsend.Metric

	hostname := args["--hostname"].(string)

	rmqc, err = makeRabbitMQClient(
		parseDSN(args["--rabbitmq"].(string)),
		args["--rabbitmq-user"].(string),
		args["--rabbitmq-secret"].(string),
		args["--ca"].(string),
		rmqTimeout,
	)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	availableVhosts, err := rmqc.ListVhosts()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	allQueues := make(map[string][]rabbithole.QueueInfo)

	for _, vhost := range availableVhosts {

		existsQueues, err := rmqc.ListQueuesIn(vhost.Name)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		allQueues[vhost.Name] = existsQueues
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

		err = discovery(
			rmqc,
			allQueues,
			aggGroup,
			args["--aggregate"].(bool),
		)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	metrics = getQueuesMetrics(
		hostname,
		allQueues,
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
