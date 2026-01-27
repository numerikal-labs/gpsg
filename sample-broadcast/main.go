package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// starting rabbitMQ using docker, should be flexible
func startRabbitMQ() {
	log.Println("Starting RabbitMQ using Docker...")

	cmd := exec.Command("docker-compose", "up", "-d")
	time.Sleep(5 * time.Second)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to start RabbitMQ: %v", err)
	}
	log.Println("RabbitMQ started successfully")
}

// very minimalist, assume data will be of a different form, so just had this as a placeholder
func readCSV(filePath string) string {
	//convert csv to string to be broadcasted
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read CSV file: %v", err)
	}
	return string(file)
}

func publishToRabbitMQ(data string) {
	// establish connection with rabbitHQ, check for errors
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	//create a channel, check for errors
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	//Declare Exchange with name "waveform_exchange" and type "fanout"
	//This function has many required parameters (after type), I did not research them much
	//But once we understand the endpoint better, they could be adjusted
	err = ch.ExchangeDeclare(
		"waveform_exchange", "fanout", true, false, false, false, nil)

	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	//publish the message
	err = ch.Publish(
		"waveform_exchange", // name
		"",                  // routing key (n/a for fanout)
		false,               // message can be routed to no queues without error
		false,               // message stays in queue until a conusmer is available
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data), //converts data into byte slice, neccesary for broadcasting
		},
	)

	//check for error

	if err != nil {
		log.Fatalf("Failed to publish data: %v", err)
	}

	log.Println("Data broadcasted successfully")
}

func main() {
	//start RabbitMQ using Docker
	startRabbitMQ()

	//read CSV file
	csvData := readCSV("sample_waveform_data.csv")

	//broadcast CSV data to RabbitMQ
	publishToRabbitMQ(csvData)

	//haven't created the queues that would recieve the broadcast, but that would be where users could access it
}
