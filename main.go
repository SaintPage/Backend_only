package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

// db es la variable global para la conexión a la base de datos
var db *sql.DB

func main() {
	// Obtener la cadena de conexión desde la variable de entorno o usar la por defecto
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		// Formato: postgres://usuario:contraseña@host:puerto/base_de_datos?sslmode=disable
		dbConnStr = "postgres://postgres:postgres@db:5432/series_tracker?sslmode=disable"
	}

	var err error
	db, err = sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal("Error al abrir la base de datos: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("No se pudo conectar a la base de datos: ", err)
	}
	log.Println("Conectado a la base de datos")

	// Configurar el router y definir los endpoints
	router := mux.NewRouter()

	router.HandleFunc("/api/series", GetAllSeries).Methods("GET")
	router.HandleFunc("/api/series/{id}", GetSeriesByID).Methods("GET")
	router.HandleFunc("/api/series", CreateSeries).Methods("POST")
	router.HandleFunc("/api/series/{id}", UpdateSeries).Methods("PUT")
	router.HandleFunc("/api/series/{id}", DeleteSeries).Methods("DELETE")
	router.HandleFunc("/api/series/{id}/status", UpdateSeriesStatus).Methods("PATCH")
	router.HandleFunc("/api/series/{id}/episode", IncrementEpisode).Methods("PATCH")
	router.HandleFunc("/api/series/{id}/upvote", UpvoteSeries).Methods("PATCH")
	router.HandleFunc("/api/series/{id}/downvote", DownvoteSeries).Methods("PATCH")

	// Configurar CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Println("Servidor corriendo en el puerto 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
