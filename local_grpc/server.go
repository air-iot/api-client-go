package local_grpc

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/air-iot/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
	"google.golang.org/grpc/status"
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

func (s *Server) Invoke(ctx context.Context, method string, args any, reply any) error {
	sm := method
	if sm != "" && sm[0] == '/' {
		sm = sm[1:]
	}
	pos := strings.LastIndex(sm, "/")
	if pos == -1 {
		errDesc := fmt.Sprintf("malformed method name: %q", method)
		return status.Error(codes.Unimplemented, errDesc)
	}
	service := sm[:pos]
	method1 := sm[pos+1:]
	srv, knownService := s.GetService(service)
	codec := encoding.GetCodecV2(proto.Name)
	if knownService {
		if md, ok := srv.GetMethods(method1); ok {
			df := func(v any) error {
				fmt.Println("codec", codec)
				fmt.Println("args", args)
				d, err := codec.Marshal(args)
				if err != nil {
					return status.Errorf(codes.Internal, "grpc: error marshalling request: %v", err)
				}
				if err := codec.Unmarshal(d, v); err != nil {
					return status.Errorf(codes.Internal, "grpc: error unmarshalling request: %v", err)
				}
				return nil
			}

			appReply, appErr := md.Handler(srv.GetServiceImpl(), ctx, df, nil)
			if appErr != nil {
				appStatus, ok := status.FromError(appErr)
				if !ok {
					// Convert non-status application error to a status error with code
					// Unknown, but handle context errors specifically.
					appStatus = status.FromContextError(appErr)
					appErr = appStatus.Err()
				}
				return appErr
			}

			b, err := codec.Marshal(appReply)
			if err != nil {
				return status.Errorf(codes.Internal, "grpc: error while marshaling: %v", err.Error())
			}

			if err := codec.Unmarshal(b, reply); err != nil {
				return status.Errorf(codes.Internal, "grpc: error while unmarshaling: %v", err.Error())
			}
			return nil
		}
	}
	errDesc := fmt.Sprintf("unknown service %v", service)
	return status.Error(codes.Unimplemented, errDesc)
}
