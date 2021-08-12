package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type codeWrapper struct {
	code int
	http.ResponseWriter
}

func (w *codeWrapper) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

type proxyHandler struct {
	config Config
	logger *logrus.Logger
	hc     *http.Client
}

func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	resp := codeWrapper{ResponseWriter: w}
	var err error

	defer r.Body.Close()
	defer func() {
		if err != nil || h.logger.Level >= logrus.DebugLevel {
			fields := logrus.Fields{
				"method":        r.Method,
				"ellapsed":      time.Since(start).String(),
				"path":          r.URL.RequestURI(),
				"proxyEndpoint": h.config.Proxy.Endpoint,
				"response":      resp.code,
			}
			for _, header := range h.config.Proxy.LogHeaders {
				if value := r.Header.Get(header); value != "" {
					fields["ReqHeader/"+header] = value
				}
			}
			entry := h.logger.WithFields(fields)
			if err != nil {
				entry.WithError(err).Error("failed to handle request")
			} else {
				entry.Debug("finished handling request")
			}
		}
	}()

	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// If GCS_HELPER_ROOT_REDIRECT is set, '/' should redirect to the given path
	if r.URL.Path == "/" && h.config.Proxy.RootRedirect != "" {
		r.URL.Path = h.config.Proxy.RootRedirect
	}

	if r.URL.Path == "/" {
		w.WriteHeader(http.StatusOK)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), h.config.Proxy.Timeout)
	defer cancel()

	host := "storage.googleapis.com"
	if !h.config.Proxy.BucketOnPath {
		host = h.config.BucketName + "." + host
	}
	url := fmt.Sprintf("https://%s%s", host, r.URL.RequestURI())
	// no support for request body, do we care? :)
	gcsReq, err := http.NewRequest(r.Method, url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	gcsReq = gcsReq.WithContext(ctx)
	for name, values := range r.Header {
		for _, value := range values {
			gcsReq.Header.Add(name, value)
		}
	}
	gcsResp, err := h.hc.Do(gcsReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer gcsResp.Body.Close()

	for name, values := range gcsResp.Header {
		for _, value := range values {
			resp.Header().Add(name, value)
		}
	}
	resp.WriteHeader(gcsResp.StatusCode)
	io.Copy(resp, gcsResp.Body)
}

// Proxy returns the proxy handler.
func Proxy(c Config, hc *http.Client) http.Handler {
	logger := c.Logger()
	return &proxyHandler{logger: logger, hc: hc, config: c}
}
