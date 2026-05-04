package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/database" // Importe ton nouveau dossier database
	"forum/handlers"
)

func main() {
	// INITIALISER LA BASE DE DONNÉES EN PREMIER
	database.InitDB()

	// 1. Gérer les fichiers statiques (CSS, images)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 2. Définir la route de la page d'accueil
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	// 3. Lancer le serveur sur le port 8080
	port := ":8080"
	fmt.Println("🌐 Serveur démarré sur http://localhost" + port)
	
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Erreur lors du démarrage du serveur : ", err)
	}
}