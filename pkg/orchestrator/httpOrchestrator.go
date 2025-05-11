package orchestrator

import (
	"log"
	"net/http"

	"github.com/alextrufmanov/asyncCalculator/pkg/config"
	"github.com/gorilla/mux"
)

// Функция создания и запуска сервера оркестратора
func StartHttpOrchestrator(cfg config.Cfg, storage *Storage, db *DBStorage) {
	// "настраиваем" мультиплексор
	router := mux.NewRouter()
	// router.HandleFunc("/", indexHandler)

	router.HandleFunc("/api/v1/register", POSTRegisterHandler(storage, db)).Methods("POST")
	router.HandleFunc("/api/v1/login", POSTLoginHandler(storage, db)).Methods("POST")

	router.HandleFunc("/api/v1/calculate", POSTCalculateHandler(storage, db)).Methods("POST")
	router.HandleFunc("/api/v1/expressions", GETExpressionsHandler(storage, db)).Methods("GET")
	router.HandleFunc("/api/v1/expressions/{id}", GETExpressionByIDHandler(storage, db)).Methods("GET")

	// запускаем сервер оркестратора
	log.Printf("http server started on %s", cfg.HttpAddr)
	go func() {
		err := http.ListenAndServe(cfg.HttpAddr, router)
		if err != nil {
			log.Fatal("... with error:", err)
		}
	}()
}
