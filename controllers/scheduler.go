package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/lucasbyte/go-clipse/scheduler"
)

var isImporting bool
var mu sync.Mutex

func CanImport(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	_, isImporting, _ = scheduler.Atualizou()
	response := map[string]bool{"canImport": isImporting}
	json.NewEncoder(w).Encode(response)
}

func CheckImportAuto(w http.ResponseWriter, r *http.Request) {
	checkboxValue := r.FormValue("check-auto-scheduler")
	fmt.Println(checkboxValue)
	bolean := checkboxValue == "on"
	scheduler.SetAuto(bolean)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func StartPageImport(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	isImporting = scheduler.ReadAuto()
	response := map[string]bool{"canScheduler": isImporting}
	json.NewEncoder(w).Encode(response)
}
