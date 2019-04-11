package binaryBuf

import "fmt"

const(
	BinaryBuf_32  int = 32
	BinaryBuf_64  int = 64
	BinaryBuf_128 int = 128
	BinaryBuf_256 int = 256
	BinaryBuf_512 int = 512
	BinaryBuf_1k  int = 1024
	BinaryBuf_2k  int = 2048
	BinaryBuf_4k  int = 4096
	BinaryBuf_8k  int = 8192
	BinaryBuf_16k int = 16384
	BinaryBuf_32k int = 32768
	BinaryBuf_64k int = 65536
)

type BinaryBufPtr interface {
	AppendBuf(arrByte []byte) error
	ClearBuf()
	GetBuf() []byte
	getType() int
	GetFreeSize()int
	Seek(nSk int) error
	GetSize() int
}

type binaryBufSilce struct {
	m_bufSilce []byte
}

type binaryBuf struct {
	m_nSindex int
	m_nEindex int
	m_nType   int
	m_pBufSilce binaryBufSilcePtr
}
type binaryBufSilcePtr interface {
	appendBuf(arrByte []byte, pSindex, pEindex *int)
	getBuf(nSIndex, nEIndex int) []byte
}

type arrByte32 struct {
	binaryBufSilce
	m_buf [BinaryBuf_32]byte
}
type arrByte64  struct {
	binaryBufSilce
	m_buf [BinaryBuf_64]byte
}
type arrByte128 struct {
	binaryBufSilce
	m_buf [BinaryBuf_128]byte
}
type arrByte256 struct {
	binaryBufSilce
	m_buf [BinaryBuf_256]byte
}
type arrByte512 struct {
	binaryBufSilce
	m_buf[BinaryBuf_512]byte
}
type arrByte1k  struct {
	binaryBufSilce
	m_buf[BinaryBuf_1k]byte
}
type arrByte2k  struct {
	binaryBufSilce
	m_buf[BinaryBuf_2k]byte
}
type arrByte4k  struct {
	binaryBufSilce
	m_buf[BinaryBuf_4k]byte
}
type arrByte8k  struct {
	binaryBufSilce
	m_buf[BinaryBuf_8k]byte
}
type arrByte16k struct {
	binaryBufSilce
	m_buf[BinaryBuf_16k]byte
}
type arrByte32k struct {
	binaryBufSilce
	m_buf[BinaryBuf_32k]byte
}
type arrByte64k struct {
	binaryBufSilce
	m_buf[BinaryBuf_64k]byte
}

type binaryBuf32 struct {
	binaryBuf
	m_arrByte arrByte32
}

type binaryBuf64 struct {
	binaryBuf
	m_arrByte arrByte64
}

type binaryBuf128 struct {
	binaryBuf
	m_arrByte arrByte128
}

type binaryBuf256 struct {
	binaryBuf
	m_arrByte arrByte256
}

type binaryBuf512 struct {
	binaryBuf
	m_arrByte arrByte512
}

type binaryBuf1k struct {
	binaryBuf
	m_arrByte arrByte1k
}

type binaryBuf2k struct {
	binaryBuf
	m_arrByte arrByte2k
}

type binaryBuf4k struct {
	binaryBuf
	m_arrByte arrByte4k
}

type binaryBuf8k struct {
	binaryBuf
	m_arrByte arrByte8k
}

type binaryBuf16k struct {
	binaryBuf
	m_arrByte arrByte16k
}

type binaryBuf32k struct {
	binaryBuf
	m_arrByte arrByte32k
}

type binaryBuf64k struct {
	binaryBuf
	m_arrByte arrByte64k
}

func (this *binaryBufSilce) appendBuf(arrByte []byte, pSindex, pEindex *int)  {
	appendBuf(this.m_bufSilce, arrByte, pSindex, pEindex)
}

func (this *binaryBufSilce)getBuf(nSIndex, nEIndex int) []byte  {
	return this.m_bufSilce[nSIndex:nEIndex]
}

