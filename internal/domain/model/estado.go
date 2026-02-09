package model

// Estado represents a Brazilian state
type Estado struct {
	ID   uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Nome string `gorm:"size:80;not null" json:"nome"`
}

func (Estado) TableName() string {
	return "estado"
}
