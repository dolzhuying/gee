package geecache
import pb "geecache/geecache/geecachepb/geecachepb"

type PeerGetter interface{
	Get(in *pb.Request, out *pb.Response) error //从对应 group 查找缓存值
}

type PeerPicker interface{
	PickPeer(key string)(peer PeerGetter,ok bool)
}

