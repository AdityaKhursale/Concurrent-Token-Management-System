package main

import (
	"flag"
	"fmt"
	"net"

	"proj_2/token"
	"proj_2/utils"

	"google.golang.org/grpc"
)

func main() {
	portPtr := flag.String("port", "50051", "port number to use")
	flag.Parse()

	fmt.Println("\nServer started on port", (*portPtr), "\n")
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", (*portPtr)))
	utils.IsSuccess(err)

	s := token.Server{}
	server := grpc.NewServer()

	token.RegisterTokenServiceServer(server, &s)

	err = server.Serve(ln)
	utils.IsSuccess(err)
}
