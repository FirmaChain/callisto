package main

//TODO: fix import firmachain
/*import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tokentypes "github.com/firmachain/firmachain/x/token/types"
)

func FirmaChainTokenMessagesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	switch msg := cosmosMsg.(type) {
	case *tokentypes.MsgCreateToken:
		return []string{msg.Owner}, nil

	case *tokentypes.MsgBurn:
		return []string{msg.Owner}, nil

	case *tokentypes.MsgUpdateTokenURI:
		return []string{msg.Owner}, nil

	case *tokentypes.MsgMint:
		return []string{msg.Owner, msg.ToAddress}, nil
	}

	return nil, nil
}*/
