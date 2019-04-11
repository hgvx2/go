package threadPool

import (
	"sync"
)

type taskPool struct {
	m_listFreeTask []TaskPtr
	m_mapUseTask map[TaskPtr]TaskPtr
	m_mutex sync.Mutex
}

func (this *taskPool)provideTask() TaskPtr {
	pTask := TaskPtr(nil)
	for ; ;  {
		this.m_mutex.Lock()
		defer this.m_mutex.Unlock()
		if 0 == len(this.m_listFreeTask) {
			break
		}
		pTask = this.m_listFreeTask[0]
		this.m_mapUseTask[pTask] = pTask
		this.m_listFreeTask = this.m_listFreeTask[1:]
		break
	}

	return pTask
}

func (this *taskPool)destoryTask(pTask TaskPtr) {
	this.m_mutex.Lock()
	defer this.m_mutex.Unlock()

	delete(this.m_mapUseTask, pTask)
	this.m_listFreeTask = append(this.m_listFreeTask, pTask)
}

func (this *taskPool)addNewTask(pTask TaskPtr) error  {
	errCode := error(nil)
	for ; ;  {
		this.m_mutex.Lock()
		defer this.m_mutex.Unlock()
		this.m_mapUseTask[pTask] = pTask
		break
	}
	return errCode
}