package connManage

import "tools/memoryPool"

type connLivePool struct {
	memoryPool.MemoryPool
}

func getConnLivePool() connLivePool {
	return connLivePool{memoryPool.GetMemoryPool()}
}

func (this *connLivePool)provideConLive() (*connLive, error){
	pRet, errCode := (*connLive)(nil), error(nil)
	for ; ;  {
		pTemp := this.MemoryPool.ProvideNode()
		if nil != pTemp {
			pRet = pTemp.(*connLive)
			break
		}

		temp := connLive{0,nil,0}
		if errCode = this.MemoryPool.AddNewNode(&temp); nil != errCode{
			break
		}
		pRet = &temp
		break
	}
	return pRet, errCode
}

func (this *connLivePool)destoryConLive(pConLive *connLive) error {
	pConLive.m_fd = 0
	pConLive.m_pCon = nil
	pConLive.m_unixTime = 0
	return this.DestoryNode(pConLive)
}
