package core

import (
	"context"
	"fmt"
	"net/http"

	"github.com/anz-bank/sysl-go/log"

	"github.com/anz-bank/sysl-go/common"
	"github.com/anz-bank/sysl-go/config"
	"github.com/anz-bank/sysl-go/core/authrules"
	"github.com/anz-bank/sysl-go/jwtauth"
	"github.com/go-chi/chi"
	"google.golang.org/grpc"
)

// RestGenCallback is used by `sysl-go` to call hand-crafted code.
type RestGenCallback interface {
	// AddMiddleware allows hand-crafted code to add middleware to the router
	AddMiddleware(ctx context.Context, r chi.Router)
	// BasePath allows hand-crafted code to set the base path for the Router
	BasePath() string
	// Config returns a structure representing the server config
	// This is returned from the status endpoint
	Config() interface{}
	// MapError maps an error to an HTTPError in instances where custom error mapping is required. Return nil to perform default error mapping; defined as:
	// 1. CustomError.HTTPError if the original error is a CustomError, otherwise
	// 2. common.MapError
	MapError(ctx context.Context, err error) *common.HTTPError
	// DownstreamTimeoutContext add the desired timeout duration to the context for downstreams
	// A separate service timeout (usually greater than the downstream) should also be in
	// place to automatically respond to callers
	DownstreamTimeoutContext(ctx context.Context) (context.Context, context.CancelFunc)
}

// GrpcGenCallback is currently a subset of RestGenCallback so is defined separately for convenience.
type GrpcGenCallback interface {
	DownstreamTimeoutContext(ctx context.Context) (context.Context, context.CancelFunc)
}

// Hooks can be used to customise the behaviour of an autogenerated sysl-go service.
type Hooks struct {

	// Logger returns the common.Logger instance to set use within Sysl-go.
	// By default, if this Logger hook is not set then an instance of the pkg logger is used.
	// This hook can also be used to define a custom logger.
	// For more information about logging see log/README.md within this project.
	// Note: The returned logger is guaranteed to have the log level from the external configuration
	// file (library: log: level) set against it.
	Logger func() log.Logger

	// MapError maps an error to an HTTPError in instances where custom error mapping is required.
	// Return nil to perform default error mapping; defined as:
	// 1. CustomError.HTTPError if the original error is a CustomError, otherwise
	// 2. common.MapError
	// By default, if this MapError hook is not customised, the default error mapping will be used.
	MapError func(ctx context.Context, err error) *common.HTTPError

	// AdditionalGrpcDialOptions can be used to append to the default grpc.DialOption configuration used by
	// an autogenerated service when it calls grpc.Dial when using a grpc.Client to connect to a gRPC server.
	// If given, AdditionalGrpcDialOptions will be appended to the list of default options created by
	// DefaultGrpcDialOptions(CommonGRPCDownstreamData).
	//
	// Use AdditionalGrpcDialOptions if you need both default and custom options. Be careful that you do
	// not specify any options that clash with the default options.
	//
	// If you need to completely override the default options, use OverrideGrpcDialOptions.
	// It is an error to set both AdditionalGrpcDialOptions and OverrideGrpcDialOptions.
	AdditionalGrpcDialOptions []grpc.DialOption

	// OverrideGrpcDialOptions can be used to override the default grpc.DialOption configuration used by an
	// an autogenerated service when it calls grpc.Dial when using a grpc.Client to connect to a gRPC server.
	//
	// The serviceName parameter will be filled with the name of the target service that we
	// are about to call grpc.Dial to connect to -- a function implementing this hook can use the
	// serviceName to customise different dial options for different targets.
	//
	// Prefer to use AdditionalGrpcDialOptions instead of OverrideGrpcDialOptions if you only need
	// to append to the default grpc.DialOption configuration instead of overriding it completely.
	//
	// It is an error to set both AdditionalGrpcDialOptions and OverrideGrpcDialOptions.
	OverrideGrpcDialOptions func(serviceName string, cfg *config.CommonGRPCDownstreamData) ([]grpc.DialOption, error)

	// AdditionalGrpcServerOptions can be used to append to the default grpc.ServerOption configuration used by
	// an autogenerated service when it creates a gRPC server. If given, AdditionalGrpcServerOptions will be
	// appended to the list of default options created by DefaultGrpcServerOptions(context.Context, CommonServerConfig).
	//
	// Use AdditionalGrpcServerOptions if you need both default and custom options. Be careful that you do
	// not specify any options that clash with the default options.
	//
	// If you need to completely override the default options, use OverrideGrpcServerOptions.
	// It is an error to set both AdditionalGrpcServerOptions and OverrideGrpcServerOptions.
	AdditionalGrpcServerOptions []grpc.ServerOption

	// OverrideGrpcServerOptions can be used to override the default grpc.ServerOption configuration used by an
	// autogenerated service when it creates a gRPC server.
	//
	// Prefer to use AdditionalGrpcServerOptions instead of OverrideGrpcServerOptions if you only need
	// to append to the default grpc.ServerOption configuration instead of overriding it completely.
	//
	// It is an error to set both AdditionalGrpcServerOptions and OverrideGrpcServerOptions.
	OverrideGrpcServerOptions func(ctx context.Context, grpcPublicServerConfig *config.CommonServerConfig) ([]grpc.ServerOption, error)

	// OverrideMakeJWTClaimsBasedAuthorizationRule can be used to customise how authorization rule
	// expressions are evaluated and used to decide if JWT claims are authorised. By default, if this
	// hook is nil, then authrules.MakeDefaultJWTClaimsBasedAuthorizationRule is used.
	OverrideMakeJWTClaimsBasedAuthorizationRule func(authorizationRuleExpression string) (authrules.JWTClaimsBasedAuthorizationRule, error)

	// AddHTTPMiddleware can be used to install additional HTTP middleware into the chi.Router
	// used to serve all (non-admin) HTTP endpoints. By default, sysl-go installs a number of
	// HTTP middleware -- refer to prepareMiddleware inside sysl-go/core. This hook can only
	// be used to add middleware, not override any of the default middleware.
	AddHTTPMiddleware func(ctx context.Context, r chi.Router)

	// AddAdminHTTPMiddleware can be used to install additional HTTP middleware into the chi.Router
	// used to serve the admin HTTP endpoints. See AddHTTPMiddleware for further details.
	AddAdminHTTPMiddleware func(ctx context.Context, r chi.Router)

	// DownstreamRoundTripper can be used to install additional HTTP RoundTrippers to the downstream clients
	DownstreamRoundTripper func(serviceName string, serviceURL string, original http.RoundTripper) http.RoundTripper

	// ValidateConfig can be used to validate (or override) values in the config.
	ValidateConfig func(ctx context.Context, cfg *config.DefaultConfig) error
}

