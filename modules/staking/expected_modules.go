package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/callisto/v4/types"
)

type SlashingModule interface {
	GetSigningInfo(height int64, consAddr sdk.ConsAddress) (types.ValidatorSigningInfo, error)
}
