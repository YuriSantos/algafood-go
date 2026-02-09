package model

// Grupo represents a user group with permissions
type Grupo struct {
	ID         uint64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nome       string      `gorm:"size:60;not null" json:"nome"`
	Permissoes []Permissao `gorm:"many2many:grupo_permissao;" json:"permissoes,omitempty"`
}

func (Grupo) TableName() string {
	return "grupo"
}

// AdicionarPermissao adds a permission to the group
func (g *Grupo) AdicionarPermissao(permissao Permissao) {
	g.Permissoes = append(g.Permissoes, permissao)
}

// RemoverPermissao removes a permission from the group
func (g *Grupo) RemoverPermissao(permissao Permissao) {
	for i, p := range g.Permissoes {
		if p.ID == permissao.ID {
			g.Permissoes = append(g.Permissoes[:i], g.Permissoes[i+1:]...)
			return
		}
	}
}
