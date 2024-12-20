package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.SayHelloAgain(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Status=%s, Task ID=%s\n", r.GetMessage(), r.GetUuid())

	// Simulate waiting for task completion (in a real-world app, this might be an event-based check)
	taskID := r.GetUuid()       // Replace with actual UUID received
	time.Sleep(6 * time.Second) // Simulate waiting for the task to finish

	// Query task status
	statusResp, err := c.GetTaskStatus(ctx, &pb.TaskStatusRequest{Uuid: taskID})
	if err != nil {
		log.Fatalf("could not get task status: %v", err)
	}
	log.Printf("Task %s status: %s - %s", statusResp.GetUuid(), statusResp.GetStatus(), statusResp.GetMessage())
}
