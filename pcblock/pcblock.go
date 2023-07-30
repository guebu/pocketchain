package pcblock

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/guebu/common-utils/errors"
	"github.com/guebu/common-utils/logger"
	"github.com/guebu/pocketchain/pctrx"
	"time"
)

type PCBlock struct {
	TimeStamp             int64
	Nonce                 int
	HashOfPreviousBCBlock [32]byte
	Transactions          []*pctrx.PCChainTransaction
}

type BCBlockHash = [32]byte

func (b *PCBlock) Print() {
	//logger.Info("... printing a block", "Status:Open", "App:bcblock")
	fmt.Printf("Timestamp:     %d\n", b.TimeStamp)
	fmt.Printf("Nonce:         %d\n", b.Nonce)
	fmt.Printf("Previous hash: %x\n", b.HashOfPreviousBCBlock)

	for _, trx := range b.Transactions {
		trx.Print()
	}
	//logger.Info("... printing a block", "Status:End", "App:bcblock")
}

func (b *PCBlock) GenerateHashForBlock() (BCBlockHash, *errors.ApplicationError) {
	json, err := b.MarshalJSON()
	var someValue BCBlockHash

	if err != nil {
		return someValue, err
	}

	//fmt.Println(string(json))

	return sha256.Sum256([]byte(json)), nil

}

func (b *PCBlock) MarshalJSON() ([]byte, *errors.ApplicationError) {
	json, err := json.Marshal(struct {
		Timestamp    int64                       `json:"timestamp"`
		Nonce        int                         `json:"nonce"`
		PreviousHash [32]byte                    `json:"previous_hash"`
		Transactions []*pctrx.PCChainTransaction `json:"transactions"`
	}{
		Timestamp:    b.TimeStamp,
		Nonce:        b.Nonce,
		PreviousHash: b.HashOfPreviousBCBlock,
		Transactions: b.Transactions,
	})

	if err != nil {
		return nil, errors.NewInternalServerError("JSON for current block couldn't be created!", err)
	}
	return json, nil
}

func NewBlock(nonce int, previousHash [32]byte, trx []*pctrx.PCChainTransaction) *PCBlock {
	logger.Info("... creating a new block", "Status:Open", "App:bcblock")
	b := new(PCBlock)
	b.TimeStamp = time.Now().UnixNano()
	b.Nonce = nonce
	b.HashOfPreviousBCBlock = previousHash
	b.Transactions = trx
	logger.Info("... new block created", "Status:End", "App:bcblock")
	return b
}
