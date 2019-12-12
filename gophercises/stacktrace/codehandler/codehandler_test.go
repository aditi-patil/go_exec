package codehandler

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeCopy struct {
	err error
}

func (f *fakeCopy) iocopy(p io.Writer, q io.Reader) (int64, error) {
	return 0, f.err
}

var mux = http.NewServeMux()

func TestHello(t *testing.T) {
	t.Run("it gives status 200 when path is right", func(t *testing.T) {
		mux.HandleFunc("/", Hello)
		request, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		handler := DevMiddleware(mux)
		handler.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, 200)
	})
}

func TestSourceCodeHandler(t *testing.T) {
	mux.HandleFunc("/debug/", SourceCodeHandler)
	t.Run("it gives internal server error when path is not present", func(t *testing.T) {
		path := ""
		request, err := http.NewRequest("GET", "/debug/"+path, nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		handler := DevMiddleware(mux)
		handler.ServeHTTP(response, request)
		assert.Equal(t, response.Code, 500)
	})

	t.Run("it gives status 200 when path is present", func(t *testing.T) {
		path := "?line=24&path=/usr/local/go/src/runtime/debug/stack.go"
		request, err := http.NewRequest("GET", "/debug/"+path, nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		handler := DevMiddleware(mux)
		handler.ServeHTTP(response, request)
		assert.Equal(t, response.Code, 200)
	})

	t.Run("it gives internal server error when file failed to copy", func(t *testing.T) {
		f := &fakeCopy{err: errors.New("Failed")}
		ioCopy = f.iocopy
		path := "?line=24&path=/usr/local/go/src/runtime/debug/stack.go"
		request, err := http.NewRequest("GET", "/debug/"+path, nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		handler := DevMiddleware(mux)
		handler.ServeHTTP(response, request)
		assert.Equal(t, response.Code, 500)
	})
}

func TestPanicDemo(t *testing.T) {
	t.Run("it gives status 500 when path is right", func(t *testing.T) {
		mux.HandleFunc("/panic/", PanicDemo)
		request, err := http.NewRequest("GET", "/panic/", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		handler := DevMiddleware(mux)
		handler.ServeHTTP(response, request)
		assert.Equal(t, response.Code, 500)
	})
}
