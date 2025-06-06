package public

import (
	"connectrpc.com/connect"
	"context"
	"github.com/jamespfennell/transiter/internal/gen/api"
	"github.com/jamespfennell/transiter/internal/gen/api/apiconnect"
)

type ConnectServer struct {
	impl *Server
}

func NewConnectWrapper(impl *Server) *ConnectServer {
	return &ConnectServer{impl: impl}
}

func wrap[TReq any, TResp any](ctx context.Context, req *connect.Request[TReq], h func(context.Context, *TReq) (*TResp, error)) (*connect.Response[TResp], error) {
	ret, err := h(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(ret), nil
}

func (c *ConnectServer) Entrypoint(ctx context.Context, req *connect.Request[api.EntrypointRequest]) (*connect.Response[api.EntrypointReply], error) {
	return wrap(ctx, req, c.impl.Entrypoint)
}

func (c *ConnectServer) ListSystems(ctx context.Context, req *connect.Request[api.ListSystemsRequest]) (*connect.Response[api.ListSystemsReply], error) {
	return wrap(ctx, req, c.impl.ListSystems)
}

func (c *ConnectServer) GetSystem(ctx context.Context, req *connect.Request[api.GetSystemRequest]) (*connect.Response[api.System], error) {
	return wrap(ctx, req, c.impl.GetSystem)
}

func (c *ConnectServer) ListAgencies(ctx context.Context, req *connect.Request[api.ListAgenciesRequest]) (*connect.Response[api.ListAgenciesReply], error) {
	return wrap(ctx, req, c.impl.ListAgencies)
}

func (c *ConnectServer) GetAgency(ctx context.Context, req *connect.Request[api.GetAgencyRequest]) (*connect.Response[api.Agency], error) {
	return wrap(ctx, req, c.impl.GetAgency)
}

func (c *ConnectServer) ListStops(ctx context.Context, req *connect.Request[api.ListStopsRequest]) (*connect.Response[api.ListStopsReply], error) {
	return wrap(ctx, req, c.impl.ListStops)
}

func (c *ConnectServer) GetStop(ctx context.Context, req *connect.Request[api.GetStopRequest]) (*connect.Response[api.Stop], error) {
	return wrap(ctx, req, c.impl.GetStop)
}

func (c *ConnectServer) ListRoutes(ctx context.Context, req *connect.Request[api.ListRoutesRequest]) (*connect.Response[api.ListRoutesReply], error) {
	return wrap(ctx, req, c.impl.ListRoutes)
}

func (c *ConnectServer) GetRoute(ctx context.Context, req *connect.Request[api.GetRouteRequest]) (*connect.Response[api.Route], error) {
	return wrap(ctx, req, c.impl.GetRoute)
}

func (c *ConnectServer) ListTrips(ctx context.Context, req *connect.Request[api.ListTripsRequest]) (*connect.Response[api.ListTripsReply], error) {
	return wrap(ctx, req, c.impl.ListTrips)
}

func (c *ConnectServer) GetTrip(ctx context.Context, req *connect.Request[api.GetTripRequest]) (*connect.Response[api.Trip], error) {
	return wrap(ctx, req, c.impl.GetTrip)
}

func (c *ConnectServer) ListAlerts(ctx context.Context, req *connect.Request[api.ListAlertsRequest]) (*connect.Response[api.ListAlertsReply], error) {
	return wrap(ctx, req, c.impl.ListAlerts)
}

func (c *ConnectServer) GetAlert(ctx context.Context, req *connect.Request[api.GetAlertRequest]) (*connect.Response[api.Alert], error) {
	return wrap(ctx, req, c.impl.GetAlert)
}

func (c *ConnectServer) ListTransfers(ctx context.Context, req *connect.Request[api.ListTransfersRequest]) (*connect.Response[api.ListTransfersReply], error) {
	return wrap(ctx, req, c.impl.ListTransfers)
}

func (c *ConnectServer) GetTransfer(ctx context.Context, req *connect.Request[api.GetTransferRequest]) (*connect.Response[api.Transfer], error) {
	return wrap(ctx, req, c.impl.GetTransfer)
}

func (c *ConnectServer) ListFeeds(ctx context.Context, req *connect.Request[api.ListFeedsRequest]) (*connect.Response[api.ListFeedsReply], error) {
	return wrap(ctx, req, c.impl.ListFeeds)
}

func (c *ConnectServer) GetFeed(ctx context.Context, req *connect.Request[api.GetFeedRequest]) (*connect.Response[api.Feed], error) {
	return wrap(ctx, req, c.impl.GetFeed)
}

func (c *ConnectServer) ListVehicles(ctx context.Context, req *connect.Request[api.ListVehiclesRequest]) (*connect.Response[api.ListVehiclesReply], error) {
	return wrap(ctx, req, c.impl.ListVehicles)
}

func (c *ConnectServer) GetVehicle(ctx context.Context, req *connect.Request[api.GetVehicleRequest]) (*connect.Response[api.Vehicle], error) {
	return wrap(ctx, req, c.impl.GetVehicle)
}

func (c *ConnectServer) ListShapes(ctx context.Context, req *connect.Request[api.ListShapesRequest]) (*connect.Response[api.ListShapesReply], error) {
	return wrap(ctx, req, c.impl.ListShapes)
}

func (c *ConnectServer) GetShape(ctx context.Context, req *connect.Request[api.GetShapeRequest]) (*connect.Response[api.Shape], error) {
	return wrap(ctx, req, c.impl.GetShape)
}

var _ apiconnect.PublicHandler = &ConnectServer{}
