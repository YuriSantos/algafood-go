package repository

import (
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"gorm.io/gorm"
)

type grupoRepositoryImpl struct {
	db *gorm.DB
}

// NewGrupoRepository creates a new GrupoRepository
func NewGrupoRepository(db *gorm.DB) *grupoRepositoryImpl {
	return &grupoRepositoryImpl{db: db}
}

func (r *grupoRepositoryImpl) FindAll() ([]model.Grupo, error) {
	var grupos []model.Grupo
	if err := r.db.Find(&grupos).Error; err != nil {
		return nil, err
	}
	return grupos, nil
}

func (r *grupoRepositoryImpl) FindByID(id uint64) (*model.Grupo, error) {
	var grupo model.Grupo
	if err := r.db.Preload("Permissoes").First(&grupo, id).Error; err != nil {
		return nil, err
	}
	return &grupo, nil
}

func (r *grupoRepositoryImpl) Save(grupo *model.Grupo) error {
	return r.db.Save(grupo).Error
}

func (r *grupoRepositoryImpl) Delete(id uint64) error {
	return r.db.Delete(&model.Grupo{}, id).Error
}

func (r *grupoRepositoryImpl) AddPermissao(grupoID, permissaoID uint64) error {
	return r.db.Exec("INSERT INTO grupo_permissao (grupo_id, permissao_id) VALUES (?, ?)", grupoID, permissaoID).Error
}

func (r *grupoRepositoryImpl) RemovePermissao(grupoID, permissaoID uint64) error {
	return r.db.Exec("DELETE FROM grupo_permissao WHERE grupo_id = ? AND permissao_id = ?", grupoID, permissaoID).Error
}
