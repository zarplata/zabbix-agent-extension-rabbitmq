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

For both installation you also should import `template_app_rabbitmq_service.xml` template into Zabbix server
