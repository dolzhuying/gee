package geecache

import (
	"fmt"
	singleflight "geecache/geecache/single-flight"
	pb "geecache/geecache/geecachepb/geecachepb"

	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

// 函数类型实现某一个接口，称之为接口型函数，
// 方便使用者在调用时既能够传入函数作为参数，也能够传入实现了该接口的结构体作为参数
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

//一个group代表一个节点，peers用于当本地缓存miss时请求其他节点
type Group struct {
	name      string
	getter    Getter //缓存未命中时获取源数据的回调
	mainCache cache
	peers PeerPicker
	loader *singleflight.Group
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group) //全局变量
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:&singleflight.Group{},
	}
	groups[name] = g
	return g

}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g*Group) RegisterPeer(peers PeerPicker){
	if g.peers!=nil{
		panic("RegisterPeerPicker called more than once")
	}
	g.peers=peers
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.Add(key, value) //添加缓存
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key) //缓存不存在时调用接口函数，获取源数据
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) load(key string) (value ByteView, err error) {
	//确保了并发场景下针对相同的 key，load 过程只会调用一次
	viewi,err:=g.loader.Do(key,func()(interface{},error){
		if g.peers!=nil{
			if peer,ok:=g.peers.PickPeer(key);ok{//判断是否远程节点
				if value,err:=g.getFromPeer(peer,key);err==nil{//从具远程节点的客户端获取缓存
					return value,nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
		//本机节点或者失败，则回退到getlocally
		return g.getLocally(key)
	})

	if err==nil{
		return viewi.(ByteView),nil
	}
	return
	
}

func (g*Group) getFromPeer(peer PeerGetter,key string)(ByteView,error){
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("empty key")
	}

	if v, ok := g.mainCache.Get(key); ok {
		log.Printf("cache hit")
		return v, nil
	}
	return g.load(key)

}
