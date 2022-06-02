package main

import (
	"Hung/Hung-Test/Test-gRPC/calculator"
	"context"
	"google.golang.org/grpc"
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

	client := calculator.NewCalculatorServiceClient(cc)
	log.Printf("service client %v", client)

	//callSum(client, 6, 9)
	//callAverage(client)
	callFindMax(client)
}

func callSum(cli calculator.CalculatorServiceClient, num1, num2 int) {
	log.Printf("callSum api is called")
	req := &calculator.SumRequest{
		Num1: int32(num1),
		Num2: int32(num2),
	}
	resp, err := cli.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("err %v", err)
	}
	log.Printf("resp %v", resp.GetResult())
}

func callAverage(cli calculator.CalculatorServiceClient) {
	log.Printf("callAverage api is called")
	stream, err := cli.Average(context.Background())
	if err != nil {
		log.Fatalf("err %v", err)
	}
	listReq := []calculator.AverageRequest{
		calculator.AverageRequest{
			Num: 21.11,
		},
		calculator.AverageRequest{
			Num: 15.22,
		},
		calculator.AverageRequest{
			Num: 7.9,
		},
		calculator.AverageRequest{
			Num: 9.3,
		},
		calculator.AverageRequest{
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

func callFindMax(cli calculator.CalculatorServiceClient) {
	log.Printf("callFindMax api is called")
	stream, err := cli.FindMax(context.Background())
	if err != nil {
		log.Fatalf("call FindMax err %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		//send requests
		listReq := []calculator.FindMaxRequest{
			calculator.FindMaxRequest{
				Num: 2,
			},
			calculator.FindMaxRequest{
				Num: 15,
			},
			calculator.FindMaxRequest{
				Num: 7,
			},
			calculator.FindMaxRequest{
				Num: 21,
			},
			calculator.FindMaxRequest{
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
