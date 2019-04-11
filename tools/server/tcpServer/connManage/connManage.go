package connManage

import (
	"fmt"
	"net"
	"sync"
	"time"
	"tools/server/serverDefine"
)

const group_nums int = 100
const try_times int = 3
const (
	CM_KEEPALIVE int = iota
	CM_CLOSE
)

type connManage struct {
	m_connPool connLivePool
	m_pTimer *time.Ticker

	m_listMapCon [group_nums]map[serverDefine.SOCKET_FD]*connLive
	m_listMutex [group_nums]sync.Mutex

	m_nKeepalive int64
	m_pCallback connManageCallback

	m_nCurConnectNums int
}

type ConnManagePtr interface {
	AddNewConn(fd serverDefine.SOCKET_FD, pConn *net.Conn) error
	DelConn(fd serverDefine.SOCKET_FD) error
	GetCurConnectNums() int
	UpdateAliveTime(fd serverDefine.SOCKET_FD)
}

var g_pThis *connManage = nil

type connManageCallback func(int, serverDefine.SOCKET_FD)

func InitConnManage(nKeepAlive int, pFun connManageCallback) error {
	errCode := error(nil)
	for ; ;  {
		if nil != g_pThis {
			errCode = fmt.Errorf("connManage 已经初始化了")
			break
		}
		g_pThis = &connManage{getConnLivePool(),
			nil,
			[group_nums]map[serverDefine.SOCKET_FD]*connLive{},
			[group_nums]sync.Mutex{},
		    0,
			pFun,
		0}

		for i := 0; i < group_nums; i ++ {
			g_pThis.m_listMapCon[i] = make(map[serverDefine.SOCKET_FD]*connLive)
			g_pThis.m_listMutex[i] = sync.Mutex{}
		}

		if 0 < nKeepAlive {
			g_pThis.m_nKeepalive = int64(nKeepAlive)
			g_pThis.m_pTimer = time.NewTicker(time.Duration(nKeepAlive) * time.Second)
			go g_pThis.keepaliveTimer()
		}
		break
	}
	return errCode
}

func GetConnManagePtr() ConnManagePtr {
	return g_pThis
}

func (this *connManage)keepaliveTimer() {
	for {
		select {
		case <-this.m_pTimer.C:
			this.clearTimeoutConn()
		}
	}
}

func (this *connManage)clearTimeoutConn()  {
	for ; ;  {
		if 0 >= this.m_nKeepalive {
			break
		}
		timeOut := int64(try_times) * this.m_nKeepalive
		for i := 0; i < group_nums; i ++ {
			this.m_listMutex[i].Lock()
			for fd, pCon := range this.m_listMapCon[i] {
				now := time.Now().Unix()
				if  timeOut < now - pCon.m_unixTime {
					this.m_pCallback(CM_CLOSE, fd)
					continue
				}
				if this.m_nKeepalive < now - pCon.m_unixTime{
					this.m_pCallback(CM_KEEPALIVE, fd)
				}
			}
			this.m_listMutex[i].Unlock()
		}
		break
	}
}

func (this *connManage)addConn(pLive *connLive) error  {
	errCode := error(nil)
	for ; ;  {
		nIndex := int(pLive.m_fd % serverDefine.SOCKET_FD(group_nums))
		this.m_listMutex[nIndex].Lock()
		defer this.m_listMutex[nIndex].Unlock()

		if _, ok := this.m_listMapCon[nIndex][pLive.m_fd]; ok{
			errCode = fmt.Errorf("fd=%llu 已经存在", pLive.m_fd)
			break
		}
		this.m_listMapCon[nIndex][pLive.m_fd] = pLive
		break
	}
	return errCode
}

func (this *connManage)delConn(fd serverDefine.SOCKET_FD) error  {
	errCode := error(nil)
	for ; ;  {
		nIndex := int(fd % serverDefine.SOCKET_FD(group_nums))
		this.m_listMutex[nIndex].Lock()
		defer this.m_listMutex[nIndex].Unlock()

		var pCon *connLive = nil
		ok := false
		if pCon, ok = this.m_listMapCon[nIndex][fd]; !ok{
			errCode = fmt.Errorf("fd=%d 不存在", fd)
			break
		}
		(*pCon.m_pCon).Close()
		delete(this.m_listMapCon[nIndex], fd)
		this.m_connPool.destoryConLive(pCon)
		this.m_nCurConnectNums -= 1
		break
	}
	return errCode
}

func (this *connManage)AddNewConn(fd serverDefine.SOCKET_FD, pConn *net.Conn) error{
	 errCode :=  error(nil)
	for ; ;  {
		var pLive *connLive = nil
		if pLive, errCode = this.m_connPool.provideConLive(); nil != errCode{
			break
		}
		pLive.m_fd = fd
		pLive.m_pCon = pConn
		pLive.m_unixTime = time.Now().Unix()

		if errCode = this.addConn(pLive); nil != errCode{
			this.m_connPool.destoryConLive(pLive)
		}
		this.m_nCurConnectNums += 1
		break
	}
	return errCode
}

func (this *connManage)DelConn(fd serverDefine.SOCKET_FD) error{
	return this.delConn(fd)
}

func (this *connManage)GetCurConnectNums() int{
	return this.m_nCurConnectNums
}

func (this *connManage) UpdateAliveTime(fd serverDefine.SOCKET_FD){
	if 0 == this.m_nKeepalive {
		return
	}

	nIndex := int(fd % serverDefine.SOCKET_FD(group_nums))
	this.m_listMutex[nIndex].Lock()
	defer this.m_listMutex[nIndex].Unlock()

	if pLive, ok := this.m_listMapCon[nIndex][fd]; ok{
		pLive.UpdateNowTime()
	}
}




