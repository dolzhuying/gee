package consistenthash

import(
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct{
	hash Hash  //hash函数
	replicas int //虚拟节点倍数
	keys []int  //哈希环
	hashMap map[int]string //虚拟节点和真实节点的映射  虚拟hash--真实name
}

func New(replicas int,fn Hash)*Map{
	m:=&Map{
		replicas: replicas,
		hash: fn,
		hashMap: make(map[int]string),
	}
	if m.hash==nil{
		m.hash=crc32.ChecksumIEEE
	}
	return m
}

//添加真实节点
func (m*Map) Add(keys...string){
	for _,key:=range keys{
		for i:=0;i<m.replicas;i++{
			hash:=int(m.hash([]byte(strconv.Itoa(i)+key)))
			m.keys=append(m.keys,hash)
			m.hashMap[hash]=key
		}
	}
	sort.Ints(m.keys)
}

//选择节点
func (m*Map)Get(key string)string{
	if len(m.keys)==0{
		return ""
	}
	hash:=int(m.hash([]byte(key)))
	idx:=sort.Search(len(m.keys),func(i int)bool{
		return m.keys[i]>=hash
	})//如果所有虚拟节点的哈希值都小于目标哈希值，idx 会等于 len(m.keys)
	return m.hashMap[m.keys[idx%len(m.keys)]]
}


