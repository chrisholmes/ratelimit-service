package main

import (
	"sort"
)

type WhiteList struct {
	whiteList []string
}

func NewWhiteList(whiteList []string) *WhiteList {
	return &WhiteList{
		whiteList: whiteList,
	}
}

func (wl *WhiteList) contains(ip string) bool {
	i := sort.SearchStrings(wl.whiteList, ip)
	return (i < len(wl.whiteList) && ip == wl.whiteList[i])
}
