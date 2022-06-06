package main

import (
	calculator2 "Hung/Hung-Test/Test-gRPC/CallCalculator/calculator"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
	"net"
)

type server struct {
}

func (s *server) Sum(ctx context.Context, req *calculator2.SumRequest) (*calculator2.SumResponse, error) {
	fmt.Println("Sum function is called")
	resp := &calculator2.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}
	return resp, nil
}

func (s *server) Average(stream calculator2.CalculatorService_AverageServer) error {
	log.Println("Average function is called")
	var total float32
	var count int
	for {
		req, err := stream.Recv()
		//client send finish request
		if err == io.EOF {
			//tinh trung binh cong và trả về kq cho client
			resp := &calculator2.AverageResponse{
				Result: total / float32(count),
			}
			return stream.SendAndClose(resp)
		}
		if err != nil {
			//log.Fatalf tuong duong voi in va exit
			log.Fatalf("send Average err %v", err)
		}
		log.Println("req.GetNum(): ", req.GetNum())
		total += req.GetNum()
		count++
	}
}

func (s *server) FindMax(stream calculator2.CalculatorService_FindMaxServer) error {
	log.Println("FindMax function is called")
	var max int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("client send finish request EOF")
			log.Println("~~~~~")
			return nil
		}
		if err != nil {
			log.Fatalf("err %v", err)
			return err
		}
		num := req.GetNum()
		log.Printf("num %v", num)
		if num > max {
			max = num
		}
		resp := &calculator2.FindMaxResponse{
			Max: max,
		}
		err = stream.Send(resp)
		if err != nil {
			log.Fatalf("send max err %v", err)
			return err
		}
	}
}

func (s *server) Square(ctx context.Context, req *calculator2.SquareRequest) (*calculator2.SquareResponse, error) {
	fmt.Println("Square function is called")
	num := req.GetNum()
	if num < 0 {
		log.Printf("req num < 0, num %v return InvalidArgument", num)
		return nil, status.Errorf(codes.InvalidArgument, "num must be non-negative")
	}
	resp := &calculator2.SquareResponse{
		SquareRoot: math.Sqrt(float64(num)),
	}
	return resp, nil
}

func main() {
	list, err := net.Listen("tcp", "0.0.0.0:2204")
	if err != nil {
		log.Fatalf("err %v", err)
	}
	s := grpc.NewServer()
	calculator2.RegisterCalculatorServiceServer(s, &server{})
	fmt.Println("server is running...")
	err = s.Serve(list)
	if err != nil {
		log.Fatalf("err %v", err)
	}
}
