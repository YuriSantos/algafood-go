package assembler

import (
	"github.com/shopspring/decimal"
	"github.com/yurisasc/algafood-go/internal/api/dto"
	"github.com/yurisasc/algafood-go/internal/domain/model"
)

// ToEstadoModel converts Estado entity to EstadoModel DTO
func ToEstadoModel(e *model.Estado) dto.EstadoModel {
	return dto.EstadoModel{
		ID:   e.ID,
		Nome: e.Nome,
	}
}

// ToEstadoModels converts slice of Estado entities to slice of EstadoModel DTOs
func ToEstadoModels(estados []model.Estado) []dto.EstadoModel {
	models := make([]dto.EstadoModel, len(estados))
	for i, e := range estados {
		models[i] = ToEstadoModel(&e)
	}
	return models
}

// ToEstadoEntity converts EstadoInput DTO to Estado entity
func ToEstadoEntity(input *dto.EstadoInput) *model.Estado {
	return &model.Estado{
		Nome: input.Nome,
	}
}

// ToCidadeModel converts Cidade entity to CidadeModel DTO
func ToCidadeModel(c *model.Cidade) dto.CidadeModel {
	return dto.CidadeModel{
		ID:     c.ID,
		Nome:   c.Nome,
		Estado: ToEstadoModel(&c.Estado),
	}
}

// ToCidadeModels converts slice of Cidade entities to slice of CidadeModel DTOs
func ToCidadeModels(cidades []model.Cidade) []dto.CidadeModel {
	models := make([]dto.CidadeModel, len(cidades))
	for i, c := range cidades {
		models[i] = ToCidadeModel(&c)
	}
	return models
}

// ToCidadeEntity converts CidadeInput DTO to Cidade entity
func ToCidadeEntity(input *dto.CidadeInput) *model.Cidade {
	return &model.Cidade{
		Nome:     input.Nome,
		EstadoID: input.Estado.ID,
	}
}

// ToCozinhaModel converts Cozinha entity to CozinhaModel DTO
func ToCozinhaModel(c *model.Cozinha) dto.CozinhaModel {
	return dto.CozinhaModel{
		ID:   c.ID,
		Nome: c.Nome,
	}
}

// ToCozinhaModels converts slice of Cozinha entities to slice of CozinhaModel DTOs
func ToCozinhaModels(cozinhas []model.Cozinha) []dto.CozinhaModel {
	models := make([]dto.CozinhaModel, len(cozinhas))
	for i, c := range cozinhas {
		models[i] = ToCozinhaModel(&c)
	}
	return models
}

// ToCozinhaEntity converts CozinhaInput DTO to Cozinha entity
func ToCozinhaEntity(input *dto.CozinhaInput) *model.Cozinha {
	return &model.Cozinha{
		Nome: input.Nome,
	}
}

// ToFormaPagamentoModel converts FormaPagamento entity to FormaPagamentoModel DTO
func ToFormaPagamentoModel(fp *model.FormaPagamento) dto.FormaPagamentoModel {
	return dto.FormaPagamentoModel{
		ID:              fp.ID,
		Descricao:       fp.Descricao,
		DataAtualizacao: fp.DataAtualizacao,
	}
}

// ToFormaPagamentoModels converts slice of FormaPagamento entities
func ToFormaPagamentoModels(formasPagamento []model.FormaPagamento) []dto.FormaPagamentoModel {
	models := make([]dto.FormaPagamentoModel, len(formasPagamento))
	for i, fp := range formasPagamento {
		models[i] = ToFormaPagamentoModel(&fp)
	}
	return models
}

// ToFormaPagamentoEntity converts FormaPagamentoInput DTO to FormaPagamento entity
func ToFormaPagamentoEntity(input *dto.FormaPagamentoInput) *model.FormaPagamento {
	return &model.FormaPagamento{
		Descricao: input.Descricao,
	}
}

// ToPermissaoModel converts Permissao entity to PermissaoModel DTO
func ToPermissaoModel(p *model.Permissao) dto.PermissaoModel {
	return dto.PermissaoModel{
		ID:        p.ID,
		Nome:      p.Nome,
		Descricao: p.Descricao,
	}
}

// ToPermissaoModels converts slice of Permissao entities
func ToPermissaoModels(permissoes []model.Permissao) []dto.PermissaoModel {
	models := make([]dto.PermissaoModel, len(permissoes))
	for i, p := range permissoes {
		models[i] = ToPermissaoModel(&p)
	}
	return models
}

// ToGrupoModel converts Grupo entity to GrupoModel DTO
func ToGrupoModel(g *model.Grupo) dto.GrupoModel {
	return dto.GrupoModel{
		ID:   g.ID,
		Nome: g.Nome,
	}
}

