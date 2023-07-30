package pcutil

import (
	"fmt"
	"math/big"
)

type PCSignature struct {
	R *big.Int
	S *big.Int
}

func (s *PCSignature) String() string {
	return fmt.Sprintf("%x%x", s.R, s.S)
}
