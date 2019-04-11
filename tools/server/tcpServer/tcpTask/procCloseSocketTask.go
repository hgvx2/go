package tcpTask

import (
	"tools/server/serverDefine"
	"tools/threadPool"
)

type closeSocketTaskCallback func (serverDefine.SOCKET_FD)

type ProcCloseSocketTaskParam struct {
	m_chClose chan serverDefine.SOCKET_FD
	m_pCallbackFun closeSocketTaskCallback
}

func GetProcCloseSocketTaskParam(chClose chan serverDefine.SOCKET_FD, pCallbackFun closeSocketTaskCallback) ProcCloseSocketTaskParam {
	return ProcCloseSocketTaskParam{chClose,pCallbackFun}
}

type procCloseSocketTask struct {
	threadPool.Task
	m_taskParam ProcCloseSocketTaskParam
}

func (this *procCloseSocketTask)Run() error {
	errCode := error(nil)
	for ; ;  {
		if errCode = this.Task.Run(); nil != errCode {
			break
		}
		// 处理自己的逻辑
		for ; ;  {
			bChClose, fd := false, serverDefine.SOCKET_FD(0)
			select {
				case fd, bChClose = <-this.m_taskParam.m_chClose:
					if bChClose {
						this.m_taskParam.m_pCallbackFun(fd)
					}
			}
			if !bChClose {
				// 通道已关闭
				this.m_taskParam.m_pCallbackFun(fd)
				break
			}
		}
		break
	}
	return errCode
}
