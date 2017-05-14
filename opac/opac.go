package opac

import (
	//"fmt"
	"github.com/panthesingh/goson"
	"io/ioutil"
	"net/http"
)

type Hold struct {
	Callno   string
	State    string
	Lib      string
	Location string
}

func Get(id string) []Hold {

        hold_list := make([]Hold, 0)

	json := get_json(id)

	g, _ := goson.Parse(json)

	holds := g.Get("holdingList")
	state_map := g.Get("holdStateMap")
	local_map := g.Get("localMap")
	lib_map := g.Get("libcodeMap")

	for i := 0; i < holds.Len(); i++ {
		hold := holds.Index(i)

		state_code := hold.Get("state").String()
		state_name := state_map.Get(state_code).Get("stateName").String()

		loc_code := hold.Get("curlocal").String()
		loc_name := local_map.Get(loc_code).String()

		lib_code := hold.Get("curlib").String()
		lib_name := lib_map.Get(lib_code).String()

		h := Hold{
			Callno:   hold.Get("callno").String(),
			State:    state_name,
			Lib:      lib_name,
			Location: loc_name,
		}
                hold_list = append(hold_list, h)
	}
        return hold_list
}

func get_json(id string) []byte {
	resp, err := http.Get("http://opac.gzlib.gov.cn/opac/api/holding/" + id)

	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	return body
}
