package mongo

import (
	"encoding/json"
	"fmt"
	"jtools/cc"
	"jtools/database/mongo/sm"
	"jtools/database/mongo/zzp"
	"jtools/dbg"
	"jtools/jmath"
	"time"
)

const (
	_hint_tag = "hint"
)

func get_col_node(
	node zzp.NodeFunc,
	cmd zzp.CmdFunc,
	brunch zzp.BrunchFunc) []interface{} {
	//node, cmd, brunch := zzp.NodeMaker()
	list := []interface{}{}
	find := node(
		cmd(
			"find",
			zzp.TagSmall,
			zzp.ArgJson,
		),

		brunch(
			cmd(
				"skip",
				zzp.TagSmall,
				zzp.ArgText,
			),
			cmd(
				"sort",
				zzp.TagSmall,
				zzp.ArgTextComma,
			),
			cmd(
				"limit",
				zzp.TagSmall,
				zzp.ArgText,
			),
			cmd(
				_hint_tag,
				zzp.TagSmall,
				zzp.ArgJson,
			),
		),

		node(
			cmd(
				"count",
				zzp.TagSmall,
			),
		),
		node(
			cmd(
				"one",
				zzp.TagSmall,
			),
		),
		node(
			cmd(
				"all",
				zzp.TagSmall,
			),
		),
		node(
			cmd(
				"sum",
				zzp.TagSmall,
				zzp.ArgTextDot,
			),
		),
	)
	findProjection := node(
		cmd(
			"findProjection",
			zzp.TagSmall,
			zzp.ArgJson,
			zzp.ArgJsonArray,
		),
		brunch(
			cmd(
				"skip",
				zzp.TagSmall,
				zzp.ArgText,
			),
			cmd(
				"sort",
				zzp.TagSmall,
				zzp.ArgTextComma,
			),
			cmd(
				"limit",
				zzp.TagSmall,
				zzp.ArgText,
			),
			cmd(
				_hint_tag,
				zzp.TagSmall,
				zzp.ArgJson,
			),
		),

		node(
			cmd(
				"count",
				zzp.TagSmall,
			),
		),
		node(
			cmd(
				"one",
				zzp.TagSmall,
			),
		),
		node(
			cmd(
				"all",
				zzp.TagSmall,
			),
		),
	)
	update := node(
		cmd(
			"update",
			zzp.TagSmall,
			zzp.ArgJson,
			zzp.ArgJson,
		),
	)
	updateAll := node(
		cmd(
			"updateAll",
			zzp.TagSmall,
			zzp.ArgJson,
			zzp.ArgJson,
		),
	)
	upsert := node(
		cmd(
			"upsert",
			zzp.TagSmall,
			zzp.ArgJson,
			zzp.ArgJson,
		),
	)
	insert := node(
		cmd(
			"insert",
			zzp.TagSmall,
			zzp.ArgJson,
		),
	)
	remove := node(
		cmd(
			"remove",
			zzp.TagSmall,
			zzp.ArgJson,
		),
	)
	removeAll := node(
		cmd(
			"removeAll",
			zzp.TagSmall,
			zzp.ArgJson,
		),
	)
	aggregate := node(
		cmd(
			"aggregate",
			zzp.TagSmall,
			zzp.ArgJsonArray,
		),

		node(
			cmd(
				"one",
				zzp.TagSmall,
			),
		),
		node(
			cmd(
				"all",
				zzp.TagSmall,
			),
		),
	)
	dropcollection := node(
		cmd(
			"dropcollection",
			zzp.TagSmall,
		),
	)
	dropindexname := node(
		cmd(
			"dropindexname",
			zzp.TagSmall,
			zzp.ArgTextComma,
		),
	)
	ensureindex := node(
		cmd(
			"ensureindex",
			zzp.TagSmall,
			zzp.ArgTextComma,
		),
	)
	dropindexAll := node(
		cmd(
			"dropindexAll",
			zzp.TagSmall,
		),
	)
	indexes := node(
		cmd(
			"indexes",
			zzp.TagSmall,
		),
	)

	FindOneAndUpdate := node(
		cmd(
			"FindOneAndUpdate",
			zzp.TagSmall,
			zzp.ArgJson,
			zzp.ArgJson,
		),
	)
	FindOneAndDelete := node(
		cmd(
			"FindOneAndDelete",
			zzp.TagSmall,
			zzp.ArgJson,
		),
	)
	FindOneAndReplace := node(
		cmd(
			"FindOneAndReplace",
			zzp.TagSmall,
			zzp.ArgJson,
			zzp.ArgJson,
		),
	)

	list = append(list,
		find,
		findProjection,

		update,
		upsert,
		updateAll,
		insert,
		remove,
		removeAll,
		aggregate,
		dropcollection,
		dropindexname,
		ensureindex,
		dropindexAll,
		indexes,

		FindOneAndUpdate,
		FindOneAndDelete,
		FindOneAndReplace,
	)
	return list
}

