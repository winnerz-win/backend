package jrpc

import (
	"jtools/cc"
	"jtools/cloud/ebcm"
	"jtools/dbg"
	"strings"
)

type IContract interface {
	Input() ebcm.PMethodIDDataMap
	cMerge(vs ...IContract)
	addEvent(vs ...IContract)
	getEvent() EventMap

	ParseEvent(tx ebcm.TransactionBlock, is_global ...bool) EventLogList
}

type cContract struct {
	input ebcm.PMethodIDDataMap
	event EventMap
}

func (my cContract) Input() ebcm.PMethodIDDataMap {
	return my.input
}

func (my *cContract) cMerge(vs ...IContract) {
	for _, v := range vs {
		if v.Input() != nil {
			my.input.Merge(v.Input())
		}
		if v.getEvent() != nil {
			my.event.merge(v.getEvent())
		}
	}
}

func (my *cContract) addEvent(vs ...IContract) {
	for _, v := range vs {
		if v.getEvent() != nil {
			my.event.merge(v.getEvent())
		}
	} //for
}

func (my cContract) getEvent() EventMap {
	return my.event
}

type EventItem struct {
	Name  string
	Parse func(log ebcm.TxLog) interface{}
}

type EventMap map[string]EventItem

func (my EventMap) getEvent() EventMap { return my }

func (my EventMap) merge(event EventMap) {
	for key, val := range event {
		if _, do := my[key]; do {
			cc.Gray("[jrpc] EvemtMap.merge(event) is duplicated key[", key, "]")
		}
		my[key] = val
	}
}

type EventLog struct {
	Contract string      `bson:"contract" json:"contract"`
	Name     string      `bson:"name" json:"name"`
	Data     interface{} `bson:"data" json:"data"`
}
type EventLogList []EventLog

func (my EventLog) String() string     { return dbg.ToJsonString(my) }
func (my EventLogList) String() string { return dbg.ToJsonString(my) }

func (my EventLog) ParseData(ptr interface{}) error {
	return dbg.ParseStruct(my.Data, ptr)
}

func (my EventLogList) MatchName(name string, f func(event EventLog)) bool {
	name = strings.TrimSpace(name)
	for _, event := range my {
		if event.Name == name {
			f(event)
		}
	}

	return false
}

func ParseEvent[T any](log EventLog, p T) T {
	log.ParseData(p)
	return p
}

//////////////////////////////////////////////////////////////////

var (
	global_event_all = EventMap{}
)

func NewContract(input ebcm.PMethodIDDataMap, event EventMap) IContract {
	nc := &cContract{
		input: input,
		event: event,
	}
	if input == nil {
		nc.input = ebcm.PMethodIDDataMap{}
	}
	if event == nil {
		nc.event = EventMap{}
	} else {
		global_event_all.merge(event)
	}
	return nc
}

func (my cContract) ParseEvent(tx ebcm.TransactionBlock, is_global ...bool) EventLogList {
	list := EventLogList{}

	var target EventMap
	if dbg.IsTrue(is_global) {
		target = global_event_all
	} else {
		target = my.event
	}

	for _, log := range tx.Logs {
		if parser, do := target[log.Topics.GetName()]; do {
			el := EventLog{
				Contract: log.Address,
				Name:     parser.Name,
				Data:     parser.Parse(log),
			}
			list = append(list, el)
		}
	} //for
	return list
}
