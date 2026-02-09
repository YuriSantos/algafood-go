-- AlgaFood Database Schema

-- Estados
CREATE TABLE IF NOT EXISTS estado (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(80) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Cidades
CREATE TABLE IF NOT EXISTS cidade (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(80) NOT NULL,
    estado_id BIGINT NOT NULL,
    CONSTRAINT fk_cidade_estado FOREIGN KEY (estado_id) REFERENCES estado(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Cozinhas
CREATE TABLE IF NOT EXISTS cozinha (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(60) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Formas de Pagamento
CREATE TABLE IF NOT EXISTS forma_pagamento (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    descricao VARCHAR(60) NOT NULL,
    data_atualizacao DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Permissoes
CREATE TABLE IF NOT EXISTS permissao (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    descricao VARCHAR(255)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Grupos
CREATE TABLE IF NOT EXISTS grupo (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(60) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Grupo Permissao (M2M)
CREATE TABLE IF NOT EXISTS grupo_permissao (
    grupo_id BIGINT NOT NULL,
    permissao_id BIGINT NOT NULL,
    PRIMARY KEY (grupo_id, permissao_id),
    CONSTRAINT fk_grupo_permissao_grupo FOREIGN KEY (grupo_id) REFERENCES grupo(id),
    CONSTRAINT fk_grupo_permissao_permissao FOREIGN KEY (permissao_id) REFERENCES permissao(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Usuarios
CREATE TABLE IF NOT EXISTS usuario (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(80) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    senha VARCHAR(255) NOT NULL,
    data_cadastro DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Usuario Grupo (M2M)
CREATE TABLE IF NOT EXISTS usuario_grupo (
    usuario_id BIGINT NOT NULL,
    grupo_id BIGINT NOT NULL,
    PRIMARY KEY (usuario_id, grupo_id),
    CONSTRAINT fk_usuario_grupo_usuario FOREIGN KEY (usuario_id) REFERENCES usuario(id),
    CONSTRAINT fk_usuario_grupo_grupo FOREIGN KEY (grupo_id) REFERENCES grupo(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Restaurantes
CREATE TABLE IF NOT EXISTS restaurante (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(80) NOT NULL,
    taxa_frete DECIMAL(10,2) NOT NULL,
    cozinha_id BIGINT NOT NULL,
    ativo BOOLEAN DEFAULT TRUE,
    aberto BOOLEAN DEFAULT FALSE,
    data_cadastro DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    data_atualizacao DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- Endereco (embedded)
    endereco_cep VARCHAR(9),
    endereco_logradouro VARCHAR(100),
    endereco_numero VARCHAR(20),
    endereco_complemento VARCHAR(60),
    endereco_bairro VARCHAR(60),
    endereco_cidade_id BIGINT,

    CONSTRAINT fk_restaurante_cozinha FOREIGN KEY (cozinha_id) REFERENCES cozinha(id),
    CONSTRAINT fk_restaurante_cidade FOREIGN KEY (endereco_cidade_id) REFERENCES cidade(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Restaurante Forma Pagamento (M2M)
CREATE TABLE IF NOT EXISTS restaurante_forma_pagamento (
    restaurante_id BIGINT NOT NULL,
    forma_pagamento_id BIGINT NOT NULL,
    PRIMARY KEY (restaurante_id, forma_pagamento_id),
    CONSTRAINT fk_rfp_restaurante FOREIGN KEY (restaurante_id) REFERENCES restaurante(id),
    CONSTRAINT fk_rfp_forma_pagamento FOREIGN KEY (forma_pagamento_id) REFERENCES forma_pagamento(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Restaurante Usuario Responsavel (M2M)
CREATE TABLE IF NOT EXISTS restaurante_usuario_responsavel (
    restaurante_id BIGINT NOT NULL,
    usuario_id BIGINT NOT NULL,
    PRIMARY KEY (restaurante_id, usuario_id),
    CONSTRAINT fk_rur_restaurante FOREIGN KEY (restaurante_id) REFERENCES restaurante(id),
    CONSTRAINT fk_rur_usuario FOREIGN KEY (usuario_id) REFERENCES usuario(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Produtos
CREATE TABLE IF NOT EXISTS produto (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(80) NOT NULL,
    descricao VARCHAR(255),
    preco DECIMAL(10,2) NOT NULL,
    ativo BOOLEAN DEFAULT TRUE,
    restaurante_id BIGINT NOT NULL,
    CONSTRAINT fk_produto_restaurante FOREIGN KEY (restaurante_id) REFERENCES restaurante(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Foto Produto
CREATE TABLE IF NOT EXISTS foto_produto (
    id BIGINT PRIMARY KEY,
    produto_id BIGINT NOT NULL UNIQUE,
    nome_arquivo VARCHAR(150) NOT NULL,
    descricao VARCHAR(150),
    content_type VARCHAR(80) NOT NULL,
    tamanho BIGINT NOT NULL,
    CONSTRAINT fk_foto_produto FOREIGN KEY (produto_id) REFERENCES produto(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Pedidos
CREATE TABLE IF NOT EXISTS pedido (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    codigo VARCHAR(36) NOT NULL UNIQUE,
    subtotal DECIMAL(10,2) NOT NULL,
    taxa_frete DECIMAL(10,2) NOT NULL,
    valor_total DECIMAL(10,2) NOT NULL,
    status VARCHAR(15) NOT NULL DEFAULT 'CRIADO',
    data_criacao DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    data_confirmacao DATETIME,
    data_cancelamento DATETIME,
    data_entrega DATETIME,

    restaurante_id BIGINT NOT NULL,
    usuario_cliente_id BIGINT NOT NULL,
    forma_pagamento_id BIGINT NOT NULL,

    -- Endereco Entrega (embedded)
    endereco_cep VARCHAR(9),
    endereco_logradouro VARCHAR(100),
    endereco_numero VARCHAR(20),
    endereco_complemento VARCHAR(60),
    endereco_bairro VARCHAR(60),
    endereco_cidade_id BIGINT,

    CONSTRAINT fk_pedido_restaurante FOREIGN KEY (restaurante_id) REFERENCES restaurante(id),
    CONSTRAINT fk_pedido_usuario FOREIGN KEY (usuario_cliente_id) REFERENCES usuario(id),
    CONSTRAINT fk_pedido_forma_pagamento FOREIGN KEY (forma_pagamento_id) REFERENCES forma_pagamento(id),
    CONSTRAINT fk_pedido_cidade FOREIGN KEY (endereco_cidade_id) REFERENCES cidade(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Item Pedido
CREATE TABLE IF NOT EXISTS item_pedido (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    pedido_id BIGINT NOT NULL,
    produto_id BIGINT NOT NULL,
    quantidade INT NOT NULL,
    preco_unitario DECIMAL(10,2) NOT NULL,
    preco_total DECIMAL(10,2) NOT NULL,
    observacao VARCHAR(255),
    CONSTRAINT fk_item_pedido_pedido FOREIGN KEY (pedido_id) REFERENCES pedido(id),
    CONSTRAINT fk_item_pedido_produto FOREIGN KEY (produto_id) REFERENCES produto(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Indexes for better performance
CREATE INDEX idx_cidade_estado ON cidade(estado_id);
CREATE INDEX idx_restaurante_cozinha ON restaurante(cozinha_id);
CREATE INDEX idx_produto_restaurante ON produto(restaurante_id);
CREATE INDEX idx_pedido_cliente ON pedido(usuario_cliente_id);
CREATE INDEX idx_pedido_restaurante ON pedido(restaurante_id);
CREATE INDEX idx_pedido_status ON pedido(status);
CREATE INDEX idx_pedido_data_criacao ON pedido(data_criacao);
CREATE INDEX idx_item_pedido_pedido ON item_pedido(pedido_id);
