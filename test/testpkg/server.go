package testpkg

import (
	"net"
)

type Server struct {
	addr string
	port int
	listener net.Listener
	isShuttingDown bool
	sessions sessionInfo
}

type sessionInfo struct {
	clients []net.Conn
}

type clientID int64
type clientList []clientID

func NewServer(addr string, port int) (*Server, error) {
	return &Server{
		addr,
		port,
		nil,
		false,
		sessionInfo{},
	}, nil
}

func (s *Server) Listen() error {
	logger.Printf("listening on port %d", s.port, "\n")

	for !s.isShuttingDown {
		conn, err := s.listener.Accept()
		s.sessions.clients = append(s.sessions.clients, conn)
		
		if err != nil {
			if !s.isShuttingDown {
				logger.Printf("Listener err: %s\n", err.Error())
			}
			return nil
		}
		s.handleClient(conn)
	}

	return nil
}

func (s *Server) handleClient(conn net.Conn) {
	return
}