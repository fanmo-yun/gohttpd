package server

import (
	"gohttp/utils"
	"net/http"
	"net/url"
	"path/filepath"
)

func Handle(response http.ResponseWriter, request *http.Request) *url.URL {
	if request.Method != http.MethodGet {
		SendHTTPErrorResponse(response, http.StatusMethodNotAllowed)
		return nil
	}

	if len(request.UserAgent()) == 0 {
		SendHTTPErrorResponse(response, http.StatusForbidden)
		return nil
	}

	parsedURL, err := url.Parse(request.URL.Path)
	if err != nil {
		SendHTTPErrorResponse(response, http.StatusBadRequest)
		return nil
	}

	return parsedURL
}

func HandleRouter(config *utils.Config) http.HandlerFunc {
	proxies := CreateProxies(config.Proxy)
	backends := NewBackend(config.Backend)
	loadBalancer := NewLoadBalancer(backends)

	return func(res http.ResponseWriter, req *http.Request) {
		urlPath := Handle(res, req)
		if urlPath == nil {
			return
		}
		path := urlPath.Path

		if proxies != nil && FindAndServeProxy(res, req, path, proxies) {
			return
		}

		if backends != nil {
			loadBalancer.ServeHTTP(res, req)
			return
		}

		if !CustomRouter(res, req, path, config.Static.Dirpath, config.Custom) {
			Router(res, req, path, config.Static)
		}
	}
}

func CustomRouter(response http.ResponseWriter, request *http.Request, URLPath string, staticDir string, cus []utils.CustomConfig) bool {
	if len(cus) == 0 {
		return false
	}

	for _, custom := range cus {
		if custom.Urlpath == URLPath {
			fullPath := filepath.Join(staticDir, filepath.Clean(custom.Filepath))
			SendStaticFile(response, request, fullPath)
			return true
		}
	}
	return false
}

func Router(res http.ResponseWriter, req *http.Request, URLPath string, h utils.HtmlConfig) {
	switch URLPath {
	case "/":
		fullPath := filepath.Join(h.Dirpath, h.Index)
		SendStaticFile(res, req, fullPath)
	default:
		fullPath := filepath.Join(h.Dirpath, filepath.Clean(URLPath))
		SendStaticFile(res, req, fullPath)
	}
}
