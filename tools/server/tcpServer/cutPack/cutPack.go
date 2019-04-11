package cutPack

import (
	"fmt"
	"sync"
	"tools/binaryBuf"
	"tools/server/dataDepot"
	"tools/server/serverDefine"
	"tools/transfromBytes"
)

const cutpack_group int = 100

const(
	CUT_PACK_OK int = iota
	CUT_PACK_ERR
)

type cutPackCallback func(int, serverDefine.SOCKET_FD, binaryBuf.BinaryBufPtr)

type cutPack struct {
	m_listMapHalfPack [cutpack_group]map[serverDefine.SOCKET_FD]binaryBuf.BinaryBufPtr
	m_listMutex [cutpack_group]sync.Mutex

	m_pCallbackFun cutPackCallback
}

var g_pThis *cutPack = nil
func InitCutPack(pFun cutPackCallback) error {
	errCode := error(nil)
	for ; ;  {
		if nil != g_pThis {
			errCode = fmt.Errorf("curpack 已经初始化了")
			break
		}
		g_pThis = &cutPack{[cutpack_group]map[serverDefine.SOCKET_FD]binaryBuf.BinaryBufPtr{},
			[cutpack_group]sync.Mutex{},
			pFun}
		for i := 0; i < cutpack_group; i++ {
			g_pThis.m_listMapHalfPack[i] = make(map[serverDefine.SOCKET_FD]binaryBuf.BinaryBufPtr)
			g_pThis.m_listMutex[i] = sync.Mutex{}
		}
		break
	}
	return errCode
}

func CutPack(fd serverDefine.SOCKET_FD) error {
	errCode := error(nil)
	for ; ;  {
		if nil == g_pThis {
			errCode = fmt.Errorf("curpack 没有初始化")
			break
		}
		pDBManage := dataDepot.GetDepotManagePtr()
		if nil == pDBManage {
			errCode = fmt.Errorf("没有获取到仓库管理指针")
			break
		}
		pHalf := g_pThis.findLocakHalfPack(fd)
		pTemp := binaryBuf.BinaryBufPtr(nil)
		for ; ;  {
			if pTemp, errCode = pDBManage.PopRecvData(fd); nil != errCode{
				break
			}
			if nil == pTemp {
				break
			}
			if 0 == pTemp.GetSize() {
				binaryBuf.DestroyBinaryBuf(pTemp)
				errCode = fmt.Errorf("数据仓库中pop出来的数据错误， 长度为0 fd=%d", fd)
				break
			}

			if errCode = g_pThis.cutPack(&pHalf, &pTemp, fd); nil != errCode{
				if nil != pTemp {
					binaryBuf.DestroyBinaryBuf(pTemp)
				}
				break
			}
		}
		break
	}
	if nil != errCode {
		g_pThis.m_pCallbackFun(CUT_PACK_ERR, fd, nil)
	}
	return errCode
}

func ClearHalfPack(fd serverDefine.SOCKET_FD) error {
	errCode := error(nil)
	for ; ;  {
		if nil == g_pThis {
			errCode = fmt.Errorf("curpack 没有初始化")
			break
		}

		nIndex := int(fd % serverDefine.SOCKET_FD(cutpack_group))
		g_pThis.m_listMutex[nIndex].Lock()
		defer g_pThis.m_listMutex[nIndex].Unlock()
		if pHalf, ok := g_pThis.m_listMapHalfPack[nIndex][fd]; ok{
			if nil != pHalf {
				errCode = binaryBuf.DestroyBinaryBuf(pHalf)
			}
			delete(g_pThis.m_listMapHalfPack[nIndex], fd)
		}
		break
	}
	return errCode
}

func (this* cutPack)findLocakHalfPack(fd serverDefine.SOCKET_FD) binaryBuf.BinaryBufPtr{
	pRet := binaryBuf.BinaryBufPtr(nil)
	for ; ;  {
		nIndex := int(fd % serverDefine.SOCKET_FD(cutpack_group))
		this.m_listMutex[nIndex].Lock()
		defer this.m_listMutex[nIndex].Unlock()
		if pHalf, bFind := this.m_listMapHalfPack[nIndex][fd]; bFind{
			pRet = pHalf
		}
		break
	}
	return pRet
}

func (this* cutPack)popLocakHalfPack(fd serverDefine.SOCKET_FD){
	for ; ;  {
		nIndex := int(fd % serverDefine.SOCKET_FD(cutpack_group))
		this.m_listMutex[nIndex].Lock()
		defer this.m_listMutex[nIndex].Unlock()
		delete(this.m_listMapHalfPack[nIndex], fd)
		break
	}
}