var (
	node, cmd, brunch = zzp.NodeMaker()
	pack              = node(
		node(cmd("collections", zzp.TagSmall)),
		node(cmd("show", zzp.TagSmall)),
		node(cmd("drop", zzp.TagSmall)),
		node(cmd("dropDatabase", zzp.TagSmall)),
		node(get_col_node(node, cmd, brunch)...),
		true,
	)
)

func queryLog(a ...interface{}) {
	if isRemoteMode() == true {
		remoteMsg(a...)
	} else {
		//fmt.Println(a...)
		cc.Println(a...)
	}
}

func jsonString(v interface{}, taps ...string) string {
	tap := "  "
	if len(taps) > 0 {
		tap = taps[0]
	}
	b, _ := json.MarshalIndent(v, "", tap)
	return string(b)
}

func text2Json(text string) (interface{}, error) {
	if text == "nil" || text == "" {
		return VOID{}, nil
	}
	// queryArray := []interface{}{}
	// if err := json.Unmarshal([]byte(text), &queryArray); err == nil {
	// 	return queryArray, nil
	// }
	queryArray := []MAP{}
	if err := json.Unmarshal([]byte(text), &queryArray); err == nil {
		for i := 0; i < len(queryArray); i++ {
			if err := StringToDocument(queryArray[i]); err != nil {
				queryLog("text2Json_StringToDocument:", err)
			}
		}
		return queryArray, nil
	}

	queryOne := map[string]interface{}{}
	if err := json.Unmarshal([]byte(text), &queryOne); err != nil {
		queryLog("text2Json_err:", text)
		return nil, err
	}

	if _, do := queryOne["_id"]; do {
		objectIDString := queryOne["_id"].(string)
		if IsObjectIdHex(objectIDString) {
			queryOne["_id"] = ObjectIDFromHex(objectIDString)
		}
	}

	if err := StringToDocument(queryOne); err != nil {
		queryLog("text2Json_StringToDocument:", err)
	}

	return queryOne, nil
}
func text2Dson(text string) (Dson, error) {
	d := Dson{}
	if text == "nil" || text == "" {
		return Dson{}, nil
	}
	v, e := text2Json(text)
	if e != nil {
		return Dson{}, e
	}
	switch arr := v.(type) {
	case []MAP:
		for _, a := range arr {
			for key, val := range a {
				d.Set(key, val)
			}
		}
		return d, nil

	case map[string]interface{}:
		for key, val := range arr {
			d.Set(key, val)
		}
		return d, nil

	case []interface{}:
		if len(arr) > 0 {
			item := arr[0]
			switch key := item.(type) {
			case map[string]interface{}: //Bson
				_, do := key["Key"]
				_, do2 := key["Value"]
				if do && do2 {

				} else {
					for _, a := range arr {
						aa := a.(map[string]interface{})
						for k, v := range aa {
							d.Set(Eson{k, v})
						}
					} //for
				}
			}
		}
		err := dbg.ParseStruct(v, &d)
		return d, err

	default:

	} //switch
	err := dbg.ParseStruct(v, &d)
	return d, err
}

