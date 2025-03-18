package singleflight

import(
	"sync"
)

//代表正在进行中，或已经结束的请求
type call struct{
	wg sync.WaitGroup
	val interface{}
	err error
}

//singleflight 的主数据结构，管理不同 key 的请求(call)
type Group struct{
	mu sync.Mutex
	m map[string]*call
}


//多个并发请求同时访问同一个未在本地cache的缓存值，先获取锁，将请求加入任务队列，此时其他请求阻塞，
// 加入完成后释放锁，此时其他请求获取锁会发现任务队列已有，则释放锁，阻塞等待第一个远程调用返回
func (g*Group)Do(key string,fn func()(interface{},error))(interface{},error){
	g.mu.Lock()
	if g.m==nil{
		g.m=make(map[string]*call)
	}
	if c,ok:=g.m[key];ok{
		g.mu.Unlock()
		c.wg.Wait() // 如果请求正在进行中，则等待,,,阻塞，直到锁被释放
		return c.val,c.err
	}
	c:=new(call)
	c.wg.Add(1) // 发起请求前加锁,锁加1
	g.m[key]=c // 添加到 g.m，表明 key 已经有对应的请求在处理
	g.mu.Unlock()

	c.val,c.err=fn()
	c.wg.Done()//锁减1

	g.mu.Lock()
	delete(g.m,key)
	g.mu.Unlock()
	
	return c.val,c.err
}