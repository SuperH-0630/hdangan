package model

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type FileSetType int64

const (
	QianRu FileSetType = iota
	ChuSheng
	QianChu
	SiWang
	BianGeng
	SuoNeiYiJu
	SuoJianYiJu
	NongZiZhuanFei
	YiZhanShiQianYiZheng
)

var FileSetTypeList = []FileSetType{QianRu, ChuSheng, QianChu, SiWang, BianGeng, SuoNeiYiJu, SuoJianYiJu, NongZiZhuanFei, YiZhanShiQianYiZheng}
var FileSetTypeName = map[FileSetType]string{
	QianRu:               "迁入",
	ChuSheng:             "出生",
	QianChu:              "迁出",
	SiWang:               "死亡",
	BianGeng:             "变更",
	SuoNeiYiJu:           "所内移居",
	SuoJianYiJu:          "所间移居",
	NongZiZhuanFei:       "农（自）转非",
	YiZhanShiQianYiZheng: "一站式迁移证",
}
var FileSetTypeID = map[string]FileSetType{
	"迁入":     QianRu,
	"出生":     ChuSheng,
	"迁出":     QianChu,
	"死亡":     SiWang,
	"变更":     BianGeng,
	"所内移居":   SuoNeiYiJu,
	"所间移居":   SuoJianYiJu,
	"农（自）转非": NongZiZhuanFei,
	"一站式迁移证": YiZhanShiQianYiZheng,
}
var FileSetTypeMaker = map[FileSetType]func() File{
	QianRu: func() File {
		return &FileQianRu{}
	},
	ChuSheng: func() File {
		return &FileChuSheng{}
	},
	QianChu: func() File {
		return &FileQianChu{}
	},
	SiWang: func() File {
		return &FileSiWang{}
	},
	BianGeng: func() File {
		return &FileBianGeng{}
	},
	SuoNeiYiJu: func() File {
		return &FileSuoNeiYiJu{}
	},
	SuoJianYiJu: func() File {
		return &FileSuoJianYiJu{}
	},
	NongZiZhuanFei: func() File {
		return &FileNongZiZhuanFei{}
	},
	YiZhanShiQianYiZheng: func() File {
		return &FileYiZhanShiQianYiZheng{}
	},
}

type FileSet struct {
	gorm.Model
	FileSetID   int64       `gorm:"column:filesetid;type:uint;not null;index:fileset"`
	FileSetType FileSetType `gorm:"column:filesettype;type:uint;not null;index:fileset"`
	PageCount   int64       `gorm:"column:pagecount;type:uint;not null"`

	File []FileAbs `gorm:"foreignKey:FileSetSQLID;references:ID"`
}

type File interface {
	GetFile() *FileAbs
}

// FileAbs 档案
type FileAbs struct {
	gorm.Model

	FileSetSQLID int64       `gorm:"column:filesetsqlid;type:uint;not null"`
	FileSetID    int64       `gorm:"column:filesetid;type:uint;not null"`
	FileSetType  FileSetType `gorm:"column:filesettype;type:uint;not null"`
	FileSet      FileSet     `gorm:"foreignKey:FileSetSQLID;references:ID"`

	FileUnionID int64 `gorm:"column:fileunionid;type:uint;not null"` // 联合ID
	FileGroupID int64 `gorm:"column:filegroupid;type:uint;not null"` // 组内id
	FileID      int64 `gorm:"column:fileid;type:uint;not null"`

	// 档案人基本信息
	Name     string         `gorm:"column:name;type:VARCHAR(50);not null"`
	OldName  sql.NullString `gorm:"column:oldname;type:VARCHAR(50);not null"`
	IDCard   sql.NullString `gorm:"column:idcard;type:VARCHAR(20)"`
	IsMan    bool           `gorm:"column:isman;type:BOOL;not null"`
	Birthday time.Time      `gorm:"column:birthday;type:DATE;not null"`
	Comment  sql.NullString `gorm:"column:comment;type:TEXT"`

	// 档案联合
	SameAsAbove bool  `gorm:"column:sameabove;type:BOOL;not null"`
	PeopleCount int64 `gorm:"column:peoplecount;type:uint;not null"`

	// 档案信息
	Time time.Time `gorm:"column:time;type:DATETIME;not null"`

	PageStart int64 `gorm:"column:pagestart;type:uint;not null"`
	PageEnd   int64 `gorm:"column:pageend;type:uint;not null"`
	PageCount int64 `gorm:"column:pagecount;type:uint;not null"`

	BeiKao   sql.NullString `gorm:"column:beikao;type:VARCHAR(50);"`
	Material sql.NullString `gorm:"column:material;type:VARCHAR(500);"`

	// 出入库
	LastMoveRecordID sql.NullInt64 `gorm:"column:lastmoverecordid;type:uint;"`
}

