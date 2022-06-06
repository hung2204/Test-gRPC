package main

import (
	calculator2 "Hung/Hung-Test/Test-gRPC/CallCalculator/calculator"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

func main() {
	cc, err := grpc.Dial("localhost:2204", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("err %v", err)
	}
	defer func(cc *grpc.ClientConn) {
		err := cc.Close()
		if err != nil {
			log.Fatalf("err %v", err)
		}
	}(cc)

	client := calculator2.NewCalculatorServiceClient(cc)
	log.Printf("service client %v", client)

	//callSum(client, 6, 9)
	//callAverage(client)
	//callFindMax(client)
	callSquareRoot(client, 9)
}

func callSum(cli calculator2.CalculatorServiceClient, num1, num2 int) {
	log.Printf("callSum api is called")
	req := &calculator2.SumRequest{
		Num1: int32(num1),
		Num2: int32(num2),
	}
	resp, err := cli.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("err %v", err)
	}
	log.Printf("resp %v", resp.GetResult())
}

func callAverage(cli calculator2.CalculatorServiceClient) {
	log.Printf("callAverage api is called")
	stream, err := cli.Average(context.Background())
	if err != nil {
		log.Fatalf("err %v", err)
	}
	listReq := []calculator2.AverageRequest{
		calculator2.AverageRequest{
			Num: 21.11,
		},
		calculator2.AverageRequest{
			Num: 15.22,
		},
		calculator2.AverageRequest{
			Num: 7.9,
		},
		calculator2.AverageRequest{
			Num: 9.3,
		},
		calculator2.AverageRequest{
			Num: 12.5,
		},
	}
	for _, req := range listReq {
		err := stream.Send(&req)
		if err != nil {
			log.Fatalf("err %v", err)
		}
		time.Sleep(1 * time.Second)
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("err %v", err)
	}
	log.Println("resp.GetResult(): ", resp.GetResult())
}

func callFindMax(cli calculator2.CalculatorServiceClient) {
	log.Printf("callFindMax api is called")
	stream, err := cli.FindMax(context.Background())
	if err != nil {
		log.Fatalf("call FindMax err %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		//send requests
		listReq := []calculator2.FindMaxRequest{
			calculator2.FindMaxRequest{
				Num: 2,
			},
			calculator2.FindMaxRequest{
				Num: 15,
			},
			calculator2.FindMaxRequest{
				Num: 7,
			},
			calculator2.FindMaxRequest{
				Num: 21,
			},
			calculator2.FindMaxRequest{
				Num: 12,
			},
		}
		for _, req := range listReq {
			err := stream.Send(&req)
			if err != nil {
				log.Fatalf("err %v", err)
				break
			}
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Println("ending...")
				break
			}
			if err != nil {
				log.Fatalf("err %v", err)
				break
			}
			log.Printf("resp %v", resp.GetMax())
		}
		close(waitc)
	}()

	<-waitc
}

func callSquareRoot(cli calculator2.CalculatorServiceClient, num int32) {
	log.Printf("callSquareRoot api is called")
	req := &calculator2.SquareRequest{
		Num: num,
	}
	resp, err := cli.Square(context.Background(), req)
	if err != nil {
		log.Printf("call square root API err %v", err)
		// lay errStatus roi tra ve
		if errStatus, ok := status.FromError(err); ok {
			log.Println("err message: %v\n", errStatus.Message())
			log.Println("err code: %v\n", errStatus.Code())
			if errStatus.Code() == codes.InvalidArgument {
				log.Printf("Invalid argument: %v\n", num)
				return
			}
		}
	}
	log.Printf("square resp %v", resp.GetSquareRoot())
}
