CREATE TABLE album (
    id SERIAL PRIMARY KEY,             -- ID do álbum, auto-incremental
    title VARCHAR(255) NOT NULL,      -- Título do álbum, não pode ser nulo
    artist VARCHAR(255) NOT NULL,     -- Artista do álbum, não pode ser nulo
    price NUMERIC(10, 2) CHECK (price >= 0) -- Preço do álbum, deve ser não negativo
);

