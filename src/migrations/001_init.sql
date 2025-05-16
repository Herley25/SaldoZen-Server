-- users
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

-- expenses
CREATE TABLE IF NOT EXISTS expenses (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  descricao TEXT NOT NULL,
  valor NUMERIC(10,2) CHECK (valor > 0),
  vencimento DATE NOT NULL,
  paga BOOLEAN DEFAULT false,
  data_pagamento DATE,
  categoria TEXT NOT NULL,
  observacoes TEXT,
  created_at TIMESTAMP DEFAULT NOW()
);

-- categories
CREATE TABLE IF NOT EXISTS categories (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
  UNIQUE (user_id, name) -- impede categorias duplicadas para o mesmo usuário
);