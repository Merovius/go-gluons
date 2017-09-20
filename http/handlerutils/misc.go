package handlerutils

import (
	"io"
	"net/http"
	"os"

	"github.com/prasannavl/mchain"
)

func Nop(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NopHandler() mchain.Handler {
	return mchain.HandlerFunc(Nop)
}

func FileHandler(filePath string, status int) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		w.WriteHeader(status)
		_, err = io.Copy(w, f)
		if err != nil {
			return err
		}
		return nil
	}
	return mchain.HandlerFunc(f)
}