func (this *cutPack)addLocalHalfPack(fd serverDefine.SOCKET_FD, ptr binaryBuf.BinaryBufPtr) error {
	errCode := error(nil)
	for ; ;  {
		nIndex := int(fd % serverDefine.SOCKET_FD(cutpack_group))
		this.m_listMutex[nIndex].Lock()
		defer this.m_listMutex[nIndex].Unlock()
		_, ok := this.m_listMapHalfPack[nIndex][fd]
		if ok {
			errCode = fmt.Errorf("socket %llu 添加半包失败， 已经有半包存在")
			break
		}
		this.m_listMapHalfPack[nIndex][fd] = ptr
		break
	}
	return errCode
}

func (this *cutPack)cutPackNotHalfPack(ppBuf *binaryBuf.BinaryBufPtr, fd serverDefine.SOCKET_FD) (binaryBuf.BinaryBufPtr, error) {
	pHalfRet, errCode:= binaryBuf.BinaryBufPtr(nil), error(nil)
	for ; ;  {
		nSize := (*ppBuf).GetSize()
		if 3 >= nSize {
			if errCode = this.addLocalHalfPack(fd, *ppBuf); nil != errCode{
				break
			}
			pHalfRet = *ppBuf
			break
		}

		// 数据长度中可以读出数据包的包长
		bufMem := (*ppBuf).GetBuf()
		packLenBuf := bufMem[:4]
		nPackLen := int(transfromBytes.BytesToInt32(packLenBuf))
		if 0 >= nPackLen || binaryBuf.BinaryBuf_64k < nPackLen {
			errCode = fmt.Errorf("数据包长度错误， len=%d", nPackLen)
			break
		}
		// 不够一个数据包
		if nSize < nPackLen {
			if errCode = this.addLocalHalfPack(fd, *ppBuf); nil != errCode{
				break
			}
			pHalfRet = *ppBuf
			break
		}
		// 刚好够一个数据包
		if nSize == nPackLen{
			this.m_pCallbackFun(CUT_PACK_OK, fd, *ppBuf)
			*ppBuf = nil
			break
		}

		// 超出一个包
		pNewBinaryTemp := binaryBuf.BinaryBufPtr(nil)
		if pNewBinaryTemp, errCode = binaryBuf.ProvideBinaryBuf(nPackLen); nil != errCode{
			break
		}
		if errCode = pNewBinaryTemp.AppendBuf(bufMem[:nPackLen]); nil != errCode{
			binaryBuf.DestroyBinaryBuf(pNewBinaryTemp)
			break
		}
		this.m_pCallbackFun(CUT_PACK_OK, fd, pNewBinaryTemp)
		if errCode = (*ppBuf).Seek(nPackLen); nil != errCode{
			break
		}
	}
	return pHalfRet, errCode
}

func (this *cutPack)cutPackHasHalfPackNoAnalyLen(ppHalf, ppNew *binaryBuf.BinaryBufPtr, fd serverDefine.SOCKET_FD) error {
	errCode := error(nil)
	for ; ;  {
		nHalfSize := (*ppHalf).GetSize()
		nHalfFreeSize := (*ppHalf).GetFreeSize()
		pHalfBuf := (*ppHalf).GetBuf()

		nNewSize := (*ppNew).GetSize()
		pNewBuf := (*ppNew).GetBuf()
		if nHalfFreeSize >= nNewSize {
			if errCode = (*ppHalf).AppendBuf(pNewBuf); nil != errCode{
				break
			}
			binaryBuf.DestroyBinaryBuf(*ppNew)
			*ppNew = nil
			break
		}

		// 半包空余长度不够
		pNewTemp := binaryBuf.BinaryBufPtr(nil)
		if pNewTemp, errCode = binaryBuf.ProvideBinaryBuf(nHalfSize + 3); nil != errCode{
			break
		}
		if errCode = pNewTemp.AppendBuf(pHalfBuf); nil != errCode{
			binaryBuf.DestroyBinaryBuf(pNewTemp)
			break
		}
		if errCode = pNewTemp.AppendBuf(pNewBuf); nil != errCode{
			binaryBuf.DestroyBinaryBuf(pNewTemp)
			break
		}
		binaryBuf.DestroyBinaryBuf(*ppNew)
		*ppNew = binaryBuf.BinaryBufPtr(nil)

		// 清空本地半包
		ClearHalfPack(fd)
		*ppHalf = binaryBuf.BinaryBufPtr(nil)
		// 添加本地半包
		if errCode = this.addLocalHalfPack(fd, pNewTemp); nil != errCode{
			binaryBuf.DestroyBinaryBuf(pNewTemp)
			break
		}
		*ppHalf = pNewTemp
		break
	}
	return errCode
}

