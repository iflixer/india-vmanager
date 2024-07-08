package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"videomanager/database"
	"videomanager/telegram"
	"videomanager/web"
)

func main() {
	log.Println("START")

	log.Println("runtime.GOMAXPROCS:", runtime.GOMAXPROCS(0))

	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Cant load .env: ", err)
	}

	mysqlURL := os.Getenv("MYSQL_URL")
	if os.Getenv("MYSQL_URL_FILE") != "" {
		mysqlURL_, err := os.ReadFile(os.Getenv("MYSQL_URL_FILE"))
		if err != nil {
			log.Fatal(err)
		}
		mysqlURL = strings.TrimSpace(string(mysqlURL_))
	}

	telegramApiToken := os.Getenv("TELEGRAM_APITOKEN")
	if os.Getenv("TELEGRAM_APITOKEN_FILE") != "" {
		telegramApiToken_, err := os.ReadFile(os.Getenv("TELEGRAM_APITOKEN_FILE"))
		if err != nil {
			log.Fatal(err)
		}
		telegramApiToken = strings.TrimSpace(string(telegramApiToken_))
	}

	telegramService, err := telegram.NewService(telegramApiToken)
	if err != nil {
		log.Println(err)
	}

	telegramService.Send(telegram.ChanVideo, fmt.Sprintf("india-vmanager started"))

	promRegistry := prometheus.NewRegistry()

	dbService, err := database.NewService(mysqlURL)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("dbService OK")
	}

	webService, err := web.NewService(dbService, promRegistry, telegramService)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("webService OK")
	}

	rtr := mux.NewRouter()
	rtr.HandleFunc("/alive", getAlive)
	rtr.HandleFunc("/job/get", webService.VideoGetJob).Methods("POST")
	rtr.HandleFunc("/job/done", webService.VideoDone).Methods("POST")
	rtr.HandleFunc("/job/progress/save", webService.VideoProgress).Methods("POST")
	rtr.HandleFunc("/job/progress/get", webService.VideoProgressGet).Methods("GET")
	// rtr.HandleFunc("/converter_alive", webService.ConverterAlive).Methods("POST")

	http.Handle("/", rtr)

	port := os.Getenv("HTTP_PORT")
	log.Println("Listening port :" + port + "...")
	if err = http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}

func getAlive(w http.ResponseWriter, r *http.Request) {
	log.Println("get " + r.RequestURI)
	_, _ = w.Write([]byte("OK"))
}
