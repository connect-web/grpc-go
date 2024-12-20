package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

// Task struct to hold task data and status
type Task struct {
	UUID    string
	Status  string
	Message string
}

type RequestTask struct {
	UUID    string
	Status  string
	Message string
}

var (
	taskQueue = make(map[string]*Task) // In-memory queue of tasks
	mu        sync.Mutex               // Mutex to protect access to taskQueue

	requestTaskQueue = make(map[string]*RequestTask) // In-memory queue of tasks

)

type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello creates a task, adds it to the queue, and returns a UUID
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	taskID := uuid.New().String()
	log.Printf("Received task: %s", taskID)

	// Simulate adding task to queue
	mu.Lock()
	taskQueue[taskID] = &Task{UUID: taskID, Status: "Pending", Message: "Task is pending"}
	mu.Unlock()

	// Asynchronously process the task in a separate goroutine
	go processTask(taskID, in.GetName())

	return &pb.HelloReply{Message: "Task started with UUID: " + taskID}, nil
}

// SayHelloAgain can also create a task (similarly) and return a UUID
func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	taskID := uuid.New().String()
	log.Printf("Received task: %s", taskID)

	// Simulate adding task to queue
	mu.Lock()
	taskQueue[taskID] = &Task{UUID: taskID, Status: "Pending", Message: "Task is pending"}
	mu.Unlock()

	// Asynchronously process the task in a separate goroutine
	go processTask(taskID, in.GetName())

	return &pb.HelloReply{Message: "Task started", Uuid: taskID}, nil
}

// GetTaskStatus allows clients to query the status of a task
func (s *server) GetTaskStatus(ctx context.Context, in *pb.TaskStatusRequest) (*pb.TaskStatusReply, error) {
	mu.Lock()
	task, exists := taskQueue[in.GetUuid()]
	mu.Unlock()

	if !exists {
		log.Printf("Task with UUID %s not found", in.GetUuid())
		return nil, fmt.Errorf("task not found")
	}

	log.Printf("Task %s found with status: %s", in.GetUuid(), task.Status)

	return &pb.TaskStatusReply{
		Uuid:    task.UUID,
		Status:  task.Status,
		Message: task.Message,
	}, nil
}

// Simulate task processing
func processTask(taskID, name string) {
	mu.Lock()
	taskQueue[taskID].Status = "In Progress"
	mu.Unlock()

	// Simulate work (e.g., time-consuming task)
	time.Sleep(5 * time.Second)

	mu.Lock()
	taskQueue[taskID].Status = "Completed"
	taskQueue[taskID].Message = "Hello " + name + ", task completed!"
	mu.Unlock()

	log.Printf("Task %s completed", taskID)
}

func main() {
	// Listen on a TCP port
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Register the Greeter service with the server
	pb.RegisterGreeterServer(s, &server{})

	log.Println("Server listening on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