// ToGrupoModels converts slice of Grupo entities
func ToGrupoModels(grupos []model.Grupo) []dto.GrupoModel {
	models := make([]dto.GrupoModel, len(grupos))
	for i, g := range grupos {
		models[i] = ToGrupoModel(&g)
	}
	return models
}

// ToGrupoEntity converts GrupoInput DTO to Grupo entity
func ToGrupoEntity(input *dto.GrupoInput) *model.Grupo {
	return &model.Grupo{
		Nome: input.Nome,
	}
}

// ToUsuarioModel converts Usuario entity to UsuarioModel DTO
func ToUsuarioModel(u *model.Usuario) dto.UsuarioModel {
	return dto.UsuarioModel{
		ID:           u.ID,
		Nome:         u.Nome,
		Email:        u.Email,
		DataCadastro: u.DataCadastro,
	}
}

// ToUsuarioModels converts slice of Usuario entities
func ToUsuarioModels(usuarios []model.Usuario) []dto.UsuarioModel {
	models := make([]dto.UsuarioModel, len(usuarios))
	for i, u := range usuarios {
		models[i] = ToUsuarioModel(&u)
	}
	return models
}

// ToUsuarioEntity converts UsuarioComSenhaInput DTO to Usuario entity
func ToUsuarioEntity(input *dto.UsuarioComSenhaInput) *model.Usuario {
	return &model.Usuario{
		Nome:  input.Nome,
		Email: input.Email,
		Senha: input.Senha,
	}
}

// ToRestauranteModel converts Restaurante entity to RestauranteModel DTO
func ToRestauranteModel(r *model.Restaurante) dto.RestauranteModel {
	model := dto.RestauranteModel{
		ID:              r.ID,
		Nome:            r.Nome,
		TaxaFrete:       r.TaxaFrete,
		Cozinha:         ToCozinhaModel(&r.Cozinha),
		Ativo:           r.Ativo,
		Aberto:          r.Aberto,
		DataCadastro:    r.DataCadastro,
		DataAtualizacao: r.DataAtualizacao,
	}

	if r.Endereco.CEP != "" {
		model.Endereco = &dto.EnderecoModel{
			CEP:         r.Endereco.CEP,
			Logradouro:  r.Endereco.Logradouro,
			Numero:      r.Endereco.Numero,
			Complemento: r.Endereco.Complemento,
			Bairro:      r.Endereco.Bairro,
			Cidade: dto.CidadeResumoModel{
				ID:     r.Endereco.Cidade.ID,
				Nome:   r.Endereco.Cidade.Nome,
				Estado: r.Endereco.Cidade.Estado.Nome,
			},
		}
	}

	return model
}

// ToRestauranteResumoModels converts slice of Restaurante entities to summary DTOs
func ToRestauranteResumoModels(restaurantes []model.Restaurante) []dto.RestauranteResumoModel {
	models := make([]dto.RestauranteResumoModel, len(restaurantes))
	for i, r := range restaurantes {
		models[i] = dto.RestauranteResumoModel{
			ID:        r.ID,
			Nome:      r.Nome,
			TaxaFrete: r.TaxaFrete,
			Cozinha:   ToCozinhaModel(&r.Cozinha),
			Ativo:     r.Ativo,
			Aberto:    r.Aberto,
		}
	}
	return models
}

// ToRestauranteEntity converts RestauranteInput DTO to Restaurante entity
func ToRestauranteEntity(input *dto.RestauranteInput) *model.Restaurante {
	r := &model.Restaurante{
		Nome:      input.Nome,
		TaxaFrete: decimal.NewFromFloat(input.TaxaFrete),
		CozinhaID: input.Cozinha.ID,
	}

	if input.Endereco != nil {
		r.Endereco = model.Endereco{
			CEP:         input.Endereco.CEP,
			Logradouro:  input.Endereco.Logradouro,
			Numero:      input.Endereco.Numero,
			Complemento: input.Endereco.Complemento,
			Bairro:      input.Endereco.Bairro,
			CidadeID:    input.Endereco.Cidade.ID,
		}
	}

	return r
}

// ToProdutoModel converts Produto entity to ProdutoModel DTO
func ToProdutoModel(p *model.Produto) dto.ProdutoModel {
	return dto.ProdutoModel{
		ID:        p.ID,
		Nome:      p.Nome,
		Descricao: p.Descricao,
		Preco:     p.Preco,
		Ativo:     p.Ativo,
	}
}

// ToProdutoModels converts slice of Produto entities
func ToProdutoModels(produtos []model.Produto) []dto.ProdutoModel {
	models := make([]dto.ProdutoModel, len(produtos))
	for i, p := range produtos {
		models[i] = ToProdutoModel(&p)
	}
	return models
}

