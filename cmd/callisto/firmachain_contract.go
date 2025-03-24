package main

//TODO: fix import firmachain
/*import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	contracttypes "github.com/firmachain/firmachain/x/contract/types"
)

func FirmaChainContractMessagesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	switch msg := cosmosMsg.(type) {
	case *contracttypes.MsgAddContractLog:
		return []string{msg.Creator, msg.OwnerAddress}, nil

	case *contracttypes.MsgCreateContractFile:
		var addresses []string
		addresses = append(addresses, msg.Creator)
		addresses = append(addresses, msg.OwnerList...)

		return addresses, nil
	}

	return nil, nil
}*/
