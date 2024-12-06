package model

import (
	"database/sql"
)

type SearchWhere struct {
	Name          string
	OldName       string
	IDCard        string
	IsMan         string
	BirthdayStart sql.NullTime
	BirthdayEnd   sql.NullTime
	Comment       string

	FileSetID   int64
	FileUnionID int64
	FileID      int64
	FileGroupID int64
}

type SearchRecord struct {
	MoveOutStart sql.NullTime
	MoveOutEnd   sql.NullTime

	MoveStatus        string
	MoveOutPeopleName string
	MoveOutPeopleUnit string
	MoveInPeopleName  string
	MoveInPeopleUnit  string
}
