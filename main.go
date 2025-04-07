package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Kommandozeilenparameter definieren:
	// -seed: Falls gesetzt, werden zufällige Produkte eingefügt.
	// -count: Anzahl der zufälligen Produkte, die eingefügt werden sollen.
	seedFlag := flag.Bool("seed", false, "Seed die Datenbank mit zufälligen Produkten")
	count := flag.Int("count", 0, "Anzahl zufälliger Produkte")
	flag.Parse()

	// Lade die .env-Datei, damit Umgebungsvariablen verfügbar sind.
	if err := godotenv.Load(); err != nil {
		log.Println("Keine .env Datei gefunden")
	}

	// Initialisiere die App mit den Umgebungsvariablen.
	a := App{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	// Wenn das Seed-Flag gesetzt ist, führe das Seeding aus.
	if *seedFlag {
		if err := SeedProducts(a.DB, *count); err != nil {
			log.Fatalf("Fehler beim Seeding: %v", err)
		}
		log.Printf("Seeding abgeschlossen. %d Produkte wurden eingefügt.", *count)
	}

	// Starte den HTTP-Server auf Port 8010.
	a.Run(":8010")
}
