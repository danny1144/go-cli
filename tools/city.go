package tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

var CityMap = make(map[string]int)

func init() {
		//CityMap["杭州"] = 330100
		//CityMap["上海"] = 310000
		//CityMap["深州"] = 131182
		//CityMap["广州"] = 440100
		//CityMap["开化"] = 330824
	err := ReadJson("D:\\gowork\\go-cli\\tools\\行政区划与名称.json")
	if err != nil {
		fmt.Println(err)
	}
}

type citySource struct {
	Xzqhbh string `json:"xzqhbh"`
	Xzqhmc string `json:"xzqhmc"`
}

func ReadJson(path string) error {

	data, err := ioutil.ReadFile(path)
	var ret []citySource
	if err = json.Unmarshal(data, &ret); err != nil {
		return err
	}
	for _, item := range ret {
		bh, _ := strconv.Atoi(item.Xzqhbh)
		CityMap[item.Xzqhmc] = bh
	}
	return nil
}