func (f *FileAbs) GetFile() *FileAbs {
	return f
}

type FileQianRu struct {
	// 迁入
	FileAbs
	Type        string `gorm:"column:type;type:VARCHAR(20);not null"`
	OldLocation string `gorm:"column:oldloc;type:VARCHAR(50);not null"`
	NewLocation string `gorm:"column:newloc;type:VARCHAR(50);not null"`
}

type FileChuSheng struct {
	// 出生
	FileAbs
	Type        string `gorm:"column:type;type:VARCHAR(20);not null"`
	NewLocation string `gorm:"column:newloc;type:VARCHAR(50);not null"`
}

type FileQianChu struct {
	// 迁出
	FileAbs
	Type        string `gorm:"column:type;type:VARCHAR(20);not null"`
	NewLocation string `gorm:"column:newloc;type:VARCHAR(50);not null"`
}

type FileSiWang struct {
	// 死亡
	FileAbs
	Type     string `gorm:"column:type;type:VARCHAR(20);not null"`
	Location string `gorm:"column:loc;type:VARCHAR(50);not null"`
}

type FileBianGeng struct {
	// 变更
	FileAbs
	Type     string `gorm:"column:type;type:VARCHAR(20);not null"`
	Location string `gorm:"column:loc;type:VARCHAR(50);not null"`
}

type FileSuoNeiYiJu struct {
	// 所内移居
	FileAbs
	Type     string `gorm:"column:type;type:VARCHAR(20);not null"`
	Location string `gorm:"column:loc;type:VARCHAR(50);not null"`
}

type FileSuoJianYiJu struct {
	// 所间移居
	FileAbs
	Type     string `gorm:"column:type;type:VARCHAR(20);not null"`
	Location string `gorm:"column:loc;type:VARCHAR(50);not null"`
}

type FileNongZiZhuanFei struct {
	// 农自转非
	FileAbs
	Type     string `gorm:"column:type;type:VARCHAR(20);not null"`
	Location string `gorm:"column:loc;type:VARCHAR(50);not null"`
}

type FileYiZhanShiQianYiZheng struct {
	// 一站式迁移证
	FileAbs
	Type     string `gorm:"column:type;type:VARCHAR(20);not null"`
	Location string `gorm:"column:loc;type:VARCHAR(50);not null"`
}

type FileMoveRecord struct {
	gorm.Model

	FileSetSQLID int64       `gorm:"column:filesetsqlid;type:uint;not null"`
	FileSetID    int64       `gorm:"column:filesetid;type:uint;not null"`
	FileSetType  FileSetType `gorm:"column:filesettype;type:uint;not null"`

	FileSQLID   int64 `gorm:"column:filesqlid;type:uint;not null"`
	FileUnionID int64 `gorm:"column:fileunionid;type:uint;not null"`

	UpRecord   sql.NullInt64 `gorm:"column:uprecord;type:uint"`
	NextRecord sql.NullInt64 `gorm:"column:nextrecord;type:uint"`
	File       FileAbs       `gorm:"foreignKey:FileSQLID;references:ID"`

	MoveStatus        string         `gorm:"column:movestatus;type:VARCHAR(20);not null"`
	MoveTime          time.Time      `gorm:"column:movetime;type:DATETIME;not null"`
	MoveOutPeopleName sql.NullString `gorm:"column:moveoutpeoplename;type:VARCHAR(10)"`
	MoveOutPeopleUnit sql.NullString `gorm:"column:moveoutpeopleunit;type:VARCHAR(50)"`
	MoveInPeopleName  sql.NullString `gorm:"column:moveinpeoplename;type:VARCHAR(10)"`
	MoveInPeopleUnit  sql.NullString `gorm:"column:moveinpeopleunit;type:VARCHAR(50)"`
	MoveComment       sql.NullString `gorm:"column:movecomment;type:TEXT"`
}
