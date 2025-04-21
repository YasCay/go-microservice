package main

import (
	"database/sql"
	"encoding/json"
	"fmt" // Für String-Formatierung
	"log"
	"net/http"
	"os" // Added for environment variables
	"strconv"
	"strings" // Für String-Funktionen

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// App enthält den Router und die DB-Verbindung.
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize baut die DB-Verbindung anhand der übergebenen Parameter auf.
// Falls ein Parameter leer ist, wird ein Fallback auf Standardwerte verwendet.
func (a *App) Initialize(user, password, dbname string) {

	// Get host and port from environment variables or use defaults
	host := os.Getenv("APP_DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("APP_DB_PORT")
	if port == "" {
		port = "5432"
	}

	// Aufbau des Connection-Strings im URL-Format
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)

	log.Printf("Connecting to database: %s:%s", host, port)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	// Test the connection
	err = a.DB.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to database")

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run startet den HTTP-Server auf dem angegebenen Port.
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

// initializeRoutes registriert alle API-Endpunkte.
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
	// Neuer Endpunkt zum Löschen aller Produkte
	a.Router.HandleFunc("/products", a.deleteAllProducts).Methods("DELETE")
}

// getProduct liefert ein einzelnes Produkt anhand der ID.
func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := product{ID: id}
	if err := p.getProduct(a.DB); err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Product not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

// getProducts liefert eine Liste von Produkten.
// Zusätzliche Funktion: Wenn count als "all" übergeben wird, werden alle Produkte zurückgeliefert.
func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	countParam := r.FormValue("count")
	var count int
	if strings.ToLower(countParam) == "all" || countParam == "" {
		// Wenn "all" angegeben oder kein Parameter vorhanden, setze einen hohen Wert
		count = 10000
	} else {
		var err error
		count, err = strconv.Atoi(countParam)
		if err != nil || count < 1 {
			count = 10 // Standardwert, falls ungültiger Wert
		}
	}

	start, _ := strconv.Atoi(r.FormValue("start"))
	if start < 0 {
		start = 0
	}

	products, err := getProducts(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

// createProduct legt ein neues Produkt in der Datenbank an.
func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.createProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, p)
}

// updateProduct aktualisiert ein bestehendes Produkt.
func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	p.ID = id

	if err := p.updateProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

// deleteProduct löscht ein einzelnes Produkt anhand der ID.
func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := product{ID: id}
	if err := p.deleteProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// deleteAllProducts löscht alle Produkte aus der Datenbank
// und setzt die Sequenz für die IDs zurück.
func (a *App) deleteAllProducts(w http.ResponseWriter, r *http.Request) {
	if _, err := a.DB.Exec("DELETE FROM products"); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if _, err := a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1"); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "all products deleted"})
}

// respondWithError sendet eine Fehlermeldung als JSON-Antwort.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON sendet das Payload als JSON-Antwort.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