func appendBuf(dstByte, srcByte []byte, pSindex, pEindex *int)  {
	temp := dstByte[*pSindex : *pEindex]
	temp = append(temp[:], srcByte...)
	*pEindex += len(srcByte)
}

func (this *binaryBuf) AppendBuf(arrByte []byte) error {
	errCode := error(nil)
	for ; ;  {
		copyLen, cacheLen := len(arrByte), this.m_nType - this.m_nEindex
		if copyLen > cacheLen {
			errCode = fmt.Errorf("缓存长度不够，无法拷贝， arrByte=%d cacheLen=%d", copyLen,
				cacheLen)
			break
		}
		if 0 >= copyLen {
			break
		}
		this.m_pBufSilce.appendBuf(arrByte, &this.m_nSindex, &this.m_nEindex)
		break
	}
	return errCode;
}

func (this *binaryBuf) ClearBuf() {
	this.m_nSindex = 0
	this.m_nEindex = 0
}

func (this *binaryBuf) getType() int{
	return this.m_nType
}

func (this *binaryBuf) GetBuf() []byte {
	return this.m_pBufSilce.getBuf(this.m_nSindex, this.m_nEindex)
}

func (this *binaryBuf)GetFreeSize() int{
	return this.m_nType - this.m_nEindex
}
func (this *binaryBuf)Seek(nSk int) error{
	errCode := error(nil)
	for ; ;  {
		if 0 == nSk {
			break
		}
		if 0 < nSk {
			if nSk + this.m_nSindex <= this.m_nEindex {
				this.m_nSindex += nSk
				break
			}
			errCode = fmt.Errorf("seek err sindex=%d seek=%d eindex=%d",
				this.m_nSindex, nSk, this.m_nEindex)
			break
		}

		if this.m_nEindex - nSk >= this.m_nSindex {
			this.m_nEindex -= nSk
			break
		}
		errCode = fmt.Errorf("seek err sindex=%d seek=%d eindex=%d",
			this.m_nSindex, nSk, this.m_nEindex)
		break
	}
	return errCode
}

func (this *binaryBuf)GetSize() int{
	return this.m_nEindex - this.m_nSindex
}

