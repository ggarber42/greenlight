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
	Movies MovieModel
	Users  UserModel
	// Movies interface { Todo
	// 	Insert(movie *Movie) error
	// 	Get(id int64) (*Movie, error)
	// 	Update(movie *Movie) error
	// 	Delete(id int64) error
	// 	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
		Users:  UserModel{DB: db},
	}
}

// func NewMockModels() Models { Todo
// 	return Models{
// 		Movies: MockMovieModel{},
// 	}
// }
