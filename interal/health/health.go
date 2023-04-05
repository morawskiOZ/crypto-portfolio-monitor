package health

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// StartHealthCheck starts the health check server
func StartHealthCheck() {
	http.HandleFunc("/health", handleHealthCheck)
	fmt.Printf("Health check listening on port %s\n", viper.GetString("APP_PORT"))
	log.Fatal(http.ListenAndServe(":"+viper.GetString("APP_PORT"), nil))
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}
