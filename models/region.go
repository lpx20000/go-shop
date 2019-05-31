package models

import (
	"encoding/json"
	"fmt"
	"shop/pkg/e"
	"shop/pkg/gredis"
)

type Region struct {
	Id         int    `json:"id"`
	Pid        int    `json:"pid"`
	Shortname  string `json:"shortname,omitempty"`
	Name       string `json:"name"`
	MergerName string `json:"merger_name,omitempty"`
	Level      int    `json:"level"`
	Pinyin     string `json:"pinyin,omitempty"`
	Code       string `json:"code,omitempty"`
	ZipCode    string `json:"zip_code,omitempty"`
	First      string `json:"first,omitempty"`
	Lng        string `json:"lng,omitempty"`
	Lat        string `json:"lat,omitempty"`
}

type Tree struct {
	CommonList
	City map[int]City `json:"city"`
}

type City struct {
	CommonList
	RegionInfo map[int]CommonList `json:"region"`
}

type CommonList struct {
	Id    int    `json:"id"`
	Pid   int    `json:"pid"`
	Name  string `json:"name"`
	Level uint8  `json:"level"`
}

type RegionInfo struct {
	Province string `json:"province"`
	City     string `json:"city"`
	Region   string `json:"region"`
}

func GetRegion() (commonList []CommonList) {
	Db.Model(&Region{}).Select("id, pid, name, level").Scan(&commonList)
	return
}

func GetRegionInfo(province, region, city int) (regionInfo RegionInfo, err error) {
	var (
		allRegion  map[int]CommonList
		dataByte   []byte
		commonList []CommonList
	)

	if gredis.Exists(e.CACHE_REGION) {
		if dataByte, err = gredis.Get(e.CACHE_REGION); err != nil {
			return
		}
		if err = json.Unmarshal(dataByte, &allRegion); err != nil {
			return
		}
		regionInfo.Region = allRegion[region].Name
		regionInfo.Province = allRegion[province].Name
		regionInfo.City = allRegion[city].Name
		return
	}
	commonList = GetRegion()
	allRegion = make(map[int]CommonList, len(commonList))
	for _, item := range commonList {
		allRegion[item.Id] = item
	}

	if err = gredis.Set(e.CACHE_REGION, allRegion, 0); err != nil {
		return
	}
	regionInfo.Region = allRegion[region].Name
	regionInfo.Province = allRegion[province].Name
	regionInfo.City = allRegion[city].Name
	return
}

func GetIdByRegionName() (regionInfo map[string]int, err error) {
	var commonList []CommonList
	commonList = GetRegion()
	regionInfo = make(map[string]int, len(commonList))
	for _, item := range commonList {
		regionInfo[fmt.Sprintf("%s:%d:%d", item.Name, item.Level, item.Pid)] = item.Id
	}

	err = gredis.Set(e.CACHA_APP_REGION_ID, regionInfo, 0)
	return
}
