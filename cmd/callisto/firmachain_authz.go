package main

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authztype "github.com/cosmos/cosmos-sdk/x/authz"
)

func FirmaChainAuthzMessagesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {

	switch msg := cosmosMsg.(type) {

	case *authztype.MsgGrant:
		var stringArray = []string{msg.Grantee, msg.Granter}
		return stringArray, nil

	case *authztype.MsgRevoke:
		var stringArray = []string{msg.Grantee, msg.Granter}
		return stringArray, nil

	case *authztype.MsgExec:
		msgs, _ := msg.GetMessages()

		var stringArray = []string{}
		stringArray = append(stringArray, msg.Grantee)

		total := len(msgs)

		for i := 0; i < total; i++ {
			msgText := msgs[i].String()

			orgTotalLength := len(msgText)

			if orgTotalLength > 0 {

				totalLength := orgTotalLength
				msgTempText := msgText

				for i := 0; i < totalLength; i++ {
					idx := strings.Index(msgTempText, "firma1")
					if idx != -1 {
						const lenghOfAddress = 44
						tempAddress := msgTempText[idx : idx+lenghOfAddress]
						msgTempText = msgTempText[idx+lenghOfAddress:]
						totalLength = len(msgTempText)
						i = 0

						isFindAddress := false
						for _, v := range stringArray {
							if v == tempAddress {
								isFindAddress = true
								break
							}
						}

						if !isFindAddress {
							stringArray = append(stringArray, tempAddress)
						}
					}
				}

				totalLength = orgTotalLength
				msgTempText = msgText

				for i := 0; i < totalLength; i++ {
					idx := strings.Index(msgTempText, "firmavaloper1")
					if idx != -1 {
						const lenghOfValidatorAddress = 51
						tempAddress := msgTempText[idx : idx+lenghOfValidatorAddress]
						msgTempText = msgTempText[idx+lenghOfValidatorAddress:]
						totalLength = len(msgTempText)
						i = 0

						isFindAddress := false
						for _, v := range stringArray {
							if v == tempAddress {
								isFindAddress = true
								break
							}
						}

						if !isFindAddress {
							stringArray = append(stringArray, tempAddress)
						}
					}
				}
			}
		}

		return stringArray, nil
	}

	return nil, nil
}
