package handler

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/streadway/amqp"
)

type VoteHandler struct {
	Conn *amqp.Connection
}

type votePayload struct {
	ParedaoId    uint32 `json:"paredao_id"`
	EmparedadoId uint32 `json:"emparedado_id"`
}

func (vp votePayload) Validate() error {
	return validation.ValidateStruct(&vp,
		validation.Field(&vp.ParedaoId, validation.Required),
		validation.Field(&vp.EmparedadoId, validation.Required),
	)
}

func (vh *VoteHandler) HandleVote(w http.ResponseWriter, r *http.Request) {
	var vp votePayload

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&vp); err != nil {
		response, _ := json.Marshal(map[string]string{"error": err.Error()})

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}

	defer r.Body.Close()

	if err := vp.Validate(); err != nil {
		response, _ := json.Marshal(map[string]string{"error": err.Error()})

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}

	ch, err := vh.Conn.Channel()

	if err != nil {
		response, _ := json.Marshal(map[string]string{"error": "Internal server error"})

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response)
		return
	}

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
		response, _ := json.Marshal(map[string]string{"error": "Internal server error"})

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response)
		return
	}

	body, _ := json.Marshal(vp)

	err = ch.Publish(
		"votes",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		response, _ := json.Marshal(map[string]string{"error": "Internal server error"})

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(nil)
}
