package model

const DefaultPageItemCount = 20

type PageSizeResult struct {
	AllCount int64 `gorm:"column:ps"`
}

type FileIDMax struct {
	FileID int64 `gorm:"column:fid"`
}
