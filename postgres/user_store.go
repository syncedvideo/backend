package postgres

import "database/sql"

// UserStore manages the database connection
type UserStore struct {
	*sql.DB
}
