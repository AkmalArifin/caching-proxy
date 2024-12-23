package proxy

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/AkmalArifin/caching-proxy/internal/cache"
)

type ProxyObject struct {
	Origin string
	Cache  map[string]*cache.CacheObject
	Mutex  sync.RWMutex
}

func NewProxy(origin string) *ProxyObject {
	return &ProxyObject{
		Origin: origin,
		Cache:  make(map[string]*cache.CacheObject),
	}
}

func (p *ProxyObject) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	CACHE_KEY := fmt.Sprintf("%v:%v", r.Method, r.URL)

	// if cache is found
	p.Mutex.RLock()
	if c, ok := p.Cache[CACHE_KEY]; ok {
		p.Mutex.RUnlock()
		ResponseWithHeader(w, *c.Response, c.ResponseBody, "HIT", CACHE_KEY)
		return
	}
	p.Mutex.RUnlock()

	originURL := p.Origin + r.URL.String()
	resp, err := http.Get(originURL)
	if err != nil {
		fmt.Fprintf(w, "cannot forwarding request to the origin url %v", originURL)
		return
	}
	defer resp.Body.Close()

	var newCache cache.CacheObject
	newCache.Response = resp
	newCache.ResponseBody, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "cannot read response body")
		return
	}

	p.Cache[CACHE_KEY] = &newCache
	ResponseWithHeader(w, *newCache.Response, newCache.ResponseBody, "MISS", CACHE_KEY)
}

func ResponseWithHeader(w http.ResponseWriter, response http.Response, body []byte, cacheHeader, key string) {
	fmt.Printf("Cache: %s %s\n", key, cacheHeader)
	w.Header().Set("X-Cache", cacheHeader)
	w.WriteHeader(response.StatusCode)
	for k, v := range response.Header {
		w.Header()[k] = v
	}
	w.Write(body)
}
