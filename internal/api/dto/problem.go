package dto

import "time"

// ProblemType represents the type of error
type ProblemType string

const (
	ProblemTypeResourceNotFound   ProblemType = "recurso-nao-encontrado"
	ProblemTypeEntityInUse        ProblemType = "entidade-em-uso"
	ProblemTypeBusinessError      ProblemType = "erro-negocio"
	ProblemTypeInvalidMessage     ProblemType = "mensagem-incompreensivel"
	ProblemTypeInvalidParameter   ProblemType = "parametro-invalido"
	ProblemTypeSystemError        ProblemType = "erro-de-sistema"
	ProblemTypeInvalidData        ProblemType = "dados-invalidos"
	ProblemTypeAccessDenied       ProblemType = "acesso-negado"
	ProblemTypeInvalidCredentials ProblemType = "credenciais-invalidas"
)

var problemTypeTitles = map[ProblemType]string{
	ProblemTypeResourceNotFound:   "Recurso nao encontrado",
	ProblemTypeEntityInUse:        "Entidade em uso",
	ProblemTypeBusinessError:      "Violacao de regra de negocio",
	ProblemTypeInvalidMessage:     "Mensagem incompreensivel",
	ProblemTypeInvalidParameter:   "Parametro invalido",
	ProblemTypeSystemError:        "Erro de sistema",
	ProblemTypeInvalidData:        "Dados invalidos",
	ProblemTypeAccessDenied:       "Acesso negado",
	ProblemTypeInvalidCredentials: "Credenciais invalidas",
}

func (p ProblemType) Title() string {
	return problemTypeTitles[p]
}

func (p ProblemType) URI() string {
	return "https://algafood.com.br/" + string(p)
}

// Problem represents the standard error response
type Problem struct {
	Status      int           `json:"status"`
	Type        string        `json:"type"`
	Title       string        `json:"title"`
	Detail      string        `json:"detail"`
	UserMessage string        `json:"userMessage"`
	Timestamp   time.Time     `json:"timestamp"`
	Objects     []ObjectError `json:"objects,omitempty"`
}

// ObjectError represents a field validation error
type ObjectError struct {
	Name        string `json:"name"`
	UserMessage string `json:"userMessage"`
}

// NewProblem creates a new Problem instance
func NewProblem(status int, problemType ProblemType, detail, userMessage string) *Problem {
	return &Problem{
		Status:      status,
		Type:        problemType.URI(),
		Title:       problemType.Title(),
		Detail:      detail,
		UserMessage: userMessage,
		Timestamp:   time.Now(),
	}
}

// NewProblemWithObjects creates a Problem with field errors
func NewProblemWithObjects(status int, problemType ProblemType, detail, userMessage string, objects []ObjectError) *Problem {
	p := NewProblem(status, problemType, detail, userMessage)
	p.Objects = objects
	return p
}
