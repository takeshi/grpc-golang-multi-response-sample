package main

import (
	"context"
	"google.golang.org/grpc"
	pb "grpc-batch-client/pd/batch.sample"
	"io"
	"log"
)

func main() {
	//sampleなのでwithInsecure
	conn, err := grpc.Dial("127.0.0.1:6565", grpc.WithInsecure())
	if err != nil {
		log.Fatal("client connection error:", err)
	}
	defer conn.Close()
	//pb.BatchExecutorClient()
	client := pb.NewBatchExecutorClient(conn)
	input := &pb.BatchRequest{}
	res, err := client.Execute(context.TODO(), input)

	if err != nil {
		panic(err)
	}

	done := make(chan bool)
	go func() {
		for {
			resp, err := res.Recv()
			if err == io.EOF {
				done <- true //means stream is finished
				return
			}
			if err != nil {
				log.Fatalf("cannot receive %v", err)
			}
			switch resp.Type {
			case "message":
				println("getMessage")
				println(resp.Output)
				break
			case "response":
				println("getReponse")
				println(resp.Output)
				res.CloseSend()
				done <- true //means stream is finished
				break
			}
		}
	}()

	<-done //we will wait until all response is received
	log.Printf("finished")
}
