package handlers

import (
	"fmt"

	"github.com/forbole/callisto/v4/modules/actions/types"

	"github.com/rs/zerolog/log"
)

func TotalSupplyHandler(ctx *types.Context, payload *types.Payload) (interface{}, error) {
	log.Debug().Str("address", payload.GetAddress()).
		Int64("height", payload.Input.Height).
		Msg("executing total supply action")

	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	balance, err := ctx.Sources.BankSource.GetSupply(height)
	if err != nil {
		return nil, fmt.Errorf("error while getting total supply: %s", err)
	}

	return types.Balance{
		Coins: types.ConvertCoins(balance),
	}, nil
}
