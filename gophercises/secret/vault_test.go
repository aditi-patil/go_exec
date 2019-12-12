package secret

import (
	"crypto/cipher"
	"errors"
	"fmt"
	Cipher "gophercises/secret/cipher"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeCipher struct {
	number int
	err    error
}

func (m *fakeCipher) ioread(r io.Reader, buf []byte) (int, error) {
	return m.number, m.err
}

func (m *fakeCipher) iowrite(w io.Writer, buf []byte) (int, error) {
	return m.number, m.err
}

func (m *fakeCipher) decryptStream(key string, iv []byte) (cipher.Stream, error) {
	return nil, m.err
}

func getFilepath() string {
	tmpfile, err := ioutil.TempFile(os.TempDir(), "temp")
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(tmpfile.Name())
	defer os.Remove(tmpfile.Name())
	os.Stdin = tmpfile

	return path
}

func InitFile() *Vault {
	return &Vault{
		encodingKey: "testencodingKey",
		filepath:    getFilepath(),
	}
}

func TestGet(t *testing.T) {
	t.Run("it gives value for the key from the file", func(t *testing.T) {
		tmpfile, err := ioutil.TempFile(os.TempDir(), "temp")
		if err != nil {
			log.Fatal(err)
		}
		path := filepath.Join(tmpfile.Name())
		defer os.Remove(tmpfile.Name())
		if _, err := tmpfile.Seek(0, 0); err != nil {
			log.Fatal(err)
		}
		os.Stdin = tmpfile

		v := &Vault{
			encodingKey: "testencodingKey",
			filepath:    path,
		}
		_, err = Cipher.EncryptWriter(v.encodingKey, tmpfile)

		er := v.Set("test_key", "testkeyvalue")
		if er != nil {
			log.Fatal(er)
		}
		value, _ := v.Get("test_key")
		assert.Equal(t, value, "testkeyvalue")
	})

	t.Run("it does not give result if no value found", func(t *testing.T) {
		_, err := InitFile().Get("test_key2")
		expected := "secret: no value for that key"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("it returns error if unable to load file", func(t *testing.T) {
		v := InitFile()
		er := v.Set("test_key", "testkeyvalue")
		if er != nil {
			log.Fatal(er)
		}

		f := &fakeCipher{err: errors.New("decrypt: unable to read the full iv"), number: 2}
		Cipher.IoRead = f.ioread
		_, err := v.Get("test_key2")
		if f.err.Error() != err.Error() {
			t.Errorf("Expected : %v\n Got: %v", f.err, err)
		}
	})

	defer func() {
		Cipher.IoRead = io.ReadFull
	}()
}

func TestRemove(t *testing.T) {
	t.Run("it removes key and value pair from the file", func(t *testing.T) {
		tmpfile, err := ioutil.TempFile(os.TempDir(), "temp")
		if err != nil {
			log.Fatal(err)
		}
		path := filepath.Join(tmpfile.Name())
		if _, err := tmpfile.Seek(0, 0); err != nil {
			log.Fatal(err)
		}
		os.Stdin = tmpfile

		v := &Vault{
			encodingKey: "testencodingKey",
			filepath:    path,
		}
		_, err = Cipher.EncryptWriter(v.encodingKey, tmpfile)

		er := v.Set("test_key", "testkeyvalue")
		if er != nil {
			log.Fatal(er)
		}
		e := v.Remove("test_key")
		if e != nil {
			log.Fatal(e)
		}

	})

	t.Run("it gives an error for invalid key", func(t *testing.T) {
		f := &fakeCipher{err: nil, number: 21}
		Cipher.IoRead = f.ioread
		s := &fakeCipher{err: errors.New("crypto/aes: invalid key size 10")}
		Cipher.NewDecryptStream = s.decryptStream
		tmpfile, err := ioutil.TempFile(os.TempDir(), "temp")
		if err != nil {
			log.Fatal(err)
		}
		path := filepath.Join(tmpfile.Name())
		if _, err := tmpfile.Seek(0, 0); err != nil {
			log.Fatal(err)
		}
		os.Stdin = tmpfile

		v := &Vault{
			encodingKey: "testencodingKey",
			filepath:    path,
		}
		err = v.Remove("test_key")
		assert.Equal(t, s.err, err)
	})

	t.Run("it does not give result if no value found", func(t *testing.T) {
		err := InitFile().Remove("test_key2")
		expected := "secret: no value for that key"
		assert.Equal(t, expected, err.Error())
	})

	defer func() {
		Cipher.IoRead = io.ReadFull
		Cipher.NewDecryptStream = Cipher.DecryptStream
	}()
}

func TestSet(t *testing.T) {
	t.Run("it saves key values in the file", func(t *testing.T) {
		tmpfile, err := ioutil.TempFile(os.TempDir(), "temp")
		if err != nil {
			log.Fatal(err)
		}
		path := filepath.Join(tmpfile.Name())
		if _, err := tmpfile.Seek(0, 0); err != nil {
			log.Fatal(err)
		}
		os.Stdin = tmpfile

		v := &Vault{
			encodingKey: "testencodingKey",
			filepath:    path,
		}
		_, err = Cipher.EncryptWriter(v.encodingKey, tmpfile)

		er := v.Set("test_key", "testkeyvalue")
		if er != nil {
			log.Fatal(er)
		}

		// Close the file
		if err := tmpfile.Close(); err != nil {
			fmt.Println(err)
		}

	})

	t.Run("it returns error if unable to load file", func(t *testing.T) {
		f := &fakeCipher{err: errors.New("decrypt: unable to read the full iv"), number: 2}
		Cipher.IoRead = f.ioread
		err := InitFile().Set("test_key", "testsetfunc")
		if f.err.Error() != err.Error() {
			t.Errorf("Expected : %v\n Got: %v", f.err, err)
		}
	})

	t.Run("it gives an error for invalid key", func(t *testing.T) {
		f := &fakeCipher{err: nil, number: 21}
		Cipher.IoRead = f.ioread
		s := &fakeCipher{err: errors.New("crypto/aes: invalid key size 10")}
		Cipher.NewDecryptStream = s.decryptStream
		tmpfile, err := ioutil.TempFile(os.TempDir(), "temp")
		if err != nil {
			log.Fatal(err)
		}
		path := filepath.Join(tmpfile.Name())
		if _, err := tmpfile.Seek(0, 0); err != nil {
			log.Fatal(err)
		}
		os.Stdin = tmpfile

		v := &Vault{
			encodingKey: "testencodingKey",
			filepath:    path,
		}
		err = v.Set("test_key", "setkeyfunc")
		assert.Equal(t, s.err, err)
	})

	defer func() {
		Cipher.IoRead = io.ReadFull
		Cipher.NewDecryptStream = Cipher.DecryptStream
	}()
}

func TestSave(t *testing.T) {

	t.Run("it fails when file not present", func(t *testing.T) {
		v := &Vault{
			encodingKey: "testencodingKey",
			filepath:    "",
		}
		err := v.Save()
		expected := "open : no such file or directory"
		assert.Equal(t, expected, err.Error())
	})
}

func TestLoad(t *testing.T) {

	t.Run("it fails when file not present", func(t *testing.T) {
		v := &Vault{
			encodingKey: "testencodingKey",
			filepath:    "",
		}
		err := v.Load()
		if err != nil {
			log.Fatal(err)
		}
	})

	t.Run("it gives an error for invalid key", func(t *testing.T) {
		f := &fakeCipher{err: nil, number: 21}
		Cipher.IoRead = f.ioread
		s := &fakeCipher{err: errors.New("crypto/aes: invalid key size 10")}
		Cipher.NewDecryptStream = s.decryptStream
		tmpfile, err := ioutil.TempFile(os.TempDir(), "temp")
		if err != nil {
			log.Fatal(err)
		}
		path := filepath.Join(tmpfile.Name())
		if _, err := tmpfile.Seek(0, 0); err != nil {
			log.Fatal(err)
		}
		os.Stdin = tmpfile

		v := &Vault{
			encodingKey: "testencodingKey",
			filepath:    path,
		}

		err = v.Load()
		assert.Equal(t, s.err, err)
	})

}

func TestFile(t *testing.T) {
	v := File("testkey", "/tmp/test")
	assert.Equal(t, v.encodingKey, "testkey")
}
