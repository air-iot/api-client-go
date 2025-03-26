package local_grpc

import (
	"reflect"
	"sync"

	"github.com/air-iot/logger"
	"google.golang.org/grpc"
)

type ServiceInfo struct {
	// Contains the implementation for the methods in this service.
	serviceImpl any
	methods     map[string]*grpc.MethodDesc
	streams     map[string]*grpc.StreamDesc
	mdata       any
}

func (s *ServiceInfo) GetMethods(md string) (*grpc.MethodDesc, bool) {
	md1, ok := s.methods[md]
	return md1, ok
}

func (s *ServiceInfo) GetStreams(md string) (*grpc.StreamDesc, bool) {
	md1, ok := s.streams[md]
	return md1, ok
}

func (s *ServiceInfo) GetServiceImpl() any {
	return s.serviceImpl
}

type Server struct {
	mu       sync.Mutex // guards following
	services map[string]*ServiceInfo
}

func NewServer() *Server {
	return &Server{
		services: make(map[string]*ServiceInfo),
	}
}

func (s *Server) GetService(serviceName string) (*ServiceInfo, bool) {
	md1, ok := s.services[serviceName]
	return md1, ok
}

func (s *Server) RegisterService(sd *grpc.ServiceDesc, ss any) {
	if ss != nil {
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) {
			logger.Fatalf("grpc: Server.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
		}
	}
	s.register(sd, ss)
}

func (s *Server) register(sd *grpc.ServiceDesc, ss any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.services[sd.ServiceName]; ok {
		logger.Fatalf("grpc: Server.RegisterService found duplicate service registration for %q", sd.ServiceName)
	}
	info := &ServiceInfo{
		serviceImpl: ss,
		methods:     make(map[string]*grpc.MethodDesc),
		streams:     make(map[string]*grpc.StreamDesc),
		mdata:       sd.Metadata,
	}
	for i := range sd.Methods {
		d := &sd.Methods[i]
		info.methods[d.MethodName] = d
	}
	for i := range sd.Streams {
		d := &sd.Streams[i]
		info.streams[d.StreamName] = d
	}
	s.services[sd.ServiceName] = info
}
