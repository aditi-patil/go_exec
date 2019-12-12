package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

type fakeCipher struct {
	number int
	err    error
}

func (m *fakeCipher) newCipherBlock(key string) (cipher.Block, error) {
	return nil, m.err
}

func (m *fakeCipher) encryptStream(key string, iv []byte) (cipher.Stream, error) {
	return nil, m.err
}

func (m *fakeCipher) decryptStream(key string, iv []byte) (cipher.Stream, error) {
	return nil, m.err
}

func (m *fakeCipher) ioread(r io.Reader, buf []byte) (int, error) {
	return m.number, m.err
}

func (m *fakeCipher) iowrite(w io.Writer, buf []byte) (int, error) {
	return m.number, m.err
}

func TestNewCipherBlock(t *testing.T) {
	_, err := cipherBlock("demo_test_key")
	if err != nil {
		log.Fatal(err)
	}
}

func TestDecryptStream(t *testing.T) {
	iv := make([]byte, aes.BlockSize)

	t.Run("it gives an error for invalid key", func(t *testing.T) {
		f := &fakeCipher{err: errors.New("crypto/aes: invalid key size 10")}
		cipherBlock = f.newCipherBlock
		_, err := NewDecryptStream("6368616e67", iv)
		if f.err != err {
			t.Errorf("Expected : %v\n Got: %v", f.err, err)
		}
	})

	t.Run("it returns decrypt stream", func(t *testing.T) {
		cipherBlock = NewCipherBlock
		stream, err := NewDecryptStream("test_key", iv)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(stream)
	})

}

func TestEncryptStream(t *testing.T) {
	iv := make([]byte, aes.BlockSize)

	t.Run("it return encrypted stream", func(t *testing.T) {
		stream, err := NewEncryptStream("test_key", iv)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(stream)
	})

	t.Run("it gives an error for invalid key", func(t *testing.T) {
		f := &fakeCipher{err: errors.New("crypto/aes: invalid key size 10")}
		cipherBlock = f.newCipherBlock
		_, err := NewEncryptStream("6368616e67", iv)
		if f.err != err {
			t.Errorf("Expected : %v\n Got: %v", f.err, err)
		}
	})
}

func TestEncryptWriter(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "secrets")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	os.Stdin = tmpfile

	t.Run("it writes key value to file for valid key", func(t *testing.T) {
		cipherBlock = NewCipherBlock
		w, error := EncryptWriter("test_key", tmpfile)
		if error != nil {
			log.Fatal(error)
		}
		fmt.Println(w)
	})

	t.Run("it returns error if encrypt stream fails", func(t *testing.T) {
		f := &fakeCipher{err: errors.New("crypto/aes: invalid key size 10")}
		NewEncryptStream = f.encryptStream
		_, err := EncryptWriter("test_key", tmpfile)
		if f.err != err {
			t.Errorf("Expected : %v\n Got: %v", f.err, err)
		}
	})

	t.Run("it returns error for reading full iv", func(t *testing.T) {
		NewEncryptStream = EncryptStream
		f := &fakeCipher{err: errors.New("encrypt: unable to write full iv to writer"), number: 2}
		IoWrite = f.iowrite
		_, err := EncryptWriter("test_key", tmpfile)
		if f.err.Error() != err.Error() {
			t.Errorf("Expected : %v\n Got: %v", f.err, err)
		}
	})

	t.Run("it returns error if there is unexpected EOF", func(t *testing.T) {
		f := &fakeCipher{err: errors.New("error: unexpected EOF")}
		IoRead = f.ioread
		_, err := EncryptWriter("test_key", tmpfile)
		if f.err != err {
			t.Errorf("Expected : %v\n Got: %v", f.err, err)
		}
	})

	if err := tmpfile.Close(); err != nil {
		fmt.Println(err)
	}

}

func TestDecryptReader(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "secrets")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())
	os.Stdin = tmpfile

	t.Run("it returns decrypted stream reader", func(t *testing.T) {
		f := &fakeCipher{err: nil, number: 21}
		IoRead = f.ioread
		w, error := DecryptReader("test_key", tmpfile)
		if error != nil {
			log.Fatal(error)
		}
		fmt.Println(w)
	})

	t.Run("it gives an error for invalid key", func(t *testing.T) {
		f := &fakeCipher{err: nil, number: 21}
		IoRead = f.ioread
		s := &fakeCipher{err: errors.New("crypto/aes: invalid key size 10")}
		NewDecryptStream = s.decryptStream
		_, err := DecryptReader("6368616e67", tmpfile)
		if s.err != err {
			t.Errorf("\nExpected : %v\n Got: %v", s.err, err)
		}
	})

	t.Run("it returns error for reading full iv", func(t *testing.T) {
		f := &fakeCipher{err: errors.New("decrypt: unable to read the full iv"), number: 2}
		IoRead = f.ioread
		_, err := DecryptReader("test_key", tmpfile)
		if f.err.Error() != err.Error() {
			t.Errorf("Expected : %v\n Got: %v", f.err, err)
		}
	})

	// Close the file
	if err := tmpfile.Close(); err != nil {
		fmt.Println(err)
	}

}
