package crypt

import (
	"crypto/cipher"
	"fmt"
	"github.com/thecloakproject/utils"
)

func AESEncryptBytes(block cipher.Block, plain []byte) (cipherBytes []byte, err error) {
	blockSize := block.BlockSize()
	plain = utils.PadBytes(plain, blockSize)
	length := len(plain)

	// Encrypt
	cipherBytes = make([]byte, length)
	for i := 0; i < length; i += blockSize {
		block.Encrypt(cipherBytes[i:i+blockSize], plain[i:i+blockSize])
	}

	return
}

func AESDecryptBytes(block cipher.Block, cipherBytes []byte) (plain []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			plain = nil
			err = fmt.Errorf("%v", e)
		}
	}()

	blockSize := block.BlockSize()
	cipherBytes = utils.PadBytes(cipherBytes, blockSize)
	length := len(cipherBytes)

	// Decrypt
	plain = make([]byte, length)
	for i := 0; i < length; i += blockSize {
		block.Decrypt(plain[i:i+blockSize], cipherBytes[i:i+blockSize])
	}

	return
}
