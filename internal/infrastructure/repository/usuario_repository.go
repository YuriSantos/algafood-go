package repository

import (
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"gorm.io/gorm"
)

type usuarioRepositoryImpl struct {
	db *gorm.DB
}

// NewUsuarioRepository creates a new UsuarioRepository
func NewUsuarioRepository(db *gorm.DB) *usuarioRepositoryImpl {
	return &usuarioRepositoryImpl{db: db}
}

func (r *usuarioRepositoryImpl) FindAll() ([]model.Usuario, error) {
	var usuarios []model.Usuario
	if err := r.db.Find(&usuarios).Error; err != nil {
		return nil, err
	}
	return usuarios, nil
}

func (r *usuarioRepositoryImpl) FindByID(id uint64) (*model.Usuario, error) {
	var usuario model.Usuario
	if err := r.db.Preload("Grupos").First(&usuario, id).Error; err != nil {
		return nil, err
	}
	return &usuario, nil
}

func (r *usuarioRepositoryImpl) FindByEmail(email string) (*model.Usuario, error) {
	var usuario model.Usuario
	if err := r.db.Where("email = ?", email).First(&usuario).Error; err != nil {
		return nil, err
	}
	return &usuario, nil
}

func (r *usuarioRepositoryImpl) Save(usuario *model.Usuario) error {
	return r.db.Save(usuario).Error
}

func (r *usuarioRepositoryImpl) AddGrupo(usuarioID, grupoID uint64) error {
	return r.db.Exec("INSERT INTO usuario_grupo (usuario_id, grupo_id) VALUES (?, ?)", usuarioID, grupoID).Error
}

func (r *usuarioRepositoryImpl) RemoveGrupo(usuarioID, grupoID uint64) error {
	return r.db.Exec("DELETE FROM usuario_grupo WHERE usuario_id = ? AND grupo_id = ?", usuarioID, grupoID).Error
}