func resultsView(results []map[string]interface{}, du time.Duration, skipEndTag ...bool) {
	for i := 0; i < len(results); i++ {
		DocumentToString(MAP(results[i]), true)
	}
	if isViewJSONFormat {
		barLog()
		for i, v := range results {
			queryLog("<", i, ">")
			queryLog(jsonString(v, "    "))
		} //for
		barLog()
		queryLog("mms :", du)
		queryLog("document total :", len(results))
		if len(skipEndTag) == 0 {
			barLog()
			queryLog("< END >")
		}
		return
	}
	barLog()
	for idx, rs := range results {
		queryLog("<", idx, ">")
		queryLog("{")
		smap := sm.MakeSortedML(rs)
		smap.Foreach(func(key string, val interface{}) {
			queryKeyValue(smap.KeyLen, key, val, 0)
		})
		queryLog("}")
	} //for

	if len(skipEndTag) == 0 {
		barLog()
		queryLog("< END >")
	}
}
func resultView(reslut map[string]interface{}, du ...time.Duration) {
	DocumentToString(MAP(reslut), true)
	if isViewJSONFormat {
		barLog()
		queryLog(jsonString(reslut, "    "))
		barLog()
		if len(du) == 0 {
			queryLog("< END >")
		} else {
			queryLog("< END", du[0], ">")
		}

		return
	}
	smap := sm.MakeSortedML(reslut)

	barLog()
	queryLog("{")
	smap.Foreach(func(key string, val interface{}) {
		queryKeyValue(smap.KeyLen, key, val, 0)
	})
	queryLog("}")
	barLog()
	queryLog("< END >")
}

func queryKeyValue(keyLen int, key string, val interface{}, depth int) {

	kbuffer := make([]byte, keyLen)
	copy(kbuffer, []byte(key))
	key = string(kbuffer)

	tab := getTab(depth)

	if sm, isDo := val.(map[string]interface{}); isDo == true {
		queryLog(tab, key, ": {")
		queryMap(sm, depth+1)

	} else if sl, isDo := val.([]interface{}); isDo == true {

		if len(sl) == 0 {
			queryLog(tab, key, ": []")
		} else {
			queryLog(tab, key, ": [")
			for _, d := range sl {
				queryVoid(d, depth+1)

				// if mm, isDo := d.(map[string]interface{}); isDo == true {
				// 	queryMap(mm, depth+1)
				// } else {
				// 	queryLog(tab, key, ":", d)
				// }
			}
			queryLog(tab, "]")
		}
	} else {
		queryLog(tab, key, ":", iString(val))

	}
}

func iString(tag interface{}) interface{} {
	if str, isDo := tag.(string); isDo == true {
		return fmt.Sprintf(`"%v"`, str)
	}
	return tag
}

func queryVoid(void interface{}, depth int) {
	tab := getTab(depth)

	if vmap, isDo := void.(map[string]interface{}); isDo == true {
		smap := sm.MakeSortedML(vmap)
		queryLog(tab, "{")
		smap.Foreach(func(key string, val interface{}) {
			queryKeyValue(smap.KeyLen, key, val, depth+1)
		})
		queryLog(tab, "}")

	} else {
		queryLog(tab, iString(void), ",")
	}

}

func queryMap(sub map[string]interface{}, depth int) {
	smap := sm.MakeSortedML(sub)

	tab := getTab(depth)
	//queryLog(tab, "{")
	smap.Foreach(func(key string, val interface{}) {
		queryKeyValue(smap.KeyLen, key, val, depth+1)
	})
	queryLog(tab, "}")
}

func barLog() {
	queryLog("------------------------------------------")
}

func getTab(depth int) string {
	tab := ""
	for i := 0; i < depth; i++ {
		tab = fmt.Sprintf("%v  ", tab)
	}
	return tab
}
func resultSumI2(fields []string, item map[string]interface{}) interface{} {
	var data interface{}

	ldx := len(fields) - 1
	for i, f := range fields {
		if i == 0 {
			if v, do := item[f]; !do {
				return "0"
			} else {
				data = v
				if ldx == i {
					return v
				}
			}
		} else {
			if data == nil {
				return "0"
			}
			pair, do := data.(map[string]interface{})
			if !do {
				return "0"
			}
			isSet := false
			for key, val := range pair {
				if f == key {
					if i == ldx {
						return val
					} else {
						data = val
						break
					}
				}
			}
			if !isSet {
				return "0"
			}

		}
	} //for
	return "0"

}

///////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////

func ConsoleQueryPS3(text string) zzp.Normalizer {
	return pack.Parse(text)
}

