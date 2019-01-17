# zabbix-agent-extension-rabbitmq
Zabbix agent extension for monitoring RabbitMQ server

This extension for monitoring RabbitMQ in standalone or cluster mode.

## Features
  - The `zabbix_sender` is no longer needed
  - Auto discovery all queues (LLD)
  - Triggers on discovered queues
 
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
