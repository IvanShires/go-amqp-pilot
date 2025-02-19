package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	amqp_user := os.Getenv("AMQP_USER")
	if amqp_user == "" {
		slog.Error("amqp_user env variable must be set: %s\n", "amqp_user", amqp_user)
		os.Exit(1)
	}

	amqp_secret := os.Getenv("AMQP_SECRET")
	if amqp_secret == "" {
		slog.Error("amqp_secret env variable must be set: %s\n", "amqp_secret", amqp_secret)
		os.Exit(1)
	}

	amqp_host := os.Getenv("AMQP_HOST")
	if amqp_host == "" {
		slog.Error("amqp_host env variable must be set: %s\n", "amqp_host", amqp_host)
		os.Exit(1)
	}

	amqp_port := os.Getenv("AMQP_PORT")
	if amqp_port == "" {
		slog.Error("amqp_port env variable must be set: %s\n", "amqp_port", amqp_port)
		os.Exit(1)
	}

	amqp_queue := os.Getenv("AMQP_QUEUE")
	if amqp_queue == "" {
		slog.Error("amqp_queue env variable must be set: %s\n", "amqp_queue", amqp_queue)
		os.Exit(1)
	}

	data := []map[string]interface{}{
		{
			"interface": "et-1/1/1",
			"allowed_vlans": map[string]interface{}{
				"vlan": []int{10, 20, 30, 40},
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		slog.Error("could not marshal json: %s\n", "error", err)
		os.Exit(1)
	}

	conn, err := amqp.Dial("amqp://" + string(amqp_user) + ":" + string(amqp_secret) + "@" + string(amqp_host) + ":" + string(amqp_port) + "/")
	if err != nil {
		slog.Error("failed to connect to AMQP instance", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		slog.Error("failed to open channel", "error", err)
		os.Exit(1)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		amqp_queue, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		slog.Error("failed to declare queue", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        jsonData,
		})
	if err != nil {
		slog.Error("failed to publish msg", "error", err)
		os.Exit(1)
	}
	log.Printf(" [x] Sent %s\n", jsonData)
}
