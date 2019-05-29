package models

type Region struct {
	Id         int    `json:"id"`
	Pid        int    `json:"pid"`
	Shortname  string `json:"shortname,omitempty"`
	Name       string `json:"name"`
	MergerName string `json:"merger_name,omitempty"`
	Level      int  `json:"level"`
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

func GetRegionInfo(provinceId, cityId, regionId int) (regionInfo RegionInfo) {
	var (
		all map[int]CommonList
	)
	all = make(map[int]CommonList)
	for _, item := range GetRegion() {
		all[item.Id] = item
	}
	regionInfo = RegionInfo{
		Province: all[provinceId].Name,
		City:     all[cityId].Name,
		Region:   all[regionId].Name,
	}
	return
}

func GetIdByRegionName(name string, level, pid int) (id int)  {
	var region Region
	Db.Where(map[string]interface{}{"name":name, "level":level, "pid":pid}).First(&region)
	id = region.Id
	return
}
