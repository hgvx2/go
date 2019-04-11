package iniConfig

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type IniConfig struct {
	m_listSection []*iniSection
	m_mapSection map[string] *iniSection
	m_pCurSection *iniSection
	m_listNote []string
}

/* public functin
1、 func NewIniConfig() *IniConfig
2、 func (this *IniConfig) AnalyzeConfigFile(strPath string) error
3、 func (this *IniConfig) ReadNumberConfigValue(strSection, strKey string, nDefault int) (int, error)
4、 func (this *IniConfig) ReadStringConfigValue(strSection, strKey string, strDefault string) (string, error)
5、 func (this *IniConfig) SaveConfigFile(filePath string) error
6、 func (this *IniConfig) AddSection(strSection string, listNote []string) error
7、 func (this *IniConfig) AddKey(strSection, strKey, strValue string, listNote []string) error
8、 func (this *IniConfig) DelSection(strSection string) error
9、 func (this *IniConfig) DelKey(strSection, strKey string) error
10、func (this *IniConfig) ModifySection(strSection, strNewSection string) error
11、func (this *IniConfig) ModifyKey(strSection, strKey, strNewKey string) error
12、func (this *IniConfig) ModifyKeyValue(strSection, strKey, strValue string) error
*/

func NewIniConfig() *IniConfig {
	return &IniConfig{make([]*iniSection, 0, 20), make(map[string] *iniSection),
		nil, make([]string, 0, 2)}
}

func (this *IniConfig) AnalyzeConfigFile(strPath string) error {
	errCode := error(nil)
	for ; ;  {
		var pFile *os.File
		pFile, errCode = os.Open(strPath)
		if  nil != errCode{
			break
		}
		defer pFile.Close()
		buf := bufio.NewReader(pFile)
		for {
			var line string  = ""
			line, errCode = buf.ReadString('\n')
			line = strings.TrimSpace(line)
			if nil != errCode {
				if io.EOF == errCode {
					errCode = error(nil)
				}
				break
			}
			if errCode = this.analyzeLine(&line); nil != errCode {
				break
			}
		}
		break
	}
	return  errCode
}

func (this *IniConfig) analyzeLine(pLine *string) error {
	errCode := error(nil)
	for ; ;  {
		strLineLen := len(*pLine)
		if 0 >= strLineLen {
			// 空行
			break
		}else if ';' == (*pLine)[0] {
			// 注释行
			this.analyzeNote(pLine)
		}else if '[' == (*pLine)[0] {
			// 节点行
			errCode = this.analyzeSection(pLine)
		} else {
			// key => value 行
			errCode = this.analyzeKeyValue(pLine)
		}
		break
	}
	return errCode
}

// 解析注释行
func (this *IniConfig) analyzeNote(pLine *string)  {
	this.m_listNote = append(this.m_listNote, string([]byte(*pLine)[1:]))
}

// 解析节行
func (this *IniConfig) analyzeSection(pLine *string) error {
	errCode := error(nil)
	for ; ;  {
		nLineLen := len(*pLine)
		strSection, bOk:= "", false
		for i := 1; i < nLineLen; i ++ {
			if ']' != (*pLine)[i]{
				strSection += string((*pLine)[i])
				continue
			}
			bOk = true
		}
		if 0 == strSection[0] || !bOk {
			errCode = fmt.Errorf("节解析错误 当前行内容=%s", *pLine)
			break
		}
		strSection = strings.TrimSpace(strSection)

		this.m_pCurSection = new_iniSection(strSection, make([]*iniKeyValue, 0),
			this.m_listNote, make(map[string] *iniKeyValue))
		if nil == this.m_pCurSection {
			break
		}
		this.m_listNote = make([]string, 0, 2)
		this.m_listSection = append(this.m_listSection, this.m_pCurSection)
		_, ok := this.m_mapSection[strSection]
		if ok {
			errCode = fmt.Errorf("节 %s 重复", strSection)
			break
		}
		this.m_mapSection[strSection] = this.m_pCurSection
		break
	}
	return errCode
}

// 解析key value 行
func (this *IniConfig) analyzeKeyValue(pLine *string) error {
	errCode := error(nil)
	for ; ;  {
		if nil == this.m_pCurSection {
			errCode = fmt.Errorf("解析错误，当前节缺失 当前行内容=%s", *pLine)
			break
		}
		nLineLen := len(*pLine)
		strKey, strValue, bIsKey := "", "", true
		for i := 0; i < nLineLen; i ++ {
			if '=' == (*pLine)[i] {
				bIsKey = false
				continue
			}
			if bIsKey {
				strKey += string((*pLine)[i])
				continue
			}
			strValue += string((*pLine)[i])
		}

		if 0 == strKey[0] || bIsKey{
			errCode = fmt.Errorf("key value 解析错误 当前行内容=%s", *pLine)
			break
		}
		strKey = strings.TrimSpace(strKey)
		strValue = strings.TrimSpace(strValue)
		errCode = this.m_pCurSection.addKey(strKey, strValue, this.m_listNote)
		this.m_listNote = make([]string, 0)
		break
	}
	return errCode
}

