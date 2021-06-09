package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/egnimos/grpc-practice/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	//connect to the IP Address
	cc, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalServiceClient(cc)

	//do unary call
	// doUnaryCalls(c)
	//server streaming
	doServerStreaming(c)

}

func doUnaryCalls(c calculatorpb.CalServiceClient) {
	fmt.Println("Startign the Unary RPC calls...")
	req := &calculatorpb.CalRequest{
		Cal: &calculatorpb.Cal{
			Intput_1: 3,
			Intput_2: 10,
		},
	}

	//get the response
	res, err := c.Calculate(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}

	log.Printf("response is: %v", res.Output)
}

func doServerStreaming(c calculatorpb.CalServiceClient) {
	fmt.Println("starting to do a server streaming RPC...")

	req := &calculatorpb.CalManyTimesRequest{
		Cal: &calculatorpb.Cal{
			Intput_1: 3,
			Intput_2: 10,
		},
	}

	//get the response streaming
	resStream, err := c.CalculateManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling the calculateManyTimes RPC: %v", err)
	}

	//run the loop
	for {
		msg, err := resStream.Recv()
		// if the error is end of file then break the loop or close the loop
		if err == io.EOF {
			//we've reach the end of stream
			break
		}
		//check other error
		if err != nil {
			log.Fatalf("error while reading the stream: %v", err)
		}
		//print the actual message
		log.Println("Response from the GreetManyTimes: ", msg.GetOutput())
	}
}
