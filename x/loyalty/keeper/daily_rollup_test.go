package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"tokenchain/x/loyalty/types"
)

func TestRunDailyRollup_OncePerLocalDay(t *testing.T) {
	f := initFixture(t)

	ctxDayStart := sdk.UnwrapSDKContext(f.ctx).
		WithBlockTime(time.Date(2026, 2, 26, 8, 0, 0, 0, time.UTC)).
		WithEventManager(sdk.NewEventManager())
	require.NoError(t, f.keeper.RunDailyRollup(ctxDayStart))

	lastDate, err := f.keeper.LastDailyRollupDate.Get(ctxDayStart)
	require.NoError(t, err)
	require.Equal(t, "2026-02-26", lastDate)
	requireRollupEvent(t, ctxDayStart.EventManager().Events(), "2026-02-26", "America/Edmonton")

	ctxSameDay := sdk.UnwrapSDKContext(f.ctx).
		WithBlockTime(time.Date(2026, 2, 26, 20, 0, 0, 0, time.UTC)).
		WithEventManager(sdk.NewEventManager())
	require.NoError(t, f.keeper.RunDailyRollup(ctxSameDay))
	require.Empty(t, rollupEvents(ctxSameDay.EventManager().Events()))
}

func TestRunDailyRollup_UsesEdmontonBoundary(t *testing.T) {
	f := initFixture(t)

	// 06:59 UTC is still previous local day in America/Edmonton during winter (UTC-7).
	ctxBeforeMidnight := sdk.UnwrapSDKContext(f.ctx).
		WithBlockTime(time.Date(2026, 2, 26, 6, 59, 0, 0, time.UTC)).
		WithEventManager(sdk.NewEventManager())
	require.NoError(t, f.keeper.RunDailyRollup(ctxBeforeMidnight))
	requireRollupEvent(t, ctxBeforeMidnight.EventManager().Events(), "2026-02-25", "America/Edmonton")

	// 07:01 UTC crosses local midnight (00:01 in America/Edmonton).
	ctxAfterMidnight := sdk.UnwrapSDKContext(f.ctx).
		WithBlockTime(time.Date(2026, 2, 26, 7, 1, 0, 0, time.UTC)).
		WithEventManager(sdk.NewEventManager())
	require.NoError(t, f.keeper.RunDailyRollup(ctxAfterMidnight))
	requireRollupEvent(t, ctxAfterMidnight.EventManager().Events(), "2026-02-26", "America/Edmonton")
}

func requireRollupEvent(t *testing.T, events sdk.Events, expectedDate string, expectedTimezone string) {
	t.Helper()

	found := rollupEvents(events)
	require.Len(t, found, 1)

	event := found[0]
	require.Equal(t, expectedDate, attrValue(event, types.AttributeKeyDate))
	require.Equal(t, expectedTimezone, attrValue(event, types.AttributeKeyTimezone))
}

func rollupEvents(events sdk.Events) []sdk.Event {
	out := make([]sdk.Event, 0, len(events))
	for _, event := range events {
		if event.Type == types.EventTypeDailyRollup {
			out = append(out, event)
		}
	}
	return out
}

func attrValue(event sdk.Event, key string) string {
	for _, attr := range event.Attributes {
		if attr.Key == key {
			return attr.Value
		}
	}
	return ""
}
