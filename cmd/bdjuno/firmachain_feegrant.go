package main

import (
	"github.com/cosmos/cosmos-sdk/codec"

	feegranttype "cosmossdk.io/x/feegrant"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func FirmaChainFeegrantMessagesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {

	switch msg := cosmosMsg.(type) {

	case *feegranttype.MsgGrantAllowance:
		var stringArray = []string{msg.Grantee, msg.Granter}
		return stringArray, nil

	case *feegranttype.MsgRevokeAllowance:
		var stringArray = []string{msg.Grantee, msg.Granter}
		return stringArray, nil
	}

	return nil, nil
}
