package mmdbwriter

import (
"fmt"
"github.com/maxmind/mmdbwriter/mmdbtype"
)

type  IsoRegion struct {
	CountryCode mmdbtype.String
	RegionCode mmdbtype.String
}

func (t *Tree) Move(
	from IsoRegion,
	to IsoRegion,
	) {
	//iso_country := from.Country
	var country_DataType mmdbtype.DataType
	to_move := map[dataMapKey]mmdbtype.Uint32{}

	for k, v := range t.dataMap.data {
		data := v.data.(mmdbtype.Map)
		
		// get geoname_id for to_country
		if country_DataType == nil {
			if data["country"] != nil {
				if data["country"].(mmdbtype.Map)["iso_code"] == to.CountryCode {
					country_DataType = data["country"]
					fmt.Printf("country: %s geoname id: %d\n",
						to.CountryCode,
						country_DataType.(mmdbtype.Map)["geoname_id"])
				}
			}
		}

		if data["subdivisions"] != nil && data["country"] != nil {
			if data["country"].(mmdbtype.Map)["iso_code"].(mmdbtype.String) == from.CountryCode {
				subdivisions := data["subdivisions"].(mmdbtype.Slice)
				for n := range subdivisions {
					if subdivisions[n].(mmdbtype.Map)["iso_code"] == from.RegionCode {
						to_move[k] = subdivisions[n].(mmdbtype.Map)["geoname_id"].(mmdbtype.Uint32)
						fmt.Printf("move: geoname id: %d, subdivision: %d\n", to_move[k], n)
					}
				}
			}
		}
	}
	
	for key, geoname_id := range to_move {
		t.dataMap.data[key].data.(mmdbtype.Map)["country"] = country_DataType
		t.dataMap.data[key].data.(mmdbtype.Map)["registered_country"] = country_DataType
		subdivisions := t.dataMap.data[key].data.(mmdbtype.Map)["subdivisions"].(mmdbtype.Slice)
		for n := range subdivisions {
			if subdivisions[n].(mmdbtype.Map)["geoname_id"] == geoname_id {
				subdivisions[n].(mmdbtype.Map)["iso_code"] = to.RegionCode
			}
		}
	}
}


