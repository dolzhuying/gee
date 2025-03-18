package lru_test

import (
	"geecache/geecache/lru"
	"reflect"
	"testing"
)

type String string

func (s String) Len() int{
	return len(s)
}

func TestGet(t*testing.T){
	t.Helper()
	lru:=lru.New(int64(0),nil)
	lru.Add("hello",String("world"))
	if v,ok:=lru.Get("hello");!ok||string(v.(String))!="world"{
		t.Fatalf("cannot find key hello for value:%s",v)
	}
}

func TestRemoveOldest(t*testing.T){
	t.Helper()
	k1,k2,k3:="11","22","33"
	v1,v2,v3:="qq","ww","pp"
	cap:=len(k1+k2+v1+v2)
	lru:=lru.New(int64(cap),nil)
	lru.Add(k1,String(v1))
	lru.Add(k2,String(v2))
	lru.Add(k3,String(v3))

	if _,ok:=lru.Get(k1);ok||lru.Len()!=2{
		t.Fatalf("failed to remove oldest")
	}
}

func TestOnEvicted(t*testing.T){
	t.Helper()
	keys:=make([]string,0)
	callback:=func(key string,value lru.Value){
		keys=append(keys,key)
	}

	k1,k2,k3:="11","22","33"
	v1,v2,v3:="qq","ww","pp"
	cap:=len(k1+k2+v1+v2)
	lru:=lru.New(int64(cap),callback)
	lru.Add(k1,String(v1))
	lru.Add(k2,String(v2))
	lru.Add(k3,String(v3))

	expect:=[]string{k1}
	if !reflect.DeepEqual(expect,keys){
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}

