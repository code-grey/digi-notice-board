package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"

    "github.com/gorilla/mux"

    "digi-notice-board/db"
    "digi-notice-board/models"
    "digi-notice-board/ws"
)

// ListAnnouncements handles GET /announcements
func ListAnnouncements(w http.ResponseWriter, r *http.Request) {
    var announcements []models.Announcement
    if err := db.DB.Find(&announcements).Error; err != nil {
        http.Error(w, "Failed to list announcements", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(announcements)
}

// GetAnnouncement handles GET /announcements/{id}
func GetAnnouncement(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    idStr := params["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var announcement models.Announcement
    if err := db.DB.First(&announcement, id).Error; err != nil {
        http.Error(w, "Announcement not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(announcement)
}

// UpdateAnnouncement handles PUT /announcements/{id}
func UpdateAnnouncement(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    idStr := params["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var existing models.Announcement
    if err := db.DB.First(&existing, id).Error; err != nil {
        http.Error(w, "Announcement not found", http.StatusNotFound)
        return
    }

    var updated models.Announcement
    if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Update fields
    existing.Title = updated.Title
    existing.Content = updated.Content
    // existing.CreatedAt stays the same, or change if needed

    // Save changes to DB
    if err := db.DB.Save(&existing).Error; err != nil {
        http.Error(w, "Failed to update announcement", http.StatusInternalServerError)
        return
    }

    // Optionally broadcast updated announcement
    // ws.Broadcast <- existing

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(existing)
}

// CreateAnnouncement handles POST /announcements
func CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
    var announcement models.Announcement
    if err := json.NewDecoder(r.Body).Decode(&announcement); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    announcement.CreatedAt = time.Now()

    if err := db.DB.Create(&announcement).Error; err != nil {
        http.Error(w, "Failed to create announcement", http.StatusInternalServerError)
        return
    }

    // Broadcast new announcement if you want real-time
    ws.Broadcast <- announcement

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(announcement)
}

// DeleteAnnouncement handles DELETE /announcements/{id}
func DeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    idStr := params["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var announcement models.Announcement
    if err := db.DB.First(&announcement, id).Error; err != nil {
        http.Error(w, "Announcement not found", http.StatusNotFound)
        return
    }

    if err := db.DB.Delete(&announcement).Error; err != nil {
        http.Error(w, "Failed to delete announcement", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent) // 204
}

