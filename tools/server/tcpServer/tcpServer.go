package tcpServer

import (
	"fmt"
	"net"
	"tools/binaryBuf"
	"tools/server/dataDepot"
	"tools/server/serverDefine"
	"tools/server/tcpServer/connManage"
	"tools/server/tcpServer/cutPack"
	"tools/server/tcpServer/tcpTask"
	"tools/threadPool"
)

const (
	min_service_fd serverDefine.SOCKET_FD = 10000

	TCP_ACCEPT int = iota
	TCP_DISCONNECT
	TCP_KEEPALIVE
	TCP_PACK

	default_list_close_size int = 3000
)

type tcpServer struct {
	m_listenSocket  net.Listener
	m_curFd         serverDefine.SOCKET_FD
	m_pDepotManage  dataDepot.DepotManagePtr
	m_nKeepAlive    int
	m_pCallbackFun  tcpServerCallback
	m_pConnManage   connManage.ConnManagePtr
	m_pTcpTaskManage threadPool.TaskManagerPtr
	m_chCloseSocket  chan serverDefine.SOCKET_FD
	m_socketThreadPool threadPool.ThreadPool
	m_pWorkThreadPool *threadPool.ThreadPool
	m_nMaxConnNums int
}

type tcpServerCallback func(int, serverDefine.SOCKET_FD, string, binaryBuf.BinaryBufPtr)
var g_pThis *tcpServer = nil

func InitTcpServerSingle(strIpPort string, nKeepAlive, maxConnNums int,
	pFun tcpServerCallback, pWorkThreadPool *threadPool.ThreadPool) error {
	errCode := error(nil)
	for ; ;  {
		if nil != g_pThis {
			errCode = fmt.Errorf("tcpServer 已经初始化过了")
			break
		}

		listen, err := net.Listen("tcp", strIpPort)
		if err != nil {
			errCode = err
			break
		}
		g_pThis = &tcpServer{listen,
			min_service_fd,
			dataDepot.GetDepotManagePtr(),
			nKeepAlive,
		pFun,
		nil,
		tcpTask.GetTcpTaskManagePtr(),
		make(chan serverDefine.SOCKET_FD, default_list_close_size),
		threadPool.GetThreadPool(maxConnNums * 2 + 1, 10),
			pWorkThreadPool,
			maxConnNums}

		if errCode = g_pThis.m_socketThreadPool.ThreadRunning(); nil != errCode{
			break
		}

		if errCode = connManage.InitConnManage(nKeepAlive, g_pThis.connManageCallbackFun); nil != errCode{
			break
		}
		g_pThis.m_pConnManage = connManage.GetConnManagePtr()
		if errCode = cutPack.InitCutPack(g_pThis.onCutPack); nil != errCode{
			break
		}
		break
	}
	return errCode
}

func StartRun() error {
	errCode := error(nil)
	for ; ;  {
		if nil == g_pThis {
			errCode = fmt.Errorf("tcpServer 没有初始化")
			break
		}
		// 创建一个处理断开socket 的任务
		if errCode = g_pThis.createProcCloseSocketTask(); nil != errCode {
			break
		}
		for ; ;  {
			cli, err := g_pThis.m_listenSocket.Accept()
			if err != nil {
				errCode = err
				break
			}
			if g_pThis.m_pConnManage.GetCurConnectNums() >= g_pThis.m_nMaxConnNums {
				cli.Close()
				continue
			}
			g_pThis.m_curFd ++
			if min_service_fd > g_pThis.m_curFd{
				g_pThis.m_curFd = min_service_fd
			}
			// 分配数据仓库
			sendChan := make(chan int)
			if err = g_pThis.m_pDepotManage.ProvideDataDepot(g_pThis.m_curFd, sendChan); nil != err{
				// 关闭socket
				close(sendChan)
				CloseSocket(g_pThis.m_curFd)
				continue
			}
			if err = g_pThis.m_pConnManage.AddNewConn(g_pThis.m_curFd, &cli); nil != err{
				// 关闭socket
				CloseSocket(g_pThis.m_curFd)
				continue
			}

			// 创建recvTask
			if err = g_pThis.createRecvTask(cli, g_pThis.m_curFd); nil != err {
				CloseSocket(g_pThis.m_curFd)
				continue
			}
			// create sendTask
			if err = g_pThis.createSendTask(sendChan, cli, g_pThis.m_curFd); nil != err {
				CloseSocket(g_pThis.m_curFd)
				continue
			}

			g_pThis.m_pCallbackFun(TCP_ACCEPT, g_pThis.m_curFd, cli.RemoteAddr().String(), nil)
		}
		break
	}
	return errCode
}

func PushSendData(fd serverDefine.SOCKET_FD, ptr binaryBuf.BinaryBufPtr) error {
	errCode := error(nil)
	for ; ;  {
		if nil == g_pThis {
			errCode = fmt.Errorf("tcpServer 没有初始化")
			break
		}
		if nil == ptr {
			errCode = fmt.Errorf("tcpServer 空指针")
			break
		}
		_, errCode = g_pThis.m_pDepotManage.PushSendData(ptr, fd)
		break
	}
	return errCode
}

func CloseSocket(fd serverDefine.SOCKET_FD)  {
	if nil != g_pThis {
		g_pThis.m_chCloseSocket <- fd
	}
}

