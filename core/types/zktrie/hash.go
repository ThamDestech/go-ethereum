package zktrie

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/iden3/go-iden3-crypto/utils"

	"github.com/scroll-tech/go-ethereum/common"
)

const numCharPrint = 8

// ElemBytesLen is the length of the Hash byte array
const ElemBytesLen = 32

var HashZero = Hash{}

// Hash is the generic type stored in the MerkleTree
type Hash [32]byte

// MarshalText implements the marshaler for the Hash type
func (h Hash) MarshalText() ([]byte, error) {
	return []byte(h.BigInt().String()), nil
}

// UnmarshalText implements the unmarshaler for the Hash type
func (h *Hash) UnmarshalText(b []byte) error {
	ha, err := NewHashFromString(string(b))
	copy(h[:], ha[:])
	return err
}

// String returns decimal representation in string format of the Hash
func (h Hash) String() string {
	s := h.BigInt().String()
	if len(s) < numCharPrint {
		return s
	}
	return s[0:numCharPrint] + "..."
}

// Hex returns the hexadecimal representation of the Hash
func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

// BigInt returns the *big.Int representation of the *Hash
func (h *Hash) BigInt() *big.Int {
	if new(big.Int).SetBytes(ReverseByteOrder(h[:])) == nil {
		return big.NewInt(0)
	}
	return new(big.Int).SetBytes(ReverseByteOrder(h[:]))
}

// Bytes returns the little endian []byte representation of the *Hash, which always is 32
// bytes length.
func (h *Hash) Bytes() []byte {
	b := [32]byte{}
	copy(b[:], h[:])
	return ReverseByteOrder(b[:])
}

// NewBigIntFromHashBytes returns a *big.Int from a byte array, swapping the
// endianness in the process. This is the intended method to get a *big.Int
// from a byte array that previously has ben generated by the Hash.Bytes()
// method.
func NewBigIntFromHashBytes(b []byte) (*big.Int, error) {
	if len(b) != ElemBytesLen {
		return nil, fmt.Errorf("expected 32 bytes, found %d bytes", len(b))
	}
	bi := new(big.Int).SetBytes(b[:ElemBytesLen])
	if !utils.CheckBigIntInField(bi) {
		return nil, fmt.Errorf("NewBigIntFromHashBytes: Value not inside the Finite Field")
	}
	return bi, nil
}

// NewHashFromBigInt returns a *Hash representation of the given *big.Int
func NewHashFromBigInt(b *big.Int) *Hash {
	r := &Hash{}
	copy(r[:], ReverseByteOrder(b.Bytes()))
	return r
}

// NewHashFromBytes returns a *Hash from a byte array, swapping the endianness
// in the process. This is the intended method to get a *Hash from a byte array
// that previously has ben generated by the Hash.Bytes() method.
func NewHashFromBytes(b []byte) (*Hash, error) {
	if len(b) != ElemBytesLen {
		return nil, fmt.Errorf("expected 32 bytes, found %d bytes", len(b))
	}
	var h Hash
	copy(h[:], ReverseByteOrder(b))
	return &h, nil
}

// NewHashFromHex returns a *Hash representation of the given hex string
func NewHashFromHex(h string) (*Hash, error) {
	return NewHashFromBytes(ReverseByteOrder(common.FromHex(h)))
}

// NewHashFromString returns a *Hash representation of the given decimal string
func NewHashFromString(s string) (*Hash, error) {
	bi, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("can not parse string to Hash")
	}
	return NewHashFromBigInt(bi), nil
}

// ReverseByteOrder swaps the order of the bytes in the slice.
func ReverseByteOrder(b []byte) []byte {
	o := make([]byte, len(b))
	for i := range b {
		o[len(b)-1-i] = b[i]
	}
	return o
}
