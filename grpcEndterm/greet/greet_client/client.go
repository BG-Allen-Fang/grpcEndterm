package main

import (
	"com.grpc/greet/greetpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	//doManyTimesFromServer for 1 task
	//doLongGreet 			for 2 task

	doLongGreet(c)
	//doManyTimesFromServer(c)
}

func doManyTimesFromServer(c greetpb.GreetServiceClient) {
	ctx := context.Background()
	req := &greetpb.GreetManyTimesRequest{Greeting: &greetpb.Greeting{
		Number: 120,
	}}

	stream, err := c.GreetManyTimes(ctx, req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC %v", err)
	}
	defer stream.CloseSend()

LOOP:
	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break LOOP
			}
			log.Fatalf("error while reciving from GreetManyTimes RPC %v", err)
		}
		log.Printf("response from GreetManyTimes:%v \n", res.GetResult())
	}

}

func doLongGreet(c greetpb.GreetServiceClient) {

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				Number: 1,
			},
		},
		{
			Greeting: &greetpb.Greeting{
				Number: 2,
			},
		},
		{
			Greeting: &greetpb.Greeting{
				Number: 3,
			},
		},
		{
			Greeting: &greetpb.Greeting{
				Number: 4,
			},
		},
	}

	ctx := context.Background()
	stream, err := c.LongGreet(ctx)
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res)
}
