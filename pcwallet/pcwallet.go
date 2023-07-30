package pcwallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"pocketchain/pctrx"
)

type PCWallet struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	pcAddress  string
	//walletTrx  *pctrx.PCWalletTransaction
}

func NewPCWallet() *PCWallet {
	// See following link:
	// https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses

	// 1.Creating ECDSA private key (32 bytes) and public key (64 bytes)
	wallet := new(PCWallet)

	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	wallet.privateKey = privateKey

	wallet.publicKey = &wallet.privateKey.PublicKey

	// 2. Perform SHA-256 hashing on pubnlic key (32 bytes)
	h2 := sha256.New()
	h2.Write(wallet.publicKey.Y.Bytes())
	h2.Write(wallet.publicKey.X.Bytes())
	digest2 := h2.Sum(nil)
	// 3. Perform RIPEMD-160 hashing on the result of SHA-256 (20 bytes)
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)
	// 4. Add version byte in front of RIPEDM-160 hash (0x00 for Main Network)
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])

	// 5. Perform SHA-256 hashing on extended RIPEMD-160 result
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)

	// 6. Perform SHA-256 hashing on the result of step-5
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)

	// 7. Take the first 4 Bytes of the secon SHA-256 for checksum
	chksum := digest6[:4]

	// 8. Add the 4 checksum bytes from 7 at the end of extende RIPEMD-160 hash from step 4 (25 bytes)
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[:21], chksum)

	// 9. Convert the result from a byte string into base58
	address := base58.Encode(dc8)
	wallet.pcAddress = address
	return wallet
}

func (pcw *PCWallet) PrivateKey() *ecdsa.PrivateKey {
	return pcw.privateKey
}

func (pcw *PCWallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", pcw.privateKey.D.Bytes())
}

func (pcw *PCWallet) PublicKey() *ecdsa.PublicKey {
	return pcw.publicKey
}

func (pcw *PCWallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", pcw.publicKey.X.Bytes(), pcw.publicKey.Y.Bytes())
}

func (pcw *PCWallet) PocketChainAddress() string {
	return pcw.pcAddress
}

func NewTransaction(senderPrivKey *ecdsa.PrivateKey, senderPubKey *ecdsa.PublicKey, senderPCAddress string, recepientPCAddress string, value float32) *pctrx.PCWalletTransaction {
	wTrx := pctrx.NewPCWalletTransaction(senderPrivKey, senderPubKey, senderPCAddress, recepientPCAddress, value)
	//pcw.walletTrx = wTrx
	return wTrx
}

/*
func (pcw *PCWallet) GetTransaction() *pctrx.PCWalletTransaction {
	return pcw.walletTrx
}
*/

