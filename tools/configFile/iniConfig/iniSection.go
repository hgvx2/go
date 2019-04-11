package iniConfig

import (
	"fmt"
	"os"
)

type iniSection struct {
	m_strName string
	m_listKeyValueNode []*iniKeyValue
	m_listNote []string
	m_mapKeyValue map[string] *iniKeyValue
}

func new_iniSection(sName string, listKeyValueNode []*iniKeyValue,
	listNote []string, mapKeyValue map[string] *iniKeyValue) *iniSection {
	return &iniSection{sName, listKeyValueNode,
		listNote, mapKeyValue}
}

func (this *iniSection) addKey(strKey, strValue string, listNote []string) error{
	errCode := error(nil)
	for ; ;  {
		if _, bFind := this.m_mapKeyValue[strKey]; bFind {
			errCode = fmt.Errorf("key = %s 已存在", strKey)
			break
		}
		pKeyValueNode := new_iniKeyValue(&strKey, &strValue, listNote)
		this.m_mapKeyValue[strKey] = pKeyValueNode
		this.m_listKeyValueNode = append(this.m_listKeyValueNode, pKeyValueNode)
		break
	}
	return errCode
}

func (this *iniSection)writeSection(pFile *os.File) error {
	errCode := error(nil)
	for ; ;  {

		if _, errCode = pFile.WriteString("\n"); nil != errCode {
			break
		}
		// 写注释
		if errCode = writeNote(pFile, this.m_listNote); nil != errCode {
			break
		}
		// 写节
		if errCode = writeSectionName(pFile, this.m_strName); nil != errCode{
			break
		}
		// 写key value
		for _, pKeyValue := range this.m_listKeyValueNode{
			if errCode = pKeyValue.writeKeyValue(pFile); nil != errCode {
				break
			}
		}
		break
	}
	return errCode
}

func writeNote(pFile *os.File, listNote []string) error {
	errCode := error(nil)
	for ; ;  {
		for  _, strNote := range listNote {
			strTemp := ";" + strNote + "\n"
			_, errCode = pFile.WriteString(strTemp)
		}
		break
	}
	return errCode
}

func writeSectionName(pFile *os.File, strSection string) error {
	strTemp := "[" + strSection + "]" + "\n"
	_, errCode := pFile.WriteString(strTemp)
	return errCode
}

func (this *iniSection)delKey(strKey string) error {
	errCode := error(nil)
	for ; ;  {
		pKeyValue, bFind := this.m_mapKeyValue[strKey]
		if !bFind {
			errCode = fmt.Errorf("mapKeyValue 中没有找到指定的key=%s", strKey)
			break
		}
		delete(this.m_mapKeyValue, strKey)
		errCode = fmt.Errorf("listKeyValue中没有找到指定的Key=%s", strKey)
		for i, pKeyValueTemp := range this.m_listKeyValueNode {
			if pKeyValue == pKeyValueTemp{
				this.m_listKeyValueNode = append(this.m_listKeyValueNode[:i], this.m_listKeyValueNode[i + 1:]...)
				errCode = error(nil)
			}
		}
		break
	}
	return errCode
}

func (this *iniSection) modifyKey(strKey, strNewKey string) error {
	errCode := error(nil)
	for ; ;  {
		pKeyValue, bFind := this.m_mapKeyValue[strKey]
		if !bFind {
			errCode = fmt.Errorf("mapKeyValue 中没有找到指定的key=%s", strKey)
			break
		}
		if  _, bFind = this.m_mapKeyValue[strNewKey]; bFind{
			errCode = fmt.Errorf("mapKeyValue中 key = %s 已存在", strNewKey)
			break
		}
		pKeyValue.m_strKey = strNewKey
		this.m_mapKeyValue[strNewKey] = pKeyValue
		delete(this.m_mapKeyValue, strKey)
		break
	}
	return errCode
}

func (this *iniSection) modifyKeyValue(strKey, strValue string) error {
	errCode := error(nil)
	for ; ;  {
		pKeyValue, bFind := this.m_mapKeyValue[strKey]
		if !bFind {
			errCode = fmt.Errorf("mapKeyValue 中没有找到指定的key=%s", strKey)
			break
		}
		pKeyValue.m_strValue = strValue
		break
	}
	return errCode
}