package log

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

type LogLevel int

const (
	Log_g_nInfo LogLevel = iota
	Log_g_nDebug
	Log_g_nWarring
	Log_g_nError
)

type log struct {
	m_nDay    int
	m_nLevel  LogLevel
	m_strPath string
	m_strName string
	m_pFile *os.File
	m_mutex sync.Mutex
}

var g_pThis *log = nil

func LogInfo(format string, argv ...interface{}) error  {
	errCode := error(nil)
	for ; ;  {
		if nil == g_pThis {
			errCode = fmt.Errorf("没有初始化log")
			break
		}
		errCode = g_pThis.writeInfo(&format, argv ...)
		break
	}
	return errCode
}

func LogDebug(format string, argv ...interface{}) error  {
	errCode := error(nil)
	for ; ;  {
		if nil == g_pThis {
			errCode = fmt.Errorf("没有初始化log")
			break
		}
		errCode = g_pThis.writeDebug(&format, argv ...)
		break
	}
	return errCode
}

func LogWarring(format string, argv ...interface{}) error  {
	errCode := error(nil)
	for ; ;  {
		if nil == g_pThis {
			errCode = fmt.Errorf("没有初始化log")
			break
		}
		errCode = g_pThis.writeWarring(&format, argv ...)
		break
	}
	return errCode
}

func LogError(format string, argv ...interface{}) error  {
	errCode := error(nil)
	for ; ;  {
		if nil == g_pThis {
			errCode = fmt.Errorf("没有初始化log")
			break
		}
		errCode = g_pThis.writeError(&format, argv ...)
		break
	}
	return errCode
}


func InitInstance(nLevel LogLevel, strPath, strName string) error {
	errCode := error(nil)
	for ; ;  {
		if nil == g_pThis {
			now := time.Now()
			date := fmt.Sprintf("%d-%02d-%02d", now.Year(), int(now.Month()), now.Day())
			pFile, err := os.OpenFile(strPath + "/" + date + strName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModeAppend)
			if nil != err {
				errCode = err
				break
			}

			g_pThis = &log{now.Day(), nLevel, strPath, strName, pFile, sync.Mutex{}}
		}
		break
	}
	return errCode
}

func ReleaseInstance()  {
	if nil != g_pThis && nil != g_pThis.m_pFile {
		g_pThis.m_pFile.Close()
		g_pThis.m_pFile = nil
	}
}

func (this *log)writeInfo(pStrFormat *string, argv ...interface{}) error {
	errCode := error(nil)
	for ; ;  {
		if !this.allowWrite(Log_g_nInfo) {
			break
		}
		if errCode = this.checkWriteFile(); nil != errCode {
			break
		}
		errCode = this.writeLog(Log_g_nInfo, pStrFormat, argv ...)
		break
	}
	return errCode
}

func (this *log)writeDebug(pStrFormat *string, argv ...interface{}) error {
	errCode := error(nil)
	for ; ;  {
		if !this.allowWrite(Log_g_nDebug) {
			break
		}
		if errCode = this.checkWriteFile(); nil != errCode {
			break
		}
		errCode = this.writeLog(Log_g_nDebug, pStrFormat, argv ...)
		break
	}
	return errCode
}

func (this *log)writeWarring(pStrFormat *string, argv ...interface{}) error {
	errCode := error(nil)
	for ; ;  {
		if !this.allowWrite(Log_g_nWarring) {
			break
		}
		if errCode = this.checkWriteFile(); nil != errCode {
			break
		}
		errCode = this.writeLog(Log_g_nWarring, pStrFormat, argv ...)
		break
	}
	return errCode
}

func (this *log)writeError(pStrFormat *string, argv ...interface{}) error {
	errCode := error(nil)
	for ; ;  {
		if !this.allowWrite(Log_g_nError) {
			break
		}
		if errCode = this.checkWriteFile(); nil != errCode {
			break
		}
		errCode = this.writeLog(Log_g_nError, pStrFormat, argv ...)
		break
	}
	return errCode
}

func (this *log) allowWrite(curWriteLevel LogLevel) bool {
	return curWriteLevel >= this.m_nLevel
}

func (this *log) checkWriteFile() error{
	errCode := error(nil)
	for ; ;  {
		if nil == g_pThis {
			errCode = fmt.Errorf("没有初始化")
			break
		}
		if time.Now().Day() == this.m_nDay {
			break
		}
		this.m_mutex.Lock()
		defer this.m_mutex.Unlock()
		this.m_pFile.Close()
		now := time.Now()
		date := fmt.Sprintf("%d-%02d-%02d", now.Year(), int(now.Month()), now.Day())
		this.m_pFile, errCode = os.OpenFile(this.m_strPath + "/" + date + this.m_strName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModeAppend)
		break
	}
	return errCode
}

func (this *log) writeLog(nLevel LogLevel, pStrFormat *string, argv ...interface{}) error {
	errCode := error(nil)
	for ; ;  {
		now := time.Now()
		date := fmt.Sprintf("%d-%02d-%02d", now.Year(), int(now.Month()), now.Day())
		s := "[" + date + "]"
		switch nLevel {
		case Log_g_nInfo:
			s += "[info]"
		case Log_g_nDebug:
			s += "[debug]"
		case Log_g_nError:
			s += "[error]"
		case Log_g_nWarring:
			s += "[warring]"
		}
		pc,_,line,_ := runtime.Caller(2)
		f := runtime.FuncForPC(pc)
		s += fmt.Sprintf("[Line:%d][func=%s]", line, f.Name())
		s += fmt.Sprintf(*pStrFormat, argv...)
		s += "\r\n"
		this.m_mutex.Lock()
		_, errCode = this.m_pFile.WriteString(s)
		this.m_mutex.Unlock()
		break
	}
	return errCode
}