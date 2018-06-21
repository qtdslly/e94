package util

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/mozillazg/go-pinyin"
)

func Atou(s string) uint {
	i, _ := strconv.Atoi(s)
	return uint(i)
}

func Atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func Itoa(i int) string {
	return strconv.Itoa(i)
}

func Atoi64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	} else {
		return i
	}
}

func Atof64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	} else {
		return f
	}
}

func ContainString(list []string, e string) bool {
	for _, a := range list {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsId(list [][]byte, e []byte) bool {
	for _, a := range list {
		if bytes.Compare(a, e) == 0 {
			return true
		}
	}
	return false
}

func ContainsInt(list []uint32, e uint32) bool {
	for _, a := range list {
		if a == e {
			return true
		}
	}
	return false
}

/*
	This method allows to compute the union of two list and remove the duplicate ids.
*/
func Union(oldIds, newIds [][]byte) (unionIds [][]byte) {
	for _, oldId := range oldIds {
		if !ContainsId(unionIds, oldId) {
			unionIds = append(unionIds, oldId)
		}
	}
	for _, newId := range newIds {
		if !ContainsId(unionIds, newId) {
			unionIds = append(unionIds, newId)
		}
	}
	return
}

/*
	This method allows to compute the intersection of two list.

	input:
		oldIds			:	a list of old ids
		newIds			:	a list of new ids

	return:
		excludedIds		:	the ids which are in old ids but not in new ids
		keptIds			:	the ids which are both in old ids and new ids
		addedIds		:	the ids which are in new ids but not in old olds
*/
func IntersectionUuid(oldIds, newIds [][]byte) (excludedIds [][]byte, keptIds [][]byte, addedIds [][]byte) {
	for _, oldId := range oldIds {
		found := false
		for _, newId := range newIds {
			if bytes.Equal(oldId, newId) {
				found = true
				keptIds = append(keptIds, oldId)
				break
			}
		}
		if !found {
			excludedIds = append(excludedIds, oldId)
		}
	}

	for _, newId := range newIds {
		found := false
		for _, keptId := range keptIds {
			if bytes.Equal(newId, keptId) {
				found = true
				break
			}
		}
		if !found {
			addedIds = append(addedIds, newId)
		}
	}
	return excludedIds, keptIds, addedIds
}

/*
	This method allows to compute the intersection of two list.

	input:
		oldIds			:	a list of old ids
		newIds			:	a list of new ids

	return:
		excludedIds		:	the ids which are in old ids but not in new ids
		keptIds			:	the ids which are both in old ids and new ids
		addedIds		:	the ids which are in new ids but not in old olds
*/
func IntersectionInt(oldIds, newIds []uint32) (excludedIds, keptIds, addedIds []uint32) {
	for _, oldId := range oldIds {
		found := false
		for _, newId := range newIds {
			if oldId == newId {
				found = true
				keptIds = append(keptIds, oldId)
				break
			}
		}
		if !found {
			excludedIds = append(excludedIds, oldId)
		}
	}

	for _, newId := range newIds {
		found := false
		for _, keptId := range keptIds {
			if newId == keptId {
				found = true
				break
			}
		}
		if !found {
			addedIds = append(addedIds, newId)
		}
	}
	return excludedIds, keptIds, addedIds
}

func DifferenceInt(arr1, arr2 []uint32) []uint32 {
	var diff []uint32
	for _, v1 := range arr1 {
		found := false
		for _, v2 := range arr2 {
			if v1 == v2 {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, v1)
		}
	}
	return diff
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// 判断是否是非0正整数,用于判断是否是id
func IsNumber(str string) bool {
	numReg := regexp.MustCompile(`^\+?[1-9][0-9]*$`)
	return numReg.MatchString(str)
}

func JSONMarshal(v interface{}, safeEncoding bool) ([]byte, error) {
	b, err := json.Marshal(v)
	if safeEncoding {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}

func TitleToPinyin(title string) string {
	str := ""
	r := []rune(title)
	for i := 0; i < len(r); i++ {
		if match, _ := regexp.MatchString("[0-9A-Za-z]", string(r[i])); match {
			str += string(r[i])
		} else {
			p := pinyin.Pinyin(string(r[i]), pinyin.NewArgs())
			for _, v := range p {
				if len(v) > 0 {
					str += v[0][0:1]
				}
			}
		}
	}

	return str
}


func TitleToFullPinyin(title string) string {
	str := ""
	r := []rune(title)
	for i := 0; i < len(r); i++ {
		if match, _ := regexp.MatchString("[0-9A-Za-z]", string(r[i])); match {
			str += string(r[i])
		} else {
			p := pinyin.Pinyin(string(r[i]), pinyin.NewArgs())
			for _, v := range p {
				str += v[0]
			}
		}
	}

	return str
}