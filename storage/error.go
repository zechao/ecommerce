package storage

import (
	"errors"
)

var (
	// ErrRecordNotFound is an error that indicates that a requested record could not be found in the database or data store.
	ErrRecordNotFound = errors.New("record not found")

	// ErrDuplicateKey is returned when attempting to create a record that already exists in storage.
	// This typically occurs when trying to insert a record with a unique key constraint
	// that conflicts with an existing entry.
	ErrDuplicateKey   = errors.New("record already exist")
)
