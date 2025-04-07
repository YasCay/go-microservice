package main

import (
	"database/sql"
)

// product repräsentiert ein Produkt in der Datenbank.
type product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// getProduct liest den Namen und Preis des Produkts anhand der ID aus der Datenbank.
func (p *product) getProduct(db *sql.DB) error {
	return db.QueryRow("SELECT name, price FROM products WHERE id=$1", p.ID).
		Scan(&p.Name, &p.Price)
}

// updateProduct aktualisiert den Namen und Preis eines existierenden Produkts.
func (p *product) updateProduct(db *sql.DB) error {
	_, err := db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3",
		p.Name, p.Price, p.ID)
	return err
}

// deleteProduct entfernt ein Produkt aus der Datenbank.
func (p *product) deleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products WHERE id=$1", p.ID)
	return err
}

// createProduct legt ein neues Produkt in der Datenbank an und setzt die ID des Produkts.
func (p *product) createProduct(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO products(name, price) VALUES($1, $2) RETURNING id",
		p.Name, p.Price).Scan(&p.ID)
	return err
}

// getProducts ruft eine Liste von Produkten ab, beginnend ab 'start' und begrenzt auf 'count' Produkte.
func getProducts(db *sql.DB, start, count int) ([]product, error) {
	rows, err := db.Query(
		"SELECT id, name, price FROM products LIMIT $1 OFFSET $2",
		count, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []product{}

	for rows.Next() {
		var p product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
