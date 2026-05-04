package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // L'underscore est vital : il charge le driver SQLite en tâche de fond
)

// DB est une variable globale qui nous permettra d'utiliser la base de données partout dans notre code
var DB *sql.DB

func InitDB() {
	var err error
	// 1. Ouvrir la connexion (cela crée le fichier forum.db s'il n'existe pas)
	DB, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal("❌ Erreur d'ouverture de la base de données :", err)
	}

	// 2. Vérifier que la base répond bien
	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ La base de données ne répond pas :", err)
	}

	// 3. Lire le fichier init.sql
	sqlScript, err := os.ReadFile("database/init.sql")
	if err != nil {
		log.Fatal("❌ Erreur de lecture du fichier init.sql :", err)
	}

	// 4. Exécuter le script pour créer les tables
	_, err = DB.Exec(string(sqlScript))
	if err != nil {
		log.Fatal("❌ Erreur lors de la création des tables :", err)
	}

	fmt.Println("📦 Base de données SQLite initialisée avec succès !")
}