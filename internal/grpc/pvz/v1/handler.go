package pvzv1

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *PVZServer) GetPVZList(ctx context.Context, req *GetPVZListRequest) (*GetPVZListResponse, error) {
	PVZsInfo, err := s.pvzService.GetPvz(nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	var res []*PVZ
	for _, pvzInfo := range PVZsInfo {
		res = append(res, &PVZ{
			Id:               pvzInfo.Pvz.Id.String(),
			City:             pvzInfo.Pvz.City,
			RegistrationDate: timestamppb.New(pvzInfo.Pvz.RegistrationDate),
		})
	}

	return &GetPVZListResponse{
		Pvzs: res,
	}, nil
}
