package handlers

import (
	"html/template"
	"net/http"
)

// HomeHandler gère la route "/"
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// On vérifie que l'URL est bien "/" et pas une page inexistante
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// On charge le fichier HTML
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		return
	}

	// On affiche la page
	tmpl.Execute(w, nil)
}