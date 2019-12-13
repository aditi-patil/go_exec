package api

import (
	"bytes"
	"errors"
	"gophercises/transform/primitive"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var router = mux.NewRouter()

type fakefile struct {
	err error
}

func (f *fakefile) tmpfile(prefix, ext string) (*os.File, error) {
	return nil, f.err
}

func (f *fakefile) copy(w io.Writer, r io.Reader) (int64, error) {
	return 0, f.err
}

func (f *fakefile) iotmpfile(dir, pattern string) (*os.File, error) {
	return nil, f.err
}

func (f *fakefile) transform(image io.Reader, ext string, numShapes int, opts ...func() []string) (io.Reader, error) {
	return bytes.NewBuffer(nil), f.err
}

func TestModifyImage(t *testing.T) {
	router.HandleFunc("/modify/{id}", ModifyImage)
	t.Run("check when primitive transform not have error", func(t *testing.T) {
		f := &fakefile{err: nil}
		primitive.NewTransform = f.transform

		t.Run("it displays images with different shapes if mode is not provided", func(t *testing.T) {
			request, err := http.NewRequest("GET", "/modify/test_image.png", nil)
			request.Header.Add("Accept", "*/*")
			if err != nil {
				t.Fatal(err)
			}
			router.HandleFunc("/modify/{id}", ModifyImage)

			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			assert.Equal(t, 200, response.Code, "StatusOk found")
		})

		t.Run("it displays images with different num shapes if number is not provided", func(t *testing.T) {
			request, err := http.NewRequest("GET", "/modify/test_image.png?mode=2", nil)
			request.Header.Add("Accept", "*/*")
			if err != nil {
				t.Fatal(err)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			assert.Equal(t, 200, response.Code, "StatusOk found")
		})

		t.Run("it failes if image is invalid", func(t *testing.T) {
			payload := strings.NewReader("------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"mode\"\r\n\r\n\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--")

			request, err := http.NewRequest("GET", "/modify/invalid.png", payload)
			request.Header.Add("Accept", "*/*")
			if err != nil {
				t.Fatal(err)
			}
			response := httptest.NewRecorder()
			router.HandleFunc("/modify/{id}", ModifyImage)
			router.ServeHTTP(response, request)
			assert.Equal(t, 400, response.Code, "BadRequest")
		})

		t.Run("it fails if mode is invalid", func(t *testing.T) {
			request, err := http.NewRequest("GET", "/modify/test_image.png?mode=invalid", nil)
			request.Header.Add("Accept", "*/*")
			if err != nil {
				t.Fatal(err)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			assert.Equal(t, 400, response.Code, "BadRequest")
		})

		t.Run("it fails if numshape is invalid", func(t *testing.T) {
			request, err := http.NewRequest("GET", "/modify/test_image.png?mode=2&n=invalid", nil)
			request.Header.Add("Accept", "*/*")
			if err != nil {
				t.Fatal(err)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			assert.Equal(t, 400, response.Code, "BadRequest")
		})

		t.Run("it redirects to image path", func(t *testing.T) {
			request, err := http.NewRequest("GET", "/modify/test_image.png?mode=2&n=10", nil)
			request.Header.Add("Accept", "*/*")
			if err != nil {
				t.Fatal(err)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			assert.Equal(t, 302, response.Code, "Status Found")
		})
	})

	t.Run("check when primitive transform fails", func(t *testing.T) {
		f := &fakefile{err: errors.New("Failed")}
		primitive.NewTransform = f.transform
		request, err := http.NewRequest("GET", "/modify/test_image.png?mode=2", nil)
		request.Header.Add("Accept", "*/*")
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assert.Equal(t, 500, response.Code, "Internal Server Error")
	})

	t.Run("it fails if tmp file is not present", func(t *testing.T) {
		s := &fakefile{err: nil}
		primitive.NewTransform = s.transform
		q := &fakefile{err: errors.New("main: failed to create temporary file")}
		newTempfile = q.tmpfile
		request, err := http.NewRequest("GET", "/modify/test_image.png", nil)
		request.Header.Add("Accept", "*/*")
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assert.Equal(t, 500, response.Code, "Internal Server Error")
	})

	defer func() {
		primitive.NewTransform = primitive.Transform
		newTempfile = tempfile
	}()

}

func TestUploadImage(t *testing.T) {
	router.HandleFunc("/upload", UploadImage)

	t.Run("it fails to redirect if image is not present", func(t *testing.T) {
		request, err := http.NewRequest("POST", "/upload", nil)
		request.Header.Add("content-type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
		request.Header.Add("Accept", "*/*")
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assert.Equal(t, 400, response.Code, "Bad Request")
	})

	t.Run("it checks when payload is present", func(t *testing.T) {
		t.Run("it redirect to modify page if image is uploaded", func(t *testing.T) {
			payload := strings.NewReader("------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"image\"; filename=\"test.png\"\r\nContent-Type: image/png\r\n\r\n\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--")
			request, err := http.NewRequest("POST", "/upload", payload)
			request.Header.Add("content-type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
			request.Header.Add("Accept", "*/*")
			if err != nil {
				t.Fatal(err)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			assert.Equal(t, 302, response.Code, "Redirect url found")
		})

		t.Run("it fails to redirect if file is not copied", func(t *testing.T) {
			f := &fakefile{err: errors.New("Failed to copy file")}
			ioCopy = f.copy
			payload := strings.NewReader("------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"image\"; filename=\"test.png\"\r\nContent-Type: image/png\r\n\r\n\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--")
			request, err := http.NewRequest("POST", "/upload", payload)
			request.Header.Add("content-type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
			request.Header.Add("Accept", "*/*")
			if err != nil {
				t.Fatal(err)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			assert.Equal(t, 500, response.Code, "Internal Server Error")
		})

		t.Run("it fails to redirect if given tmp file is path not present", func(t *testing.T) {
			f := &fakefile{err: errors.New("Tempfile path directory is not found")}
			ioTempFile = f.iotmpfile
			payload := strings.NewReader("------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"image\"; filename=\"test.png\"\r\nContent-Type: image/png\r\n\r\n\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--")
			request, err := http.NewRequest("POST", "/upload", payload)
			request.Header.Add("content-type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
			request.Header.Add("Accept", "*/*")
			if err != nil {
				t.Fatal(err)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			assert.Equal(t, 500, response.Code, "Internal Server Error")
		})

		t.Run("it fails to redirect if tmp file is not present", func(t *testing.T) {
			f := &fakefile{err: errors.New("main: failed to create temporary file")}
			newTempfile = f.tmpfile
			payload := strings.NewReader("------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"image\"; filename=\"test.png\"\r\nContent-Type: image/png\r\n\r\n\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--")
			request, err := http.NewRequest("POST", "/upload", payload)
			request.Header.Add("content-type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
			request.Header.Add("Accept", "*/*")
			if err != nil {
				t.Fatal(err)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			assert.Equal(t, 500, response.Code, "Internal Server Error")
		})

	})

	defer func() {
		ioCopy = io.Copy
		newTempfile = tempfile
		ioTempFile = ioutil.TempFile
	}()

}
