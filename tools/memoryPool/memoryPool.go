package memoryPool

import (
	"fmt"
	"sync"
)

const default_pool_nums  = 100

type MemoryPoolPtr interface {
	ProvideNode() interface{}
	DestoryNode(vKey interface{}) error
	AddNewNode(vKey interface{}) error
}

type MemoryPool struct {
	m_liseFreeNode []interface{}
	m_mapUseNode map[interface{}]interface{}
	m_mutex sync.Mutex
}

func GetMemoryPool() MemoryPool {
	return MemoryPool{make([]interface{}, 0, default_pool_nums),
		make(map[interface{}]interface{}),
	sync.Mutex{}}
}

func (this *MemoryPool)ProvideNode() interface{} {
	vRet := interface{}(nil)
	for ; ;  {
		this.m_mutex.Lock()
		defer this.m_mutex.Unlock()

		if 0 < len(this.m_liseFreeNode){
			vRet = this.m_liseFreeNode[0]
			this.m_liseFreeNode = this.m_liseFreeNode[1:]
			this.m_mapUseNode[vRet] = vRet
		}
		break
	}
	return vRet
}

func (this *MemoryPool)DestoryNode(vKey interface{}) error{
	errCode := error(nil)
	for ; ;  {
		this.m_mutex.Lock()
		defer this.m_mutex.Unlock()

		if _, ok := this.m_mapUseNode[vKey]; ok{
			this.m_liseFreeNode = append(this.m_liseFreeNode, vKey)
			delete(this.m_mapUseNode, vKey)
		}else{
			errCode = fmt.Errorf("没有找到需要归还的节点")
		}
		break
	}
	return errCode
}

func (this *MemoryPool)AddNewNode(vKey interface{}) error {
	errCode := error(nil)
	for ; ;  {
		this.m_mutex.Lock()
		defer this.m_mutex.Unlock()
		if _, ok := this.m_mapUseNode[vKey]; ok{
			errCode = fmt.Errorf("在使用节点中找到重复的")
			break
		}
		this.m_mapUseNode[vKey] = vKey
		break
	}
	return errCode
}