// ToProdutoEntity converts ProdutoInput DTO to Produto entity
func ToProdutoEntity(input *dto.ProdutoInput) *model.Produto {
	return &model.Produto{
		Nome:      input.Nome,
		Descricao: input.Descricao,
		Preco:     decimal.NewFromFloat(input.Preco),
		Ativo:     input.Ativo,
	}
}

// ToFotoProdutoModel converts FotoProduto entity to FotoProdutoModel DTO
func ToFotoProdutoModel(f *model.FotoProduto) dto.FotoProdutoModel {
	return dto.FotoProdutoModel{
		NomeArquivo: f.NomeArquivo,
		Descricao:   f.Descricao,
		ContentType: f.ContentType,
		Tamanho:     f.Tamanho,
	}
}

// ToPedidoModel converts Pedido entity to PedidoModel DTO
func ToPedidoModel(p *model.Pedido) dto.PedidoModel {
	itens := make([]dto.ItemPedidoModel, len(p.Itens))
	for i, item := range p.Itens {
		itens[i] = dto.ItemPedidoModel{
			ProdutoID:     item.ProdutoID,
			ProdutoNome:   item.Produto.Nome,
			Quantidade:    item.Quantidade,
			PrecoUnitario: item.PrecoUnitario,
			PrecoTotal:    item.PrecoTotal,
			Observacao:    item.Observacao,
		}
	}

	return dto.PedidoModel{
		Codigo:           p.Codigo,
		Subtotal:         p.Subtotal,
		TaxaFrete:        p.TaxaFrete,
		ValorTotal:       p.ValorTotal,
		Status:           string(p.Status),
		DataCriacao:      p.DataCriacao,
		DataConfirmacao:  p.DataConfirmacao,
		DataCancelamento: p.DataCancelamento,
		DataEntrega:      p.DataEntrega,
		Restaurante: dto.RestauranteApenasNomeModel{
			ID:   p.Restaurante.ID,
			Nome: p.Restaurante.Nome,
		},
		Cliente:        ToUsuarioModel(&p.Cliente),
		FormaPagamento: ToFormaPagamentoModel(&p.FormaPagamento),
		EnderecoEntrega: dto.EnderecoModel{
			CEP:         p.EnderecoEntrega.CEP,
			Logradouro:  p.EnderecoEntrega.Logradouro,
			Numero:      p.EnderecoEntrega.Numero,
			Complemento: p.EnderecoEntrega.Complemento,
			Bairro:      p.EnderecoEntrega.Bairro,
			Cidade: dto.CidadeResumoModel{
				ID:     p.EnderecoEntrega.Cidade.ID,
				Nome:   p.EnderecoEntrega.Cidade.Nome,
				Estado: p.EnderecoEntrega.Cidade.Estado.Nome,
			},
		},
		Itens: itens,
	}
}

// ToPedidoResumoModels converts slice of Pedido entities to summary DTOs
func ToPedidoResumoModels(pedidos []model.Pedido) []dto.PedidoResumoModel {
	models := make([]dto.PedidoResumoModel, len(pedidos))
	for i, p := range pedidos {
		models[i] = dto.PedidoResumoModel{
			Codigo:      p.Codigo,
			Subtotal:    p.Subtotal,
			TaxaFrete:   p.TaxaFrete,
			ValorTotal:  p.ValorTotal,
			Status:      string(p.Status),
			DataCriacao: p.DataCriacao,
			Restaurante: dto.RestauranteApenasNomeModel{
				ID:   p.Restaurante.ID,
				Nome: p.Restaurante.Nome,
			},
			Cliente: ToUsuarioModel(&p.Cliente),
		}
	}
	return models
}

// ToPedidoEntity converts PedidoInput DTO to Pedido entity
func ToPedidoEntity(input *dto.PedidoInput, clienteID uint64) *model.Pedido {
	itens := make([]model.ItemPedido, len(input.Itens))
	for i, item := range input.Itens {
		itens[i] = model.ItemPedido{
			ProdutoID:  item.ProdutoID,
			Quantidade: item.Quantidade,
			Observacao: item.Observacao,
		}
	}

	return &model.Pedido{
		RestauranteID:    input.Restaurante.ID,
		FormaPagamentoID: input.FormaPagamento.ID,
		ClienteID:        clienteID,
		EnderecoEntrega: model.EnderecoEntrega{
			CEP:         input.EnderecoEntrega.CEP,
			Logradouro:  input.EnderecoEntrega.Logradouro,
			Numero:      input.EnderecoEntrega.Numero,
			Complemento: input.EnderecoEntrega.Complemento,
			Bairro:      input.EnderecoEntrega.Bairro,
			CidadeID:    input.EnderecoEntrega.Cidade.ID,
		},
		Itens: itens,
	}
}
