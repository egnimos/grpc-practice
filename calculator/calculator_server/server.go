package main

import (
	"context"
	"fmt"
	"log"
	"net"

	// "strconv"
	"time"

	"github.com/egnimos/grpc-practice/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	startServer()
}

type server struct {
	//unimplemented
	calculatorpb.UnimplementedCalServiceServer
}

func (s *server) Calculate(ctx context.Context, req *calculatorpb.CalRequest) (*calculatorpb.CalResponse, error) {
	//get the inputs from the request
	fmt.Println("Cal function is called or invoke")
	input1 := req.GetCal().GetIntput_1()
	input2 := req.GetCal().GetIntput_2()
	//add
	result := input1 + input2
	res := &calculatorpb.CalResponse{
		Output: result,
	}

	return res, nil
}

func (s *server) CalculateManyTimes(req *calculatorpb.CalManyTimesRequest, stream calculatorpb.CalService_CalculateManyTimesServer) error {
	fmt.Printf("CalculateManyTimes function has been invoked with %v\n", req)
	input1 := req.GetCal().GetIntput_1()
	input2 := req.GetCal().GetIntput_2()
	for i := 0; i < 1000; i++ {
		result := input1 + input2 + int32(i)
		res := &calculatorpb.CalManyTimesResponse{
			Output: result,
		}
		//send the response in the stream
		stream.Send(res)
		//delay the time
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

// startServer
func startServer() {
	fmt.Println("server is started")

	//open the port and connection
	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalln("failed to listen: ", err)
	}

	//create a new GRPC server
	s := grpc.NewServer()
	//register your server to the newly created server
	calculatorpb.RegisterCalServiceServer(s, &server{})

	//bind the port
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
}
