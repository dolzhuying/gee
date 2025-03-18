package lru

import "container/list"

type Cache struct{
	maxBytes int64 //最大允许内存
	nBytes int64 //已使用内存
	ll *list.List  
	cache map[string]*list.Element //值为指向双向链表的节点
	OnEvicted func(key string,value Value) //某条记录被移除时的回调函数
}

//为了通用性，允许值类型value是实现了该接口的任意类型
type Value interface{
	Len() int //返回值所占用内存大小
}

//双向链表节点的数据类型
//在链表中仍保存每个值对应的 key 的好处在于，淘汰队首节点时，需要用 key 从字典中删除对应的映射
type entry struct{ 
	key string
	value Value
}

func New(maxBytes int64,onEvicted func(string,Value))*Cache{
	return &Cache{
		maxBytes: maxBytes,
		ll:list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c*Cache) Get(key string) (value Value,ok bool){
	if ele,ok:=c.cache[key];ok{
		c.ll.MoveToFront(ele) 
		kv:=ele.Value.(*entry)
		return kv.value,ok
	}
	return
}

func (c*Cache) RemoveOldest(){
	ele:=c.ll.Back()
	if ele!=nil{
		c.ll.Remove(ele)
		kv:=ele.Value.(*entry)
		delete(c.cache,kv.key) //delete(c.cache, kv.key)，从字典中 c.cache 删除该节点的映射关系
		c.nBytes-=int64(len(kv.key))+int64(kv.value.Len())
		if c.OnEvicted!=nil{
			c.OnEvicted(kv.key,kv.value)
		}

	}
}

//新增/修改，判断是否已有键值对，有则更新键值对和占用内存大小，没有则创建键值对并加到链表尾，更新键值对大小，，然后删除溢出内存
func (c*Cache) Add(key string,value Value){
	if ele,ok:=c.cache[key];ok{
		c.ll.MoveToFront(ele)
		kv:=ele.Value.(*entry)
		c.nBytes+=int64(value.Len())-int64(kv.value.Len())
		kv.value=value
	}else{
		ele:=c.ll.PushFront(&entry{key,value})
		c.cache[key]=ele
		c.nBytes+=int64(len(key))+int64(value.Len())
	}
	for c.maxBytes!=0&&c.maxBytes<c.nBytes{
		c.RemoveOldest()
	}
}

func (c*Cache) Len() int{
	return c.ll.Len()
}

