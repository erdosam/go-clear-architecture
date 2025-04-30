package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Postgres -.
type Postgres struct {
	*sqlx.DB
}

func New(url string) (*Postgres, error) {
	conn, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("postgres.New: %w", err)
	}
	return &Postgres{conn}, nil
}

func (p *Postgres) Disconnect() {
	if p.DB == nil {
		fmt.Println("No DB connection")
	}
	err := p.Close()
	if err != nil {
		_ = fmt.Errorf("closing failed : %w", err)
	}
}
