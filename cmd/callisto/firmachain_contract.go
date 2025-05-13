package main

import (
	cdc "github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	contracttypes "github.com/firmachain/firmachain/x/contract/types"
	"github.com/forbole/juno/v6/types"
)

func FirmaChainContractMessagesParser(tx *types.Transaction) ([]string, error) {
	// Create a slice to hold the addresses.
	var addresses []string

	// Loop over each message in the transaction body.
	for _, anyMsg := range tx.Body.Messages {
		// Unpack the Any message into a concrete sdk.Msg.
		anyMsgByte := anyMsg.GetBytes()
		var cosmosMsg sdk.Msg
		if err := cdc.NewLegacyAmino().Unmarshal(anyMsgByte, &cosmosMsg); err != nil {
			return nil, err
		}
		switch msg := cosmosMsg.(type) {
		case *contracttypes.MsgAddContractLog:
			addresses = append(addresses, msg.Creator, msg.OwnerAddress)
		case *contracttypes.MsgCreateContractFile:
			addresses = append(addresses, msg.Creator)
			addresses = append(addresses, msg.OwnerList...)
		}
	}
	return addresses, nil
}
