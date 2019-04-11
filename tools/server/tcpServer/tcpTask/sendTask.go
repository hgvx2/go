package tcpTask

import (
	"net"
	"tools/binaryBuf"
	"tools/server/dataDepot"
	"tools/server/serverDefine"
	"tools/threadPool"
)

const(
	SendTask_msg_close int = iota
	SendTask_msg_pop_data
)

type sendTaskCallback func(msgType int, fd serverDefine.SOCKET_FD) binaryBuf.BinaryBufPtr

type SendTaskParam struct {
	m_chSend chan int
	m_conn net.Conn
	m_fd serverDefine.SOCKET_FD
	m_pCallbackFun sendTaskCallback
}

func GetSendTask(chSend chan int, conn net.Conn, fd serverDefine.SOCKET_FD, pCallbackFun sendTaskCallback) SendTaskParam {
	return SendTaskParam{chSend, conn, fd, pCallbackFun}
}

type sendTask struct {
	threadPool.Task
	m_taskParam SendTaskParam
}

func (this *sendTask)Run() error {
	errCode := error(nil)
	for ; ;  {
		if errCode = this.Task.Run(); nil != errCode {
			break
		}
		// 处理自己的逻辑
		defer close(this.m_taskParam.m_chSend)
		bClose := false
		for ; ;  {
			select {
			case i := <-this.m_taskParam.m_chSend:
				if dataDepot.SD_QUIT == i{
					bClose = true
				}else{
					// 发送数据
					for ; ;  {
						binPtr, err := binaryBuf.BinaryBufPtr(nil), error(nil)
						binPtr = this.m_taskParam.m_pCallbackFun(SendTask_msg_pop_data, this.m_taskParam.m_fd)
						if nil == binPtr {
							// 发送队列已空
							break
						}
						arrByte := binPtr.GetBuf()
						nLen := binPtr.GetSize()
						if 0 < nLen{
							nSend := 0
							if nSend, err = this.m_taskParam.m_conn.Write(arrByte); nil != err{
								// 出错 关闭socket
								binaryBuf.DestroyBinaryBuf(binPtr)
								this.m_taskParam.m_pCallbackFun(SendTask_msg_close, this.m_taskParam.m_fd)
								break
							}
							binaryBuf.DestroyBinaryBuf(binPtr)
							if nSend != nLen{
								// 出错 关闭socket
								this.m_taskParam.m_pCallbackFun(SendTask_msg_close, this.m_taskParam.m_fd)
								break
							}
						}else{
							this.m_taskParam.m_pCallbackFun(SendTask_msg_close, this.m_taskParam.m_fd)
						}
					}
				    }
			}
			if bClose {
				break
			}
		}
		break
	}
	return errCode
}