func newBinaryBuf(nType int) (BinaryBufPtr, error)  {
	pRet, errCode := BinaryBufPtr(nil), error(nil)
	for ; ;  {
		switch nType {
		case BinaryBuf_32:
			pTemp := &binaryBuf32{binaryBuf{0, 0, BinaryBuf_32, nil},
				arrByte32{binaryBufSilce{nil},[BinaryBuf_32]byte{}}}
			pTemp.binaryBuf.m_pBufSilce = &pTemp.m_arrByte
			pTemp.m_arrByte.m_bufSilce = pTemp.m_arrByte.m_buf[:]
			pRet = pTemp
		case BinaryBuf_64:
			pTemp := &binaryBuf64{binaryBuf{0, 0, BinaryBuf_64, nil},
				arrByte64{binaryBufSilce{nil},[BinaryBuf_64]byte{}}}
			pTemp.binaryBuf.m_pBufSilce = &pTemp.m_arrByte
			pTemp.m_arrByte.m_bufSilce = pTemp.m_arrByte.m_buf[:]
			pRet = pTemp
		case BinaryBuf_128:
			pTemp := &binaryBuf128{binaryBuf{0, 0, BinaryBuf_128, nil},
				arrByte128{binaryBufSilce{nil},[BinaryBuf_128]byte{}}}
			pTemp.binaryBuf.m_pBufSilce = &pTemp.m_arrByte
			pTemp.m_arrByte.m_bufSilce = pTemp.m_arrByte.m_buf[:]
			pRet = pTemp
		case BinaryBuf_256:
			pTemp := &binaryBuf256{binaryBuf{0, 0, BinaryBuf_256, nil},
				arrByte256{binaryBufSilce{nil},[BinaryBuf_256]byte{}}}
			pTemp.binaryBuf.m_pBufSilce = &pTemp.m_arrByte
			pTemp.m_arrByte.m_bufSilce = pTemp.m_arrByte.m_buf[:]
			pRet = pTemp
		case BinaryBuf_512:
			pTemp := &binaryBuf512{binaryBuf{0, 0, BinaryBuf_512, nil},
				arrByte512{binaryBufSilce{nil},[BinaryBuf_512]byte{}}}
			pTemp.binaryBuf.m_pBufSilce = &pTemp.m_arrByte
			pTemp.m_arrByte.m_bufSilce = pTemp.m_arrByte.m_buf[:]
			pRet = pTemp
		case BinaryBuf_1k:
			pTemp := &binaryBuf1k{binaryBuf{0, 0, BinaryBuf_1k, nil},
				arrByte1k{binaryBufSilce{nil},[BinaryBuf_1k]byte{}}}
			pTemp.binaryBuf.m_pBufSilce = &pTemp.m_arrByte
			pTemp.m_arrByte.m_bufSilce = pTemp.m_arrByte.m_buf[:]
			pRet = pTemp
		case BinaryBuf_2k:
			pTemp := &binaryBuf2k{binaryBuf{0, 0, BinaryBuf_2k, nil},
				arrByte2k{binaryBufSilce{nil},[BinaryBuf_2k]byte{}}}
			pTemp.binaryBuf.m_pBufSilce = &pTemp.m_arrByte
			pTemp.m_arrByte.m_bufSilce = pTemp.m_arrByte.m_buf[:]
			pRet = pTemp
		case BinaryBuf_4k:
			pTemp := &binaryBuf4k{binaryBuf{0, 0, BinaryBuf_4k, nil},
				arrByte4k{binaryBufSilce{nil},[BinaryBuf_4k]byte{}}}
			pTemp.binaryBuf.m_pBufSilce = &pTemp.m_arrByte
			pTemp.m_arrByte.m_bufSilce = pTemp.m_arrByte.m_buf[:]
			pRet = pTemp
		case BinaryBuf_8k:
			pTemp := &binaryBuf8k{binaryBuf{0, 0, BinaryBuf_8k, nil},
				arrByte8k{binaryBufSilce{nil},[BinaryBuf_8k]byte{}}}
			pTemp.binaryBuf.m_pBufSilce = &pTemp.m_arrByte
			pTemp.m_arrByte.m_bufSilce = pTemp.m_arrByte.m_buf[:]
			pRet = pTemp
		case BinaryBuf_16k:
			pTemp := &binaryBuf16k{binaryBuf{0, 0, BinaryBuf_16k, nil},
				arrByte16k{binaryBufSilce{nil},[BinaryBuf_16k]byte{}}}
			pTemp.binaryBuf.m_pBufSilce = &pTemp.m_arrByte
			pTemp.m_arrByte.m_bufSilce = pTemp.m_arrByte.m_buf[:]
			pRet = pTemp
		case BinaryBuf_32k:
			pTemp := &binaryBuf32k{binaryBuf{0, 0, BinaryBuf_32k, nil},
				arrByte32k{binaryBufSilce{nil},[BinaryBuf_32k]byte{}}}
			pTemp.binaryBuf.m_pBufSilce = &pTemp.m_arrByte
			pTemp.m_arrByte.m_bufSilce = pTemp.m_arrByte.m_buf[:]
			pRet = pTemp
		case BinaryBuf_64k:
			pTemp := &binaryBuf64k{binaryBuf{0, 0, BinaryBuf_64k, nil},
				arrByte64k{binaryBufSilce{nil},[BinaryBuf_64k]byte{}}}
			pTemp.binaryBuf.m_pBufSilce = &pTemp.m_arrByte
			pTemp.m_arrByte.m_bufSilce = pTemp.m_arrByte.m_buf[:]
			pRet = pTemp
		default:
			errCode = fmt.Errorf("没有找到指定的类型， type=%d", nType)
			pRet = nil
		}
		break
	}
	return pRet, errCode
}