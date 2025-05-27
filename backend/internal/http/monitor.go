package http

import (
	"encoding/json"
	"log"
	"net/http"
	"uptime/internal/constants"
	"uptime/internal/models"
	"uptime/internal/postgres"
)

func (s *Server) handleGetAllMonitors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	monitors, err := postgres.GetAllMonitors(ctx, s.PostgresDB)
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

func (s *Server) handleGetSingleMonitor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monitorId := r.PathValue("monitorId")

	monitor, err := postgres.GetSingleMonitor(ctx, s.PostgresDB, monitorId)
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

func (s *Server) handleCreateMonitor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var monitor models.MonitorCreateDTO

	if err := json.NewDecoder(r.Body).Decode(&monitor); err != nil {
		log.Printf("🔌 API - 🔴 JSON decoder error: %v\n", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	monitorId, err := postgres.CreateMonitor(ctx, s.PostgresDB, monitor)
	if err != nil {
		log.Printf("🔌 API - 🔴 Failed to create monitor: %v\n", err)
		http.Error(w, "Failed to create monitor", http.StatusInternalServerError)
		return
	}

	if monitor.Active {
		if err := s.kafkaProducer.ProduceMessage(
			constants.KafkaMonitorScheduleTopic,
			monitorId,
			models.MonitorEvent{
				Action:    models.MonitorCreate,
				MonitorId: monitorId,
			},
		); err != nil {
			log.Printf("🔌 API - 🔴 Failed to schedule monitor: %v\n", err)
			http.Error(w, "Failed to schedule monitor", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(monitorId))
}

func (s *Server) handleUpdateMonitor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var monitor models.MonitorUpdateDTO

	if err := json.NewDecoder(r.Body).Decode(&monitor); err != nil {
		log.Printf("🔌 API - 🔴 JSON decoder error: %v\n", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	rowsAffected, err := postgres.UpdateMonitor(ctx, s.PostgresDB, monitor)
	if err != nil {
		log.Printf("🔌 API - 🔴 Failed to update monitor: %v\n", err)
		http.Error(w, "Failed to update monitor", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Monitor not found", http.StatusNotFound)
		return
	}

	if err := s.kafkaProducer.ProduceMessage(
		constants.KafkaMonitorScheduleTopic,
		monitor.Id,
		models.MonitorEvent{
			Action:    models.MonitorUpdate,
			MonitorId: monitor.Id,
		},
	); err != nil {
		log.Printf("🔌 API - 🔴 Failed to schedule monitor: %v\n", err)
		http.Error(w, "Failed to schedule monitor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleDeleteMonitor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monitorId := r.PathValue("monitorId")

	rowsAffected, err := postgres.DeleteMonitor(ctx, s.PostgresDB, monitorId)
	if err != nil {
		log.Printf("🔌 API - 🔴 Failed to delete monitor: %v\n", err)
		http.Error(w, "Failed to update monitor", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Monitor not found", http.StatusNotFound)
		return
	}

	if err := s.kafkaProducer.ProduceMessage(
		constants.KafkaMonitorScheduleTopic,
		monitorId,
		models.MonitorEvent{
			Action:    models.MonitorDelete,
			MonitorId: monitorId,
		},
	); err != nil {
		log.Printf("🔌 API - 🔴 Failed to schedule monitor: %v\n", err)
		http.Error(w, "Failed to schedule monitor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
