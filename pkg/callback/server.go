package callback

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
)

func oidcCallbackHandlerFactory(cancelFunc context.CancelFunc, result chan CallbackResult) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("got request", slog.String("url", r.URL.String()), slog.String("method", r.Method))

		params, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			slog.Error("error parsing request query", slog.String("error", err.Error()))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			cancelFunc()
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("<html><script type=\"text/javascript\">window.close()</script></html>"))
		if err != nil {
			slog.Error("writing response", slog.String("error", err.Error()))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			cancelFunc()
		}
		resultData := CallbackResult{
			Code:  params["code"][0],
			State: params["state"][0],
		}
		result <- resultData

	}
}

func StartHttpServer(cancel context.CancelFunc, port int, callbackPath string, result chan CallbackResult) func() {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	mux := http.NewServeMux()
	mux.HandleFunc(callbackPath, oidcCallbackHandlerFactory(cancel, result))
	srv.Handler = mux

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("error running http server", slog.String("error", err.Error()))
			cancel()
		}
	}()

	return func() {
		err := srv.Close()
		if err != nil {
			slog.Error("error closing http server", slog.String("error", err.Error()))
		}
	}

}
