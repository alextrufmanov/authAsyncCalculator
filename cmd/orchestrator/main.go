package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alextrufmanov/asyncCalculator/pkg/config"
	"github.com/alextrufmanov/asyncCalculator/pkg/orchestrator"
)

func main() {

	// получаем настройки
	config := *config.NewCfg()

	// создаем БД хранилище
	dbStorage, err := orchestrator.NewDBStorage()
	if err != nil {
		log.Printf("Init DBStorage error %v", err)
	}
	defer dbStorage.Close()

	// создаем хранилище арифметических выражений и задач агента
	storage := orchestrator.NewStorage(config, dbStorage)

	osSignalCh := make(chan os.Signal, 1)
	signal.Notify(osSignalCh, os.Interrupt, syscall.SIGTERM)

	log.Printf("Orchestrator started")

	orchestrator.StartHttpOrchestrator(config, storage, dbStorage)
	orchestrator.StartGrpcOrchestrator(config, storage)

	<-osSignalCh

}
