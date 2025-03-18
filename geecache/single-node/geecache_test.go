package geecache

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetter(t*testing.T){
	t.Helper()
	var f =GetterFunc(func(key string)([]byte,error){
		return []byte(key),nil
	})

	expect:=[]byte("key")
	if v,_:=f.Get("key");!reflect.DeepEqual(expect,v){
		t.Errorf("callback failed")
	}
}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t*testing.T){
	t.Helper()
	localCounts:=make(map[string]int,len(db))
	gee:=NewGroup("scores",2<<10,GetterFunc(
		func(key string)([]byte,error){
			if v,ok:=db[key];ok{
				if _,ok:=localCounts[key];!ok{
					localCounts[key]=0
				}
				localCounts[key]+=1
				return []byte(v),nil
			}
			return nil,fmt.Errorf("%s not exist",key)
		}))

		for k,v:=range db{
			if view,err:=gee.Get(k);err!=nil||view.String()!=v{
				t.Fatal("failed to get value of Tom")
			}
			if _,err:=gee.Get(k);err!=nil||localCounts[k]>1{
				t.Fatalf("cache %s miss", k)
			}
		}

		if view, err := gee.Get("unknown"); err == nil {
			t.Fatalf("the value of unknow should be empty, but %s got", view)
		}

}