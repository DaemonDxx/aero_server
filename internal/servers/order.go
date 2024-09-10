package servers

import (
	"context"
	orderpb "github.com/daemondxx/lks_back/gen/pb/go/order"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderServer struct {
	orderpb.UnimplementedOrderServiceServer
	oServ OrderService
}

func NewOrderServer(o OrderService) *OrderServer {
	return &OrderServer{oServ: o}
}

func (o *OrderServer) GetActualOrder(ctx context.Context, r *orderpb.ActualOrderRequest) (*orderpb.ActualOrderResponse, error) {
	order, err := o.oServ.GetActualOrder(ctx, uint(r.UserId))
	if err != nil {
		return nil, ErrInternal
	}
	res := orderpb.Order{
		Id:     uint64(order.ID),
		UserId: uint64(order.UserID),
	}
	for _, i := range order.Items {
		flights := make([]*orderpb.Flight, 0, len(i.Flights))
		for _, f := range i.Flights {
			flights = append(flights, &orderpb.Flight{
				Id:           uint64(f.ID),
				FlightNumber: f.FlightNumber,
				Airplane:     f.Airplane,
				Status:       int64(f.Status),
			})
		}
		oItem := orderpb.OrderItem{
			Id:      uint64(i.ID),
			Flights: nil,
			Departure: &timestamppb.Timestamp{
				Seconds: int64(i.Departure.Unix()),
				Nanos:   int32(i.Departure.Nanosecond()),
			},
			Arrival: &timestamppb.Timestamp{
				Seconds: int64(i.Arrival.Unix()),
				Nanos:   int32(i.Arrival.Nanosecond()),
			},
			Description: i.Description,
			Route:       i.Route,
		}

		if i.ConfirmDate != nil {
			oItem.ConfirmDate = &timestamppb.Timestamp{
				Seconds: int64(i.ConfirmDate.Unix()),
				Nanos:   int32(i.ConfirmDate.Nanosecond()),
			}
		}
		oItem.Flights = flights
		res.Items = append(res.Items, &oItem)
	}

	return &orderpb.ActualOrderResponse{
		Order: &res,
	}, nil
}

func (o *OrderServer) mustEmbedUnimplementedOrderServiceServer() {
	panic("implement me")
}
