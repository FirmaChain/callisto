package main

import (
	cdc "github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tokentypes "github.com/firmachain/firmachain/x/token/types"
	"github.com/forbole/juno/v6/types"
)

func FirmaChainTokenMessagesParser(tx *types.Transaction) ([]string, error) {
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
		case *tokentypes.MsgCreateToken:
			addresses = append(addresses, msg.Owner)
		case *tokentypes.MsgBurn:
			addresses = append(addresses, msg.Owner)
		case *tokentypes.MsgUpdateTokenURI:
			addresses = append(addresses, msg.Owner)
		case *tokentypes.MsgMint:
			addresses = append(addresses, msg.Owner, msg.ToAddress)
		}
	}
	return addresses, nil
}
