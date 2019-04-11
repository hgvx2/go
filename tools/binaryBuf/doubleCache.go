package binaryBuf

import (
	"sync"
)

type DoubleCachePtr interface {
	Push(ptr BinaryBufPtr)
	Pop()  BinaryBufPtr
	Clear()
}

type doubleCache struct {
	m_inputCache []BinaryBufPtr
	m_outputCache []BinaryBufPtr

	m_pInputCache *[]BinaryBufPtr
	m_pOutPutCache *[]BinaryBufPtr

	m_inputMutex sync.Mutex
	m_outputMutex sync.Mutex
}

func GetDoubleCachePtr() (DoubleCachePtr, error)  {
	pRet, errCode := DoubleCachePtr(nil), error(nil)
	for ; ;  {
		pTemp := &doubleCache{make([]BinaryBufPtr, 0, 100),
			make([]BinaryBufPtr, 0, 100), nil, nil,
		sync.Mutex{}, sync.Mutex{}}

		pTemp.m_pInputCache = &pTemp.m_inputCache
		pTemp.m_pOutPutCache = &pTemp.m_outputCache
		pRet = pTemp
		break
	}
	return pRet, errCode
}

func (this *doubleCache) Push(ptr BinaryBufPtr) {
	this.m_inputMutex.Lock()
	defer this.m_inputMutex.Unlock()
	*this.m_pInputCache = append(*this.m_pInputCache, ptr)
}

func (this *doubleCache) Pop() BinaryBufPtr{
	pRet:= BinaryBufPtr(nil)
	for ; ;  {
		this.m_outputMutex.Lock()
		defer this.m_outputMutex.Unlock()

		nLen := len(*this.m_pOutPutCache)
		if  0 < nLen{
			pRet = (*this.m_pOutPutCache)[0]
			*this.m_pOutPutCache = (*this.m_pOutPutCache)[1:]
			break
		}
		this.m_inputMutex.Lock()
		defer this.m_inputMutex.Unlock()
		nLen = len(*this.m_pInputCache)
		if 0 < nLen{
			this.m_pInputCache, this.m_pOutPutCache = this.m_pOutPutCache, this.m_pInputCache
			pRet = (*this.m_pOutPutCache)[0]
			*this.m_pOutPutCache = (*this.m_pOutPutCache)[1:]
			break
		}
		break
	}
	return pRet
}

func (this *doubleCache) Clear() {
	this.m_outputMutex.Lock()
	this.m_inputMutex.Lock()
	defer func() {
		this.m_outputMutex.Unlock()
		this.m_inputMutex.Unlock()
	}()

	for _, pBinaryPtr := range this.m_inputCache {
		DestroyBinaryBuf(pBinaryPtr)
	}
	this.m_inputCache = this.m_inputCache[0:0]

	for _, pBinaryPtr := range this.m_outputCache {
		DestroyBinaryBuf(pBinaryPtr)
	}
	this.m_outputCache = this.m_outputCache[0:0]

}