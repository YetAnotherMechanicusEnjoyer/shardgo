package network

import (
	"bufio"
	"errors"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/YetAnotherMechanicusEnjoyer/shardgo/internal/cache"
)

type Server struct {
	addr  string
	cache *cache.Cache
}

func NewServer(addr string, cache *cache.Cache) *Server {
	return &Server{
		addr:  addr,
		cache: cache,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Printf("ShardGo server listening on %s", s.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		req, err := ParseRequest(reader)
		if err != nil {
			_, _ = conn.Write([]byte(FormatError(err)))
			return
		}

		resp, err := s.processRequest(req)
		if err != nil {
			_, _ = conn.Write([]byte(FormatError(err)))
			continue
		}
		_, _ = conn.Write([]byte(resp))
	}
}

func (s *Server) processRequest(req *Request) (string, error) {
	switch req.Command {
	case "SET":
		if len(req.Args) < 3 {
			return "", errors.New("SET requires key and value")
		}
		key := req.Args[1]
		value := []byte(req.Args[2])
		var ttl time.Duration
		if len(req.Args) >= 4 {
			ttlSecs, err := strconv.Atoi(req.Args[3])
			if err != nil {
				return "", err
			}
			ttl = time.Duration(ttlSecs) * time.Second
		}
		if err := s.cache.Set(key, value, ttl); err != nil {
			return "", err
		}
		return FormatResponse("OK"), nil

	case "GET":
		if len(req.Args) < 2 {
			return "", errors.New("GET requires key")
		}
		value, err := s.cache.Get(req.Args[1])
		if err != nil {
			return "", err
		}
		return FormatBulkString(string(value)), nil

	case "DEL":
		if len(req.Args) < 2 {
			return "", errors.New("DEL requires key")
		}
		if err := s.cache.Delete(req.Args[1]); err != nil {
			return "", err
		}
		return FormatResponse("OK"), nil

	default:
		return "", errors.New("unknown command")
	}
}
