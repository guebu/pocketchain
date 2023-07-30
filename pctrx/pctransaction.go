package pctrx

import (
	"encoding/json"
	"fmt"
	"github.com/guebu/common-utils/errors"
	"github.com/guebu/common-utils/logger"
	"strings"
)

type PCTransaction struct {
	SenderBlockChainAddress    string
	RecipientBlockChainAddress string
	Value                      float32
}

func NewTransaction(senderAddress string, recipient string, value float32) *PCTransaction {
	return &PCTransaction{senderAddress, recipient, value}
}

func (trx *PCTransaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 60))
	fmt.Printf("... sender    PCBC address:     %s\n", trx.SenderBlockChainAddress)
	fmt.Printf("... recipient PCBC address:     %s\n", trx.RecipientBlockChainAddress)
	fmt.Printf("... value:                      %v\n", trx.Value)
}

func (trx *PCTransaction) MarshalJSON() ([]byte, *errors.ApplicationError) {
	logger.Info("Creating JSON for common transation", "Status:Open", "App:Trx")
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
	logger.Info("Creating JSON for common transation", "Status:End", "App:Trx")
	return json, nil
}
