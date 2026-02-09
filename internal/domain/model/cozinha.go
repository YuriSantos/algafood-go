package model

// Cozinha represents a type of cuisine
type Cozinha struct {
	ID   uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Nome string `gorm:"size:60;not null" json:"nome"`
}

func (Cozinha) TableName() string {
	return "cozinha"
}
