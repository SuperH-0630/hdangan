package model

const DefaultPageItemCount = 20
const MaxLimit = 10000

type PageSizeResult struct {
	AllCount int64 `gorm:"column:ps"`
}

type FileIDMax struct {
	FileID int64 `gorm:"column:fid"`
}
