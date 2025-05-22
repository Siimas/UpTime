package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"uptime/internal/models"
	"uptime/internal/postgres"
)

func GetAllMonitors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := postgres.NewConnection(ctx)

	monitors, err := postgres.GetAllMonitors(ctx, db)
	if err != nil {
		log.Printf("🔌 API - 🔴 Failed to fetch monitors: %v\n", err)
		http.Error(w, "Failed to fetch monitors", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(monitors)
	if err != nil {
		log.Printf("🔌 API - 🔴 JSON marshal error: %v\n", err)
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func GetSingleMonitor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := postgres.NewConnection(ctx)
	monitorId := r.PathValue("monitorId")

	monitor, err := postgres.GetSingleMonitor(ctx, db, monitorId)
	if err != nil {
		log.Printf("🔌 API - 🔴 Failed to fetch monitor: %v\n", err)
		http.Error(w, "Failed to fetch monitor", http.StatusNotFound)
		return
	}

	response, err := json.Marshal(monitor)
	if err != nil {
		log.Printf("🔌 API - 🔴 JSON marshal error: %v\n", err)
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func CreateMonitor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := postgres.NewConnection(ctx)

	var monitor models.MonitorCreateDTO

	if err := json.NewDecoder(r.Body).Decode(&monitor); err != nil {
		log.Printf("🔌 API - 🔴 JSON decoder error: %v\n", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	response, err := postgres.CreateMonitor(ctx, db, monitor)
	if err != nil {
		log.Printf("🔌 API - 🔴 Failed to create monitor: %v\n", err)
		http.Error(w, "Failed to create monitor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(response))
}

func UpdateMonitor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := postgres.NewConnection(ctx)

	var monitor models.MonitorUpdateDTO

	if err := json.NewDecoder(r.Body).Decode(&monitor); err != nil {
		log.Printf("🔌 API - 🔴 JSON decoder error: %v\n", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	rowsAffected, err := postgres.UpdateMonitor(ctx, db, monitor)
	if err != nil {
		log.Printf("🔌 API - 🔴 Failed to update monitor: %v\n", err)
		http.Error(w, "Failed to update monitor", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Monitor not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteMonitor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := postgres.NewConnection(ctx)
	monitorId := r.PathValue("monitorId")

	rowsAffected, err := postgres.DeleteMonitor(ctx, db, monitorId)
	if err != nil {
		log.Printf("🔌 API - 🔴 Failed to delete monitor: %v\n", err)
		http.Error(w, "Failed to update monitor", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Monitor not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
