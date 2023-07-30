package pctrx

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/guebu/common-utils/errors"
	"github.com/guebu/common-utils/logger"
	"github.com/guebu/pocketchain/pcutil"
	"strings"
)

type PCWalletTransaction struct {
	PCTransaction
	SenderPrivateKey *ecdsa.PrivateKey
	SenderPublicKey  *ecdsa.PublicKey
	//SenderBlockChainAddress    string
	//RecipientBlockChainAddress string
	//Value                      float32
}

func NewPCWalletTransaction(senderPrivateKey *ecdsa.PrivateKey, senderPublicKey *ecdsa.PublicKey, senderAddress string, recipient string, value float32) *PCWalletTransaction {
	trx := NewTransaction(senderAddress, recipient, value)

	return &PCWalletTransaction{*trx, senderPrivateKey, senderPublicKey}
}

func (trx *PCWalletTransaction) GenerateSignature() *pcutil.PCSignature {
	t := NewTransaction(trx.SenderBlockChainAddress, trx.RecipientBlockChainAddress, trx.Value)

	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	r, s, err := ecdsa.Sign(rand.Reader, trx.SenderPrivateKey, h[:])
	if err != nil {
		logger.Error("Error during generating signature!", err, "Status: Error", "App:PCWalletTransaction")
	}
	return &pcutil.PCSignature{r, s}
}

func (trx *PCWalletTransaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 60))
	fmt.Printf("... sender    PCBC address:     %s\n", trx.SenderBlockChainAddress)
	fmt.Printf("... recipient PCBC address:     %s\n", trx.RecipientBlockChainAddress)
	fmt.Printf("... value:                      %v\n", trx.Value)
}

func (trx *PCWalletTransaction) MarshalJSON() ([]byte, error) {
	logger.Info("Creating JSON for wallet transation", "Status:Open", "App:WalletTrx")
	trxJSON, err := json.Marshal(struct {
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
	logger.Info("Creating JSON for wallet transation", "Status:Open", "App:WalletTrx")
	//fmt.Printf(trxJSON.Byte)
	return trxJSON, nil
}
