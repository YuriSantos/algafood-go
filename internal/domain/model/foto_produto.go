package model

// FotoProduto represents a product photo
type FotoProduto struct {
	ID          uint64  `gorm:"primaryKey" json:"id"`
	ProdutoID   uint64  `gorm:"uniqueIndex" json:"produtoId"`
	Produto     Produto `gorm:"foreignKey:ProdutoID" json:"-"`
	NomeArquivo string  `gorm:"size:150;not null" json:"nomeArquivo"`
	Descricao   string  `gorm:"size:150" json:"descricao"`
	ContentType string  `gorm:"size:80;not null" json:"contentType"`
	Tamanho     int64   `gorm:"not null" json:"tamanho"`
}

func (FotoProduto) TableName() string {
	return "foto_produto"
}

// GetRestauranteID returns the restaurant ID from the associated product
func (f *FotoProduto) GetRestauranteID() uint64 {
	if f.Produto.ID != 0 {
		return f.Produto.RestauranteID
	}
	return 0
}
