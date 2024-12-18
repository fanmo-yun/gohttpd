package internal

import (
	"compress/gzip"
	"errors"
	"gohttpd/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

func SendHTTPErrorResponse(response http.ResponseWriter, status int) {
	msg := "gohttpd: " + http.StatusText(status)
	zap.L().Error("gohttpd: HTTP error response", zap.Int("status", status), zap.String("message", msg))
	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response.WriteHeader(status)
	io.WriteString(response, msg)
}

func SendGzipOrFile(response http.ResponseWriter, request *http.Request, path string, gzipTransport bool) {
	if strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") && gzipTransport {
		gz := gzip.NewWriter(response)
		defer gz.Close()

		response.Header().Set("Content-Encoding", "gzip")

		gzipResponseWriter := &GzipResponseWriter{
			ResponseWriter: response,
			Writer:         gz,
		}

		http.ServeFile(gzipResponseWriter, request, path)
	} else {
		http.ServeFile(response, request, path)
	}
}

func SendStaticFile(response http.ResponseWriter, request *http.Request, path string, gzipTransport bool) {
	handleError := func(err error) {
		if errors.Is(err, os.ErrNotExist) {
			zap.L().Warn("gohttpd: File not found", zap.String("path", path))
			SendHTTPErrorResponse(response, http.StatusNotFound)
		} else if os.IsPermission(err) {
			zap.L().Warn("gohttpd: Permission denied", zap.String("path", path))
			SendHTTPErrorResponse(response, http.StatusForbidden)
		} else {
			zap.L().Error("gohttpd: Internal server error", zap.Error(err))
			SendHTTPErrorResponse(response, http.StatusInternalServerError)
		}
	}

	isDir, err := utils.VerifyPath(path)
	if err != nil {
		handleError(err)
		return
	}

	if isDir {
		indexPath := filepath.Join(path, "index.html")
		_, err := utils.VerifyPath(indexPath)
		if err != nil {
			handleError(err)
			return
		}
		path = indexPath
	}

	SendGzipOrFile(response, request, path, gzipTransport)
}

func SendTryRootFile(response http.ResponseWriter, request *http.Request, path string, h utils.HtmlConfig, gzipTransport bool) {
	zap.L().Info("gohttpd: Serving file with try-root strategy", zap.String("path", path), zap.String("request_url", request.URL.String()), zap.String("dirpath", h.Dirpath), zap.String("index", h.Index))
	handleError := func(err error) {
		if errors.Is(err, os.ErrNotExist) {
			fullpath := filepath.Join(h.Dirpath, h.Index)
			zap.L().Warn("gohttpd: File not found, trying root file", zap.String("path", path), zap.String("fallback_path", fullpath))
			http.ServeFile(response, request, fullpath)
		} else if os.IsPermission(err) {
			zap.L().Warn("gohttpd: Permission denied", zap.String("path", path))
			SendHTTPErrorResponse(response, http.StatusForbidden)
		} else {
			zap.L().Error("gohttpd: Internal server error", zap.Error(err))
			SendHTTPErrorResponse(response, http.StatusInternalServerError)
		}
	}

	_, err := utils.VerifyPath(path)
	if err != nil {
		handleError(err)
		return
	}

	SendGzipOrFile(response, request, path, gzipTransport)
}
