## Project Title

Concurrent Token Management System

## Project Description

Implemented a client-server solution for managing tokens. Server maintains an initially empty (non-persistent) collection of tokens. Clients issue RPC calls to the server to execute create, drop, read-write methods on tokens. The server executes such RPCs and returns an appropriate response to each call.

#### Token
A token is an abstract data type, with the following properties: id, name, domain, and state. Tokens are uniquely identified by their id, which is a string. The name of a token is another string. The domain of a token consists of three uint64 integers: low, mid, and high. The state of a token consists of two unit64 integers: a partial value and a final value, which is defined at the integer x in the range [low, mid) and [low, high), respectively, that minimizes h(name, x) for a hash function h. Hash function used in SHA-256.

#### Supported Operations
- **create(id):** create a token with the given id. Return a success or fail response.
- **drop(id):** to destroy/delete the token with the given id
- **write(id, name, low, high, mid):**
    1. set the properties name, low, mid, and high for the token with the given id. Assume
       uint64 integers low <= mid < high.
    2. compute the partial value of the token as min H(name, x) for x in [low, mid),
       and reset the final value of the token.
   3. return the partial value on success or fail response
- **read(id):**
    1. find min H(name, x) for x in [mid, high)
    2. set the token’s final value to the minimum of the value in step#1 and its partial
       value
    3. return the token’s final value on success or fail response

## Setup Environment

**Go Installation**
Follow: [Download and install Go](https://go.dev/dl/)
Use Version: go1.17.7
**Protocol Buffers Installation**
Follow: [Protocol Buffer Compiler Installation](https://grpc.io/docs/protoc-installation/)
Use Version: 3
**Install gRPC plugins**
```sh
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
```
**Update PATH for protoc**
```sh
$ export PATH="$PATH:$(go env GOPATH)/bin"
```

## Run

Server:
```sh
# go to project directory
cd <project_directory>
e.g. $ cd /Users/aditya/Documents/Courses/AOS/CMSC621_project2

# start server
go run server.go -port <port_number>
e.g. $ go run server.go -port 50051
```

Client:
```sh
# go to project directory
cd <project_directroy>
e.g. $ cd /Users/aditya/Documents/Courses/AOS/CMSC621_project2

# create request
go run client.go -create -id <id_num> -host <host_addreess> -port <port_number>
e.g. $ go run client.go -create -id 1 -host localhost -port 50051

# write request
go run client.go -write -id <id_num> -name <token_name> -low <low> -mid <mid> -high <high> -host <host_address> -port <port_number>
e.g. $ go run client.go -write -id 1 -name abcd -low 1 -mid 5 -high 10 -host localhost -port 50051

# read request
go run client.go -read -id <id_num> -host <host_address> -port <port_number>
e.g. $ go run client.go -read -id 1 -host localhost -port 50051

# drop request
go run client.go -drop -id <id_num> -host <host_address> -port <port_number>
e.g. $ go run client.go -drop -id 1 -host localhost -port 50051
```

To run the demo:
```sh
# go to project directory
cd <project_directroy>
e.g. $ cd /Users/aditya/Documents/Courses/AOS/CMSC621_project2

# set executable permission for demo script
$ chmod +x ./demo.sh

# execute demo script
./demo.sh <port number>
e.g. $ ./demo.sh 50051
```

To close the server press Ctrl+C

## Project Files and Directories

- **server.go:** Code for the server operations
- **client.go:** Code for the client operations
- **token:** Directory containing code related to token management like proto definitions and logic for each operation that can be performed on tokens
- **utils:** Directory containing utilities and helper functions
- **go.mod:** Root dependecny managment
- **go.sum:** Checksum for dependencies
- **demo.sh:** Shell script demonstrating demo of the project
- **analysis_helper.sh:** Lists commands I generally used to analyze my outputs and it is sometimes to difficult to navigate stdout and stderr through large pile of text
- **output:** Directory which stores redirected logs from server and client ran via demo script
- **demo_screenshots:** Screenshots of the demo I ran

## Code Description - What Did I do/Assumptions Made/Deviations

I believe the code itself is very readable. I tried to note few more things below:

- One of the major deviation I took from what is asked in the project is to use min values instead of argmin for read and write operations
- Reason being argmin will always result in same partial and final values. I discussed this already with the TA
- While supporting concurrency, my code supports following type of operations
-- Operations with different id - Parallel execution
-- Read operations with same id - Parallel execution (Although, there is write involved in the last step of read operation, it does not corrupt data as low, mid, high remains the same)
-- Any other combination of operations with same id - Serial execution
-- I used 0 as default values for partial and final values. Sign-offed by TA

## Demo
![demo_sh](https://drive.google.com/uc?export=view&id=1injx0tmpBF0gL36Q6F7nRI5hD0zsOwXH)
Above screenshot is of running the demo.sh, Sequence of server and clients launched is logged above

![requests_seq](https://drive.google.com/uc?export=view&id=16uJn1CEU9H8X1zm6kDSp_RMA0mEdiuEN)
In the above order requests are raised towards the server

![processed_seq](https://drive.google.com/uc?export=view&id=1NKjCBiCbfcuhTszBSVmnUPYLKnMgwlC6)
Server processed non conflicting requests in parallel and conflicting in serial, and correctly executed everything without error. Random processing from above screenshot clearly demonstrates the concurrency. For e.g. observe that inexpensive opeartions like create or drop are performed quite fast, and their processing was done before read and write completes with different ids even though they were launched before the read and write. Also, notice that read waits for the write if they have same ids for e.g. request 3 and 4, but parallely processed requests 5 and 6 which are of different ids.

![tokenstore_states](https://drive.google.com/uc?export=view&id=1Suqn2iDJC7GsCtzfYiZmeGZxPRWg_DiN)
Above screenshot demostrates the tokenstore after every request is processed. All 4 requests for each token are launched one after the other, However, multiple tokens are their the store at the same time which concludes the concurrency. Since, these operations are not too expensive there are only 2 tokens at the same time, if we lauch requests in thousands with expensive operations, I am sure that tokenstore will have more tokens in the store. I have also observed 3 tokens with the same script, however, since this is concurrent there can not be a predictable behavior.

![all_clients_op](https://drive.google.com/uc?export=view&id=138e7mq5ZHVjTZO2g0yBWqo0kz8SI__06)
I realized that the output of clients are stored in different files, and it is logical. However, it might be cumbersome to read all those files. Hence I printed all outputs of clients toghether. This is to demonstrate that clients are receiving correct responses. Don't get mislead by the drop responses coming before read that is because of the bash redirection delay. Observe server side output for this confusion from screenshot 3

## References

- https://forum.golangbridge.org/t/cannot-import-package-variables-to-main-package/21193/2
- https://github.com/evilsocket/opensnitch/issues/373#issuecomment-803663343
- https://github.com/grpc/grpc-go/issues/3794#issuecomment-720599532
- https://stackoverflow.com/questions/15178088/create-global-map-variables
- https://tutorialedge.net/golang/go-grpc-beginners-tutorial/
- https://yourbasic.org/golang/errors-explained/
- https://go.dev/blog/maps
- https://www.geeksforgeeks.org/math-inf-function-in-golang-with-examples/
- https://learnandlearn.com/golang-programming/golang-reference/golang-find-the-minimum-value-min-function-examples-explanation
- https://yourbasic.org/golang/multiline-string/
