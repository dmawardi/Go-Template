package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dmawardi/Go-Template/internal/models"
)

func GetJobs(w http.ResponseWriter, r *http.Request) {
	var jobs []models.Job

	jobs = append(jobs, models.Job{ID: 1, Name: "Accounting"})
	jobs = append(jobs, models.Job{ID: 2, Name: "Programming"})

	// Set header
	w.Header().Set("Content-Type", "application/json")

	// Build new JSON encoder to write to, then write jobs data
	json.NewEncoder(w).Encode(jobs)
}
