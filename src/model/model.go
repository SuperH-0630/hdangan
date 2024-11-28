package model

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

// File 档案
type File struct {
	gorm.Model

	// 档案人基本信息
	Name     string `gorm:"column:name;type:VARCHAR(50);not null;index:people"`
	IDCard   string `gorm:"column:idcard;type:VARCHAR(20);not null;index:people"`
	Location string `gorm:"column:loc;type:VARCHAR(150);not null;index:people"`

	// 档案基本inxi
	FileID      int64          `gorm:"column:fileid;type:uint;not null;index:file"`
	FileTitle   string         `gorm:"column:filetitle;type:VARCHAR(50);not null;index:file"`
	FileType    string         `gorm:"column:filetype;type:VARCHAR(20);not null;index:file"`
	FileComment sql.NullString `gorm:"column:filecomment;type:TEXT"`

	// 迁入迁出
	FirstMoveIn time.Time `gorm:"column:firstmovein;type:DATETIME;not null"`
	LastMoveIn  time.Time `gorm:"column:lastmovein;type:DATETIME;not null"`

	// 最新一次迁出记录
	LastMoveRecordID  sql.NullInt64    `gorm:"column:lastmoverecordid;type:uint"`
	MoveStatus        string           `gorm:"column:movestatus;type:VARCHAR(20);not null"`
	MoveOutPeopleName sql.NullString   `gorm:"column:moveoutpeoplename;type:VARCHAR(10)"`
	MoveOutPeopleUnit sql.NullString   `gorm:"column:moveoutpeopleunit;type:VARCHAR(50)"`
	MoveComment       sql.NullString   `gorm:"column:movecomment;type:TEXT"`
	MoveRecord        []FileMoveRecord `gorm:"foreignKey:FileSQLID;references:ID"`
}

type FileMoveRecord struct {
	gorm.Model

	FileSQLID  int64         `gorm:"column:filesqlid;type:uint;not null"`
	UpRecord   sql.NullInt64 `gorm:"column:uprecord;type:uint"`
	NextRecord sql.NullInt64 `gorm:"column:nextrecord;type:uint"`
	File       File          `gorm:"foreignKey:FileSQLID;references:ID"`

	MoveStatus        string         `gorm:"column:movestatus;type:VARCHAR(20);not null"`
	MoveTime          time.Time      `gorm:"column:movetime;type:DATETIME;not null"`
	MoveOutPeopleName sql.NullString `gorm:"column:moveoutpeoplename;type:VARCHAR(10)"`
	MoveOutPeopleUnit sql.NullString `gorm:"column:moveoutpeopleunit;type:VARCHAR(50)"`
	MoveComment       sql.NullString `gorm:"column:movecomment;type:TEXT"`
}
