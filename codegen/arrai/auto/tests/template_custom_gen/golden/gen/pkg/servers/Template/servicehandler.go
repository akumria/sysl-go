// Code generated by sysl DO NOT EDIT.
package template

import (
	"context"
	"fmt"
	"net/http"

	"github.com/anz-bank/sysl-go/common"
	"github.com/anz-bank/sysl-go/config"
	"github.com/anz-bank/sysl-go/core"
	"github.com/anz-bank/sysl-go/core/authrules"
	"github.com/anz-bank/sysl-go/restlib"
	"github.com/anz-bank/sysl-go/validator"
)

// Handler interface for Template
type Handler interface {
	GetEndpointListHandler(w http.ResponseWriter, r *http.Request)
}

// ServiceHandler for Template API
type ServiceHandler struct {
	genCallback        core.RestGenCallback
	serviceInterface   *ServiceInterface
	authorizationRules map[string]authrules.Rule
}

// NewServiceHandler for Template
func NewServiceHandler(
	ctx context.Context,
	cfg *config.DefaultConfig,
	hooks *core.Hooks,
	genCallback core.RestGenCallback,
	serviceInterface *ServiceInterface,
) (*ServiceHandler, error) {

	authorizationRules := make(map[string]authrules.Rule)

	return &ServiceHandler{
		genCallback,
		serviceInterface,
		authorizationRules,
	}, nil
}

// GetEndpointListHandler ...
func (s *ServiceHandler) GetEndpointListHandler(w http.ResponseWriter, r *http.Request) {
	if s.serviceInterface.GetEndpointList == nil {
		common.HandleError(r.Context(), w, common.InternalError, "not implemented", nil, s.genCallback.MapError)
		return
	}

	ctx := common.RequestHeaderToContext(r.Context(), r.Header)
	ctx = common.RespHeaderAndStatusToContext(ctx, make(http.Header), http.StatusOK)
	var req GetEndpointListRequest

	ctx, cancel := s.genCallback.DownstreamTimeoutContext(ctx)
	defer cancel()
	valErr := validator.Validate(&req)
	if valErr != nil {
		common.HandleError(ctx, w, common.BadRequestError, "Invalid request", valErr, s.genCallback.MapError)
		return
	}

	client := GetEndpointListClient{}

	defer func() {
		if rec := recover(); rec != nil {
			var err error
			switch rec := rec.(type) {
			case error:
				err = rec
			default:
				err = fmt.Errorf("Unknown error: %v", rec)
			}
			common.HandleError(ctx, w, common.InternalError, "Unexpected panic", err, s.genCallback.MapError)
		}
	}()
	err := s.serviceInterface.GetEndpointList(ctx, &req, client)
	if err != nil {
		common.HandleError(ctx, w, common.InternalError, "Handler error", err, s.genCallback.MapError)
		return
	}

	headermap, httpstatus := common.RespHeaderAndStatusFromContext(ctx)
	if headermap.Get("Content-Type") == "" {
		headermap.Set("Content-Type", "application/json")
	}
	restlib.SetHeaders(w, headermap)
	restlib.SendHTTPResponse(w, httpstatus, nil)
}
