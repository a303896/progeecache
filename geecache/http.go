package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_geecache/"

type HttpPool struct {
	self string
	basePath string
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		self: self,
		basePath: defaultBasePath,
	}
}

func (p *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HttpPool) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	if !strings.HasPrefix(req.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + req.URL.Path)
	}
	p.Log("%s %s", req.Method, req.URL.Path)
	parts := strings.SplitN(req.URL.Path[len(p.basePath):], "/", 2)
	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view,err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//w.Header().Set("Content-type", "application/octet-stream")
	w.Write(view.b)
}