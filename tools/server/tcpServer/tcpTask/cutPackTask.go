package tcpTask

import (
	"tools/server/serverDefine"
	"tools/server/tcpServer/cutPack"
	"tools/threadPool"
)

type CutPackTaskParam struct {
	m_fd serverDefine.SOCKET_FD
}
type cutPackTask struct {
	threadPool.Task
	m_taskParam CutPackTaskParam
}

func GetCutPackTaskParam(fd serverDefine.SOCKET_FD) CutPackTaskParam {
	return CutPackTaskParam{fd}
}

func (this *cutPackTask)Run() error {
	errCode := error(nil)
	for ; ;  {
		if errCode = this.Task.Run(); nil != errCode {
			break
		}
		// 处理自己的逻辑
		cutPack.CutPack(this.m_taskParam.m_fd)
		break
	}
	return errCode
}