func (this *IniConfig) ReadNumberConfigValue(strSection, strKey string, nDefault int) (int, error) {
	nRet, errCode := nDefault, error(nil)
	for ; ;  {
		var strValue string;
		strValue, errCode= this.ReadStringConfigValue(strSection, strKey, "")
		if nil != errCode {
			break
		}
		var err error
		nRet, err = strconv.Atoi(strValue)
		if nil != err {
			errCode = fmt.Errorf("节=%s key=%s err=%s", strSection, strKey, err.Error())
			nRet = nDefault
			break
		}
		break
	}
	return  nRet, errCode
}

func (this *IniConfig) ReadStringConfigValue(strSection, strKey string, strDefault string) (string, error) {
	strRet, errCode := strDefault, error(nil)
	for ; ;  {
		pSection, ok := this.m_mapSection[strSection]
		if !ok {
			errCode = fmt.Errorf("没有找到节=%s", strSection)
			break
		}
		var pKeyValue *iniKeyValue = nil
		pKeyValue , ok = pSection.m_mapKeyValue[strKey]
		if !ok {
			errCode = fmt.Errorf("没有找到节=%s 下面的key=%s", strSection, strKey)
			break
		}
		strRet = pKeyValue.m_strValue
		break
	}
	return  strRet, errCode
}

func (this *IniConfig) SaveConfigFile(filePath string) error  {
	errCode := error(nil)
	for ; ;  {
		var pFile *os.File
		pFile, errCode = os.Create(filePath)
		if nil != errCode {
			break
		}
		defer pFile.Close()
		for _, pSection := range this.m_listSection {
			errCode = pSection.writeSection(pFile)
			if nil != errCode {
				break
			}
		}
		break
	}
	return errCode
}

func (this *IniConfig) AddSection(strSection string, listNote []string) error {
	errCode := error(nil)
	for ; ;  {
		if _, bFind := this.m_mapSection[strSection]; bFind {
			errCode = fmt.Errorf("section 重复 section = %s", strSection)
			break
		}
		pSection := new_iniSection(strSection, make([]*iniKeyValue, 0, 5), listNote, make(map[string] *iniKeyValue))
		this.m_pCurSection = pSection
		this.m_mapSection[strSection] = pSection
		this.m_listSection = append(this.m_listSection, pSection)
		break
	}
	return errCode
}

func (this *IniConfig) AddKey(strSection, strKey, strValue string, listNote []string) error {
	errCode := error(nil)
	for ; ;  {
		var pSection *iniSection = nil
		var bFind bool
		if pSection, bFind = this.m_mapSection[strSection] ; !bFind {
			errCode = fmt.Errorf("没有找到 %s 节", strSection)
			break
		}
		errCode = pSection.addKey(strKey, strValue, listNote)

		break
	}
	return errCode
}

func (this *IniConfig) DelSection(strSection string) error  {
	errCode := error(nil)
	for ; ;  {
		pSection, bFind := this.m_mapSection[strSection]
		if !bFind {
			errCode = fmt.Errorf("mapSection中没有找到对于的section=%s", strSection)
			break
		}
		delete(this.m_mapSection, strSection)
		errCode = fmt.Errorf("listSection中没有找到对于的section=%s", strSection)
		for i, pSectionTemp := range this.m_listSection {
			if pSection == pSectionTemp {
				this.m_listSection = append(this.m_listSection[:i], this.m_listSection[i + 1:]...)
				errCode = error(nil)
			}
		}
		break
	}
	return errCode
}

func (this *IniConfig) DelKey(strSection, strKey string) error {
	errCode := error(nil)
	for ; ;  {
		pSection, bFind := this.m_mapSection[strSection]
		if !bFind {
			errCode = fmt.Errorf("mapSection中没有找到对于的section=%s", strSection)
			break
		}
		errCode = pSection.delKey(strKey)
		break
	}
	return errCode
}


func (this *IniConfig) ModifySection(strSection, strNewSection string) error  {
	errCode := error(nil)
	for ; ;  {
		pSection, bFind := this.m_mapSection[strSection]
		if !bFind {
			errCode = fmt.Errorf("mapSection中没有找到对于的section=%s", strSection)
			break
		}
		if  _, bFind = this.m_mapSection[strNewSection]; bFind{
			errCode = fmt.Errorf("mapSection中节%s 已存在", strNewSection)
			break
		}
		pSection.m_strName = strNewSection
		this.m_mapSection[strNewSection] = pSection
		delete(this.m_mapSection, strSection)
		break
	}
	return errCode
}

func (this *IniConfig) ModifyKey(strSection, strKey, strNewKey string) error  {
	errCode := error(nil)
	for ; ;  {
		pSection, bFind := this.m_mapSection[strSection]
		if !bFind {
			errCode = fmt.Errorf("mapSection中没有找到对于的section=%s", strSection)
			break
		}
		errCode = pSection.modifyKey(strKey, strNewKey)
		break
	}
	return errCode
}

func (this *IniConfig) ModifyKeyValue(strSection, strKey, strValue string) error  {
	errCode := error(nil)
	for ; ;  {
		pSection, bFind := this.m_mapSection[strSection]
		if !bFind {
			errCode = fmt.Errorf("mapSection中没有找到对于的section=%s", strSection)
			break
		}
		errCode = pSection.modifyKeyValue(strKey, strValue)
		break
	}
	return errCode
}