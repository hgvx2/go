package dataDepot

import (
	"tools/binaryBuf"
)

const (
	SD_CONTINUE int = 1
	SD_QUIT     int = 2
)

type sendDepot struct {
	depotBase
	m_sendChan chan <-int
}

func (this *sendDepot)clearData()  {
	this.depotBase.clearData()
	this.m_sendChan <- SD_QUIT
}

func (this *sendDepot)pushData(ptr binaryBuf.BinaryBufPtr) int {
	nRet := DM_PUSH_ERR
	for ; ;  {
		if nRet = this.depotBase.pushData(ptr); DM_PUSH_ET == nRet{
			this.m_sendChan <- SD_CONTINUE
		}
		break
	}
	return nRet
}