package models

type Region struct {
	Id         int    `json:"id"`
	Pid        int    `json:"pid"`
	Shortname  string `json:"shortname,omitempty"`
	Name       string `json:"name"`
	MergerName string `json:"merger_name,omitempty"`
	Level      uint8  `json:"level"`
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

func GetRegion() (commonList []CommonList) {
	Db.Model(&Region{}).Select("id, pid, name, level").Scan(&commonList)
	return
}

func GetRegionInfo() (all map[int]CommonList, tree map[int]Tree) {
	return
}
