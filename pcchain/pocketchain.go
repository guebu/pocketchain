package pcchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/guebu/common-utils/logger"
	"github.com/guebu/pocketchain/pcblock"
	"github.com/guebu/pocketchain/pctrx"
	"github.com/guebu/pocketchain/pcutil"
	"strings"
	"sync"
)

const INIT_HASH = "Init Hash"
const MINING_DIFFICULTY = 3
const MINING_SENDER = "THE Pocket Chain!"
const MINING_REWARD = 1.0

type PocketChain struct {
	transactionPool    []*pctrx.PCChainTransaction
	chainOfBlocks      []*pcblock.PCBlock
	pocketChainAddress string
}

func NewPocketChain(pcChainAddress string) *PocketChain {
	b := &pcblock.PCBlock{}
	pc := new(PocketChain)
	h, _ := b.GenerateHashForBlock()
	pc.CreateBlock(0, h)
	pc.pocketChainAddress = pcChainAddress
	return pc
}

func (pc *PocketChain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey,
	trxSignature *pcutil.PCSignature, trx *pctrx.PCChainTransaction) bool {

	m, err := json.Marshal(trx)

	if err != nil {
		logger.Error("Error during marshalling the transaction!", err, "Status:Error", "App:PocketChain")
	}

	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], trxSignature.R, trxSignature.S)
}

func (pc *PocketChain) GetLastBlock() *pcblock.PCBlock {
	return pc.chainOfBlocks[len(pc.chainOfBlocks)-1]
}

func (pc *PocketChain) CreateBlock(nonce int, previousHash [32]byte) *pcblock.PCBlock {
	logger.Info("... creating a block", "Status:Open", "App:PocketChain")
	b := pcblock.NewBlock(nonce, previousHash, pc.transactionPool)
	pc.chainOfBlocks = append(pc.chainOfBlocks, b)
	// reset of transaction pool
	pc.transactionPool = []*pctrx.PCChainTransaction{}
	logger.Info("... creating a block", "Status:End", "App:PocketChain")
	return b
}

func (pc *PocketChain) AddTransaction(sender string, recipient string, value float32,
	senderPublicKey *ecdsa.PublicKey, trxSig *pcutil.PCSignature) bool {
	trx := pctrx.NewPCChainTransaction(sender, recipient, value)

	if sender == MINING_SENDER {
		pc.transactionPool = append(pc.transactionPool, trx)
		return true
	}

	sigWasOK := pc.VerifyTransactionSignature(senderPublicKey, trxSig, trx)

	if sigWasOK {
		// check if sender has enough money...
		balOfSender := pc.CalculateTotalAmount(sender)
		if balOfSender < value {
			// Sender has not enough money!!
			logger.Error("Sender has not enough money to execute transaction!", nil, "Status:End", "App:PocketChain")
			return false
		}
		pc.transactionPool = append(pc.transactionPool, trx)
		return true
	} else {
		logger.Error("Transaction couldn't be verified!", nil, "Status:Error", "App:PocketChain")
	}
	return false
}

func (pc *PocketChain) Print() {
	for i, block := range pc.chainOfBlocks {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s \n", strings.Repeat("*", 60))
}

func (pc *PocketChain) CopyTrxFromPoolToBlock() []*pctrx.PCChainTransaction {
	trx := make([]*pctrx.PCChainTransaction, 0)

	for _, t := range pc.transactionPool {
		trx = append(trx,
			pctrx.NewPCChainTransaction(
				t.SenderBlockChainAddress,
				t.RecipientBlockChainAddress,
				t.Value))

	}
	return trx
}

func (pc *PocketChain) ValidProof(nonce int, prevHash [32]byte, trx []*pctrx.PCChainTransaction, difficulty int) bool {
	//defer w.Done()
	zeros := strings.Repeat("0", difficulty)
	guessBlock := pcblock.PCBlock{0, nonce, prevHash, trx}

	h, _ := guessBlock.GenerateHashForBlock()
	guessHashString := fmt.Sprintf("%x", h)

	return guessHashString[:difficulty] == zeros
}

func (pc *PocketChain) ProofOfWork() int {
	trx := pc.CopyTrxFromPoolToBlock()
	prevHash, _ := pc.GetLastBlock().GenerateHashForBlock()

	nonce := 0
	var wg sync.WaitGroup

	validProof := false
	for nonce = 0; !validProof; nonce++ {
		wg.Add(1)
		//validProof = go pc.ValidProof(nonce, prevHash, trx, MINING_DIFFICULTY, &wg)
		go func() {
			defer wg.Done()
			zeros := strings.Repeat("0", MINING_DIFFICULTY)
			guessBlock := pcblock.PCBlock{0, nonce, prevHash, trx}

			h, _ := guessBlock.GenerateHashForBlock()
			guessHashString := fmt.Sprintf("%x", h)

			if guessHashString[:MINING_DIFFICULTY] == zeros {
				validProof = true
			}
		}()
	}
	wg.Wait()
	return nonce
}

func (pc *PocketChain) Mining() bool {
	logger.Info("Mining starts!", "Status:End", "App:PocketChain")
	// because as a miner, the transactions mustn't be verified, we set senderPublicKey and signature to nil!
	pc.AddTransaction(MINING_SENDER, pc.pocketChainAddress, MINING_REWARD, nil, nil)
	nonce := pc.ProofOfWork()
	prevHash, _ := pc.GetLastBlock().GenerateHashForBlock()
	// after proof of work was finished and we know the nonce fitting to the reward, we can create a new block
	// which also includes the transaction
	pc.CreateBlock(nonce, prevHash)
	logger.Info("Mining was successfull!", "Status:End", "App:PocketChain")
	return true
}

func (pc *PocketChain) CalculateTotalAmount(blockChainAdress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range pc.chainOfBlocks {
		for _, t := range b.Transactions {
			value := t.Value
			if blockChainAdress == t.RecipientBlockChainAddress {
				// given address receives money. So we increase the total amount
				totalAmount += value
			}

			if blockChainAdress == t.SenderBlockChainAddress {
				// given address sended money. So we have to decrease the total amount
				totalAmount -= value
			}
		}
	}

	return totalAmount
}
