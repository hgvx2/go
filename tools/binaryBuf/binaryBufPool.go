package binaryBuf

import (
	"fmt"
	"tools/memoryPool"
)

const(
	mutex_32 int = iota
	mutex_64
	mutex_128
	mutex_256
	mutex_512
	mutex_1k
	mutex_2k
	mutex_4k
	mutex_8k
	mutex_16k
	mutex_32k
	mutex_64k
	mutex_total
)

var g_pThis *binaryBufPool = nil

type binaryBufPool struct {
	m_binaryBufPool_32 memoryPool.MemoryPool
	m_binaryBufPool_64 memoryPool.MemoryPool
	m_binaryBufPool_128 memoryPool.MemoryPool
	m_binaryBufPool_256 memoryPool.MemoryPool
	m_binaryBufPool_512 memoryPool.MemoryPool
	m_binaryBufPool_1k memoryPool.MemoryPool
	m_binaryBufPool_2k memoryPool.MemoryPool
	m_binaryBufPool_4k memoryPool.MemoryPool
	m_binaryBufPool_8k memoryPool.MemoryPool
	m_binaryBufPool_16k memoryPool.MemoryPool
	m_binaryBufPool_32k memoryPool.MemoryPool
	m_binaryBufPool_64k memoryPool.MemoryPool
}

func InitBinaryBufPool(nMin int) error{
	errCode := error(nil)
	for ; ;  {
		if 0 > nMin {
			fmt.Errorf("nMin的值不能小于0")
			break
		}
		if nil == g_pThis {
			g_pThis = &binaryBufPool{memoryPool.GetMemoryPool(),
				memoryPool.GetMemoryPool(),
				memoryPool.GetMemoryPool(),
				memoryPool.GetMemoryPool(),
				memoryPool.GetMemoryPool(),
				memoryPool.GetMemoryPool(),
				memoryPool.GetMemoryPool(),
				memoryPool.GetMemoryPool(),
				memoryPool.GetMemoryPool(),
				memoryPool.GetMemoryPool(),
				memoryPool.GetMemoryPool(),
				memoryPool.GetMemoryPool()}
		}

		for i := 0; i < nMin; i ++ {
			ptr, err := BinaryBufPtr(nil), error(nil)
			if ptr, err = newBinaryBuf(BinaryBuf_32); nil == err{
				errCode = err
				break
			}
			g_pThis.m_binaryBufPool_32.AddNewNode(ptr)

			if ptr, err = newBinaryBuf(BinaryBuf_64); nil == err{
				errCode = err
				break
			}
			g_pThis.m_binaryBufPool_64.AddNewNode(ptr)

			if ptr, err = newBinaryBuf(BinaryBuf_128); nil == err{
				errCode = err
				break
			}
			g_pThis.m_binaryBufPool_128.AddNewNode(ptr)

			if ptr, err = newBinaryBuf(BinaryBuf_256); nil == err{
				errCode = err
				break
			}
			g_pThis.m_binaryBufPool_256.AddNewNode(ptr)

			if ptr, err = newBinaryBuf(BinaryBuf_512); nil == err{
				errCode = err
				break
			}
			g_pThis.m_binaryBufPool_512.AddNewNode(ptr)

			if ptr, err = newBinaryBuf(BinaryBuf_1k); nil == err{
				errCode = err
				break
			}
			g_pThis.m_binaryBufPool_1k.AddNewNode(ptr)

			if ptr, err = newBinaryBuf(BinaryBuf_2k); nil == err{
				errCode = err
				break
			}
			g_pThis.m_binaryBufPool_2k.AddNewNode(ptr)

			if ptr, err = newBinaryBuf(BinaryBuf_4k); nil == err{
				errCode = err
				break
			}
			g_pThis.m_binaryBufPool_4k.AddNewNode(ptr)

			if ptr, err = newBinaryBuf(BinaryBuf_8k); nil == err{
				errCode = err
				break
			}
			g_pThis.m_binaryBufPool_8k.AddNewNode(ptr)

			if ptr, err = newBinaryBuf(BinaryBuf_16k); nil == err{
				errCode = err
				break
			}
			g_pThis.m_binaryBufPool_16k.AddNewNode(ptr)

			if ptr, err = newBinaryBuf(BinaryBuf_32k); nil == err{
				errCode = err
				break
			}
			g_pThis.m_binaryBufPool_32k.AddNewNode(ptr)

			if ptr, err = newBinaryBuf(BinaryBuf_64k); nil == err{
				errCode = err
				break
			}
			g_pThis.m_binaryBufPool_64k.AddNewNode(ptr)
		}
		break
	}
	return errCode
}

