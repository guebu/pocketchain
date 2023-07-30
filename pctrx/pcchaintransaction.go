package pctrx

import (
	"encoding/json"
	"fmt"
	"github.com/guebu/common-utils/errors"
	"github.com/guebu/common-utils/logger"
	"strings"
)

type PCChainTransaction struct {
	//SenderBlockChainAddress    string
	//RecipientBlockChainAddress string
	//Value                      float32
	PCTransaction
}

func NewPCChainTransaction(senderAddress string, recipient string, value float32) *PCChainTransaction {
	//return &PCChainTransaction{senderAddress, recipient, value}
	trx := NewTransaction(senderAddress, recipient, value)
	return &PCChainTransaction{*trx}
}

func (trx *PCChainTransaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 60))
	fmt.Printf("... sender    PCBC address:     %s\n", trx.SenderBlockChainAddress)
	fmt.Printf("... recipient PCBC address:     %s\n", trx.RecipientBlockChainAddress)
	fmt.Printf("... value:                      %v\n", trx.Value)
}

func (trx *PCChainTransaction) MarshalJSON() ([]byte, *errors.ApplicationError) {
	logger.Info("Creating JSON for chain transation", "Status:Open", "App:ChainTrx")
	json, err := json.Marshal(struct {
		Sender    string  `json:"sender_pocketchain_address"`
		Recipient string  `json:"recipient_pocketchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    trx.SenderBlockChainAddress,
		Recipient: trx.RecipientBlockChainAddress,
		Value:     trx.Value,
	})

	if err != nil {
		return nil, errors.NewInternalServerError("JSON for current trx couldn't be created!", err)
	}
	logger.Info("Creating JSON for chain transation", "Status:End", "App:ChainTrx")
	return json, nil
}
