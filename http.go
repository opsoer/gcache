package gcache

import (
	"fmt"
	"gcache/consistenthash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/gcache"
	defaultReplicas = 30
)

// HTTPPool 为HTTP对等体池实现PeerPicker。
type HTTPPool struct {
	self        string //自己的url+port
	basePath    string
	mu          sync.Mutex // 防止并发访问peers和httpGetters
	peers       *consistenthash.Map
	httpGetters map[string]*httpGetter
}

// NewHTTPPool 初始化HTTP对等体池。
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		//自己的IP
		self:     self,
		basePath: defaultBasePath,
	}
}

// ServeHTTP http服务器
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		log.Panicf("HTTPPool serving unexpected path: %s, expect %s\n", r.URL.Path, p.basePath)
	}
	log.Printf("[Server %s] %s\n", p.self, fmt.Sprintf("%s %s", r.Method, r.URL.Path))
	// url:port/<basepath>/<key>
	key := r.URL.Path[len(p.basePath)+1:]

	peer, ok := p.PickPeer(key)
	if !ok {
		log.Fatalf("PickPeer %s not find peer\n", key)
	}
	view, err := peer.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view)
}

// AddPeers 将节点虚拟化多个并且放入HTTPPool，peer为ip+port
func (p *HTTPPool) AddPeers(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.peers == nil {
		p.peers = consistenthash.New(defaultReplicas, nil)
	}
	p.peers.AddPeers(peers...)
	if p.httpGetters == nil {
		p.httpGetters = make(map[string]*httpGetter)
	}
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

// PickPeer 根据key选择对等体
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.GetPeer(key); peer != "" && peer != p.self {
		log.Printf("Pick peer %s\n", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

var _ PeerPicker = (*HTTPPool)(nil)

//http客户端
type httpGetter struct {
	baseURL string
}

// Get 实现了PeerGetter 接口
func (h *httpGetter) Get(key string) ([]byte, error) {
	u := fmt.Sprintf(
		"%v/%v\n",
		h.baseURL,
		url.QueryEscape(key),
	)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v\n", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v\n", err)
	}

	return bytes, nil
}

var _ PeerGetter = (*httpGetter)(nil)
