package tcpTask

import (
	"fmt"
	"tools/threadPool"
)

const(
	Task_type_cutpack int = iota
	Task_type_proc_closesocket
	Task_type_recv
	Task_type_send
)

type tcpTaskManage struct {
	threadPool.TaskManage
}


func GetTcpTaskManagePtr() threadPool.TaskManagerPtr {
	return &tcpTaskManage{threadPool.GetTaskManage([]int{Task_type_cutpack,
		Task_type_proc_closesocket,
		Task_type_recv,
		Task_type_send})}
}

func (this *tcpTaskManage) ProvideTask(taskType int, pVoid interface{})  (threadPool.TaskPtr, error) {
	pRet, errCode := threadPool.TaskPtr(nil), error(nil)
	for ; ;  {
		if pRet, errCode = this.TaskManage.ProvideTask(taskType, pVoid); nil != errCode{
			break
		}
		if nil != pRet {
			break
		}

		switch taskType {
		case Task_type_cutpack:
			pRet, errCode = this.newCutpackTask(pVoid)
		case Task_type_proc_closesocket:
			pRet, errCode = this.newProcCloseSocketTask(pVoid)
		case Task_type_recv:
			pRet, errCode = this.newRecvTask(pVoid)
		case Task_type_send:
			pRet, errCode = this.newSendTask(pVoid)
		default:
			errCode = fmt.Errorf("没有找到指定的任务类型 type=%d", taskType)
		}
		break
	}
	return pRet, errCode
}

func (this *tcpTaskManage)newCutpackTask(pVoid interface{}) (threadPool.TaskPtr, error) {
	pRet, errCode := threadPool.TaskPtr(nil), error(nil)
	for ; ;  {
		param := pVoid.(CutPackTaskParam)
		pRet = &cutPackTask{threadPool.GetTask(this, Task_type_cutpack), param}
		errCode = this.AddNewTask(pRet)
		if nil != errCode {
			pRet = nil
		}
		break
	}
	return pRet, errCode
}

func (this *tcpTaskManage)newProcCloseSocketTask(pVoid interface{}) (threadPool.TaskPtr, error) {
	pRet, errCode := threadPool.TaskPtr(nil), error(nil)
	for ; ;  {
		param := pVoid.(ProcCloseSocketTaskParam)
		pRet = &procCloseSocketTask{threadPool.GetTask(this, Task_type_proc_closesocket), param}
		errCode = this.AddNewTask(pRet)
		if nil != errCode {
			pRet = nil
		}
		break
	}
	return pRet, errCode
}

func (this *tcpTaskManage)newRecvTask(pVoid interface{}) (threadPool.TaskPtr, error) {
	pRet, errCode := threadPool.TaskPtr(nil), error(nil)
	for ; ;  {
		param := pVoid.(RecvTaskParam)
		pRet = &recvTask{threadPool.GetTask(this, Task_type_recv), param}
		errCode = this.AddNewTask(pRet)
		if nil != errCode {
			pRet = nil
		}
		break
	}
	return pRet, errCode
}

func (this *tcpTaskManage)newSendTask(pVoid interface{}) (threadPool.TaskPtr, error) {
	pRet, errCode := threadPool.TaskPtr(nil), error(nil)
	for ; ;  {
		param := pVoid.(SendTaskParam)
		pRet = &sendTask{threadPool.GetTask(this, Task_type_send), param}
		errCode = this.AddNewTask(pRet)
		if nil != errCode {
			pRet = nil
		}
		break
	}
	return pRet, errCode
}

