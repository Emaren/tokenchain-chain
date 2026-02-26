package keeper

import (
	"context"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"tokenchain/x/loyalty/types"
)

func (q queryServer) FilterRecoveryoperation(ctx context.Context, req *types.QueryFilterRecoveryoperationRequest) (*types.QueryFilterRecoveryoperationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	filterStatus := strings.TrimSpace(req.Status)
	switch filterStatus {
	case "", types.RecoveryStatusQueued, types.RecoveryStatusExecuted, types.RecoveryStatusCancelled:
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid status filter")
	}

	filterDenom := strings.TrimSpace(req.Denom)
	filterRequestedBy := strings.TrimSpace(req.RequestedBy)
	filterFromAddress := strings.TrimSpace(req.FromAddress)
	filterToAddress := strings.TrimSpace(req.ToAddress)

	filtered := make([]types.Recoveryoperation, 0)
	if err := q.k.Recoveryoperation.Walk(ctx, nil, func(_ uint64, op types.Recoveryoperation) (bool, error) {
		if filterStatus != "" && op.Status != filterStatus {
			return false, nil
		}
		if filterDenom != "" && op.Denom != filterDenom {
			return false, nil
		}
		if filterRequestedBy != "" && op.RequestedBy != filterRequestedBy {
			return false, nil
		}
		if filterFromAddress != "" && op.FromAddress != filterFromAddress {
			return false, nil
		}
		if filterToAddress != "" && op.ToAddress != filterToAddress {
			return false, nil
		}
		filtered = append(filtered, op)
		return false, nil
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	start := uint64(0)
	limit := uint64(len(filtered))
	needTotal := true
	if req.Pagination != nil {
		needTotal = req.Pagination.CountTotal
		if len(req.Pagination.Key) > 0 {
			keyStart, err := strconv.ParseUint(string(req.Pagination.Key), 10, 64)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, "invalid pagination key")
			}
			start = keyStart
		} else {
			start = req.Pagination.Offset
		}
		if req.Pagination.Limit > 0 {
			limit = req.Pagination.Limit
		}
	}

	total := uint64(len(filtered))
	if start > total {
		start = total
	}
	end := total
	if limit < end-start {
		end = start + limit
	}

	items := make([]types.Recoveryoperation, end-start)
	copy(items, filtered[start:end])

	pageRes := &query.PageResponse{}
	if end < total {
		pageRes.NextKey = []byte(strconv.FormatUint(end, 10))
	}
	if needTotal {
		pageRes.Total = total
	}

	return &types.QueryFilterRecoveryoperationResponse{
		Recoveryoperation: items,
		Pagination:        pageRes,
	}, nil
}
