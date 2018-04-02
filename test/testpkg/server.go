package testpkg

type Server struct {
	addr string
	port int
}

type clientID int64
type clientList []clientID

func NewServer(addr string, port int) (*Server, error) {
	return &Server{
		addr,
		port,
	}, nil
}

func (s *Server) Listen() error {
	// logger.Printf("listening on port %d", s.port, "\n")
	return nil
}