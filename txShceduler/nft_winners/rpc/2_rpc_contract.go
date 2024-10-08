package rpc

import (
	"jtools/cc"
	"jtools/cloud/ebcm"
	"txscheduler/brix/tools/dbg"
)

type IContract interface {
	Input() ebcm.PMethodIDDataMap
	cMerge(vs ...IContract)
	addEvent(vs ...IContract)
	getEvent() event_map_list

	ParseEvent(tx ebcm.TransactionBlock, isGlobal ...bool) EventLogList
}
type cContract struct {
	input ebcm.PMethodIDDataMap
	event event_map_list
}

func newContract(input ebcm.PMethodIDDataMap, event event_map_list) IContract {
	nc := &cContract{
		input: input,
		event: event,
	}
	if input == nil {
		nc.input = ebcm.PMethodIDDataMap{}
	}
	if event == nil {
		nc.event = event_map_list{}
	} else {
		global_event_all.merge(event)

	}
	return nc
}

func (my cContract) Input() ebcm.PMethodIDDataMap {
	return my.input
}
func (my cContract) getEvent() event_map_list {
	return my.event
}

func (my *cContract) cMerge(vs ...IContract) {
	for _, v := range vs {
		if v.Input() != nil {
			my.input.Merge(v.Input())
		}
		if v.getEvent() != nil {
			my.event.merge(v.getEvent())
		}
	} //for
}

func (my *cContract) addEvent(vs ...IContract) {
	for _, v := range vs {
		if v.getEvent() != nil {
			my.event.merge(v.getEvent())
		}
	} //for
}

////////////////////////////////////////////////////////////////////////

type EventItem map[string]interface{}

type eventParseBox struct {
	name  string
	parse func(log ebcm.TxLog) interface{}
}

var (
	global_event_all = event_map_list{}
)

type event_map_list map[string][]eventParseBox

type event_map map[string]eventParseBox

func MakeEventMap(em event_map) event_map_list {
	map_list := event_map_list{}
	if em == nil {
		return map_list
	}
	for k, v := range em {
		map_list[k] = []eventParseBox{v}
	}
	return map_list
}

func (my event_map_list) getEvent() event_map_list { return my }

func (my event_map_list) merge(event event_map_list) {
	for key, val_list := range event {
		if my_list, do := my[key]; do {

			add_list := []eventParseBox{}
			for _, a := range val_list {
				is_append := true
				for _, v := range my_list {
					if v.name == a.name {
						is_append = false
						break
					}
				} //for
				if is_append {
					add_list = append(add_list, a)
				}
			}
			my_list = append(my_list, add_list...)
			my[key] = my_list

		} else {
			my[key] = val_list
		}
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

func (my EventLog) ChangeStructData(p interface{}) error {
	return dbg.ParseStruct(my.Data, p)
}

func (my EventLogList) MatchName(name string, f func(event EventLog)) {
	for _, event := range my {
		if event.Name == name {
			f(event)
			break
		}
	} //for
}
func FindEventOne[T any](events EventLogList, name string) T {
	var item T
	for _, event := range events {
		if event.Name == name {
			item = event.Data.(T)
			break
		}
	} //for
	return item
}

func ParseEvent[T any](log EventLog, p T) T {
	log.ChangeStructData(p)
	return p
}

func (my cContract) ParseEvent(tx ebcm.TransactionBlock, isGlobal ...bool) EventLogList {
	list := EventLogList{}

	var event_list event_map_list
	if dbg.IsTrue(isGlobal) {
		event_list = global_event_all
	} else {
		event_list = my.event
	}

	for _, log := range tx.Logs {
		if parser_list, do := event_list[log.Topics.GetName()]; do {
			for _, parser := range parser_list {
				make_data := func(parser eventParseBox) (el *EventLog) {
					defer func() {
						if e := recover(); e != nil {
							el = nil
						}
					}()
					el = &EventLog{
						Contract: log.Address,
						Name:     parser.name,
						Data:     parser.parse(log),
					}
					return el
				}

				if elp := make_data(parser); elp != nil {
					list = append(list, *elp)
					break
				}

				cc.Gray("[ParseEvent]", parser.name, "is skip!")

			} //for
		}
	} //for

	return list
}

////////////////////////////////////////////////////////////////////////

type increaser int

func (my *increaser) N() int {
	v := int(*my)
	(*my)++
	return v
}
