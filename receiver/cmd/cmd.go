package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"

	"github.com/xsadia/bbb_queue/receiver/internal/entity"
	"github.com/xsadia/bbb_queue/receiver/internal/storage"
)

var tableCreationQuery = `CREATE TABLE IF NOT EXISTS votes
(
	id uuid DEFAULT uuid_generate_v4() PRIMARY KEY, 
	paredao_id INTEGER NOT NULL,
	emparedado_id INTEGER NOT NULL
)
`

func main() {
	godotenv.Load(".env")

	db, err := storage.NewDB(
		os.Getenv("APP_DB_HOST"),
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	if err != nil {
		log.Fatalf("Error: %q", err)
	}

	defer db.Close()

	if _, err = db.Exec(tableCreationQuery); err != nil {
		log.Fatalf("Error: %q", err)
	}

	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URI"))

	if err != nil {
		log.Fatalf("Error: %q", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("Error: %q", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"votes",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Error: %q", err)
	}

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Error: %q", err)
	}

	err = ch.QueueBind(
		q.Name,
		"",
		"votes",
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Error: %q", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Error: %q", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var v entity.Vote
			json.Unmarshal(d.Body, &v)
			v.Register(db)
			log.Printf("[x] %v\n", string(d.Body))
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C\n")
	<-forever
}
