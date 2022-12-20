# Cena de shushi

Práctica 3 de la asignatura de programación concurrente.

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

# Ejecución

1. Iniciar el servidor de RabbitMQ.
2. Ejecutar N clientes, donde N es el número de clientes, N gangsters y 1 cocinero.

## Ejemplo

```bash
# Ejecutar N clientes -> 1 cada cliente en una terminal
go run cliente.go

# Ejecutar N gangsters -> 1 cada gangster en una terminal
go run gangster.go

# Ejecutar 1 cocinero
go run cocinero.go
```
