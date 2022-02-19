package beer

import (
	"database/sql"
	"errors"
)


type PostgresRepository struct {
	DB *sql.DB
}


func (r *PostgresRepository) GetBeer(id string) (*RepoBeer, error) {
	return nil, errors.New("not found")
}