func (this *cutPack)cutPackHasHalfPackHasAnalyLen(ppHalf, ppNew *binaryBuf.BinaryBufPtr, fd serverDefine.SOCKET_FD) error {
	errCode := error(nil)
	for ; ;  {
		nHalfSize := (*ppHalf).GetSize()
		pHalfBuf := (*ppHalf).GetBuf()
		nHalfFreeSize := (*ppHalf).GetFreeSize()
		pNewBuf := (*ppNew).GetBuf()
		nNewSize := (*ppNew).GetSize()

		nPackLen := 0
		if 3 >= nHalfSize {
			arrBuf := [4]byte{0,0,0,0}
			for i := 0; i < nHalfSize; i ++{
				arrBuf[i] = pHalfBuf[i]
			}
			for i := nHalfSize; i < 4; i ++{
				arrBuf[i] = pNewBuf[i - nHalfSize]
			}
			nPackLen = int(transfromBytes.BytesToInt32(arrBuf[:]))
		}else{
			nPackLen = int(transfromBytes.BytesToInt32(pHalfBuf[:4]))
		}
		if 0 >= nPackLen || binaryBuf.BinaryBuf_64k < nPackLen {
			errCode = fmt.Errorf("数据包长度错误， len=%d", nPackLen)
			break
		}

		if nPackLen > nHalfSize + nHalfFreeSize{
			// 半包不可以容纳整个数据包
			pNewTemp := binaryBuf.BinaryBufPtr(nil)
			if pNewTemp, errCode = binaryBuf.ProvideBinaryBuf(nPackLen); nil != errCode{
				break
			}
			if errCode = pNewTemp.AppendBuf(pHalfBuf); nil != errCode {
				binaryBuf.DestroyBinaryBuf(pNewTemp)
				break
			}
			if errCode = ClearHalfPack(fd); nil != errCode{
				binaryBuf.DestroyBinaryBuf(pNewTemp)
				break
			}
			*ppHalf = binaryBuf.BinaryBufPtr(nil)
			if errCode = this.addLocalHalfPack(fd, pNewTemp); nil != errCode{
				binaryBuf.DestroyBinaryBuf(pNewTemp)
				break
			}
			*ppHalf = pNewTemp
		}

		// 两份数据加起来不够一个整包
		if nPackLen > nHalfSize + nNewSize {
			if errCode = (*ppHalf).AppendBuf(pNewBuf); nil != errCode{
				break
			}
			binaryBuf.DestroyBinaryBuf(*ppNew)
			*ppNew = binaryBuf.BinaryBufPtr(nil)
			break
		}
		// 两份数据加起来刚好够一个整包
		if nPackLen == nHalfSize + nNewSize {
			if errCode = (*ppHalf).AppendBuf(pNewBuf); nil != errCode{
				break
			}
			binaryBuf.DestroyBinaryBuf(*ppNew)
			*ppNew = binaryBuf.BinaryBufPtr(nil)
			this.m_pCallbackFun(CUT_PACK_OK, fd, *ppHalf)
			this.popLocakHalfPack(fd)
			*ppHalf = binaryBuf.BinaryBufPtr(nil)
			break
		}

		// 两份数据加起来超过一个整包
		if errCode = (*ppHalf).AppendBuf(pNewBuf[:nPackLen - nHalfSize]); nil != errCode {
			break
		}
		this.m_pCallbackFun(CUT_PACK_OK, fd, *ppHalf)
		this.popLocakHalfPack(fd)
		*ppHalf = binaryBuf.BinaryBufPtr(nil)
		if errCode = (*ppNew).Seek(nPackLen - nHalfSize); nil != errCode {
			break
		}
		*ppHalf, errCode = this.cutPackNotHalfPack(ppNew, fd)
 		break
	}
	return errCode
}

func (this *cutPack)cutPack(ppHalf, ppNew *binaryBuf.BinaryBufPtr, fd serverDefine.SOCKET_FD) error  {
	errCode := error(nil)
	for ; ;  {
		if nil == (*ppHalf) {
			(*ppHalf), errCode= this.cutPackNotHalfPack(ppNew, fd)
			break
		}
		// half pack is not nil
		nHalfSize := (*ppHalf).GetSize()
		nNewSize := (*ppNew).GetSize()
		if 3 >= nHalfSize + nNewSize {
			errCode = this.cutPackHasHalfPackNoAnalyLen(ppHalf, ppNew, fd)
			break
		}
		errCode = this.cutPackHasHalfPackHasAnalyLen(ppHalf, ppNew, fd)
		break
	}
	return errCode
}

