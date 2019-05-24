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

func GetRegionInfo() (all map[int]CommonList, tree map[int]Tree) {
	var (
		commonList []CommonList
	)

	all = make(map[int]CommonList)
	tree = make(map[int]Tree)

	Db.Model(&Region{}).Select("id, pid, name, level").Scan(&commonList)

	for pKey, province := range commonList {
		all[province.Id] = province
		if province.Level == 1 {
			tree[province.Id] = Tree{
				CommonList: province,
				City:       make(map[int]City, 0),
			}
			commonList[pKey] = CommonList{}
			for cKey, city := range commonList {
				if city.Level == 2 && city.Pid == province.Id {
					tree[province.Id].City[city.Id] = City{
						CommonList: city,
						RegionInfo: make(map[int]CommonList, 0),
					}
					commonList[cKey] = CommonList{}
					for rKey, region := range commonList {
						if region.Level == 3 && region.Pid == city.Id {
							tree[province.Id].City[city.Id].RegionInfo[region.Id] = region
							commonList[rKey] = CommonList{}
						}
					}
				}
			}
		}
	}
	return
}
