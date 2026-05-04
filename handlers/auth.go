package handlers

import (
	"html/template"
	"log"
	"net/http"

	"forum/database" // Pour parler à ta DB

	"github.com/gofrs/uuid"           // Pour générer l'ID unique
	"golang.org/x/crypto/bcrypt"      // Pour sécuriser le mot de passe
	"time"                            // Pour gérer l'expiration des sessions
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Si la méthode est GET, on affiche simplement le formulaire
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	// 2. Si la méthode est POST, on traite les données du formulaire
	if r.Method == http.MethodPost {
		// Récupération des valeurs tapées par l'utilisateur
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Vérification basique (pour éviter les champs vides)
		if username == "" || email == "" || password == "" {
			http.Error(w, "Tous les champs sont obligatoires", http.StatusBadRequest)
			return
		}

		// A. Génération de l'UUID
		userID, err := uuid.NewV4()
		if err != nil {
			log.Println("Erreur génération UUID :", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		// B. Hachage du mot de passe avec Bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Erreur hachage mot de passe :", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		// C. Insertion dans la base de données SQLite
		query := `INSERT INTO users (id, username, email, password) VALUES (?, ?, ?, ?)`
		_, err = database.DB.Exec(query, userID.String(), username, email, string(hashedPassword))
		if err != nil {
			log.Println("Erreur insertion BDD :", err)
			// En réalité, il faudrait vérifier si l'erreur vient d'un email déjà existant (UNIQUE constraint)
			http.Error(w, "Cet utilisateur ou cet email existe déjà", http.StatusConflict)
			return
		}

		// D. Si tout est bon, on redirige vers l'accueil (ou plus tard vers la page de connexion)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Si ce n'est ni GET ni POST
	http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("templates/login.html")
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// 1. Chercher l'utilisateur dans la base de données
		var storedHash, userID string
		query := `SELECT id, password FROM users WHERE email = ?`
		err := database.DB.QueryRow(query, email).Scan(&userID, &storedHash)
		if err != nil {
			http.Error(w, "Email ou mot de passe incorrect", http.StatusUnauthorized)
			return
		}

		// 2. Comparer le mot de passe tapé avec le hash de la BDD
		err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
		if err != nil {
			http.Error(w, "Email ou mot de passe incorrect", http.StatusUnauthorized)
			return
		}

		// 3. Créer une session avec un UUID
		sessionID, _ := uuid.NewV4()
		expiresAt := time.Now().Add(24 * time.Hour) // La session expire dans 24h

		// 4. Sauvegarder la session dans la base de données
		insertSession := `INSERT INTO sessions (uuid, user_id, expires_at) VALUES (?, ?, ?)`
		_, err = database.DB.Exec(insertSession, sessionID.String(), userID, expiresAt)
		if err != nil {
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		// 5. Envoyer le cookie au navigateur de l'utilisateur
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionID.String(),
			Expires:  expiresAt,
			HttpOnly: true, // Sécurité : empêche le JS de lire le cookie
			Path:     "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}