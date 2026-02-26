package keeper

import (
	"context"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"tokenchain/x/loyalty/types"
)

func (q queryServer) FilterMerchantallocation(ctx context.Context, req *types.QueryFilterMerchantallocationRequest) (*types.QueryFilterMerchantallocationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	filterDate := strings.TrimSpace(req.Date)
	if filterDate != "" {
		if _, err := time.Parse(rollupDateLayout, filterDate); err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid date filter")
		}
	}

	filterDenom := strings.TrimSpace(req.Denom)
	if filterDenom != "" {
		if err := sdk.ValidateDenom(filterDenom); err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid denom filter")
		}
	}

	filtered := make([]types.Merchantallocation, 0)
	if err := q.k.Merchantallocation.Walk(ctx, nil, func(_ string, record types.Merchantallocation) (bool, error) {
		if filterDate != "" && record.Date != filterDate {
			return false, nil
		}
		if filterDenom != "" && record.Denom != filterDenom {
			return false, nil
		}
		filtered = append(filtered, record)
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

	items := make([]types.Merchantallocation, end-start)
	copy(items, filtered[start:end])

	pageRes := &query.PageResponse{}
	if end < total {
		pageRes.NextKey = []byte(strconv.FormatUint(end, 10))
	}
	if needTotal {
		pageRes.Total = total
	}

	return &types.QueryFilterMerchantallocationResponse{
		Merchantallocation: items,
		Pagination:         pageRes,
	}, nil
}
