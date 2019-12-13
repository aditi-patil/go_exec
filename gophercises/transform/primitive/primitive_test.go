package primitive

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakefile struct {
	err error
}

func (f *fakefile) copy(w io.Writer, r io.Reader) (int64, error) {
	return 0, f.err
}

func (f *fakefile) mockIoTmpfile(prefix, ext string) (*os.File, error) {
	return nil, f.err
}

func (f *fakefile) mockPrimitive(inFile, outFile string, numShapes int, args ...string) (string, error) {
	return "", f.err
}

func TestTransform(t *testing.T) {
	validImage, _ := os.Open("../api/img/test_image.png")
	invalidImage, _ := os.Open("./img/invalid.png")

	t.Run("it returns bytes buffer for valid image", func(t *testing.T) {
		out, _ := Transform(validImage, ".png", 12, WithMode(ModeCircle))
		assert.Equal(t, reflect.TypeOf(out).String(), "*bytes.Buffer")
	})

	t.Run("it returns error if primitive is failed to run", func(t *testing.T) {
		f := &fakefile{err: errors.New("Failed")}
		newPrimitive = f.mockPrimitive
		_, err := Transform(validImage, ".png", 12, WithMode(ModeCircle))
		assert.NotEqual(t, nil, err.Error())

		defer func() {
			newPrimitive = primitive
		}()
	})

	t.Run("it returns error for invalid image", func(t *testing.T) {
		out, err := Transform(invalidImage, ".png", 12, WithMode(ModeCircle))
		assert.Equal(t, out, nil)
		assert.Equal(t, "primitive: failed to copy image into temp input file", err.Error())
	})

	t.Run("it returns error if failed to create temp file", func(t *testing.T) {
		f := &fakefile{err: errors.New("main: failed to create temporary file")}
		ioTempFile = f.mockIoTmpfile
		_, err := Transform(validImage, ".png", 12, WithMode(ModeCircle))
		assert.Equal(t, "primitive: failed to create temporary input file", err.Error())
		defer func() {
			ioTempFile = ioutil.TempFile
		}()
	})

	t.Run("it returns error if primitive is failed to copy output file", func(t *testing.T) {
		f := &fakefile{err: nil}
		newPrimitive = f.mockPrimitive
		s := &fakefile{err: errors.New("Failed to copy file")}
		ioCopy = s.copy
		_, err := Transform(validImage, ".png", 12, WithMode(ModeCircle))
		assert.Equal(t, "primitive: Failed to copy output file into byte buffer", err.Error())

		defer func() {
			newPrimitive = primitive
			ioCopy = io.Copy
		}()
	})

}
