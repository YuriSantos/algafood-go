-- Drop tables in reverse order of creation (respecting foreign key constraints)

DROP TABLE IF EXISTS item_pedido;
DROP TABLE IF EXISTS pedido;
DROP TABLE IF EXISTS foto_produto;
DROP TABLE IF EXISTS produto;
DROP TABLE IF EXISTS restaurante_usuario_responsavel;
DROP TABLE IF EXISTS restaurante_forma_pagamento;
DROP TABLE IF EXISTS restaurante;
DROP TABLE IF EXISTS usuario_grupo;
DROP TABLE IF EXISTS usuario;
DROP TABLE IF EXISTS grupo_permissao;
DROP TABLE IF EXISTS grupo;
DROP TABLE IF EXISTS permissao;
DROP TABLE IF EXISTS forma_pagamento;
DROP TABLE IF EXISTS cozinha;
DROP TABLE IF EXISTS cidade;
DROP TABLE IF EXISTS estado;
