package data

import (
	"io"
	"net/http"
	"net/textproto"

	"github.com/BBVA/kapow/internal/server/model"
	"github.com/gorilla/mux"
)

func getRequestBody(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	n, err := io.Copy(w, h.Request.Body)
	if err != nil {
		if n == 0 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			// Only way to abort current connection as of go 1.13
			// https://github.com/golang/go/issues/16542
			panic("Truncated body")
		}
	}
}

func getRequestMethod(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	_, _ = w.Write([]byte(h.Request.Method))
}

func getRequestHost(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	_, _ = w.Write([]byte(h.Request.Host))
}

func getRequestPath(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	// TODO: Discuss a how to obtain URL.EscapedPath() instead
	_, _ = w.Write([]byte(h.Request.URL.Path))
}

func getRequestMatches(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	name := mux.Vars(r)["name"]
	vars := mux.Vars(h.Request)
	if value, ok := vars[name]; ok {
		_, _ = w.Write([]byte(value))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func getRequestParams(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	name := mux.Vars(r)["name"]
	if values, ok := h.Request.URL.Query()[name]; ok {
		_, _ = w.Write([]byte(values[0]))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func getRequestHeaders(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	name := mux.Vars(r)["name"]
	if values, ok := h.Request.Header[textproto.CanonicalMIMEHeaderKey(name)]; ok {
		_, _ = w.Write([]byte(values[0]))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// TODO: Add to the note section of the specification the fact that
// Cookie keys are CaseSensitive
func getRequestCookies(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	name := mux.Vars(r)["name"]
	if cookie, err := h.Request.Cookie(name); err == nil {
		_, _ = w.Write([]byte(cookie.Value))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
