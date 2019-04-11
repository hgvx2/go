package dataDepot

import (
	"fmt"
	"sync"
	"tools/binaryBuf"
	"tools/server/serverDefine"
)

const (
	max_socket_group int = 100
)

const (
	DM_PUSH_ERR int = iota
	DM_PUSH_OK
	DM_PUSH_ET
)

type depotManage struct {
	m_depotPool depotPool

	m_arrMutex [max_socket_group]sync.Mutex
	m_arrMapSocket [max_socket_group]map[serverDefine.SOCKET_FD] *depot
}

type DepotManagePtr interface {
	ProvideDataDepot(fd serverDefine.SOCKET_FD, sendChan chan <-int) error
	DestoryDataDepot(fd serverDefine.SOCKET_FD) error
	PushRecvData(binPtr binaryBuf.BinaryBufPtr, fd serverDefine.SOCKET_FD) (int, error)
	PopRecvData(fd serverDefine.SOCKET_FD) (binaryBuf.BinaryBufPtr, error)
	PushSendData(binPtr binaryBuf.BinaryBufPtr, fd serverDefine.SOCKET_FD) (int, error)
	PopSendData(fd serverDefine.SOCKET_FD) (binaryBuf.BinaryBufPtr, error)
}

var g_pThis *depotManage = nil

func InitDepotManage() error {
	errCode := error(nil)
	for ; ;  {
		if nil != g_pThis {
			errCode = fmt.Errorf("depotManage 已经初始化过了")
			break
		}
		depotMan := depotManage{}
		depotMan.m_depotPool = getDepotPool()
		for i := 0; i < max_socket_group; i ++ {
			depotMan.m_arrMutex[i] = sync.Mutex{}
			depotMan.m_arrMapSocket[i] = make(map[serverDefine.SOCKET_FD]*depot)
		}
		g_pThis = &depotMan
		break
	}
	return errCode
}

func GetDepotManagePtr() DepotManagePtr  {
	return g_pThis
}


func (this *depotManage)ProvideDataDepot(fd serverDefine.SOCKET_FD, sendChan chan <-int) error {
	errCode := error(nil)
	for ; ;  {
		pDepot := (*depot)(nil)
		if pDepot, errCode = this.m_depotPool.provideDepot(fd, sendChan); nil != errCode{
			break
		}
		if nil == pDepot {
			errCode = fmt.Errorf("分配数据仓库失败")
			break
		}
		errCode = this.addDepot(fd, pDepot)
		break
	}
	return errCode
}

func (this *depotManage)DestoryDataDepot(fd serverDefine.SOCKET_FD) error {
	errCode := error(nil)
	for ; ;  {
		pDepot := this.popDepot(fd)
		if nil == pDepot {
			errCode = fmt.Errorf("没有找到仓库 fd=", fd)
			break
		}
		errCode = this.m_depotPool.destoryDepot(pDepot)
		break
	}
	return errCode

}

func (this *depotManage)findDepot(fd serverDefine.SOCKET_FD) (*depot, error) {
	pRet, errCode:= (*depot)(nil), error(nil)
	for ; ;  {
		nIndex := fd % serverDefine.SOCKET_FD(max_socket_group)
		this.m_arrMutex[nIndex].Lock()
		defer  this.m_arrMutex[nIndex].Unlock()

		if pDepot, ok := this.m_arrMapSocket[nIndex][fd]; ok {
			pRet = pDepot
		}else{
			errCode = fmt.Errorf("没有找到 %lu 对于的socket连接", fd)
		}
		break
	}
	return pRet, errCode
}
func (this *depotManage)popDepot(fd serverDefine.SOCKET_FD) *depot {
	pRet := (*depot)(nil)
	for ; ;  {
		nIndex := fd % serverDefine.SOCKET_FD(max_socket_group)
		this.m_arrMutex[nIndex].Lock()
		defer  this.m_arrMutex[nIndex].Unlock()
		if pDepot, ok := this.m_arrMapSocket[nIndex][fd]; ok {
			pRet = pDepot
			delete(this.m_arrMapSocket[nIndex], fd)
		}
		break
	}
	return pRet
}

func (this *depotManage)addDepot(fd serverDefine.SOCKET_FD, pDepot *depot) error {
	errCode := error(nil)
	for ; ;  {
		nIndex := fd % serverDefine.SOCKET_FD(max_socket_group)
		this.m_arrMutex[nIndex].Lock()
		defer  this.m_arrMutex[nIndex].Unlock()
		if _, ok := this.m_arrMapSocket[nIndex][fd]; ok {
			errCode = fmt.Errorf("仓库已经存在 fd=", fd)
			break
		}
		this.m_arrMapSocket[nIndex][fd] = pDepot
		break
	}
	return errCode
}

func (this *depotManage)PushRecvData(binPtr binaryBuf.BinaryBufPtr, fd serverDefine.SOCKET_FD) (int, error) {
	nRet, errCode := DM_PUSH_ERR, error(nil)
	for ; ;  {
		pDepot := (*depot)(nil)
		if pDepot, errCode = this.findDepot(fd); nil == errCode{
			nRet = pDepot.pushRecvDepot(fd, binPtr)
			if DM_PUSH_ERR == nRet {
				errCode = fmt.Errorf("数据添加到接收仓库失败")
			}
		}
		break
	}
	return nRet, errCode
}

func (this *depotManage)PopRecvData(fd serverDefine.SOCKET_FD) (binaryBuf.BinaryBufPtr, error)  {
	pRet,errCode := binaryBuf.BinaryBufPtr(nil), error(nil)
	for ; ;  {
		pDepot := (*depot)(nil)
		if pDepot, errCode = this.findDepot(fd); nil == errCode{
			pRet = pDepot.popRecvDepot(fd)
		}
		break
	}
	return pRet, errCode
}

func (this *depotManage)PushSendData(binPtr binaryBuf.BinaryBufPtr, fd serverDefine.SOCKET_FD) (int, error) {
	nRet, errCode := DM_PUSH_ERR, error(nil)
	for ; ;  {
		pDepot := (*depot)(nil)
		if pDepot, errCode = this.findDepot(fd); nil == errCode{
			nRet = pDepot.pushSendDepot(fd, binPtr)
			if DM_PUSH_ERR == nRet {
				errCode = fmt.Errorf("数据添加到发送仓库失败")
			}
		}
		break
	}
	return nRet, errCode
}

func (this *depotManage)PopSendData(fd serverDefine.SOCKET_FD) (binaryBuf.BinaryBufPtr, error)  {
	pRet,errCode := binaryBuf.BinaryBufPtr(nil), error(nil)
	for ; ;  {
		pDepot := (*depot)(nil)
		if pDepot, errCode = this.findDepot(fd); nil == errCode{
			pRet = pDepot.popSendDepot(fd)
		}
		break
	}
	return pRet, errCode
}