func ResolveGrpcDialOptions(ctx context.Context, serviceName string, h *Hooks, grpcDownstreamConfig *config.CommonGRPCDownstreamData) ([]grpc.DialOption, error) {
	switch {
	case len(h.AdditionalGrpcDialOptions) > 0 && h.OverrideGrpcDialOptions != nil:
		return nil, fmt.Errorf("Hooks.AdditionalGrpcDialOptions and Hooks.OverrideGrpcDialOptions cannot both be set")
	case h.OverrideGrpcDialOptions != nil:
		return h.OverrideGrpcDialOptions(serviceName, grpcDownstreamConfig)
	default:
		opts, err := config.DefaultGrpcDialOptions(ctx, grpcDownstreamConfig)
		if err != nil {
			return nil, err
		}
		opts = append(opts, h.AdditionalGrpcDialOptions...)
		return opts, nil
	}
}

func ResolveGrpcServerOptions(ctx context.Context, h *Hooks, grpcPublicServerConfig *config.CommonServerConfig) ([]grpc.ServerOption, error) {
	switch {
	case len(h.AdditionalGrpcServerOptions) > 0 && h.OverrideGrpcServerOptions != nil:
		return nil, fmt.Errorf("Hooks.AdditionalGrpcServerOptions and Hooks.OverrideGrpcServerOptions cannot both be set")
	case h.OverrideGrpcServerOptions != nil:
		return h.OverrideGrpcServerOptions(ctx, grpcPublicServerConfig)
	default:
		opts, err := DefaultGrpcServerOptions(ctx, grpcPublicServerConfig)
		if err != nil {
			return nil, err
		}
		opts = append(opts, h.AdditionalGrpcServerOptions...)
		return opts, nil
	}
}

func ResolveGRPCAuthorizationRule(ctx context.Context, h *Hooks, endpointName string, authRuleExpression string) (authrules.Rule, error) {
	return resolveAuthorizationRule(ctx, h, endpointName, authRuleExpression, authrules.MakeGRPCJWTAuthorizationRule)
}

func ResolveRESTAuthorizationRule(ctx context.Context, h *Hooks, endpointName string, authRuleExpression string) (authrules.Rule, error) {
	return resolveAuthorizationRule(ctx, h, endpointName, authRuleExpression, authrules.MakeRESTJWTAuthorizationRule)
}

func resolveAuthorizationRule(ctx context.Context, h *Hooks, endpointName string, authRuleExpression string, ruleFactory func(authRule authrules.JWTClaimsBasedAuthorizationRule, authenticator jwtauth.Authenticator) (authrules.Rule, error)) (authrules.Rule, error) {
	cfg := config.GetDefaultConfig(ctx)
	if cfg.Development != nil && cfg.Development.DisableAllAuthorizationRules {
		log.Info(ctx, "warning: development.disableAllAuthorizationRules is set, all authorization rules are disabled, this is insecure and should not be used in production.")
		return authrules.InsecureAlwaysGrantAccess, nil
	}
	var claimsBasedAuthRuleFactory func(authorizationRuleExpression string) (authrules.JWTClaimsBasedAuthorizationRule, error)
	switch {
	case h.OverrideMakeJWTClaimsBasedAuthorizationRule != nil:
		claimsBasedAuthRuleFactory = h.OverrideMakeJWTClaimsBasedAuthorizationRule
	default:
		claimsBasedAuthRuleFactory = authrules.MakeDefaultJWTClaimsBasedAuthorizationRule
	}
	claimsBasedAuthRule, err := claimsBasedAuthRuleFactory(authRuleExpression)
	if err != nil {
		return nil, err
	}

	// TODO(fletcher) inject custom http client instrumented with monitoring
	httpClient, err := config.DefaultHTTPClient(ctx, nil)
	if err != nil {
		return nil, err
	}
	httpClientFactory := func(_ string) *http.Client {
		return httpClient
	}

	// Note: this will start a new jwtauth.Authenticator with its own cache & threads running for each of our service's endpoints, we usually want a shared one.
	if cfg == nil || cfg.Library.Authentication == nil || cfg.Library.Authentication.JWTAuth == nil {
		return nil, fmt.Errorf("method/endpoint %s requires a JWT-based authorization rule, but there is no config for library.authentication.jwtauth", endpointName)
	}
	authenticator, err := jwtauth.AuthFromConfig(ctx, cfg.Library.Authentication.JWTAuth, httpClientFactory)
	if err != nil {
		return nil, err
	}
	return ruleFactory(claimsBasedAuthRule, authenticator)
}
