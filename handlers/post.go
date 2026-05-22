package handlers

import (
	"html/template"
	"net/http"

	"forum/database"

	"github.com/gofrs/uuid"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// 1. On vérifie que l'utilisateur est connecté via son cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		// S'il n'a pas de cookie, on le renvoie vers la page de connexion
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 2. On récupère son user_id grâce à sa session
	var userID string
	err = database.DB.QueryRow("SELECT user_id FROM sessions WHERE uuid = ? AND expires_at > CURRENT_TIMESTAMP", cookie.Value).Scan(&userID)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 3. Affichage du formulaire
	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("templates/create_post.html")
		tmpl.Execute(w, nil)
		return
	}

	// 4. Traitement du formulaire
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")

		if title == "" || content == "" {
			http.Error(w, "Le titre et le contenu sont obligatoires", http.StatusBadRequest)
			return
		}

		postID, _ := uuid.NewV4()

		// Insertion du post en liant l'ID de l'utilisateur connecté
		query := `INSERT INTO posts (id, user_id, title, content) VALUES (?, ?, ?, ?)`
		_, err = database.DB.Exec(query, postID.String(), userID, title, content)
		if err != nil {
			http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}