package structurederror_test

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adampresley/structurederror"
	"github.com/adampresley/structurederror/errorparsers/httperrorparser"
	"github.com/stretchr/testify/assert"
)

func TestStructuredError(t *testing.T) {
	w := httptest.NewRecorder()

	http.Error(w, "That didn't go well", http.StatusInternalServerError)
	resp := w.Result()

	logger := slog.New(slog.Default().Handler())

	maker := structurederror.New(
		structurederror.WithDelimiter(";;"),
		structurederror.WithSlog(logger),
	)

	f := func() []structurederror.ErrorArg {
		return []structurederror.ErrorArg{
			{
				Key:   "key1",
				Value: "value1",
			},
			{
				Key:   "key2",
				Value: "value2",
			},
		}
	}

	want := fmt.Errorf("example error ;; age: 123 ;; body: That didn't go well\n ;; key1: value1 ;; key2: value2 ;; statusCode: 500")
	got := maker("example error", "age", 123, httperrorparser.Parse(resp, httperrorparser.WithResponseBody()), f())

	assert.Equal(t, want.Error(), got.Error())
}
