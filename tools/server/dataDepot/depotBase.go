package dataDepot

import (
	"tools/binaryBuf"
)

type depotBase struct {
	m_ch chan binaryBuf.BinaryBufPtr
	m_bPopRunning bool
}

func (this *depotBase) pushData(ptr binaryBuf.BinaryBufPtr) int {
	nRet := DM_PUSH_ERR
	for ; ;  {
		bPush := true
		select {
		case this.m_ch <- ptr:
		default:
			bPush = false
		}
		if !bPush {
			break
		}
		if !this.m_bPopRunning {
			this.m_bPopRunning = true
			nRet = DM_PUSH_ET
			break
		}
		nRet = DM_PUSH_OK
		break
	}
	return nRet

}

func (this *depotBase) popData() binaryBuf.BinaryBufPtr{
	pRet := binaryBuf.BinaryBufPtr(nil)
	select {
	case pRet = <-this.m_ch:
	default:
		this.m_bPopRunning = false
	}
	return pRet
}

func (this *depotBase) clearData() {
	for ; ;  {
		pRet, bContinue := binaryBuf.BinaryBufPtr(nil), false
		select {
		case pRet ,bContinue= <-this.m_ch:
			if nil !=  pRet{
				binaryBuf.DestroyBinaryBuf(pRet)
			}
			if !bContinue {
				this.m_bPopRunning = false
			}
		default:
			bContinue = false
			this.m_bPopRunning = false
		}
		if !bContinue {
			break
		}
	}
}

