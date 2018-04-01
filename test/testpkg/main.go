package testpkg

type Server struct {
	addr string
	port int
}

func NewServer(addr string, port int) *Server {
	return &Server{
		addr,
		port,
	}
}

func (s *Server) Listen() {
	
}