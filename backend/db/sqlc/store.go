package db

import (
	"context"
	"database/sql"
	"fmt"
)

/*
Custom Transaction Flow
1. Begin tx
2. Transfer money
3. enter entry 1 in
4. enter entry 2 out
5. update balance
6. commit transaction
*/

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db, Queries: New(db)}
}

func (s *Store) execTx(c context.Context, fq func(q *Queries) error) error {
	// initial trxn
	tx, err := s.db.BeginTx(c, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fq(q)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return fmt.Errorf("encountered rollback error: %v", txErr)
		}
		return err
	}

	return tx.Commit()
}
