package main

import (
	"context"
	"flag"
	"fmt"

	"proj_2/token"
	"proj_2/utils"

	"google.golang.org/grpc"
)

func main() {
	portPtr := flag.String("port", "50051",
		"port number where server is running")
	hostPtr := flag.String("host", "localhost",
		"host where server is running")

	createPtr := flag.Bool("create", false, "set to create token")
	dropPtr := flag.Bool("drop", false, "set to drop token")
	writePtr := flag.Bool("write", false, "set to write token")
	readPtr := flag.Bool("read", false, "set to read ptr")

	idPtr := flag.String("id", "undefined", "id of the token")
	namePtr := flag.String("name", "undefined", "name of the token")
	lowPtr := flag.Uint64("low", 1, "low value of the domain of token")
	midPtr := flag.Uint64("mid", 1, "mid value of the domain of token")
	highPtr := flag.Uint64("high", 1, "high value of the domain of token")
	flag.Parse()

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", *hostPtr, *portPtr),
		grpc.WithInsecure())
	utils.IsSuccess(err)

	defer conn.Close()

	c := token.NewTokenServiceClient(conn)

	req := token.Request{}
	req.Domain = &token.Request_Domain{}
	req.TokenState = &token.Request_State{}

	resp := &token.Response{}
	if *createPtr {
		req.Id = *idPtr
		resp, err = c.Create(context.Background(), &req)
	} else if *dropPtr {
		req.Id = *idPtr
		resp, err = c.Drop(context.Background(), &req)
	} else if *writePtr {
		req.Id = *idPtr
		req.Name = *namePtr
		req.Domain.Low = *lowPtr
		req.Domain.Mid = *midPtr
		req.Domain.High = *highPtr
		resp, err = c.Write(context.Background(), &req)
	} else if *readPtr {
		req.Id = *idPtr
		resp, err = c.Read(context.Background(), &req)
	}
	fmt.Println("Server Response:", (resp.GetBody()))
	utils.IsSuccess(err)

}
