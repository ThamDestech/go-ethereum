package trie

import (
	"math/big"
	"testing"

	zkt "github.com/scroll-tech/go-ethereum/core/types/zktrie"

	"github.com/iden3/go-iden3-crypto/constants"
	cryptoUtils "github.com/iden3/go-iden3-crypto/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/scroll-tech/go-ethereum/ethdb/memorydb"
)

type Fatalable interface {
	Fatal(args ...interface{})
}

func newTestingMerkle(f Fatalable, numLevels int) *ZkTrieImpl {
	mt, err := NewZkTrieImpl(NewZktrieDatabase((memorydb.New())), numLevels)
	if err != nil {
		f.Fatal(err)
		return nil
	}
	return mt
}

func TestHashParsers(t *testing.T) {
	h0 := zkt.NewHashFromBigInt(big.NewInt(0))
	assert.Equal(t, "0", h0.String())
	h1 := zkt.NewHashFromBigInt(big.NewInt(1))
	assert.Equal(t, "1", h1.String())
	h10 := zkt.NewHashFromBigInt(big.NewInt(10))
	assert.Equal(t, "10", h10.String())

	h7l := zkt.NewHashFromBigInt(big.NewInt(1234567))
	assert.Equal(t, "1234567", h7l.String())
	h8l := zkt.NewHashFromBigInt(big.NewInt(12345678))
	assert.Equal(t, "12345678...", h8l.String())

	b, ok := new(big.Int).SetString("4932297968297298434239270129193057052722409868268166443802652458940273154854", 10) //nolint:lll
	assert.True(t, ok)
	h := zkt.NewHashFromBigInt(b)
	assert.Equal(t, "4932297968297298434239270129193057052722409868268166443802652458940273154854", h.BigInt().String()) //nolint:lll
	assert.Equal(t, "49322979...", h.String())
	assert.Equal(t, "265baaf161e875c372d08e50f52abddc01d32efc93e90290bb8b3d9ceb94e70a", h.Hex())

	b1, err := zkt.NewBigIntFromHashBytes(b.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, new(big.Int).SetBytes(b.Bytes()).String(), b1.String())

	b2, err := zkt.NewHashFromBytes(b.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, b.String(), b2.BigInt().String())

	h2, err := zkt.NewHashFromHex(h.Hex())
	assert.Nil(t, err)
	assert.Equal(t, h, h2)
	_, err = zkt.NewHashFromHex("0x12")
	assert.NotNil(t, err)

	// check limits
	a := new(big.Int).Sub(constants.Q, big.NewInt(1))
	testHashParsers(t, a)
	a = big.NewInt(int64(1))
	testHashParsers(t, a)
}

func testHashParsers(t *testing.T, a *big.Int) {
	require.True(t, cryptoUtils.CheckBigIntInField(a))
	h := zkt.NewHashFromBigInt(a)
	assert.Equal(t, a, h.BigInt())
	hFromBytes, err := zkt.NewHashFromBytes(h.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, h, hFromBytes)
	assert.Equal(t, a, hFromBytes.BigInt())
	assert.Equal(t, a.String(), hFromBytes.BigInt().String())
	hFromHex, err := zkt.NewHashFromHex(h.Hex())
	assert.Nil(t, err)
	assert.Equal(t, h, hFromHex)

	aBIFromHBytes, err := zkt.NewBigIntFromHashBytes(h.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, a, aBIFromHBytes)
	assert.Equal(t, new(big.Int).SetBytes(a.Bytes()).String(), aBIFromHBytes.String())
}

func TestMerkleTree_AddUpdateGetWord(t *testing.T) {
	mt := newTestingMerkle(t, 10)
	err := mt.AddWord(&zkt.Byte32{1}, &zkt.Byte32{2})
	assert.Nil(t, err)
	err = mt.AddWord(&zkt.Byte32{3}, &zkt.Byte32{4})
	assert.Nil(t, err)
	err = mt.AddWord(&zkt.Byte32{5}, &zkt.Byte32{6})
	assert.Nil(t, err)

	node, err := mt.GetLeafNodeByWord(&zkt.Byte32{1})
	assert.Nil(t, err)
	assert.Equal(t, (&zkt.Byte32{2})[:], node.ValuePreimage)
	node, err = mt.GetLeafNodeByWord(&zkt.Byte32{3})
	assert.Nil(t, err)
	assert.Equal(t, (&zkt.Byte32{4})[:], node.ValuePreimage)
	node, err = mt.GetLeafNodeByWord(&zkt.Byte32{5})
	assert.Nil(t, err)
	assert.Equal(t, (&zkt.Byte32{6})[:], node.ValuePreimage)

	err = mt.UpdateWord(&zkt.Byte32{1}, &zkt.Byte32{7})
	assert.Nil(t, err)
	err = mt.UpdateWord(&zkt.Byte32{3}, &zkt.Byte32{8})
	assert.Nil(t, err)
	err = mt.UpdateWord(&zkt.Byte32{5}, &zkt.Byte32{9})
	assert.Nil(t, err)

	node, err = mt.GetLeafNodeByWord(&zkt.Byte32{1})
	assert.Nil(t, err)
	assert.Equal(t, (&zkt.Byte32{7})[:], node.ValuePreimage)
	node, err = mt.GetLeafNodeByWord(&zkt.Byte32{3})
	assert.Nil(t, err)
	assert.Equal(t, (&zkt.Byte32{8})[:], node.ValuePreimage)
	node, err = mt.GetLeafNodeByWord(&zkt.Byte32{5})
	assert.Nil(t, err)
	assert.Equal(t, (&zkt.Byte32{9})[:], node.ValuePreimage)
}