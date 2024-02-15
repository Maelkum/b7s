package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/request"
)

type ExecuteControlRequest struct {
	RequestID string `json:"request_id,omitempty"`
	Action    string `json:"action,omitempty"`
}

type ExecuteControlResponse struct {
	RequestID string     `json:"request_id,omitempty"`
	Action    string     `json:"action,omitempty"`
	Code      codes.Code `json:"code,omitempty"`
}

func (a *API) ExecControl(ctx echo.Context) error {

	var req ExecuteControlRequest
	err := ctx.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("could not unpack request: %w", err))
	}

	// TODO: Support other actions.
	// action := parseAction(req.Action)

	action := request.ExecWait
	err = a.Node.ExecutionControl(req.RequestID, action)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("could not control execution: %w", err))
	}

	res := ExecuteControlResponse{
		RequestID: req.RequestID,
		Action:    action.String(),
		Code:      codes.OK,
	}

	return ctx.JSON(http.StatusOK, res)
}