// ConsoleQueryCmd3 :
func ConsoleQueryCmd3(cdb *CDB, node zzp.Normalizer) (resultMessage string) {
	defer func() {
		if err := recover(); err != nil {
			queryLog("CDB Query error:", err)
		}
		resultMessage = remoteMessage
		remoteCallback()
	}()

	if err := node.Error(); err != nil {
		return err.Error()
	}

	startRemote()

	cmd := node.Cmd()
	switch cmd {
	case "databasenames", "show":
		dbs, err := cdb.DatabaseNames()
		if err == nil {
			queryLog("===============================")
			queryLog(" Database List")
			queryLog("-------------------------------")
			for _, d := range dbs {
				queryLog(" " + d)
			}
			queryLog("===============================")
			queryLog("< END >")
		} else {
			queryLog(fmt.Sprintf("%v_error : %v", cmd, err))
		}
		return
	}

	dbName := cmd
	switch dbName {
	case "drop", "dropDatabase":
		systemDBList := []string{
			"admin",
			"config",
			"local",
		}

		for _, system := range systemDBList {
			if dbName == system {
				queryLog("---------------------------------------------")
				queryLog("system db is not drop. [admin/config/local]")
				queryLog("---------------------------------------------")
				return
			}
		} //for

		if err := cdb.DropDatabase(dbName); err != nil {
			queryLog(dbg.Error(err))
		} else {
			queryLog("success.")
			queryLog("< END >")
		}
		return
	}

	next := node.Next()
	if next == nil {
		queryLog("missing db next", dbName)
		return
	}

	name := next.Cmd()
	switch name {
	case "collections", "show":
		colNames, err := cdb.CollectionNames(dbName)
		if err != nil {
			queryLog("err:", err)
		} else {
			queryLog("======================================================")
			queryLog(" DB :", dbName)
			queryLog("------------------------------------------------------")
			for _, v := range colNames {
				queryLog(" ", v)
			}
			queryLog("======================================================")
			queryLog("< END >")
		}
		return
	}
	colName := name
	col := next.Next()
	if col == nil {
		dbg.Cat("minssing collection query : ", dbName, ".", colName)
		return
	}

	col_action := func(f func(c Collection)) {
		cdb.Action(dbName, func(db DATABASE) {
			f(db.C(colName))
		})
	}

	command := col.Cmd()
	params := col.Params()
	switch command {
	default:
		queryLog("command error :", command)

	case "ensureindex":
		col_action(func(c Collection) {
			key := params[0]
			order := params[1]
			unique := dbg.IsTrue(params[2])
			if err := c.EnsureIndex(SingleIndex(key, order, unique)); err != nil {
				queryLog(err)
			} else {
				queryLog("success.")
				queryLog("< END >")
			}
		})

	case "dropindexname":
		col_action(func(c Collection) {
			if err := c.DropIndex(params...); err != nil {
				queryLog(err)
			} else {
				queryLog("success.")
				queryLog("< END >")
			}
		})

	case "dropindexAll":
		col_action(func(c Collection) {
			err := c.DropIndexAll()
			if err != nil {
				queryLog(dbg.Error(err))
			} else {
				queryLog(fmt.Sprintf("success."))
				queryLog("< END >")
			}
		})

	case "indexes":
		col_action(func(c Collection) {
			idxs, err := c.Indexes()
			if err != nil {
				queryLog("indexes_error:", err)
			} else {
				queryLog("===============================")
				for i, idx := range idxs {
					queryLog("<", i, ">")
					queryLog(jsonString(idx))
					queryLog("")
				}
				queryLog("===============================")
				queryLog("< END >")
			}
		})

	case "dropcollection":
		col_action(func(c Collection) {
			if err := c.DropCollection(); err != nil {
				queryLog(dbg.Error(err))
			} else {
				queryLog(fmt.Sprintf("success."))
				queryLog("< END >")
			}
		})

	case "removeAll":
		col_action(func(c Collection) {
			if fQuery, err := text2Json(params[0]); err == nil {
				qr := c.RemoveAll(fQuery)
				queryLog(jsonString(qr))
			}
		})

	case "remove":
		col_action(func(c Collection) {
			if fQuery, err := text2Json(params[0]); err == nil {
				qr := c.Remove(fQuery)
				if qr.Error == nil {
					queryLog("remove_success")
					queryLog("< END >")
				} else {
					queryLog("remove_error :", qr.Error)
				}
			}
		})

	case "insert":
		col_action(func(c Collection) {
			fQuery, err := text2Json(params[0])
			if err != nil {
				queryLog(err)
				return
			}
			if err := c.Insert(fQuery); err != nil {
				queryLog("insert_error :", err)
			} else {
				queryLog("insert_success!")
				queryLog("< END >")
			}
		})

	case "updateAll":
		col_action(func(c Collection) {
			fQuery, err := text2Json(params[0])
			if err != nil {
				queryLog(err)
				return
			}
			dQuery, err := text2Json(params[1])
			if err != nil {
				queryLog(err)
				return
			}
			if qr := c.UpdateAll(fQuery, dQuery); qr.Error != nil {
				queryLog("updateAll_error :", qr.Error)
			} else {
				queryLog("updateAll_result :", jsonString(qr))
				queryLog("< END >")
			}
		})

	case "update":
		col_action(func(c Collection) {
			fQuery, err := text2Json(params[0])
			if err != nil {
				queryLog(err)
				return
			}
			dQuery, err := text2Json(params[1])
			if err != nil {
				queryLog(err)
				return
			}
			if qr := c.Update(fQuery, dQuery); qr.Error != nil {
				queryLog("update_error :", qr.Error)
			} else {
				queryLog("update_success!")
				queryLog("< END >")
			}
		})

	case "upsert":
		col_action(func(c Collection) {
			fQuery, err := text2Json(params[0])
			if err != nil {
				queryLog(err)
				return
			}
			dQuery, err := text2Json(params[1])
			if err != nil {
				queryLog(err)
				return
			}
			if qr := c.Upsert(fQuery, dQuery); qr.Error != nil {
				queryLog("upsert_error :", qr.Error)
			} else {
				queryLog("upsert_result :", jsonString(qr))
				queryLog("< END >")
			}
		})

	case "aggregate":
		col_action(func(c Collection) {
			fQuery, err := text2Json(params[0])
			if err != nil {
				queryLog(err)
				return
			}

			iter := c.Aggregate(fQuery)

			at := time.Now()
			switch col.Next().Cmd() {
			case "one":
				var result map[string]interface{}
				if err := iter.One(&result); err != nil {
					queryLog("aggregate_error:", err)
				} else {
					resultView(result, time.Now().Sub(at))
				}

			case "all":
				var results []map[string]interface{}
				if err := iter.All(&results); err != nil {
					queryLog("aggregate_error:", err)
				} else {
					resultsView(results, time.Now().Sub(at))
				}
			}
		})

	case "FindOneAndUpdate":
		col_action(func(c Collection) {
			fQuery, err := text2Json(params[0])
			if err != nil {
				queryLog(err)
				return
			}
			dQuery, err := text2Json(params[1])
			if err != nil {
				queryLog(err)
				return
			}

			var result map[string]interface{}
			if err := c.FindOneAndUpdate(fQuery, dQuery, &result); err != nil {
				queryLog("FindOneAndUpdate_error :", err)
			} else {
				resultView(result)
			}
		})

	case "FindOneAndDelete":
		col_action(func(c Collection) {
			fQuery, err := text2Json(params[0])
			if err != nil {
				queryLog(err)
				return
			}
			var result map[string]interface{}
			if err := c.FindOneAndDelete(fQuery, &result); err != nil {
				queryLog("FindOneAndDelete_error :", err)
			} else {
				resultView(result)
			}
		})

	case "FindOneAndReplace":
		col_action(func(c Collection) {
			fQuery, err := text2Json(params[0])
			if err != nil {
				queryLog(err)
				return
			}
			dQuery, err := text2Json(params[1])
			if err != nil {
				queryLog(err)
				return
			}

			var result map[string]interface{}
			if err := c.FindOneAndReplace(fQuery, dQuery, &result); err != nil {
				queryLog("FindOneAndReplace_error :", err)
			} else {
				resultView(result)
			}
		})

	case "find":
		col_action(func(c Collection) {
			if fQuery, err := text2Json(params[0]); err == nil {

				mgoQuery := c.Find(fQuery)
				brunch := col.Brunch()
				brunch.Callback("skip", func(s []string) {
					skip := jmath.Int(s[0])
					mgoQuery = mgoQuery.Skip(skip)
				})
				brunch.Callback("limit", func(s []string) {
					limit := jmath.Int(s[0])
					mgoQuery = mgoQuery.Limit(limit)
				})
				brunch.Callback("sort", func(s []string) {
					mgoQuery = mgoQuery.Sort(s...)
				})
				brunch.Callback(_hint_tag, func(s []string) {
					if hint_query, err := text2Dson(s[0]); err == nil {
						mgoQuery = mgoQuery.Hint(hint_query)
					} else {
						queryLog("hint_query :", err)
					}
				})

				at := time.Now()
				action := col.Next()
				switch action.Cmd() {
				case "all":
					var results []map[string]interface{}

					if err := mgoQuery.All(&results); err != nil {
						queryLog("find.all_error:", err)
					} else {
						resultsView(results, time.Since(at))
					}

				case "one":
					var result map[string]interface{}
					if err := mgoQuery.One(&result); err != nil {
						queryLog("find.one_error:", err)
					} else {
						resultView(result, time.Since(at))
						//queryLog(jsonString(result))
					}

				case "count":
					if cnt, err := mgoQuery.Count(); err != nil {
						queryLog("find.one_error:", err)
					} else {
						queryLog("find.count :", cnt)
						queryLog("< END", time.Since(at), ">")
					}

				case "sum": //sumi

					iter := mgoQuery.Iter()

					sumstrs := action.Params()

					sumVal := "0"
					result := map[string]interface{}{}
					for iter.Next(&result) {
						v := resultSumI2(sumstrs, result)
						sumVal = jmath.ADD(sumVal, v)
					} //for
					iter.Close()

					queryLog("==========================================")
					queryLog("sumVal :", sumVal)
					queryLog("du : ", time.Since(at))
					queryLog("==========================================")

				} //switch
			}
		})

	case "findProjection":
		col_action(func(c Collection) {
			dp, err := text2Dson(params[1])
			if err != nil {
				queryLog("find.text2Dson:", err)

			} else {
				if fQuery, err := text2Json(params[0]); err == nil {

					mgoQuery := c.FindProjection(fQuery, dp)

					brunch := col.Brunch()
					brunch.Callback("skip", func(s []string) {
						skip := jmath.Int(s[0])
						mgoQuery = mgoQuery.Skip(skip)
					})
					brunch.Callback("limit", func(s []string) {
						limit := jmath.Int(s[0])
						mgoQuery = mgoQuery.Limit(limit)
					})
					brunch.Callback("sort", func(s []string) {
						mgoQuery = mgoQuery.Sort(s...)
					})
					brunch.Callback(_hint_tag, func(s []string) {
						if hint_query, err := text2Dson(s[0]); err == nil {
							mgoQuery = mgoQuery.Hint(hint_query)
						} else {
							queryLog("hint_query :", err)
						}
					})

					at := time.Now()
					action := col.Next()
					switch action.Cmd() {
					case "all":
						var results []map[string]interface{}

						if err := mgoQuery.All(&results); err != nil {
							queryLog("find.all_error:", err)
						} else {
							resultsView(results, time.Since(at))
						}

					case "one":
						var result map[string]interface{}
						if err := mgoQuery.One(&result); err != nil {
							queryLog("find.one_error:", err)
						} else {
							resultView(result, time.Since(at))
							//queryLog(jsonString(result))
						}

					case "count":
						if cnt, err := mgoQuery.Count(); err != nil {
							queryLog("find.one_error:", err)
						} else {
							queryLog("find.count :", cnt)
							queryLog("< END", time.Since(at), ">")
						}

					case "sum": //sumi

						iter := mgoQuery.Iter()

						sumstrs := action.Params()

						sumVal := "0"
						result := map[string]interface{}{}
						for iter.Next(&result) {
							v := resultSumI2(sumstrs, result)
							sumVal = jmath.ADD(sumVal, v)
						} //for
						iter.Close()

						queryLog("==========================================")
						queryLog("sumVal :", sumVal)
						queryLog("du : ", time.Since(at))
						queryLog("==========================================")

					} //switch
				}
			}

		})
	} //switch

	return resultMessage
}
