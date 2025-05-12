package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alextrufmanov/asyncCalculator/pkg/agent"
	"github.com/alextrufmanov/asyncCalculator/pkg/config"
	"github.com/alextrufmanov/asyncCalculator/pkg/models"
	"github.com/alextrufmanov/asyncCalculator/pkg/orchestrator"
)

var token string = ""
var id int32 = 0

func testRegistr(cfg config.Cfg) {
	requestBody, err := json.Marshal(models.POSTRegisterRequestBody{
		Login:    "Test",
		Password: "Test",
	})
	if err != nil {
		log.Fatalf("FAIL : /api/v1/register - marshal request body error: %v", err)
	}
	response, err := http.Post(fmt.Sprintf("http://%s/api/v1/register", cfg.HttpAddr), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("FAIL : /api/v1/register - request error : %v", err)
	}
	if response.StatusCode != http.StatusOK {
		// 	log.Fatalf("FAIL : /api/v1/login - response status code : %v", response.StatusCode)
	}
}

func testLogin(cfg config.Cfg) {
	var answerBody models.POSTLoginAnswerBody
	requestBody, err := json.Marshal(models.POSTLoginRequestBody{
		Login:    "Test",
		Password: "Test",
	})
	if err != nil {
		log.Fatalf("FAIL : /api/v1/login - marshal request body error: %v", err)
	}
	response, err := http.Post(fmt.Sprintf("http://%s/api/v1/login", cfg.HttpAddr), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("FAIL : /api/v1/login - request error : %v", err)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatalf("FAIL : /api/v1/login - response status code : %v", response.StatusCode)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("FAIL : /api/v1/login - read response body error: %v", err)
	}
	if json.Unmarshal(bodyBytes, &answerBody) != nil {
		log.Fatalf("FAIL : /api/v1/login - unmarshal response body error: %v", err)
	}
	token = answerBody.Token
}

func testCalculate(cfg config.Cfg) {
	var answerBody models.POSTCalculateAnswerBody
	requestBody, err := json.Marshal(models.POSTCalculateRequestBody{
		Token:      token,
		Expression: "2+2",
	})
	if err != nil {
		log.Fatalf("FAIL : /api/v1/calculate - marshal request body error: %v", err)
	}
	response, err := http.Post(fmt.Sprintf("http://%s/api/v1/calculate", cfg.HttpAddr), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("FAIL : /api/v1/calculate - request error : %v", err)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatalf("FAIL : /api/v1/calculate - response status code : %v", response.StatusCode)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("FAIL : /api/v1/calculate - read response body error: %v", err)
	}
	if json.Unmarshal(bodyBytes, &answerBody) != nil {
		log.Fatalf("FAIL : /api/v1/calculate - unmarshal response body error: %v", err)
	}
	id = answerBody.Id
}

func testExpressionByID(cfg config.Cfg) {
	requestBody, err := json.Marshal(models.GETExpressionByIDRequestBody{
		Token: token,
	})
	if err != nil {
		log.Fatalf("FAIL : /api/v1/expressions/{id} - marshal request body error: %v", err)
	}
	response, err := http.Post(fmt.Sprintf("http://%s/api/v1/expressions/%d", cfg.HttpAddr, id), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("FAIL : /api/v1/expressions/{id} - request error : %v", err)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatalf("FAIL : /api/v1/expressions/{id} - response status code : %v", response.StatusCode)
	}
}

func testExpressions(cfg config.Cfg) {
	requestBody, err := json.Marshal(models.GETExpressionsRequestBody{
		Token: token,
	})
	if err != nil {
		log.Fatalf("FAIL : /api/v1/expressions - marshal request body error: %v", err)
	}
	response, err := http.Post(fmt.Sprintf("http://%s/api/v1/expressions", cfg.HttpAddr), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("FAIL : /api/v1/expressions - request error : %v", err)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatalf("FAIL : /api/v1/expressions - response status code : %v", response.StatusCode)
	}
}

func main() {

	osSignalCh := make(chan os.Signal, 1)
	signal.Notify(osSignalCh, os.Interrupt, syscall.SIGTERM)

	// получаем настройки
	cfg := *config.NewCfg()
	// создаем БД хранилище
	dbStorage, err := orchestrator.NewDBStorage()
	if err != nil {
		log.Printf("Init DBStorage error %v", err)
	}
	defer dbStorage.Close()

	// создаем хранилище арифметических выражений и задач агента
	storage := orchestrator.NewStorage(cfg, dbStorage)

	// запуск оркестратора
	log.Printf("Orchestrator started")
	orchestrator.StartHttpOrchestrator(cfg, storage, dbStorage)
	orchestrator.StartGrpcOrchestrator(cfg, storage)

	// запуск агента
	go agent.StartGrpcAgents(cfg)

	testRegistr(cfg)
	testLogin(cfg)
	testCalculate(cfg)

	<-osSignalCh
}
