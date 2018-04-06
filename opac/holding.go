package opac

import (
	//"fmt"
	"io/ioutil"
	"net/http"

	"github.com/panthesingh/goson"
)

type holding struct {
	Callno   string
	State    string
	Lib      string
	Location string
}

const holdURL = "http://opac.gzlib.gov.cn/opac/api/holding/"

func (b *book) checkHolding() {

	b.holdings = make([]holding, 0)

	json := getJSON(b.Bookrecno)

	g, _ := goson.Parse(json)

	holds := g.Get("holdingList")
	stateMap := g.Get("holdStateMap")
	localMap := g.Get("localMap")
	libMap := g.Get("libcodeMap")

	for i := 0; i < holds.Len(); i++ {
		hold := holds.Index(i)

		stateCode := hold.Get("state").String()
		stateName := stateMap.Get(stateCode).Get("stateName").String()

		locCode := hold.Get("curlocal").String()
		locName := localMap.Get(locCode).String()

		libCode := hold.Get("curlib").String()
		libName := libMap.Get(libCode).String()

		h := holding{
			Callno:   hold.Get("callno").String(),
			State:    stateName,
			Lib:      libName,
			Location: locName,
		}

		b.holdings = append(b.holdings, h)
	}
}

func getJSON(id string) []byte {
	resp, err := http.Get(holdURL + id)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	return body
}
