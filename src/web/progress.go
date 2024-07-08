package web

import (
	"encoding/json"
	"net/http"
)

// VideoProgressGet returns video progress for ajax calls in admin
func (s *Service) VideoProgressGet(w http.ResponseWriter, _ *http.Request) {
	// taskId, _ := strconv.Atoi(r.URL.Query().Get("task_id"))
	s.tasksProgressMu.RLock()
	body, _ := json.Marshal(s.tasksProgress)
	s.tasksProgressMu.RUnlock()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write(body)
}
