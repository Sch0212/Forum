package handlers

import (
	"net/http"

	"forum/database"
	"github.com/gofrs/uuid"
)

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// 1. On vérifie que l'utilisateur est bien connecté
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 2. On récupère l'ID de l'utilisateur
	var userID string
	err = database.DB.QueryRow("SELECT user_id FROM sessions WHERE uuid = ? AND expires_at > CURRENT_TIMESTAMP", cookie.Value).Scan(&userID)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 3. On récupère les données du formulaire
	postID := r.FormValue("post_id")
	content := r.FormValue("content")

	if postID == "" || content == "" {
		http.Error(w, "Données manquantes", http.StatusBadRequest)
		return
	}

	// 4. On génère un UUID pour le commentaire et on l'enregistre
	commentID, _ := uuid.NewV4()
	_, err = database.DB.Exec(`
		INSERT INTO comments (id, post_id, user_id, content) 
		VALUES (?, ?, ?, ?)
	`, commentID.String(), postID, userID, content)

	if err != nil {
		http.Error(w, "Erreur lors de l'ajout du commentaire", http.StatusInternalServerError)
		return
	}

	// 5. On redirige l'utilisateur vers la page du post pour qu'il voie son commentaire
	http.Redirect(w, r, "/post?id="+postID, http.StatusSeeOther)
}