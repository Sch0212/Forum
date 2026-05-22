package handlers

import (
	"html/template"
	"net/http"

	"forum/database"
)

// Structure représentant un Post
type Post struct {
	ID      string // <-- Ajout de l'ID ici
	Title   string
	Content string
	Author  string
}

// Structure envoyée à la page HTML
type PageData struct {
	Username string
	Posts    []Post
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := PageData{}

	// 1. Vérification de la session
	cookie, err := r.Cookie("session_token")
	if err == nil {
		query := `
			SELECT u.username 
			FROM users u 
			INNER JOIN sessions s ON u.id = s.user_id 
			WHERE s.uuid = ? AND s.expires_at > CURRENT_TIMESTAMP
		`
		var username string
		err = database.DB.QueryRow(query, cookie.Value).Scan(&username)
		
		if err == nil {
			data.Username = username
		}
	}

	// 2. Récupération des posts avec l'ID inclus
	rows, err := database.DB.Query(`
		SELECT p.id, p.title, p.content, u.username 
		FROM posts p 
		INNER JOIN users u ON p.user_id = u.id 
		ORDER BY p.created_at DESC
	`)
	if err == nil {
		defer rows.Close()
		var posts []Post
		for rows.Next() {
			var p Post
			// On scanne l'ID en premier pour correspondre au SELECT
			rows.Scan(&p.ID, &p.Title, &p.Content, &p.Author)
			posts = append(posts, p)
		}
		data.Posts = posts
	}

	// 3. Affichage du template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}