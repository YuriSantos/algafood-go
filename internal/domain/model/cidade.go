package model

// Cidade represents a city
type Cidade struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Nome     string `gorm:"size:80;not null" json:"nome"`
	EstadoID uint64 `gorm:"not null" json:"estadoId"`
	Estado   Estado `gorm:"foreignKey:EstadoID" json:"estado,omitempty"`
}

func (Cidade) TableName() string {
	return "cidade"
}
