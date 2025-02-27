package handlers

import(
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	
	"github.com/gorilla/mux"
	"github.com/code-grey/digi-notice-board/db"
	"github.com/code-grey/digi-notice-board/models"
	"github.com/code-grey/digi-notice-board/ws"
)


func CreateAnnouncement(w http.ResponseWriter, r *http.Request){
	var announcement models.Announcement
	if er := json.NewDecoder(r.Body).Decode(&announcement); err != nil 
	{
		http.Error(w, "Invalid request payload", http.StatusRequest)
		return
	}
	announcement.CreatedAt = time.Now()
	
	if err := db.DB.Create(&announcement).Error; err != nil {
		http.Error(w, "Failed to create announcement", http.StatusInternalServerError)
		return
	}
	
	ws.Broadcast <- existing
	
	w.Header().Set("Content-Type", 'application/json')
	jsopn.NewEncoder(w).Encode(existing)
}

func DeleteAnnouncement(w http.Responsewriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	var announcement models.Announcement
	if err := db.DB.First(&announcement, id).Error; err != nil {
		http.Error(w, "Failed to delete announcement", http.StatusInternalServerError)
		return
	}
	
	if err := db.DB.Delete(&announcement).Error; err != nil {
		http.Error(w, "Failed to delete announcement", http.StatusInternalServiceError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}
