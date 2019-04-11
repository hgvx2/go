package connManage

import (
	"net"
	"time"
	"tools/server/serverDefine"
)

type connLive struct{
	m_fd   serverDefine.SOCKET_FD
	m_pCon *net.Conn
	m_unixTime int64
}

type ConnLivePtr interface {
	UpdateNowTime()
}

func (this *connLive)UpdateNowTime()  {
	this.m_unixTime = time.Now().Unix()
}