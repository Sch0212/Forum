package database
package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler gère la route d'inscription
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Si la méthode est GET, on affiche simplement la page HTML
		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "templates/register.html")
			return
		}

		// Si la méthode est POST, on traite le formulaire
		if r.Method == http.MethodPost {
			email := r.FormValue("email")
			username := r.FormValue("username")
			password := r.FormValue("password")

			// Vérification que les champs ne sont pas vides
			if email == "" || username == "" || password == "" {
				http.Error(w, "Tous les champs sont obligatoires", http.StatusBadRequest)
				return
			}

			// Hachage du mot de passe avec bcrypt
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
				return
			}

			// Insertion sécurisée dans la base de données
			stmt, err := db.Prepare("INSERT INTO users (email, username, password) VALUES (?, ?, ?)")
			if err != nil {
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
				return
			}
			defer stmt.Close()

			_, err = stmt.Exec(email, username, hashedPassword)
			if err != nil {
				// Vérifie si l'erreur vient de la contrainte UNIQUE (email ou pseudo déjà pris)
				if strings.Contains(err.Error(), "UNIQUE constraint failed") {
					http.Error(w, "Cet email ou ce nom d'utilisateur est déjà utilisé", http.StatusConflict)
					return
				}
				http.Error(w, "Erreur lors de la création du compte", http.StatusInternalServerError)
				return
			}

			// Inscription réussie, redirection vers la page de connexion
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Si ce n'est ni GET ni POST
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}