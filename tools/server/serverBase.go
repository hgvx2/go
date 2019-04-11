package server

import (
	"fmt"
	"tools/binaryBuf"
	"tools/server/dataDepot"
	"tools/server/serverDefine"
	"tools/server/tcpServer"
	"tools/threadPool"
)

const (
	default_work_thread_nums int = 1000           //  默认的工作线程的数量
	default_work_thread_cache_nums int = 10000    //  每个线程对应的任务队列的长度
)

type serverBaseCallback func(int, serverDefine.SOCKET_FD, string, binaryBuf.BinaryBufPtr)

type serverBase struct {
	m_pCallback serverBaseCallback
	m_workThreadPool threadPool.ThreadPool
}

var g_pThis *serverBase = nil

func InitServerBase(strIpPort string, nKeepalive int, maxConnNums int, pFun serverBaseCallback) (*threadPool.ThreadPool, error) {
	pThreadPool, errCode := (*threadPool.ThreadPool)(nil), error(nil)
	for ; ;  {
		if nil != g_pThis {
			errCode = fmt.Errorf("serverBase 已经初始化过了")
			break
		}
		if errCode = binaryBuf.InitBinaryBufPool(1000); nil != errCode {
			break
		}
		if errCode = dataDepot.InitDepotManage(); nil != errCode {
			break
		}
		g_pThis = &serverBase{pFun, threadPool.GetThreadPool(default_work_thread_nums, default_work_thread_cache_nums)}
		if errCode = g_pThis.m_workThreadPool.ThreadRunning(); nil != errCode{
			break
		}
		pThreadPool = &g_pThis.m_workThreadPool
		if errCode = tcpServer.InitTcpServerSingle(strIpPort, nKeepalive, maxConnNums, onTcpServer, pThreadPool); nil != errCode{
			break
		}
		break
	}
	return pThreadPool, errCode
}

func ServerRunning()  {
	tcpServer.StartRun()
}

func PushSendData(fd serverDefine.SOCKET_FD, ptr binaryBuf.BinaryBufPtr) error{
	return tcpServer.PushSendData(fd, ptr)
}

func CloseSocket(fd serverDefine.SOCKET_FD){
	tcpServer.CloseSocket(fd)
}

func onTcpServer(pt int, fd serverDefine.SOCKET_FD, msg string, ptr binaryBuf.BinaryBufPtr)  {
	if nil != g_pThis {
		g_pThis.m_pCallback(pt, fd, msg, ptr)
	}
}
