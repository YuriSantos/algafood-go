package model

// Permissao represents a system permission
type Permissao struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Nome      string `gorm:"size:100;not null" json:"nome"`
	Descricao string `gorm:"size:255" json:"descricao"`
}

func (Permissao) TableName() string {
	return "permissao"
}
