package handlers

import (
	"html/template"
	"net/http"

	"forum/database"
)

// Structure pour un commentaire
type Comment struct {
	Content string
	Author  string
}

// Structure qui regroupe toutes les données nécessaires pour la page post.html
type ViewPostData struct {
	Username string
	Post     Post      // On réutilise la structure Post définie dans home.go
	Comments []Comment
}

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	// 1. On récupère l'ID du post dans l'URL (ex: /post?id=12345)
	postID := r.URL.Query().Get("id")
	if postID == "" {
		http.Error(w, "Post introuvable", http.StatusBadRequest)
		return
	}

	data := ViewPostData{}

	// 2. On vérifie si l'utilisateur est connecté (comme sur la page d'accueil)
	cookie, err := r.Cookie("session_token")
	if err == nil {
		var username string
		database.DB.QueryRow(`
			SELECT u.username FROM users u 
			INNER JOIN sessions s ON u.id = s.user_id 
			WHERE s.uuid = ? AND s.expires_at > CURRENT_TIMESTAMP
		`, cookie.Value).Scan(&username)
		data.Username = username
	}

	// 3. On récupère les détails du Post principal
	err = database.DB.QueryRow(`
		SELECT p.id, p.title, p.content, u.username 
		FROM posts p 
		INNER JOIN users u ON p.user_id = u.id 
		WHERE p.id = ?
	`, postID).Scan(&data.Post.ID, &data.Post.Title, &data.Post.Content, &data.Post.Author)
	
	if err != nil {
		http.Error(w, "Post introuvable", http.StatusNotFound)
		return
	}

	// 4. On récupère tous les commentaires associés à ce post
	rows, err := database.DB.Query(`
		SELECT c.content, u.username 
		FROM comments c 
		INNER JOIN users u ON c.user_id = u.id 
		WHERE c.post_id = ? 
		ORDER BY c.created_at ASC
	`, postID)
	
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var c Comment
			rows.Scan(&c.Content, &c.Author)
			data.Comments = append(data.Comments, c)
		}
	}

	// 5. On affiche la page
	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}