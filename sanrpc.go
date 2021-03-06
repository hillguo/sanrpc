package sanrpc

import (
	"github.com/hillguo/sanrpc/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	services map[string]service.Service
}

func (s *Server) AddService(name string, svr service.Service) {
	if s.services == nil {
		s.services = make(map[string]service.Service)
	}
	s.services[name] =svr
}


func (s *Server) Serve() error{
	for _,svr := range s.services{
		go svr.Serve()
	}

	ch := make(chan os.Signal)
	signal.Notify(ch,syscall.SIGINT, syscall.SIGTERM,syscall.SIGSEGV)
	<-ch
	return nil
}

func NewServer() *Server{
	s := &Server{
		services:make(map[string]service.Service),
	}
	return s
}

func (s *Server) Register(serviceDesc interface{}) error {
	for _,s := range s.services {
		err := s.Register(serviceDesc)
		if err != nil {
			log.Fatalln("Register serviceDesc error: ", err)
		}
	}
	return nil
}