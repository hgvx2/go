package dataDepot

import (
	"sync"
	"tools/binaryBuf"
	"tools/server/serverDefine"
)

type depot struct {
	m_recvDepot recvDepot
	m_sendDepot sendDepot

	m_fd serverDefine.SOCKET_FD
	m_mutex sync.Mutex
}

func (this *depot)pushRecvDepot(fd serverDefine.SOCKET_FD, ptr binaryBuf.BinaryBufPtr) int {
	nRet := DM_PUSH_ERR
	for ; ;  {
		this.m_mutex.Lock()
		defer this.m_mutex.Unlock()
		if fd != this.m_fd {
			break
		}
		nRet = this.m_recvDepot.pushData(ptr)
		break
	}
	return nRet
}

func (this *depot)popRecvDepot(fd serverDefine.SOCKET_FD) binaryBuf.BinaryBufPtr  {
	pRet := binaryBuf.BinaryBufPtr(nil)
	for ; ;  {
		this.m_mutex.Lock()
		defer this.m_mutex.Unlock()
		if fd != this.m_fd {
			break
		}
		pRet = this.m_recvDepot.popData()
		break
	}
	return pRet
}

func (this *depot)pushSendDepot(fd serverDefine.SOCKET_FD,ptr binaryBuf.BinaryBufPtr) int {
	nRet := DM_PUSH_ERR
	for ; ;  {
		this.m_mutex.Lock()
		defer this.m_mutex.Unlock()
		if fd != this.m_fd {
			break
		}
		nRet = this.m_sendDepot.pushData(ptr)
		break
	}
	return nRet
}

func (this *depot)popSendDepot(fd serverDefine.SOCKET_FD) binaryBuf.BinaryBufPtr  {
	pRet := binaryBuf.BinaryBufPtr(nil)
	for ; ;  {
		this.m_mutex.Lock()
		defer this.m_mutex.Unlock()
		if fd != this.m_fd {
			break
		}
		pRet = this.m_sendDepot.popData()
		break
	}
	return pRet
}

func (this *depot)clearDepot()  {
	this.m_mutex.Lock()
	defer this.m_mutex.Unlock()
	if serverDefine.SOCKET_FD_ERROR == this.m_fd {
		return
	}
	this.m_fd = serverDefine.SOCKET_FD_ERROR
	this.m_recvDepot.clearData()
	this.m_sendDepot.clearData()
}