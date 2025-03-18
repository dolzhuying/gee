package geecache

import (
	"fmt"
	consistenthash "geecache/geecache/consistent-hash"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const(
	 defaultBasePath="/_geecache/"
	 defaultReplicas =50
)

type HTTPPool struct{
	self string //自己地址
	basePath string //节点间通讯地址的前缀
	mu sync.Mutex
	peers *consistenthash.Map //根据具体的 key 选择节点
	httpGetters map[string]*httpGetter //映射远程节点与对应的 httpGetter
}

func NewHTTPPool(self string)*HTTPPool{
	return &HTTPPool{
		self:self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p*HTTPPool) ServeHTTP(w http.ResponseWriter,r *http.Request){
	if !strings.HasPrefix(r.URL.Path,p.basePath){
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)

	parts:=strings.SplitN(r.URL.Path[len(p.basePath):],"/",2)
	if len(parts)!=2{
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName:=parts[0]
	key:=parts[1]
	
	group:=GetGroup(groupName)
	if group==nil{
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view,err:=group.Get(key)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlices())

}

func (p*HTTPPool) Set(peers...string){
	p.mu.Lock()
	defer p.mu.Unlock()

	p.peers=consistenthash.New(defaultReplicas,nil)
	p.peers.Add(peers...)
	p.httpGetters=make(map[string]*httpGetter,len(peers))
	for _,peer:=range peers{
		p.httpGetters[peer]=&httpGetter{baseURL: peer+p.basePath} //为每一个节点创建了一个 HTTP 客户端
	}
}

//根据具体的 key，选择节点，返回节点对应的 HTTP 客户端
func (p*HTTPPool) PickPeer(key string)(PeerGetter,bool){
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer:=p.peers.Get(key);peer!=""&&peer!=p.self{
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer],true
	}
	return nil,false
}


type httpGetter struct{
	baseURL string
}

//向远程节点客户端发送请求
func (h*httpGetter)Get(group,key string)([]byte,error){
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)

	res,err:=http.Get(u)
	if err!=nil{
		return nil,err
	}

	defer res.Body.Close()

	if res.StatusCode!=http.StatusOK{
		return nil,fmt.Errorf("server returned: %v", res.Status)
	}

	bytes,err:=io.ReadAll(res.Body)
	if err!=nil{
		return nil,fmt.Errorf("reading response body: %v", err)
	}

	return bytes,nil
}