package ecs

import (
	"jtools/cloud/ebcm"
	"jtools/dbg"
	"jtools/jmath"
	"jtools/jnet/cnet"
)

// nGasStation :
type nGasStation struct {
	Fast    float64 `json:"fast"`
	Fastest float64 `json:"fastest"`
	SafeLow float64 `json:"safeLow"`
	Average float64 `json:"average"`
	Begger  float64 `json:"begger"`
}

func (my nGasStation) String() string { return dbg.ToJsonString(my) }

func (my nGasStation) GasStationWei() GasStationWei {
	data := GasStationWei{
		Fast:    ebcm.GWEI(jmath.VALUE(my.Fast)).ToWEI().String(),
		Fastest: ebcm.GWEI(jmath.VALUE(my.Fastest)).ToWEI().String(),
		SafeLow: ebcm.GWEI(jmath.VALUE(my.SafeLow)).ToWEI().String(),
		Average: ebcm.GWEI(jmath.VALUE(my.Average)).ToWEI().String(),
		Begger:  ebcm.GWEI(jmath.VALUE(my.Begger)).ToWEI().String(),
	}
	return data
}

type GasStationWei struct {
	Fast    string `json:"fast"`
	Fastest string `json:"fastest"`
	SafeLow string `json:"safeLow"`
	Average string `json:"average"`
	Begger  string `json:"begger"`
}

func (my GasStationWei) String() string { return dbg.ToJsonString(my) }

func GasStation() nGasStation {
	pGas := nGasStation{}

	_, buf, _ := cnet.GET("https://ethgasstation.info/json/ethgasAPI.json")

	dbg.ParseStruct(buf, &pGas)
	return pGas
}
