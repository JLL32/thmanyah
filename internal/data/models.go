package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Videos VideoModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Videos: VideoModel{DB: db},
	}
}
