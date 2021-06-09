package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/egnimos/grpc-practice/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("hello I'm a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	//close the client
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	//do unary call
	// doUnaryCalls(c)
	//do server streaming
	doServerStreaming(c)
}

//doUnaryCalls
func doUnaryCalls(c greetpb.GreetServiceClient) {
	fmt.Println("starting the Unary RPC call...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Niteesh Kumar",
			LastName:  "Dubey",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

//doServerStreaming
func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do a server streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Niteesh Kumar",
			LastName:  "Dubey",
		},
	}

	//get the response streaming
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling the greetManyTimes RPC: %v", err)
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
		log.Println("Response from the GreetManyTimes: ", msg.GetResult())
	}
}
