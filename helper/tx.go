package helper

import (
	"database/sql"
)

func CommitOrRollback(tx *sql.Tx, err *error) {
	if *err != nil {
		_ = tx.Rollback()
		return
	}

	if commitErr := tx.Commit(); commitErr != nil {
		*err = commitErr
	}
}
