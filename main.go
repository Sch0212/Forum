package main

import (
	"log"
	"net/http"

	// Attention: Remplace "forum" par le nom exact écrit à la première ligne de ton fichier go.mod
	"forum/internal/database"
	"forum/internal/handlers"
)

func main() {
	// 1. Initialisation de la base de données SQLite
	// Le fichier "forum.db" sera créé à la racine du projet s'il n'existe pas
	err := database.InitDB("forum.db")
	if err != nil {
		log.Fatalf("Erreur critique : impossible d'initialiser la base de données : %v", err)
	}
	// Important : Dans ton fichier database.db.go, n'oublie pas de fermer la DB quand le programme s'arrête
	// (ex: defer database.DB.Close() si tu gères la fermeture ici)

	// 2. Gestion des fichiers statiques (CSS, JS, images uploadées)
	// Cela permet au navigateur d'accéder aux fichiers du dossier "static" via l'URL "/static/..."
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 3. Déclaration des routes (Endpoints)
	// Route d'inscription que nous avons créée précédemment
	http.HandleFunc("/register", handlers.RegisterHandler(database.DB))

	// Routes futures (à décommenter et implémenter plus tard)
	// http.HandleFunc("/login", handlers.LoginHandler(database.DB))
	// http.HandleFunc("/logout", handlers.LogoutHandler(database.DB))
	// http.HandleFunc("/", handlers.HomeHandler(database.DB)) // Page d'accueil avec les posts

	// 4. Lancement du serveur HTTP
	port := ":8080"
	log.Printf("Serveur démarré avec succès. Accès via : http://localhost%s", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur HTTP : %v", err)
	}
}
