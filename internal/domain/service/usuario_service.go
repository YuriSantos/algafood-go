package service

import (
	"errors"

	"github.com/yurisasc/algafood-go/internal/domain/exception"
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"github.com/yurisasc/algafood-go/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UsuarioService struct {
	repo     repository.UsuarioRepository
	grupoSvc *GrupoService
}

func NewUsuarioService(repo repository.UsuarioRepository, grupoSvc *GrupoService) *UsuarioService {
	return &UsuarioService{
		repo:     repo,
		grupoSvc: grupoSvc,
	}
}

func (s *UsuarioService) Authenticate(email, password string) (*model.Usuario, error) {
	usuario, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewAuthenticationException("Usuario ou senha invalidos")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Senha), []byte(password)); err != nil {
		return nil, exception.NewAuthenticationException("Usuario ou senha invalidos")
	}

	return usuario, nil
}

func (s *UsuarioService) FindAll() ([]model.Usuario, error) {
	return s.repo.FindAll()
}

func (s *UsuarioService) FindByID(id uint64) (*model.Usuario, error) {
	usuario, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewUsuarioNaoEncontradoException(id)
		}
		return nil, err
	}
	return usuario, nil
}

func (s *UsuarioService) FindByEmail(email string) (*model.Usuario, error) {
	return s.repo.FindByEmail(email)
}

func (s *UsuarioService) Save(usuario *model.Usuario) error {
	// Check if email is already in use by another user
	existing, err := s.repo.FindByEmail(usuario.Email)
	if err == nil && existing.ID != usuario.ID {
		return exception.NewNegocioException("Ja existe um usuario cadastrado com o e-mail informado")
	}

	// Hash password if it's a new user or password was changed
	if usuario.ID == 0 && usuario.Senha != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usuario.Senha), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		usuario.Senha = string(hashedPassword)
	}

	return s.repo.Save(usuario)
}

func (s *UsuarioService) AlterarSenha(id uint64, senhaAtual, novaSenha string) error {
	usuario, err := s.FindByID(id)
	if err != nil {
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Senha), []byte(senhaAtual)); err != nil {
		return exception.NewNegocioException("Senha atual informada nao coincide com a senha do usuario")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(novaSenha), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	usuario.Senha = string(hashedPassword)

	return s.repo.Save(usuario)
}

func (s *UsuarioService) AssociarGrupo(usuarioID, grupoID uint64) error {
	if _, err := s.FindByID(usuarioID); err != nil {
		return err
	}
	if _, err := s.grupoSvc.FindByID(grupoID); err != nil {
		return err
	}
	return s.repo.AddGrupo(usuarioID, grupoID)
}

func (s *UsuarioService) DesassociarGrupo(usuarioID, grupoID uint64) error {
	if _, err := s.FindByID(usuarioID); err != nil {
		return err
	}
	if _, err := s.grupoSvc.FindByID(grupoID); err != nil {
		return err
	}
	return s.repo.RemoveGrupo(usuarioID, grupoID)
}
