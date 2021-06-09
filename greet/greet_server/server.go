package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/egnimos/grpc-practice/greetpb"
	"google.golang.org/grpc"
)

//server is used to implement the greetpb.GreeterServer
type server struct {
	//embed the unimplemented server
	greetpb.UnimplementedGreetServiceServer
}

//greet: this method makes the unary calls
func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet Function was Invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	result := "Hello you fucker " + firstName + " " + lastName
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function has been invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	for i := 0; i < 1000; i++ {
		result := "Hello you fucker " + firstName + " " + lastName + " " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		//send the response in the stream
		stream.Send(res)
		//delay the time
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func main() {
	fmt.Println("helllo world")

	//define the ports and connection
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	//create a new grpc server
	s := grpc.NewServer()
	//register your server greet proto server to the newly created server
	greetpb.RegisterGreetServiceServer(s, &server{})

	//bind the port
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to listen %v", err)
	}
}
