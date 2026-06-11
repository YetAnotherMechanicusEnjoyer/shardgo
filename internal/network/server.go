package network

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/YetAnotherMechanicusEnjoyer/shardgo/internal/cache"
	"github.com/YetAnotherMechanicusEnjoyer/shardgo/internal/cluster"
)

type Server struct {
	addr     string
	cache    *cache.Cache
	cluster  *cluster.Manager
	isMaster bool
	nodeAddr string
}

func NewServer(addr string, cache *cache.Cache, cluster *cluster.Manager, isMaster bool) *Server {
	return &Server{
		addr:     addr,
		cache:    cache,
		cluster:  cluster,
		isMaster: isMaster,
		nodeAddr: addr,
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
	if req.Command == "ADD_NODE" && len(req.Args) >= 2 {
		return s.handleClusterCommand(req)
	}

	key := ""
	if len(req.Args) >= 2 {
		key = req.Args[1]
	}

	if key != "" && s.cluster != nil {
		node, ok := s.cluster.GetNodeForKey(key)
		if !ok || node != s.nodeAddr {
			return s.forwardRequest(req, node)
		}
	}

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

	case "STATS":
		size := s.cache.Size()
		return fmt.Sprintf(":%d\r\n", size), nil

	default:
		return "", errors.New("unknown command")
	}
}

func (s *Server) forwardRequest(req *Request, targetNode string) (string, error) {
	if targetNode == "" {
		return "", errors.New("no node responsible for key")
	}

	client := NewClient(targetNode)

	var rawCmd strings.Builder

	rawCmd.WriteString(req.Command)
	for _, arg := range req.Args[1:] {
		rawCmd.WriteString(" ")
		rawCmd.WriteString(arg)
	}

	return client.SendRequest(rawCmd.String())
}

func (s *Server) handleClusterCommand(req *Request) (string, error) {
	if req.Command == "ADD_NODE" && len(req.Args) >= 2 {
		node := req.Args[1]
		s.cluster.AddNode(node)
		return FormatResponse("OK"), nil
	}
	return "", errors.New("unknown cluster command")
}
