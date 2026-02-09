-- Seed data for testing

-- Estados
INSERT INTO estado (id, nome) VALUES
(1, 'Minas Gerais'),
(2, 'Sao Paulo'),
(3, 'Ceara');

-- Cidades
INSERT INTO cidade (id, nome, estado_id) VALUES
(1, 'Uberlandia', 1),
(2, 'Belo Horizonte', 1),
(3, 'Sao Paulo', 2),
(4, 'Campinas', 2),
(5, 'Fortaleza', 3);

-- Cozinhas
INSERT INTO cozinha (id, nome) VALUES
(1, 'Tailandesa'),
(2, 'Indiana'),
(3, 'Argentina'),
(4, 'Brasileira');

-- Formas de Pagamento
INSERT INTO forma_pagamento (id, descricao) VALUES
(1, 'Cartao de credito'),
(2, 'Cartao de debito'),
(3, 'Dinheiro'),
(4, 'PIX');

-- Permissoes
INSERT INTO permissao (id, nome, descricao) VALUES
(1, 'CONSULTAR_COZINHAS', 'Permite consultar cozinhas'),
(2, 'EDITAR_COZINHAS', 'Permite editar cozinhas'),
(3, 'CONSULTAR_FORMAS_PAGAMENTO', 'Permite consultar formas de pagamento'),
(4, 'EDITAR_FORMAS_PAGAMENTO', 'Permite editar formas de pagamento'),
(5, 'CONSULTAR_CIDADES', 'Permite consultar cidades'),
(6, 'EDITAR_CIDADES', 'Permite editar cidades'),
(7, 'CONSULTAR_ESTADOS', 'Permite consultar estados'),
(8, 'EDITAR_ESTADOS', 'Permite editar estados'),
(9, 'CONSULTAR_USUARIOS_GRUPOS_PERMISSOES', 'Permite consultar usuarios'),
(10, 'EDITAR_USUARIOS_GRUPOS_PERMISSOES', 'Permite criar ou editar usuarios'),
(11, 'CONSULTAR_RESTAURANTES', 'Permite consultar restaurantes'),
(12, 'EDITAR_RESTAURANTES', 'Permite criar, editar ou gerenciar restaurantes'),
(13, 'CONSULTAR_PRODUTOS', 'Permite consultar produtos'),
(14, 'EDITAR_PRODUTOS', 'Permite criar ou editar produtos'),
(15, 'CONSULTAR_PEDIDOS', 'Permite consultar pedidos'),
(16, 'GERENCIAR_PEDIDOS', 'Permite gerenciar pedidos'),
(17, 'GERAR_RELATORIOS', 'Permite gerar relatorios');

-- Grupos
INSERT INTO grupo (id, nome) VALUES
(1, 'Gerente'),
(2, 'Vendedor'),
(3, 'Secretaria'),
(4, 'Cadastrador');

-- Grupo Permissoes
INSERT INTO grupo_permissao (grupo_id, permissao_id) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7), (1, 8), (1, 9), (1, 10),
(1, 11), (1, 12), (1, 13), (1, 14), (1, 15), (1, 16), (1, 17),
(2, 1), (2, 3), (2, 5), (2, 7), (2, 11), (2, 13), (2, 15),
(3, 1), (3, 3), (3, 5), (3, 7), (3, 11), (3, 13), (3, 15), (3, 17),
(4, 1), (4, 5), (4, 7), (4, 11), (4, 12), (4, 13), (4, 14);

