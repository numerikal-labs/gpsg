// In handlers/rabbitmq.go
package handlers

import (
    "log"
    "github.com/streadway/amqp"
)

var (
    rabbitMQConn *amqp.Connection
    rabbitMQChannel *amqp.Channel
    queueName = "gpsg_queue"  // Define queue name as a constant
)

func InitRabbitMQ() error {
    var err error
    
    // Connect to RabbitMQ
    rabbitMQConn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        return err
    }

    // Create a channel
    rabbitMQChannel, err = rabbitMQConn.Channel()
    if err != nil {
        return err
    }

    // Declare the queue
    _, err = rabbitMQChannel.QueueDeclare(
        queueName, // name
        true,      // durable
        false,     // delete when unused
        false,     // exclusive
        false,     // no-wait
        nil,       // arguments
    )
    if err != nil {
        return err
    }

    log.Printf("Successfully connected to RabbitMQ and declared queue: %s", queueName)
    return nil
}

func PublishToRabbitMQ(message string) error {
    log.Printf("Publishing message to queue: %s", message)
    
    err := rabbitMQChannel.Publish(
        "",        // exchange
        queueName, // routing key (queue name)
        false,     // mandatory
        false,     // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(message),
        },
    )
    
    if err != nil {
        log.Printf("Error publishing message: %v", err)
        return err
    }
    
    log.Printf("Successfully published message to queue")
    return nil
}