# Cena de shushi

# Prerrequisitos

- [RabbitMQ](https://www.rabbitmq.com/download.html)
- [Go client](https://github.com/rabbitmq/amqp091-go)

## RabbitMQ en Docker

Inicializar imagen

```bash
docker run --rm -it -p 15672:15672 -p 5672:5672 rabbitmq:3-management
```

Activación del pluggin manager

```bash
rabbitmq-plugins enable rabbitmq_management
```
