package admin

import (
	"net/http"
	"strings"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/jnet/doc"
	"txscheduler/txm/ack"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

const (
	queryLimit = 200
)

//CAdminQueryFind :
type CAdminQueryFind struct {
	COL      string                 `json:"col"`
	Selector map[string]interface{} `json:"selector,omitempty"`
	Sort     string                 `json:"sort,omitempty"`
	Offset   int                    `json:"offset,omitempty"`
	Limit    int                    `json:"limit,omitempty"`
	Cmd      string                 `json:"cmd"`
}

//TagString :
func (CAdminQueryFind) TagString() []string {
	return []string{
		"col", "컬렉션 이름",
		"selector,omitempty", "조건",
		"sort,omitempty", "소팅순서 (ex) timestamp(오름차순), -timestamp(내림차순)",
		"offset,omitempty", "쿼리조건의 갯수 제한에서의 오프셋",
		"limit,omitempty", "요청 갯수 제한",
		"cmd", "one(한개), all(리스트), count(갯수)",
	}
}

//Valid :
func (my *CAdminQueryFind) Valid() bool {
	my.COL = strings.TrimSpace(my.COL)
	my.Sort = strings.TrimSpace(my.Sort)
	my.Cmd = dbg.TrimToLower(my.Cmd)

	if my.Limit > queryLimit || my.Limit < 0 {
		my.Limit = queryLimit
	}
	if my.COL == "" {
		return false
	}

	switch my.Cmd {
	case "one", "all", "count":
	default:
		return false
	}
	return true
}

type cAdminQueryFindResult struct {
	TotalCount int         `json:"total_count"`
	Data       interface{} `json:"data"`
}

//TagString :
func (cAdminQueryFindResult) TagString() []string {
	return []string{
		"total_count", "요청한 쿼리의 총갯수",
		"data", "요청타입에따른 데이타 (one이면 한개 또는 null, list이면 여러개 또는 [], count이면 null)",
	}
}

//____
func init() {
	hQueryFind()
}

func hQueryFind() {
	errorBadParam := ack.BadParam
	errorDBJob := ack.DBJob

	method := chttp.POST
	url := model.V2 + "/query.find"

	Doc().Comment("[ DB쿼리 ] 컬렉션 데이타 DB쿼리 요청").
		Method(method).URL(url).
		JParam(CAdminQueryFind{}, CAdminQueryFind{}.TagString()...).
		JResultOK(cAdminQueryFindResult{}, cAdminQueryFindResult{}.TagString()...).
		ResultERRR(errorBadParam).
		ResultERRR(errorDBJob).
		Apply(doc.Blue)

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := CAdminQueryFind{}
			if chttp.BindingJSON(req, &cdata) != nil {
				chttp.Fail(w, errorBadParam, "BindingJSON")
				return
			}
			if cdata.Valid() == false {
				chttp.Fail(w, errorBadParam, "Valid")
				return
			}

			inf.DB().Action(inf.DBName, func(db mongo.DATABASE) {
				c := db.C(cdata.COL)

				selector := mongo.MapToBson(cdata.Selector)
				cQuery := c.Find(selector)
				if cdata.Sort != "" {
					cQuery = cQuery.Sort(cdata.Sort)
				}

				totalCount, err := cQuery.Count()
				if err != nil {
					chttp.Fail(w, errorDBJob, err)
					return
				}

				switch cdata.Cmd {
				case "one":
					data := map[string]interface{}{}
					if err := cQuery.One(&data); err != nil {
						chttp.Fail(w, errorDBJob, err)
					} else {
						chttp.OK(w, cAdminQueryFindResult{
							TotalCount: totalCount,
							Data:       data,
						})
					}

				case "all":

					datalist := []map[string]interface{}{}

					if cdata.Limit > 0 {
						if cdata.Offset > 0 {
							iter := cQuery.Iter()
							item := map[string]interface{}{}
							for iter.Next(&item) {
								if cdata.Offset > 0 {
									item = map[string]interface{}{}
									cdata.Offset--
									continue
								}
								datalist = append(datalist, item)
								item = map[string]interface{}{}
								cdata.Limit--
								if cdata.Limit <= 0 {
									break
								}
							} //for
						} else {
							cQuery = cQuery.Limit(cdata.Limit)
							cQuery.All(&datalist)
						}

					} else {
						cQuery.All(&datalist)
					}

					chttp.OK(w, cAdminQueryFindResult{
						TotalCount: totalCount,
						Data:       datalist,
					})

				case "count":
					chttp.OK(w, cAdminQueryFindResult{
						TotalCount: totalCount,
						Data:       nil,
					})
				}

			})
		},
	)
}
