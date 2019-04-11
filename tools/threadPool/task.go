package threadPool

import "fmt"

type TaskPtr interface {
	Run() error
	DestoryTaskPtr() error
	getTaskType() int
}

type Task struct {
	m_pTaskManage TaskManagerPtr
	m_nType int
}

func GetTask(ptr TaskManagerPtr, taskType int) Task {
	return Task{ptr, taskType}
}

func (this *Task)Run() error {
	return error(nil)
}

func (this *Task)DestoryTaskPtr() error  {
	errCode := error(nil)
	for ; ;  {
		if nil == this.m_pTaskManage {
			errCode = fmt.Errorf("this.m_pTaskManage is null")
			break
		}
		//errCode = this.m_pTaskManage.(this)
		break
	}
	return errCode
}

func (this *Task)getTaskType() int {
	return this.m_nType
}