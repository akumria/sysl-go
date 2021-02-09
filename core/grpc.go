package core

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/anz-bank/sysl-go/log"

	"github.com/anz-bank/sysl-go/config"
	"github.com/anz-bank/sysl-go/handlerinitialiser"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

// Deprecated: prefer GrpcServerManager.
type GrpcManager interface {
	Interceptors() []grpc.UnaryServerInterceptor
	EnabledGrpcHandlers() []handlerinitialiser.GrpcHandlerInitialiser
	GrpcAdminServerConfig() *config.CommonServerConfig
	GrpcPublicServerConfig() *config.CommonServerConfig
}

type GrpcServerManager struct {
	GrpcServerOptions      []grpc.ServerOption
	EnabledGrpcHandlers    []handlerinitialiser.GrpcHandlerInitialiser
	GrpcPublicServerConfig *config.CommonServerConfig
}

func DefaultGrpcServerOptions(ctx context.Context, grpcPublicServerConfig *config.CommonServerConfig) ([]grpc.ServerOption, error) {
	opts, err := config.ExtractGrpcServerOptions(grpcPublicServerConfig)
	if err != nil {
		return nil, err
	}

	logger := log.GetLogger(ctx)
	// Inject the logger into the ctx so we can log when we're serving rpc calls.
	opts = append(opts, grpc.ChainUnaryInterceptor(makeLoggerInterceptor(logger)))

	opts = append(opts, grpc.ChainUnaryInterceptor(TraceidLogInterceptor))
	return opts, nil
}

func newGrpcServerManagerFromGrpcManager(hl GrpcManager) (*GrpcServerManager, error) {
	opts, err := extractGrpcServerOptionsFromGrpcManager(hl)
	if err != nil {
		return nil, err
	}
	return &GrpcServerManager{
		GrpcServerOptions:      opts,
		EnabledGrpcHandlers:    hl.EnabledGrpcHandlers(),
		GrpcPublicServerConfig: hl.GrpcPublicServerConfig(),
	}, nil
}

func extractGrpcServerOptionsFromGrpcManager(hl GrpcManager) ([]grpc.ServerOption, error) {
	opts, err := config.ExtractGrpcServerOptions(hl.GrpcPublicServerConfig())
	if err != nil {
		return nil, err
	}
	opts = append(opts, grpc.ChainUnaryInterceptor(hl.Interceptors()...))
	opts = append(opts, grpc.ChainUnaryInterceptor(TraceidLogInterceptor)) // seems wrong to have this last in chain, but that was old behaviour.
	return opts, nil
}

func configurePublicGrpcServerListener(ctx context.Context, m GrpcServerManager) StoppableServer {
	server := grpc.NewServer(m.GrpcServerOptions...)
	// Not sure if it is possible to register multiple servers
	for _, h := range m.EnabledGrpcHandlers {
		h.RegisterServer(ctx, server)
	}

	return prepareGrpcServerListener(ctx, server, *m.GrpcPublicServerConfig, "gRPC Public server")
}

type grpcServer struct {
	ctx    context.Context
	cfg    config.CommonServerConfig
	server *grpc.Server
	name   string
}

func (s grpcServer) Start() error {
	if s.cfg.TLS != nil {
		log.Infof(s.ctx, "TLS configuration present. Preparing to serve gRPC/HTTPS for address: %s:%d", s.cfg.HostName, s.cfg.Port)
	} else {
		log.Infof(s.ctx, "TLS configuration NOT present. Preparing to serve gRPC/HTTP for address: %s:%d", s.cfg.HostName, s.cfg.Port)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.HostName, s.cfg.Port))
	if err != nil {
		return err
	}
	return s.server.Serve(lis)
}

func (s grpcServer) GracefulStop() error {
	s.server.GracefulStop()
	return nil
}

func (s grpcServer) Stop() error {
	s.server.Stop()
	return nil
}

func (s grpcServer) GetName() string {
	return s.name
}

type logWriterInfo struct {
	logger log.Logger
}

func (lw *logWriterInfo) Write(p []byte) (n int, err error) {
	lw.logger.Info(string(p))
	return len(p), nil
}

type logWriterError struct {
	logger log.Logger
}

func (lw *logWriterError) Write(p []byte) (n int, err error) {
	lw.logger.Error(errors.New(string(p)), "gRPC error")
	return len(p), nil
}

func prepareGrpcServerListener(ctx context.Context, server *grpc.Server, commonConfig config.CommonServerConfig, name string) StoppableServer {
	logger := log.GetLogger(ctx)
	grpclog.SetLoggerV2(
		grpclog.NewLoggerV2(
			&logWriterInfo{logger: logger},
			&logWriterInfo{logger: logger},
			&logWriterError{logger: logger}))

	log.Infof(ctx, "configured gRPC listener for address: %s:%d", commonConfig.HostName, commonConfig.Port)
	return grpcServer{ctx: ctx, cfg: commonConfig, server: server, name: name}
}

func makeLoggerInterceptor(logger log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = log.PutLogger(ctx, logger)
		return handler(ctx, req)
	}
}

func TraceidLogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx = log.WithStr(ctx, "traceid", "traceid")
	return handler(ctx, req)
}
