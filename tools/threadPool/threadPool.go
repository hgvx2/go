package threadPool

import (
	"fmt"
	"sync"
)

const (
	tp_max_free_thread_nums int = 15
	tp_min_free_thread_nums int = 10
)

type ThreadPool struct {
	m_chWork chan TaskPtr
	m_nMaxWorkNums  int
	m_nFreeWorkNums int
	m_nUseWorkNums  int
	m_bIsRunning bool
	m_wait sync.WaitGroup
	m_mutex sync.Mutex
}

func GetThreadPool(nMaxWorkNums, nCacheTaskNums int) ThreadPool {
	return ThreadPool{make(chan TaskPtr, nCacheTaskNums),
		nMaxWorkNums,
		0,
		0,
		false,
		sync.WaitGroup{},
		sync.Mutex{}}
}

func (this *ThreadPool)ThreadRunning() error {
	errCode := error(nil)
	for ; ;  {
		this.m_mutex.Lock()
		if this.m_bIsRunning {
			this.m_mutex.Unlock()
			errCode = fmt.Errorf("threadPool 已经运行了， 不能再次运行")
			break
		}
		this.m_bIsRunning = true;
		this.m_mutex.Unlock()
		break
	}
	return errCode
}

func (this *ThreadPool)StopThreadPool()  {
	this.m_mutex.Lock()
	this.m_bIsRunning = false
	this.m_mutex.Unlock()
	close(this.m_chWork)

	this.m_wait.Wait()
}

func (this *ThreadPool)PushTaskToBack(ptr TaskPtr) error {
	errCode := error(nil)
	for ; ;  {
		this.m_mutex.Lock()
		if !this.m_bIsRunning {
			this.m_mutex.Unlock()
			errCode = fmt.Errorf("threadPool 没有运行")
			break
		}
		this.m_mutex.Unlock()
		if nil == ptr {
			errCode = fmt.Errorf("PushTaskToBack ptr 空指针")
			break
		}

		this.m_mutex.Lock()
		if tp_min_free_thread_nums >= this.m_nFreeWorkNums &&
			this.m_nFreeWorkNums + this.m_nUseWorkNums < this.m_nMaxWorkNums{
			go this.runWork()
			this.m_wait.Add(1)
		}
		this.m_mutex.Unlock()

		this.m_chWork <-ptr
		break
	}
	return errCode
}

func (this *ThreadPool)runWork() {
	bOut := false
	for ; ;  {
		this.m_mutex.Lock()
		if !this.m_bIsRunning || bOut{
			this.m_mutex.Unlock()
			this.m_wait.Done()
			break;
		}
		this.m_mutex.Unlock()

		for ; ;  {
			pTask, bClose := TaskPtr(nil), false

			this.m_mutex.Lock()
			this.m_nFreeWorkNums += 1
			if this.m_nUseWorkNums > 0 {
				this.m_nUseWorkNums -= 1
			}
			if tp_max_free_thread_nums < this.m_nFreeWorkNums {
				this.m_nFreeWorkNums -= 1
				bOut = true
				this.m_mutex.Unlock()
				break
			}
			this.m_mutex.Unlock()

			select {
			case pTask, bClose = <- this.m_chWork:
				this.m_mutex.Lock()
				this.m_nFreeWorkNums -= 1
				this.m_nUseWorkNums += 1
				this.m_mutex.Unlock()

				pTask.Run()
				pTask.DestoryTaskPtr()
			}
			if !bClose {
				// channel is close
				break;
			}
		}
	}
}


