package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetAllSeries obtiene todas las series de la base de datos.
func GetAllSeries(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, description, status, current_episode, score FROM series")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	seriesList := []Series{}
	for rows.Next() {
		var s Series
		if err := rows.Scan(&s.ID, &s.Title, &s.Description, &s.Status, &s.CurrentEpisode, &s.Score); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		seriesList = append(seriesList, s)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(seriesList)
}

// GetSeriesByID obtiene una serie por su ID.
func GetSeriesByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var s Series
	err := db.QueryRow("SELECT id, title, description, status, current_episode, score FROM series WHERE id = $1", id).
		Scan(&s.ID, &s.Title, &s.Description, &s.Status, &s.CurrentEpisode, &s.Score)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Serie no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

// CreateSeries crea una nueva serie.
func CreateSeries(w http.ResponseWriter, r *http.Request) {
	var s Series
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = db.QueryRow(
		"INSERT INTO series(title, description, status, current_episode, score) VALUES($1, $2, $3, $4, $5) RETURNING id",
		s.Title, s.Description, s.Status, s.CurrentEpisode, s.Score).Scan(&s.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

// UpdateSeries actualiza una serie existente.
func UpdateSeries(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var s Series
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.ID, _ = strconv.Atoi(id)

	result, err := db.Exec("UPDATE series SET title = $1, description = $2, status = $3, current_episode = $4, score = $5 WHERE id = $6",
		s.Title, s.Description, s.Status, s.CurrentEpisode, s.Score, s.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

// DeleteSeries elimina una serie por ID.
func DeleteSeries(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	result, err := db.Exec("DELETE FROM series WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UpdateSeriesStatus actualiza el estado de una serie.
func UpdateSeriesStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var payload struct {
		Status string `json:"status"`
	}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := db.Exec("UPDATE series SET status = $1 WHERE id = $2", payload.Status, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// IncrementEpisode incrementa el contador de episodios vistos en una serie.
func IncrementEpisode(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	result, err := db.Exec("UPDATE series SET current_episode = current_episode + 1 WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UpvoteSeries incrementa la puntuación de una serie.
func UpvoteSeries(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	result, err := db.Exec("UPDATE series SET score = score + 1 WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DownvoteSeries decrementa la puntuación de una serie.
func DownvoteSeries(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	result, err := db.Exec("UPDATE series SET score = score - 1 WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}