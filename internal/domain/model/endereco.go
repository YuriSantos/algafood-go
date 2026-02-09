package model

// Endereco represents an address (embedded in Restaurante and Pedido)
type Endereco struct {
	CEP         string `gorm:"column:endereco_cep;size:9" json:"cep"`
	Logradouro  string `gorm:"column:endereco_logradouro;size:100" json:"logradouro"`
	Numero      string `gorm:"column:endereco_numero;size:20" json:"numero"`
	Complemento string `gorm:"column:endereco_complemento;size:60" json:"complemento"`
	Bairro      string `gorm:"column:endereco_bairro;size:60" json:"bairro"`
	CidadeID    uint64 `gorm:"column:endereco_cidade_id" json:"cidadeId"`
	Cidade      Cidade `gorm:"foreignKey:CidadeID" json:"cidade,omitempty"`
}
