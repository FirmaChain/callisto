package handlers

import (
	"fmt"

	"github.com/forbole/callisto/v4/modules/actions/types"

	"github.com/rs/zerolog/log"
)

func InflationHandler(ctx *types.Context, payload *types.Payload) (interface{}, error) {
	log.Debug().Str("address", payload.GetAddress()).
		Int64("height", payload.Input.Height).
		Msg("executing inflation action")

	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	inflation, err := ctx.Sources.MintSource.GetInflation(height)
	if err != nil {
		return nil, fmt.Errorf("error while getting inflation: %s", err)
	}

	return types.Inflation{
		Amount: inflation.String(),
	}, nil
}
