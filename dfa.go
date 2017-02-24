package dfa

import (
	"regexp"
	"strings"
)

type Node struct {
	childrens map[string]*Node
	terminate bool
	reg       *regexp.Regexp
	parent    *Node
}

type Dfa struct {
	root Node
}

//添加关键词
//@word 关键词
func (p *Dfa) addWord(word string) {
	node := &p.root
	str := []rune(word)
	lastPos := len(str) - 1
	recordingFlag := false
	tempRegx := ""
	for i, s := range str {
		sKey := string(s)
		if sKey == "R" {
			if recordingFlag == false {
				recordingFlag = true
				continue
			} else {
				node.childrens["regex"] = &Node{map[string]*Node{}, i == lastPos, regexp.MustCompile(tempRegx), node}
				recordingFlag = false
			}
			tempRegx = ""
			node = node.childrens["regex"]
		} else {
			if recordingFlag {
				tempRegx += sKey
			} else {
				if node.childrens[sKey] == nil {
					node.childrens[sKey] = &Node{map[string]*Node{}, i == lastPos, nil, node}
				} else {
					if i == lastPos {
						node.childrens[sKey] = &Node{node.childrens[sKey].childrens, true, nil, node}
					}
				}
				node = node.childrens[sKey]
			}
		}
	}
}

// 载入关键词列表
// @words 关键词切片
func (p *Dfa) LoadWords(words []string) {
	p.root = Node{map[string]*Node{}, false, nil, nil}
	for _, v := range words {
		p.addWord(v)
	}
}

// 判断是否包含关键词,并提取匹配到的第一个关键词
// @text 待检测文本
func (p *Dfa) IsContain(text string) (bool, string) {
	keywork := make([]string, 0) //包含的关键词
	str := []rune(text)
	strLen := len(str)
	for i, _ := range str {
		node := &p.root
		for j := i; j < strLen && node != nil; j++ {
			currentWord := string(str[j])

			if node.childrens[currentWord] == nil { // 树中的第一层节点中，是否存在此节点
				keywork = make([]string, 0)         // 只要当前字符没找到，并且没到最终节点，将重新开始
				if node.childrens["regex"] == nil { // 当前节点是否存在，正则
					break
				} else {
					// 如果存在正则，通过正则匹配是否存在此节点
					str := getLastString(str, j)

					pos := node.childrens["regex"].reg.FindStringIndex(str)
					if pos != nil {
						j = pos[1] + j
						if node.childrens["regex"].terminate {

							return true, ""
						} else if pos[0] == pos[1] {
							j--
						}
					} else {
						break
					}
					node = node.childrens["regex"]
				}
			} else {
				keywork = append(keywork, (currentWord))
				if node.childrens[currentWord].terminate {
					return true, strings.Join(keywork, "")
				}
				node = node.childrens[currentWord]
			}
		}
	}
	return false, ""
}

// 过滤关键词
// @text 输入文本
// @replaceChar 替换的字符
func (p *Dfa) FilterWords(text, replaceChar string) string {
	str := []rune(text)
	strLen := len(str)
	strs := []string{}
	//是否继续添加
	isAppend := true
	for i := 0; i < strLen; i++ {
		node := &p.root
		tryTimes := 0
		for j := i; j < strLen && node != nil; j++ {
			currentWord := string(str[j])
			if node.childrens[currentWord] == nil {
				if node.childrens["regex"] == nil {
					if tryTimes == 1 {
						tryTimes = 2
						node = node.parent
						j -= 2
					} else {
						break
					}
				} else {
					str := getLastString(str, j)
					pos := node.childrens["regex"].reg.FindStringIndex(str)
					if pos != nil {
						j = pos[1] + j
						if node.childrens["regex"].terminate {
							if len(node.childrens["regex"].childrens) > 0 && tryTimes == 0 {
								tryTimes = 1
							} else {
								isAppend = false
								temp := []string{}
								for k := 0; k < j-i; k++ {
									temp = append(temp, replaceChar)
								}
								strs = append(strs, strings.Join(temp, ""))
								i = j - 1
							}
						} else if pos[0] == pos[1] {
							j--
						}
					} else {
						if tryTimes == 1 {
							tryTimes = 2
							node = node.parent
							j -= 2
							continue
						} else {
							break
						}
					}
					node = node.childrens["regex"]
				}
			} else {
				if node.childrens[currentWord].terminate {
					if len(node.childrens[currentWord].childrens) > 0 && tryTimes == 0 {
						tryTimes = 1
					} else {
						isAppend = false
						temp := []string{}
						for k := 0; k < j-i+1; k++ {
							temp = append(temp, replaceChar)
						}
						strs = append(strs, strings.Join(temp, ""))
						i = j
					}
				}
				node = node.childrens[currentWord]
			}
		}
		if isAppend {
			strs = append(strs, string(str[i]))
		} else {
			isAppend = true
		}
	}
	return strings.Join(strs, "")
}

func getLastString(strs []rune, j int) string {
	str := []string{}
	for i, s := range strs {
		if i >= j {
			str = append(str, string(s))
		}
	}
	return strings.Join(str, "")
}
