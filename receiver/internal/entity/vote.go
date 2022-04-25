package entity

import "database/sql"

type Vote struct {
	Id           string `json:"id"`
	ParedaoId    uint32 `json:"paredao_id"`
	EmparedadoId uint32 `json:"emparedado_id"`
}

func (v *Vote) Register(db *sql.DB) error {
	_, err := db.Exec(
		"INSERT INTO votes (paredao_id, emparedado_id) VALUES ($1, $2)",
		&v.ParedaoId, &v.EmparedadoId,
	)

	return err
}