func (this *tcpServer)createProcCloseSocketTask() error {
	errCode := error(nil)
	for ; ;  {
		pCloseTask := threadPool.TaskPtr(nil)
		if pCloseTask, errCode = this.m_pTcpTaskManage.ProvideTask(tcpTask.Task_type_proc_closesocket,
			tcpTask.GetProcCloseSocketTaskParam(this.m_chCloseSocket,
				this.onProcCloseSocketTask)); nil != errCode{
				break
			}
		if errCode = this.m_socketThreadPool.PushTaskToBack(pCloseTask); nil != errCode{
			this.m_pTcpTaskManage.DestoryTask(pCloseTask)
			break
		}
		break
	}
	return errCode
}

func (this *tcpServer)createRecvTask(con net.Conn, fd serverDefine.SOCKET_FD) error {
	errCode := error(nil)
	for ; ;  {
		pRecvTask := threadPool.TaskPtr(nil)
		if pRecvTask, errCode = this.m_pTcpTaskManage.ProvideTask(tcpTask.Task_type_recv,
			tcpTask.GetRecvTask(con, fd, this.m_nKeepAlive, this.onRecvTask)); nil != errCode{
			break
		}
		if errCode = this.m_socketThreadPool.PushTaskToBack(pRecvTask); nil != errCode{
			this.m_pTcpTaskManage.DestoryTask(pRecvTask)
			break
		}
		break
	}
	return errCode
}

func (this *tcpServer)createSendTask(ch chan int, con net.Conn, fd serverDefine.SOCKET_FD) error {
	errCode := error(nil)
	for ; ;  {
		pSendTask := threadPool.TaskPtr(nil)
		if pSendTask, errCode = this.m_pTcpTaskManage.ProvideTask(tcpTask.Task_type_send,
			tcpTask.GetSendTask(ch, con, fd, this.onSendTask)); nil != errCode{
			break
		}
		if errCode = this.m_socketThreadPool.PushTaskToBack(pSendTask); nil != errCode{
			this.m_pTcpTaskManage.DestoryTask(pSendTask)
			break
		}
		break
	}
	return errCode
}

func (this *tcpServer)createCutPackTask(fd serverDefine.SOCKET_FD) error {
	errCode := error(nil)
	for ; ;  {
		pCutPackTask := threadPool.TaskPtr(nil)
		if pCutPackTask, errCode = this.m_pTcpTaskManage.ProvideTask(tcpTask.Task_type_cutpack,
			tcpTask.GetCutPackTaskParam(fd)); nil != errCode{
				break
		}
		if errCode = this.m_pWorkThreadPool.PushTaskToBack(pCutPackTask); nil != errCode{
			this.m_pTcpTaskManage.DestoryTask(pCutPackTask)
			break
		}
		break
	}
	return errCode
}

func (this *tcpServer)onProcCloseSocketTask(fd serverDefine.SOCKET_FD)  {
	// 处理关闭连接 清空本地缓存
	bCallback := true
	if err := this.m_pConnManage.DelConn(fd); nil != err{
		bCallback = false
	}
	if err := this.m_pDepotManage.DestoryDataDepot(fd); nil != err{
		bCallback = false
	}
	if err := cutPack.ClearHalfPack(fd); nil != err{
		bCallback = false
	}
	if bCallback {
		this.m_pCallbackFun(TCP_DISCONNECT, fd, "", nil)
	}
}

func (this *tcpServer)onRecvTask(msgType int, fd serverDefine.SOCKET_FD, ptr binaryBuf.BinaryBufPtr)  {
	errCode := error(nil)
	for ; ;  {
		switch msgType {
		case tcpTask.RecvTask_msg_close:
			CloseSocket(fd)
		case tcpTask.RecvTask_msg_push_data:
			nState := 0
			if nState, errCode = this.m_pDepotManage.PushRecvData(ptr, fd); nil != errCode{
				binaryBuf.DestroyBinaryBuf(ptr)
				CloseSocket(fd)
			}else{
				if dataDepot.DM_PUSH_ET == nState {
					// 创建 切包任务
					if errCode = this.createCutPackTask(fd); nil != errCode {
						CloseSocket(fd)
					}
				}
			}
		case tcpTask.RecvTask_msg_update_live:
			this.m_pConnManage.UpdateAliveTime(fd)
		default:
			binaryBuf.DestroyBinaryBuf(ptr)
			errCode = fmt.Errorf("onRecvTask 错误的类型 type=%d", msgType)
		}
		break
	}
}
func (this *tcpServer)onSendTask(msgType int, fd serverDefine.SOCKET_FD) binaryBuf.BinaryBufPtr {
	pRet:= binaryBuf.BinaryBufPtr(nil)
	for ; ;  {
		switch msgType {
		case tcpTask.SendTask_msg_close:
			CloseSocket(fd)
		case tcpTask.SendTask_msg_pop_data:
			err := error(nil)
			if pRet, err = this.m_pDepotManage.PopSendData(fd); nil != err{
				CloseSocket(fd)
			}
		default:
		}
		break
	}
	return pRet
}

func (this *tcpServer)onCutPack(cutMsgType int, fd serverDefine.SOCKET_FD, ptr binaryBuf.BinaryBufPtr) {
	switch cutMsgType {
	case cutPack.CUT_PACK_OK:
		if nil == ptr {
			CloseSocket(fd)
		}else{
			this.m_pCallbackFun(TCP_PACK, fd, "", ptr)
		}
	case cutPack.CUT_PACK_ERR:
		CloseSocket(fd)
	default:
	}
}

func (this *tcpServer)connManageCallbackFun(tp int, fd serverDefine.SOCKET_FD)  {
	switch tp {
	case connManage.CM_KEEPALIVE:
		this.m_pCallbackFun(TCP_KEEPALIVE, fd, "", nil)
	case connManage.CM_CLOSE:
		CloseSocket(fd)
	default:
	}
}
