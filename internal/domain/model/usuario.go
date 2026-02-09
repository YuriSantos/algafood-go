package model

import "time"

// Usuario represents a system user
type Usuario struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Nome         string    `gorm:"size:80;not null" json:"nome"`
	Email        string    `gorm:"size:255;not null;uniqueIndex" json:"email"`
	Senha        string    `gorm:"size:255;not null" json:"-"`
	DataCadastro time.Time `gorm:"autoCreateTime" json:"dataCadastro"`
	Grupos       []Grupo   `gorm:"many2many:usuario_grupo;" json:"grupos,omitempty"`
}

func (Usuario) TableName() string {
	return "usuario"
}

// AdicionarGrupo adds a group to the user
func (u *Usuario) AdicionarGrupo(grupo Grupo) {
	u.Grupos = append(u.Grupos, grupo)
}

// RemoverGrupo removes a group from the user
func (u *Usuario) RemoverGrupo(grupo Grupo) {
	for i, g := range u.Grupos {
		if g.ID == grupo.ID {
			u.Grupos = append(u.Grupos[:i], u.Grupos[i+1:]...)
			return
		}
	}
}

// SenhaCoincideCom checks if the provided password matches the user's password
// Note: In production, this should use bcrypt.CompareHashAndPassword
func (u *Usuario) SenhaCoincideCom(senha string) bool {
	return u.Senha == senha
}

// SenhaNaoCoincideCom checks if the provided password doesn't match
func (u *Usuario) SenhaNaoCoincideCom(senha string) bool {
	return !u.SenhaCoincideCom(senha)
}
