package main

import (
	"com.grpc/greet/greetpb"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (s *Server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with a streaming request\n")
	var res float32
	var cnt float32
	cnt = 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			var r float32
			r = res / cnt
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: r,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		num := req.Greeting.GetNumber()
		res += float32(num)
		cnt++
	}
}

func (s *Server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	number := int(req.GetGreeting().GetNumber())
	fmt.Printf("GreetManyTimes function was invoked with %v \n", req)
	var str string
	for i := 2; number > i; i++ {
		for number%i == 0 {
			a := strconv.Itoa(int(i))
			str = str + a + ", "
			number = number / i
		}
	}
	if number > 2 {
		num := strconv.Itoa(number)
		str = str + num
	}
	newStr := strings.Split(str, ", ")
	for i := 0; i < len(newStr); i++ {
		res := &greetpb.GreetManyTimesResponse{Result: fmt.Sprintf("%d) prime multipliers %v\n", i, newStr[i])}
		if err := stream.Send(res); err != nil {
			log.Fatalf("error while sending greet many times responses: %v", err.Error())
		}
		time.Sleep(time.Second)
	}

	return nil
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &Server{})
	log.Println("Server is running on port:50051")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}
}
