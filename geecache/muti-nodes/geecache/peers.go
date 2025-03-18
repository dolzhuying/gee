package geecache

type PeerGetter interface{
	Get(group,key string)([]byte,error) //从对应 group 查找缓存值
}

type PeerPicker interface{
	PickPeer(key string)(peer PeerGetter,ok bool)
}