-- Usuarios (senha: 123)
INSERT INTO usuario (id, nome, email, senha, data_cadastro) VALUES
(1, 'Joao da Silva', 'joao.ger@algafood.com', '$2a$10$K2p3rLSzT.GxVdVYkVdXc.OjG1HNWqm7EqJj9x5rKjL2bMQrXnSNu', NOW()),
(2, 'Maria Joaquina', 'maria.vnd@algafood.com', '$2a$10$K2p3rLSzT.GxVdVYkVdXc.OjG1HNWqm7EqJj9x5rKjL2bMQrXnSNu', NOW()),
(3, 'Jose Souza', 'jose.aux@algafood.com', '$2a$10$K2p3rLSzT.GxVdVYkVdXc.OjG1HNWqm7EqJj9x5rKjL2bMQrXnSNu', NOW()),
(4, 'Sebastiao Martins', 'sebastiao.cad@algafood.com', '$2a$10$K2p3rLSzT.GxVdVYkVdXc.OjG1HNWqm7EqJj9x5rKjL2bMQrXnSNu', NOW()),
(5, 'Manoel Lima', 'manoel.loja@gmail.com', '$2a$10$K2p3rLSzT.GxVdVYkVdXc.OjG1HNWqm7EqJj9x5rKjL2bMQrXnSNu', NOW()),
(6, 'Debora Mendonca', 'debora@algafood.com', '$2a$10$K2p3rLSzT.GxVdVYkVdXc.OjG1HNWqm7EqJj9x5rKjL2bMQrXnSNu', NOW());

-- Usuario Grupos
INSERT INTO usuario_grupo (usuario_id, grupo_id) VALUES
(1, 1), (1, 2),
(2, 2),
(3, 3), (3, 4),
(4, 4);

-- Restaurantes
INSERT INTO restaurante (id, nome, taxa_frete, cozinha_id, ativo, aberto, endereco_cep, endereco_logradouro, endereco_numero, endereco_bairro, endereco_cidade_id) VALUES
(1, 'Thai Gourmet', 10, 1, true, true, '38400-999', 'Rua Joao Pinheiro', '1000', 'Centro', 1),
(2, 'Thai Delivery', 9.50, 1, true, true, '38400-111', 'Rua Floriano Peixoto', '500', 'Centro', 1),
(3, 'Tuk Tuk Comida Indiana', 15, 2, true, true, '38400-222', 'Rua Antonio Thomaz', '321', 'Centro', 1),
(4, 'Java Steakhouse', 12, 3, true, true, '38400-333', 'Rua Olegario Maciel', '789', 'Centro', 2),
(5, 'Lanchonete do Tio Sam', 11, 4, true, true, '38400-444', 'Av. Brasil', '100', 'Centro', 3),
(6, 'Bar da Maria', 6, 4, true, true, '38400-555', 'Rua das Flores', '50', 'Centro', 3);

-- Restaurante Formas Pagamento
INSERT INTO restaurante_forma_pagamento (restaurante_id, forma_pagamento_id) VALUES
(1, 1), (1, 2), (1, 3), (1, 4),
(2, 1), (2, 2), (2, 4),
(3, 1), (3, 3), (3, 4),
(4, 1), (4, 2), (4, 3), (4, 4),
(5, 1), (5, 2), (5, 4),
(6, 3), (6, 4);

-- Restaurante Responsaveis
INSERT INTO restaurante_usuario_responsavel (restaurante_id, usuario_id) VALUES
(1, 5),
(2, 5),
(3, 5),
(4, 6),
(5, 6),
(6, 6);

-- Produtos
INSERT INTO produto (id, nome, descricao, preco, ativo, restaurante_id) VALUES
(1, 'Porco com curry', 'Deliciosa carne suina ao molho picante', 78.90, true, 1),
(2, 'Camarao tailandes', '16 camaroes grandes ao molho picante', 110, true, 1),
(3, 'Salada picante com carne grelhada', 'Salada de folhas com cortes finos de carne bovina grelhada e temperado', 87.20, true, 1),
(4, 'Garlic Naan', 'Pao tradicional indiano com cobertura de alho', 21, true, 3),
(5, 'Murg Curry', 'Cubos de frango preparados com molho curry', 43, true, 3),
(6, 'Bife Ancho', 'Corte macio e suculento, com aproximadamente 400g', 79, true, 4),
(7, 'T-Bone', 'Corte muito saboroso, com osso em formato de T, com aproximadamente 700g', 89, true, 4),
(8, 'Sanduiche X-Tudo', 'Sandube com muito bacon, ovo, presunto, queijo, hamburguer bovino, cebola, maionese, ketchup', 19, true, 5),
(9, 'Espetinho de Cupim', 'Acompanha farinha, vinagrete e uma porcao de pao de alho', 8, true, 6);
