package dataDepot

import (
	"fmt"
	"sync"
	"tools/binaryBuf"
	"tools/memoryPool"
	"tools/server/serverDefine"
)

const (
	dp_default_queue_size int = 50000
)

type depotPool struct {
	memoryPool.MemoryPool
}

func getDepotPool() depotPool {
	return  depotPool{memoryPool.GetMemoryPool()}
}

func (this *depotPool)provideDepot(fd serverDefine.SOCKET_FD, sendChan chan <-int) (*depot, error){
	pRet, errCode:= (*depot)(nil), error(nil)
	for ; ;  {
		pTemp := this.ProvideNode()
		if nil == pTemp {
			pRet = &depot{recvDepot{depotBase{make(chan binaryBuf.BinaryBufPtr, dp_default_queue_size),
				false}},
				sendDepot{depotBase{make(chan binaryBuf.BinaryBufPtr, dp_default_queue_size),
					false},
					nil},
				serverDefine.SOCKET_FD_ERROR,
				sync.Mutex{}}
			if errCode = this.AddNewNode(pRet); nil != errCode{
				pRet = (*depot)(nil)
			}
		}else{
			bOk := false
			if pRet, bOk = pTemp.(*depot); !bOk{
				errCode = fmt.Errorf("断言 *depot 失败")
				pRet = (*depot)(nil)
			}
		}
		if nil != pRet {
			pRet.m_fd = fd
			pRet.m_sendDepot.m_sendChan = sendChan
		}
		break
	}
	return pRet, errCode
}

func (this *depotPool)destoryDepot(ptr *depot) error {
	errCode := error(nil)
	for ; ;  {
		if nil == ptr {
			errCode = fmt.Errorf("destoryDepot 空指针")
			break
		}
		ptr.clearDepot()
		errCode = this.DestoryNode(ptr)
		break
	}
	return errCode
}