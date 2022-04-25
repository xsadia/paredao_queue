package routes

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"github.com/xsadia/bbb_queue/producer/cmd/server"
	"github.com/xsadia/bbb_queue/producer/internal/handler"
)

var s server.Server
var vh handler.VoteHandler

func Listen() {
	godotenv.Load(".env")
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URI"))

	if err != nil {
		log.Fatalf("Error: %q", err)
	}

	vh.Conn = conn

	s = server.NewServer(":8000")
	s.Router.HandleFunc("/vote", vh.HandleVote).Methods("POST")
	s.Run()
}
