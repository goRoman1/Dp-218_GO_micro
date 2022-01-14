package httpserver

import (
	"Dp-218_GO_micro/protos"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net"
	"net/http"
	"time"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultIdleTimeout     = 30 * time.Second
	defaultShutdownTimeout = 3 * time.Second
	defaultAddr            = ":8080"
)

//Client is a struct a client who connects to the "scooter-run" page.
type Client struct {
	w    io.Writer
	done chan struct{}
}

//Server is a struct of the http-server which has a channel for gRPC connection.
type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
	client          map[int]*Client
	taken           map[int]bool
	codes           map[int]int
	in              chan *protos.ClientMessage
	*protos.UnimplementedScooterServiceServer
}

type Option func(*Server)

//New creates and starts the http-server
func New(handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
		Addr:         defaultAddr,
	}

	server := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
		client:          make(map[int]*Client),
		taken:           make(map[int]bool),
		codes:           make(map[int]int),
		in:              make(chan *protos.ClientMessage),
	}

	for _, opt := range opts {
		opt(server)
	}

	server.run()

	return server
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}

func Port(port string) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}

//ScooterHandler is a special handler which adds a new stream client to the server.
func (s *Server) ScooterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new client connected")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	client := &Client{
		w:    w,
		done: make(chan struct{}),
	}
	s.AddClient(client)

	<-client.done
	fmt.Println("connection closed")
}

//AddClient is a Server's function for adding attached Client.
func (s *Server) AddClient(c *Client) {
	s.client[1] = c
}

//Register is a function for implementing gRPC-service.
func (s *Server) Register(msg *protos.ClientRequest, stream protos.ScooterService_RegisterServer) error {
	return nil
}

//Receive is the function which receive a message from the gRPC stream and direct it to the Server's 'in' channel.
func (s *Server) Receive(stream protos.ScooterService_ReceiveServer) error {
	var err error

	for {
		msg, err := stream.Recv()
		if err != nil {
			fmt.Println(err)
			err = status.Errorf(codes.Internal, "unexpected error %v", err)
			break
		}

		s.in <- msg

	}

	return err
}

//run runs the Server and wait for messages into the channel. Then encode them and print to the console.
func (s *Server) run() {
	go func() {
		for {
			select {
			case msg := <-s.in:
				var buf bytes.Buffer
				json.NewEncoder(&buf).Encode(msg)

				for _, client := range s.client {

					go func(c *Client) {
						if _, err := fmt.Fprintf(c.w, "data: %v\n\n", buf.String()); err != nil {
							fmt.Println(err)
							delete(s.client, 1)
							close(c.done)
							return
						}

						if f, ok := c.w.(http.Flusher); ok {
							f.Flush()
						}
						fmt.Printf("data: %v\n", buf.String())
					}(client)
				}
			}
		}
	}()
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

func IdleTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.IdleTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
