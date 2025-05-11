package orchestrator

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/alextrufmanov/asyncCalculator/pkg/config"
	"github.com/alextrufmanov/asyncCalculator/proto"
	"google.golang.org/grpc"
)

type Server struct {
	proto.AsyncCalculatorServiceServer
	storage *Storage
}

func NewServer(s *Storage) *Server {
	return &Server{
		storage: s,
	}
}

func (s *Server) GetTask(ctx context.Context, in *proto.Undef) (*proto.Task, error) {
	task, res := s.storage.GetTask()
	if res {
		log.Println("grps: GetTask")
		return &proto.Task{
			Id:            task.Id,
			Arg1:          task.Arg1,
			Arg2:          task.Arg2,
			Operation:     task.Operation,
			OperationTime: task.OperationTime,
		}, nil
	}
	return &proto.Task{}, fmt.Errorf("empty tasks list")
}

func (s *Server) PostTaskResult(ctx context.Context, in *proto.TaskResult) (*proto.Undef, error) {
	res := s.storage.SetTaskResult(in.Id, in.Result, in.Success)
	if res {
		log.Printf("grps: PostTaskResult")
		return &proto.Undef{}, nil
	}
	return &proto.Undef{}, fmt.Errorf("error post task resul")
}

// Функция создания и запуска сервера оркестратора
func StartGrpcOrchestrator(cfg config.Cfg, storage *Storage) {
	// запускаем сервер оркестратора
	log.Printf("Grps server started on %s", cfg.GrpcAddr)
	lis, err := net.Listen("tcp", cfg.GrpcAddr) // будем ждать запросы по этому адресу
	if err != nil {
		log.Fatal("... with error:", err)
		os.Exit(1)
	}
	// создадим сервер grpc
	grpcServer := grpc.NewServer()
	rpcServer := NewServer(storage)
	// зарегистрируем нашу реализацию сервера
	proto.RegisterAsyncCalculatorServiceServer(grpcServer, rpcServer)
	// запустим grpc сервер
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Println("error serving grpc: ", err)
			os.Exit(1)
		}
	}()
}
