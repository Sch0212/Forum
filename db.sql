-- Activer les clés étrangères dans SQLite
PRAGMA foreign_keys = ON;

-- Table des utilisateurs (Gère l'inscription classique et OAuth, plus la modération)
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password TEXT, -- Peut être NULL si l'utilisateur se connecte via OAuth
    role TEXT DEFAULT 'user', -- Rôles possibles : user, moderator, admin
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Table des sessions (Basée sur des UUID pour une sécurité "under-pressure")
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY, -- Stockera l'UUID de la session
    user_id INTEGER NOT NULL,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Table des catégories (Sous-forums)
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

-- Table des posts (Inclut l'upload d'image du Bonus)
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    image_path TEXT, -- Chemin vers l'image uploadée (Bonus)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Table de liaison (Many-to-Many) entre Posts et Catégories
CREATE TABLE IF NOT EXISTS post_categories (
    post_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

-- Table des commentaires
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Table des Likes/Dislikes (Gère à la fois les posts et les commentaires)
CREATE TABLE IF NOT EXISTS likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER,    -- NULL si c'est un like sur un commentaire
    comment_id INTEGER, -- NULL si c'est un like sur un post
    is_like BOOLEAN NOT NULL, -- TRUE pour Like, FALSE pour Dislike
    UNIQUE(user_id, post_id, comment_id), -- Un user ne peut liker/disliker qu'une seule fois la même cible
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE
);