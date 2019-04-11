package iniConfig

import "os"

type iniKeyValue struct {
	m_strKey string
	m_strValue string
	m_listNote []string
}

func new_iniKeyValue(pStrKey, pStrValue *string, listNote []string) *iniKeyValue {
	return &iniKeyValue{*pStrKey, *pStrValue, listNote}
}

func (this *iniKeyValue) writeKeyValue(pFile *os.File) error  {
	errCode := error(nil)
	for ; ;  {
		if errCode = writeNote(pFile, this.m_listNote); nil != errCode {
			break
		}
		strTemp := this.m_strKey + "=" + this.m_strValue + "\n"
		_, errCode = pFile.WriteString(strTemp)
		break
	}
	return errCode
}