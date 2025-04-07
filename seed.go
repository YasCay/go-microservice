package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v6"
)

// SeedProducts fügt 'count' zufällige Produkte in die Datenbank ein.
// db: Datenbankverbindung
// count: Anzahl der zu erstellenden zufälligen Produkte
func SeedProducts(db *sql.DB, count int) error {
	// Initialisiere gofakeit, um Zufallsdaten zu generieren.
	gofakeit.Seed(0)

	// Schleife für die gewünschte Anzahl an Produkten
	for i := 0; i < count; i++ {
		// Generiere einen zufälligen Produktnamen
		name := gofakeit.ProductName()
		// Generiere einen zufälligen Preis zwischen 1 und 100
		price := gofakeit.Price(1.0, 100.0)
		// Führe eine SQL-Anweisung aus, um das Produkt einzufügen
		_, err := db.Exec("INSERT INTO products(name, price) VALUES($1, $2)", name, price)
		if err != nil {
			return fmt.Errorf("Fehler beim Einfügen des Produkts: %v", err)
		}
		log.Printf("Produkt #%d eingefügt: %s, Preis: %.2f", i+1, name, price)
	}
	return nil
}
