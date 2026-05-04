package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // Driver SQLite indispensable
)

var DB *sql.DB

// InitDB initialise la connexion et crée les tables si elles n'existent pas
func InitDB(dataSourceName string) error {
	var err error

	// Connexion à la base de données
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}

	// Vérifie que la connexion est bien établie
	if err = DB.Ping(); err != nil {
		return err
	}

	// Activation des clés étrangères (spécifique à SQLite)
	_, err = DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return err
	}

	// Lecture et exécution du fichier SQL de création
	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		log.Printf("Attention: Impossible de lire schema.sql, assurez-vous qu'il est à la racine : %v", err)
		return err
	}

	_, err = DB.Exec(string(schema))
	if err != nil {
		log.Printf("Erreur lors de l'exécution du schéma : %v", err)
		return err
	}

	log.Println("Base de données initialisée avec succès.")
	return nil
}
