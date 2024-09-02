package model

type Base struct {
	TotalRows uint64 `gorm:"type:bigint UNSIGNED" json:"-"`
}
