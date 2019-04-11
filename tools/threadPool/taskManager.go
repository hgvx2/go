package threadPool

import (
	"fmt"
	"sync"
)

type TaskManagerPtr interface {
	ProvideTask(taskType int, pVoid interface{}) (TaskPtr, error)
	DestoryTask(pTask TaskPtr) error
}

type TaskManage struct {
	m_mapTaskPool map[int]*taskPool
}

func GetTaskManage(arrTaskTypes []int) TaskManage {
	tm := TaskManage{make(map[int]*taskPool)}
	for _, taskType := range arrTaskTypes{
		tp:= taskPool{make([]TaskPtr, 0, 1000),
			make(map[TaskPtr]TaskPtr),
		sync.Mutex{}}
		tm.m_mapTaskPool[taskType] = &tp
	}
	return tm
}

func (this *TaskManage) ProvideTask(taskType int, pVoid interface{}) (TaskPtr, error) {
	pTask, errCode := TaskPtr(nil), error(nil)
	for ; ;  {
		pTaskPool, bFind := (*taskPool)(nil), false
		if pTaskPool, bFind = this.m_mapTaskPool[taskType]; !bFind {
			errCode = fmt.Errorf("ProvideTask 没有找到指定的任务类型 type=%d", taskType)
			break
		}
		if nil == pTaskPool {
			errCode = fmt.Errorf("ProvideTask 任务池空指针， type=%d ", taskType)
			break
		}
		pTask = pTaskPool.provideTask()
		break
	}
	return pTask, errCode
}

func (this *TaskManage) DestoryTask(pTask TaskPtr) error  {
	errCode := error(nil)
	for ; ;  {
		if nil == pTask {
			errCode = fmt.Errorf("DestoryTask 空指针")
			break
		}
		taskType := pTask.getTaskType()
		pTaskPool, bFind := (*taskPool)(nil), false
		if pTaskPool, bFind = this.m_mapTaskPool[taskType]; !bFind {
			errCode = fmt.Errorf("ProvideTask 没有找到指定的任务类型 type=%d", taskType)
			break
		}
		if nil == pTaskPool {
			errCode = fmt.Errorf("ProvideTask 任务池空指针， type=%d ", taskType)
			break
		}
		pTaskPool.destoryTask(pTask)
		break
	}
	return errCode
}

func (this *TaskManage)AddNewTask(pTask TaskPtr) error {
	errCode := error(nil)
	for ; ;  {
		if nil == pTask {
			errCode = fmt.Errorf("AddNewTask 空指针")
			break
		}
		taskType := pTask.getTaskType()
		pTaskPool, bFind := (*taskPool)(nil), false
		if pTaskPool, bFind = this.m_mapTaskPool[taskType]; !bFind {
			errCode = fmt.Errorf("AddNewTask 没有找到指定的任务类型 type=%d", taskType)
			break
		}
		if nil == pTaskPool {
			errCode = fmt.Errorf("AddNewTask 任务池空指针， type=%d ", taskType)
			break
		}
		errCode = pTaskPool.addNewTask(pTask)
		break
	}
	return errCode
}
