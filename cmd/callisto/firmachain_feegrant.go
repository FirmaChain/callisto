package main

import (
	cdc "github.com/cosmos/cosmos-sdk/codec"

	feegranttype "cosmossdk.io/x/feegrant"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v6/types"
)

func FirmaChainFeegrantMessagesParser(tx *types.Transaction) ([]string, error) {
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
		case *feegranttype.MsgGrantAllowance:
			addresses = append(addresses, msg.Grantee, msg.Granter)
		case *feegranttype.MsgRevokeAllowance:
			addresses = append(addresses, msg.Grantee, msg.Granter)
		}
	}
	return addresses, nil
}