func ProvideBinaryBuf(nSize int) (BinaryBufPtr, error)  {
	pRet, errCode := BinaryBufPtr(nil), error(nil)
	for ; ;  {
		if nil == g_pThis {
			errCode = fmt.Errorf("BinaryBufPool 没有初始化")
			break
		}
		if nSize <= BinaryBuf_32 {
			pRet, errCode = g_pThis.provideBinaryBuf(BinaryBuf_32, &g_pThis.m_binaryBufPool_32)
			break
		}
		if nSize <= BinaryBuf_64 {
			pRet, errCode = g_pThis.provideBinaryBuf(BinaryBuf_64, &g_pThis.m_binaryBufPool_64)
			break
		}
		if nSize <= BinaryBuf_128 {
			pRet, errCode = g_pThis.provideBinaryBuf(BinaryBuf_128, &g_pThis.m_binaryBufPool_128)
			break
		}
		if nSize <= BinaryBuf_256 {
			pRet, errCode = g_pThis.provideBinaryBuf(BinaryBuf_256, &g_pThis.m_binaryBufPool_256)
			break
		}
		if nSize <= BinaryBuf_512 {
			pRet, errCode = g_pThis.provideBinaryBuf(BinaryBuf_512, &g_pThis.m_binaryBufPool_512)
			break
		}
		if nSize <= BinaryBuf_1k {
			pRet, errCode = g_pThis.provideBinaryBuf(BinaryBuf_1k, &g_pThis.m_binaryBufPool_1k)
			break
		}
		if nSize <= BinaryBuf_2k {
			pRet, errCode = g_pThis.provideBinaryBuf(BinaryBuf_2k, &g_pThis.m_binaryBufPool_2k)
			break
		}
		if nSize <= BinaryBuf_4k {
			pRet, errCode = g_pThis.provideBinaryBuf(BinaryBuf_4k, &g_pThis.m_binaryBufPool_4k)
			break
		}
		if nSize <= BinaryBuf_8k {
			pRet, errCode = g_pThis.provideBinaryBuf(BinaryBuf_8k, &g_pThis.m_binaryBufPool_8k)
			break
		}
		if nSize <= BinaryBuf_16k {
			pRet, errCode = g_pThis.provideBinaryBuf(BinaryBuf_16k, &g_pThis.m_binaryBufPool_16k)
			break
		}
		if nSize <= BinaryBuf_32k {
			pRet, errCode = g_pThis.provideBinaryBuf(BinaryBuf_32k, &g_pThis.m_binaryBufPool_32k)
			break
		}
		if nSize <= BinaryBuf_64k {
			pRet, errCode = g_pThis.provideBinaryBuf(BinaryBuf_64k, &g_pThis.m_binaryBufPool_64k)
			break
		}else{
			errCode = fmt.Errorf("没有找到对于的类型 nSize=%d", nSize)
		}
		break
	}
	if nil != pRet {
		pRet.ClearBuf()
	}
	return pRet, errCode
}

func DestroyBinaryBuf(ptr BinaryBufPtr) error{
	errCode := error(nil)
	for ; ;  {
		if nil == g_pThis {
			errCode = fmt.Errorf("BinaryBufPool 没有初始化")
			break
		}
		if nil == ptr {
			errCode = fmt.Errorf("BinaryBufPool 空指针")
			break
		}
		switch ptr.getType() {
		case BinaryBuf_32:
			errCode = g_pThis.m_binaryBufPool_32.DestoryNode(ptr)
		case BinaryBuf_64:
			errCode = g_pThis.m_binaryBufPool_64.DestoryNode(ptr)
		case BinaryBuf_128:
			errCode = g_pThis.m_binaryBufPool_128.DestoryNode(ptr)
		case BinaryBuf_256:
			errCode = g_pThis.m_binaryBufPool_256.DestoryNode(ptr)
		case BinaryBuf_512:
			errCode = g_pThis.m_binaryBufPool_512.DestoryNode(ptr)
		case BinaryBuf_1k:
			errCode = g_pThis.m_binaryBufPool_1k.DestoryNode(ptr)
		case BinaryBuf_2k:
			errCode = g_pThis.m_binaryBufPool_2k.DestoryNode(ptr)
		case BinaryBuf_4k:
			errCode = g_pThis.m_binaryBufPool_4k.DestoryNode(ptr)
		case BinaryBuf_8k:
			errCode = g_pThis.m_binaryBufPool_8k.DestoryNode(ptr)
		case BinaryBuf_16k:
			errCode = g_pThis.m_binaryBufPool_64k.DestoryNode(ptr)
		case BinaryBuf_32k:
			errCode = g_pThis.m_binaryBufPool_32k.DestoryNode(ptr)
		case BinaryBuf_64k:
			errCode = g_pThis.m_binaryBufPool_64k.DestoryNode(ptr)
		default:
			errCode = fmt.Errorf("没有找到对于的binarybuf类型 type=%d", ptr.getType())
		}
		break
	}
	return errCode
}

func (this *binaryBufPool)provideBinaryBuf(nType int, pPool memoryPool.MemoryPoolPtr) (BinaryBufPtr, error) {
	pRet, errCode := BinaryBufPtr(nil), error(nil)
	for ; ;  {
		pVoid := pPool.ProvideNode()
		if nil == pVoid {
			if pRet, errCode = newBinaryBuf(nType); nil != errCode{
				break
			}
			if errCode = pPool.AddNewNode(pRet); nil != errCode{
				break
			}
		}else{
			pRet = pVoid.(BinaryBufPtr)
		}
		break
	}
	return pRet, errCode
}
