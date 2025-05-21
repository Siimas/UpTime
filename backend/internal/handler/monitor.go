package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"uptime/internal/postgres"
)

func GetMonitors(w http.ResponseWriter, r *http.Request)  {
	
}

func GetSingleMonitor(w http.ResponseWriter, r *http.Request)  {
	db := postgres.NewConnection(r.Context())
	monitorId := r.PathValue("monitorId")

	m, err := postgres.GetSingleMonitor(r.Context(), db, monitorId)
	if err != nil {
		log.Println("🔌 API: 🔴 Get Single Monitor Error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	mJson, err := json.Marshal(m)
	if err != nil {
		log.Println("🔌 API: 🔴 Marshal Error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(mJson)
}

func CreateMonitor(w http.ResponseWriter, r *http.Request)  {
	
}

func UpdateMonitor(w http.ResponseWriter, r *http.Request)  {
	
}

func DeletMonitor(w http.ResponseWriter, r *http.Request)  {
	
}