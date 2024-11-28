package model

import (
	"database/sql"
)

type SearchWhere struct {
	Name      string
	IDCard    string
	Location  string
	FileID    int64
	FileTitle string
	FileType  string

	FirstMoveInStart sql.NullTime
	FirstMoveInEnd   sql.NullTime

	LastMoveInStart sql.NullTime
	LastMoveInEnd   sql.NullTime

	LastMoveOutStart sql.NullTime
	LastMoveOutEnd   sql.NullTime

	MoveStatus string

	MoveOutPeopleName string
	MoveOutPeopleUnit string
}

type SearchRecord struct {
	MoveOutStart sql.NullTime
	MoveOutEnd   sql.NullTime

	MoveStatus        string
	MoveOutPeopleName string
	MoveOutPeopleUnit string
}
