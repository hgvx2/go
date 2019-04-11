package tcpTask

import (
	"net"
	"time"
	"tools/binaryBuf"
	"tools/server/serverDefine"
	"tools/threadPool"
)

const(
	RecvTask_msg_close int = iota
	RecvTask_msg_push_data
	RecvTask_msg_update_live
)

type recvTaskCallback func(msgType int, fd serverDefine.SOCKET_FD, ptr binaryBuf.BinaryBufPtr)

type RecvTaskParam struct {
	m_conn net.Conn
	m_fd serverDefine.SOCKET_FD
	m_pCallbackFun recvTaskCallback
	m_nKeepAlive int
}

func GetRecvTask(conn net.Conn, fd serverDefine.SOCKET_FD, nKeepAlive int, pCallbackFun recvTaskCallback) RecvTaskParam {
	return RecvTaskParam{conn, fd, pCallbackFun, nKeepAlive}
}

type recvTask struct {
	threadPool.Task
	m_taskParam RecvTaskParam
}

func (this *recvTask)Run() error {
	errCode := error(nil)
	for ; ; {
		if errCode = this.Task.Run(); nil != errCode {
			break
		}

		curUnixTime := time.Now().Unix()
		// 处理自己的逻辑
		bufTemp := [binaryBuf.BinaryBuf_64k]byte{}
		nR, err := 0, error(nil)
		for ; ;  {
			if nR, err = this.m_taskParam.m_conn.Read(bufTemp[:]); nil != err{
				// 出错 关闭socket
				this.m_taskParam.m_pCallbackFun(RecvTask_msg_close, this.m_taskParam.m_fd, nil)
				break
			}

			// 添加到数据仓库
			binPtr := binaryBuf.BinaryBufPtr(nil)
			if binPtr, err = binaryBuf.ProvideBinaryBuf(nR); nil != err{
				// 出错 关闭socket
				this.m_taskParam.m_pCallbackFun(RecvTask_msg_close, this.m_taskParam.m_fd, nil)
				break
			}
			if err = binPtr.AppendBuf(bufTemp[:nR]); nil != err{
				// 出错 关闭socket
				binaryBuf.DestroyBinaryBuf(binPtr)
				this.m_taskParam.m_pCallbackFun(RecvTask_msg_close, this.m_taskParam.m_fd, nil)
				break
			}
			this.m_taskParam.m_pCallbackFun(RecvTask_msg_push_data, this.m_taskParam.m_fd, binPtr)
			if 0 < this.m_taskParam.m_nKeepAlive &&
				curUnixTime + int64(this.m_taskParam.m_nKeepAlive) < time.Now().Unix() {
				this.m_taskParam.m_pCallbackFun(RecvTask_msg_update_live, this.m_taskParam.m_fd, nil)
			}
		}
		break
	}
	return errCode
}
