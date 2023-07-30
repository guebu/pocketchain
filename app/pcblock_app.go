package main

import (
	"fmt"
	"github.com/guebu/pocketchain/pcchain"
	"github.com/guebu/pocketchain/pcwallet"
)

func main() {

	pcWMiner := pcwallet.NewPCWallet()
	pcWA := pcwallet.NewPCWallet()
	pcWB := pcwallet.NewPCWallet()

	//fmt.Println(fmt.Sprintf("Private Key of PC-Wallet: %s\n", pcW.PrivateKeyStr()))
	//fmt.Println(fmt.Sprintf("Public  Key of PC-Wallet: %s\n", pcW.PublicKeyStr()))
	//fmt.Println(fmt.Sprintf("PC Address  of PC-Wallet: %s\n", pcW.PocketChainAddress()))

	// Let's create a wallet for Person A, who wants to send to person B an amount of 1.0
	walletTrx := pcwallet.NewTransaction(pcWA.PrivateKey(), pcWA.PublicKey(), pcWA.PocketChainAddress(), pcWB.PocketChainAddress(), 1.0)

	// Let's create the Pocket-Chain
	pc := pcchain.NewPocketChain(pcWMiner.PocketChainAddress())

	sig := walletTrx.GenerateSignature()
	// Now let's try to put it to the Block-Chain!
	wasAddedToPC := pc.AddTransaction(pcWA.PocketChainAddress(), pcWB.PocketChainAddress(), 1.0, pcWA.PublicKey(), sig)

	if wasAddedToPC {
		fmt.Println("Wallet-Trx was added to the Pocket Chain!")
		fmt.Printf("Signature of Wallet-Trx: %s \n", sig)
		pc.Mining()
		pc.Print()
	} else {
		fmt.Println("Wallet-Trx was NOT added to the Pocket Chain!")
	}

	fmt.Printf("my %.1f\n", pc.CalculateTotalAmount(pcWMiner.PocketChainAddress()))
	fmt.Printf("A  %.1f\n", pc.CalculateTotalAmount(pcWA.PocketChainAddress()))
	fmt.Printf("B  %.1f\n", pc.CalculateTotalAmount(pcWB.PocketChainAddress()))

	//myPocketChainAddress := "my_pocket_chain_address"

	/*
		pc := pcchain.NewPocketChain(pcW.PocketChainAddress())
		pc.Print()

		pc.AddTransaction("A", "B", 1.5)
		pc.Mining()
		pc.Print()

		pc.AddTransaction("B", "C", 1.0)
		pc.AddTransaction("C", "A", 3.0)
		pc.Mining()
		pc.Print()

		fmt.Printf("my %.1f\n", pc.CalculateTotalAmount(pcW.PocketChainAddress()))
		fmt.Printf("A  %.1f\n", pc.CalculateTotalAmount("A"))
		fmt.Printf("B  %.1f\n", pc.CalculateTotalAmount("B"))
		fmt.Printf("C  %.1f\n", pc.CalculateTotalAmount("C"))

		w := pcwallet.NewPCWallet()
		fmt.Println(w.PrivateKey())
		fmt.Println(w.PublicKey())
		fmt.Println(w.PrivateKeyStr())
		fmt.Println(w.PublicKeyStr())
		fmt.Println(w.PocketChainAddress())

	*/

	/*
		pc.AddTransaction("C", "B", 2.5)
		h, _ = pc.GetLastBlock().GenerateHashForBlock()
		nonce = pc.ProofOfWork()
		pc.CreateBlock(nonce, h)
	*/

}
