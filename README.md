# zabbix-agent-extension-rabbitmq
Zabbix agent extension for monitoring RabbitMQ server

This extension for monitoring RabbitMQ in standalone or cluster mode.

## Features
  - The `zabbix_sender` is no longer needed
  - Auto discovery all queues (LLD)
  - Triggers on discovered queues
  - Configurable triggers of queue max capacity by macro
 
#### This extension provides following information about each RabbitMQ queue:
  
  - consumers
  - memory usage
  - messages
  - ready messages
  - unacknowledged messages
  - deliver rate
  - deliver_get rate
  - deliver_no_ack rate
  - get rate
  - get_no_ack rate
  - publish rate
  - redeliver rate
  
#### Also extension by default have triggers which detect following situations:
  - `infinity or very long redeliver` - it happens when some message/s don't ACKed by a consumer and 
  they always move back to queue. It's not always a bad situation but for us, it means a problem with messages or consumer.
  - `zero deliver rate when queue not empty` - I think this is also a bad situation when you have some messages in the queue 
  but have not anyone consumer for handle that messages

## Installation

#### ArchLinux package

```sh
git clone https://github.com/zarplata/zabbix-agent-extension-rabbitmq.git
cd zabbix-agent-extension-rabbitmq
./build-archlinux.sh
sudo pacman -U *.tar.xz --noconfirm
systemctl restart zabbix-agent
```

#### From source

```sh
git clone https://github.com/zarplata/zabbix-agent-extension-rabbitmq.git
cd zabbix-agent-extension-rabbitmq
make

sudo cp .out/zabbix-agent-extension-rabbitmq /usr/bin/
sudo cp zabbix-agent-extension-rabbitmq.conf /etc/zabbix/zabbix_agentd.conf.d/
systemctl restart zabbix-agent
```

**Be note!**
  - You should add a global macros (Administration -> General -> Macros) - `{$ZABBIX_SERVER_IP}` with your Zabbix server IP. 
  - For both installation you also should import `template_app_rabbitmq_service.xml` template into Zabbix server.
  - Zabbix agent extensible directory path depends on Linux distribution and can be mismatch with the directory in this manual
  
  
## Configuration


#### Cluster specificity

If you have cluster you should known that each metric of queue will be the same on each node, so you will receive so many alerts, how many nodes you have in the cluster because extension should be installed on all nodes. To avoid this behave and have only one alert and have only one aggregate metric this extension use feature of Zabbix server - `aggregate checks` https://www.zabbix.com/documentation/3.4/manual/config/items/itemtypes/aggregate

So, you should create special host group in your Zabbix server, for example - `rabbitmq_aggregate` and add each node to this group, then you can add macro to one of cluster node (don't add macro on each node!!!) - `{$GROUPNAME} = rabbitmq_aggregate`, after that you can see new items and triggers only on node with defined macro.
