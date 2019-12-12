package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

var cipherBlock = NewCipherBlock
var NewEncryptStream = EncryptStream
var NewDecryptStream = DecryptStream
var IoRead = io.ReadFull
var IoWrite = io.Writer.Write
var NewDecryptReader = DecryptReader

// EncryptStream gives encrypted stream for the given key
func EncryptStream(key string, iv []byte) (cipher.Stream, error) {
	block, err := cipherBlock(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewCFBEncrypter(block, iv), nil
}

// EncryptWriter will return a writer that will write encrypted data to
// the original writer.
func EncryptWriter(key string, w io.Writer) (*cipher.StreamWriter, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := IoRead(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream, err := NewEncryptStream(key, iv)
	if err != nil {
		return nil, err
	}
	n, err := IoWrite(w, iv)
	if n != len(iv) || err != nil {
		return nil, errors.New("encrypt: unable to write full iv to writer")
	}
	return &cipher.StreamWriter{S: stream, W: w}, nil
}

// DecryptStream gives decrypted stream for the given key
func DecryptStream(key string, iv []byte) (cipher.Stream, error) {
	block, err := cipherBlock(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewCFBDecrypter(block, iv), nil
}

// NewCipherBlock creates and returns a new cipher.Block.
// The key argument should be the AES key, either 16, 24, or 32 bytes.
func NewCipherBlock(key string) (cipher.Block, error) {
	h := md5.New()
	fmt.Fprint(h, key)
	cipherkey := h.Sum(nil)
	return aes.NewCipher(cipherkey)
}

// DecryptReader will return a reader that will decrypt data from the
// provided reader and give the user a way to read that data as it if was
// not encrypted.
func DecryptReader(key string, r io.Reader) (*cipher.StreamReader, error) {
	iv := make([]byte, aes.BlockSize)
	n, err := IoRead(r, iv)
	if n < len(iv) || err != nil {
		return nil, errors.New("decrypt: unable to read the full iv")
	}
	stream, err := NewDecryptStream(key, iv)
	if err != nil {
		return nil, err
	}
	
	return &cipher.StreamReader{S: stream, R: r}, nil
}
