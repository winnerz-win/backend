package jaggregation

import (
	"fmt"
	"jtools/jlog"
	"math"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

var Client *Aggregation = &Aggregation{}

type ConfigAggregation struct {
	IsAggregation  bool     `yaml:"is_aggregation"`
	StartLogString string   `yaml:"start_log_string"`
	AccuracySec    int      `yaml:"accuracy_sec"`
	ExcelName      string   `yaml:"excel_name"`
	ReqNames       []string `yaml:"req_names"`
}

type Aggregation struct {
	is_init        bool
	is_aggregation bool
	accuracy       time.Duration // 1회당 측정 시간
	accuracySec    int           // accuracy 를 초로 환산한 값

	vUserCnt            int32             // vuser 수
	reqNames            []string          // res 이름 모음
	resAggCntMap        map[string]*int32 // 1회 측정동안 응답받은 횟수 (api 당 따로 체크)
	resAggMsMap         map[string]*int32 // 1회 측정동안 응답시간 총합 (api 당 따로 체크)
	resAggErrCntMap     map[string]*int32 // 1회 측정동안 응답에러 횟수 (api 당 따로 체크)
	resAggManyReqCntMap map[string]*int32 // 1회 측정동안 많은요청 횟수 (api 당 따로 체크)

	accuracyCnt          int             // 현재까지 측정한 횟수 (vUserCnt 1 이상부터 집계)
	sb                   strings.Builder // 로그 찍기 위한 스트링빌더
	excelName            string          // 엑셀 파일 이름
	startLogString       string          // 첫 로그에 작성할 내용
	startAggregationTime time.Time       // 집계 처음 시작한 시간 (accuracy 1회 이상부터)
}

func (my *Aggregation) isValid() bool {
	if my == nil || !my.is_aggregation {
		return false
	}
	return my.is_init
}

// 사용할 reqNames 전부 미리 등록해야함 (도중에 추가하려면 락 걸어야되서 처음에 전부 받음)
func (my *Aggregation) Init(
	isAggregation bool,
	startLogString string,
	accuracySec int,
	excelName string,
	reqNames ...string,
) {
	if my == nil {
		return
	}
	my.is_aggregation = isAggregation
	if !my.is_aggregation {
		return
	}
	if len(reqNames) == 0 {
		return
	}
	my.is_init = true
	my.accuracy = time.Second * time.Duration(accuracySec)
	my.accuracySec = accuracySec
	my.attachReq(reqNames...)
	my.sb.Grow(len(reqNames) * 100)
	my.excelName = excelName
	my.startLogString = startLogString
}
func (my *Aggregation) InitConfig(config ConfigAggregation) {
	if my == nil {
		return
	}
	my.Init(
		config.IsAggregation,
		config.StartLogString,
		config.AccuracySec,
		config.ExcelName,
		config.ReqNames...,
	)
}
func (my *Aggregation) attachReq(reqNames ...string) {
	if !my.isValid() {
		return
	}
	if my.reqNames == nil {
		my.reqNames = []string{}
	}
	if my.resAggCntMap == nil {
		my.resAggCntMap = map[string]*int32{}
	}
	if my.resAggMsMap == nil {
		my.resAggMsMap = map[string]*int32{}
	}
	if my.resAggErrCntMap == nil {
		my.resAggErrCntMap = map[string]*int32{}
	}
	if my.resAggManyReqCntMap == nil {
		my.resAggManyReqCntMap = map[string]*int32{}
	}
	notDuplicateReqNames := []string{}
	m := make(map[string]struct{})
	for _, val := range reqNames {
		if _, ok := m[val]; !ok {
			m[val] = struct{}{}
			notDuplicateReqNames = append(notDuplicateReqNames, val)
		}
	}
	for _, reqName := range notDuplicateReqNames {
		my.reqNames = append(my.reqNames, reqName)
		var resCnt int32 = 0
		my.resAggCntMap[reqName] = &resCnt
		var resMs int32 = 0
		my.resAggMsMap[reqName] = &resMs
		var resErrCnt int32 = 0
		my.resAggErrCntMap[reqName] = &resErrCnt
		var resManyReqCnt int32 = 0
		my.resAggManyReqCntMap[reqName] = &resManyReqCnt
	}
}
func (my *Aggregation) IsFirst() bool {
	if !my.isValid() {
		return false
	}
	return my.accuracyCnt == 1
}
func (my *Aggregation) IncUserCnt() {
	if !my.isValid() {
		return
	}
	atomic.AddInt32(&my.vUserCnt, 1)
}
func (my *Aggregation) DecUserCnt() {
	if !my.isValid() {
		return
	}
	atomic.AddInt32(&my.vUserCnt, -1)
}
func (my *Aggregation) IncAgg(reqName string, resDu time.Duration) {
	if !my.isValid() {
		return
	}
	if _, has := my.resAggCntMap[reqName]; !has {
		return
	}
	atomic.AddInt32(my.resAggCntMap[reqName], 1)
	resMs := float64(resDu / time.Millisecond)
	atomic.AddInt32(my.resAggMsMap[reqName], int32(math.Max(resMs, 1)))
}
func (my *Aggregation) IncError(reqName string) {
	if !my.isValid() {
		return
	}
	if _, has := my.resAggCntMap[reqName]; !has {
		return
	}
	atomic.AddInt32(my.resAggErrCntMap[reqName], 1)
}
func (my *Aggregation) IncManyReq(reqName string) {
	if !my.isValid() {
		return
	}
	if _, has := my.resAggManyReqCntMap[reqName]; !has {
		return
	}
	atomic.AddInt32(my.resAggManyReqCntMap[reqName], 1)
}

func (my *Aggregation) StartAggregation() {
	if !my.isValid() {
		return
	}
	go func() {
		for {
			// 측정기간마다 한타임씩 기록
			time.Sleep(my.accuracy)

			// 메모리에 먼저 저장
			vUserCnt := int(atomic.LoadInt32(&my.vUserCnt))

			// vUser 1 이상부터 집계 시작
			if vUserCnt == 0 {
				my.startAggregationTime = time.Now()
				continue
			}

			my.accuracyCnt++

			tpsMap := map[string]float32{}
			resMsMap := map[string]int{}
			epsMap := map[string]int{}
			mrpsMap := map[string]float32{}
			for _, reqName := range my.reqNames {
				resAggCntLoad := int(atomic.LoadInt32(my.resAggCntMap[reqName]))
				resAggMsLoad := int(atomic.LoadInt32(my.resAggMsMap[reqName]))
				resAggErrCntLoad := int(atomic.LoadInt32(my.resAggErrCntMap[reqName]))
				resAggManyReqCntMap := int(atomic.LoadInt32(my.resAggManyReqCntMap[reqName]))

				// 이름에 Agg 붙은 것을은 매 측정마다 다시 집계하므로 초기화
				atomic.SwapInt32(my.resAggCntMap[reqName], 0)
				atomic.SwapInt32(my.resAggMsMap[reqName], 0)
				atomic.SwapInt32(my.resAggErrCntMap[reqName], 0)
				atomic.SwapInt32(my.resAggManyReqCntMap[reqName], 0)

				tps := float32(resAggCntLoad) / float32(my.accuracySec)
				ms := 0
				if resAggCntLoad > 0 {
					ms = resAggMsLoad / resAggCntLoad
				}
				eps := resAggErrCntLoad / my.accuracySec
				mrps := float32(resAggManyReqCntMap) / float32(my.accuracySec)

				tpsMap[reqName] = tps
				resMsMap[reqName] = ms
				epsMap[reqName] = eps
				mrpsMap[reqName] = mrps
			}

			// 로그에 씀
			if my.IsFirst() {
				my.sb.WriteString("\n================ Test Start ==============\n")
				my.sb.WriteString(my.startLogString)
			}
			my.sb.WriteString("\nvUserCnt: ")
			my.sb.WriteString(strconv.Itoa(vUserCnt))
			my.sb.WriteString("  |  Running Time : ")
			my.sb.WriteString(time.Since(my.startAggregationTime).String())
			my.sb.WriteString("\n")
			for _, reqName := range my.reqNames {
				tps := tpsMap[reqName]
				resMs := resMsMap[reqName]
				eps := epsMap[reqName]
				mrps := mrpsMap[reqName]
				my.sb.WriteString(fmt.Sprintf("%-42v : [tps:%7.2f] [resMs:%7v] [eps:%7v] [mrps:%7.2f]\n", reqName, tps, resMs, eps, mrps))
			}
			jlog.Info(my.sb.String())
			my.sb.Reset()

			// 엑셀에 씀
			if my.excelName != "" {
				now := time.Now().UTC()
				const dataStartCell = 1
				dataStartCellStr := strconv.Itoa(dataStartCell)
				cellNumStr := strconv.Itoa(dataStartCell + my.accuracyCnt)

				file, err := excelize.OpenFile(my.excelName)
				if err != nil {
					file = excelize.NewFile()
				}

				if my.IsFirst() {
					title := "Params"
					sheet := file.NewSheet(title)
					file.SetActiveSheet(sheet)
					file.SetCellValue(title, "A"+cellNumStr, my.startLogString)
				}

				title := "Running VUser"
				sheet := file.NewSheet(title)
				file.SetActiveSheet(sheet)
				if my.IsFirst() {
					file.SetCellValue(title, "A"+dataStartCellStr, "Elapsed Time")
					file.SetCellValue(title, "B"+dataStartCellStr, "VUser Count")
				}
				file.SetCellValue(title, "A"+cellNumStr, now)
				file.SetCellValue(title, "B"+cellNumStr, vUserCnt)

				title = "TPS"
				sheet = file.NewSheet(title)
				file.SetActiveSheet(sheet)
				if my.IsFirst() {
					file.SetCellValue(title, "A"+dataStartCellStr, "Elapsed Time")
					for i, reqName := range my.reqNames {
						startCell := PlusCol("B", i)
						file.SetCellValue(title, string(startCell)+dataStartCellStr, reqName)
					}
				}
				file.SetCellValue(title, "A"+cellNumStr, now)
				for i, reqName := range my.reqNames {
					startCell := PlusCol("B", i)
					file.SetCellValue(title, string(startCell)+cellNumStr, tpsMap[reqName])
				}

				title = "Response Time (ms)"
				sheet = file.NewSheet(title)
				file.SetActiveSheet(sheet)
				if my.IsFirst() {
					file.SetCellValue(title, "A"+dataStartCellStr, "Elapsed Time")
					for i, reqName := range my.reqNames {
						startCell := PlusCol("B", i)
						file.SetCellValue(title, string(startCell)+dataStartCellStr, reqName)
					}
				}
				file.SetCellValue(title, "A"+cellNumStr, now)
				for i, reqName := range my.reqNames {
					startCell := PlusCol("B", i)
					file.SetCellValue(title, string(startCell)+cellNumStr, resMsMap[reqName])
				}

				title = "Error Per Sec"
				sheet = file.NewSheet(title)
				file.SetActiveSheet(sheet)
				if my.IsFirst() {
					file.SetCellValue(title, "A"+dataStartCellStr, "Elapsed Time")
					for i, reqName := range my.reqNames {
						startCell := PlusCol("B", i)
						file.SetCellValue(title, string(startCell)+dataStartCellStr, reqName)
					}
				}
				file.SetCellValue(title, "A"+cellNumStr, now)
				for i, reqName := range my.reqNames {
					startCell := PlusCol("B", i)
					file.SetCellValue(title, string(startCell)+cellNumStr, epsMap[reqName])
				}

				title = "ManyReq Per Sec"
				sheet = file.NewSheet(title)
				file.SetActiveSheet(sheet)
				if my.IsFirst() {
					file.SetCellValue(title, "A"+dataStartCellStr, "Elapsed Time")
					for i, reqName := range my.reqNames {
						startCell := PlusCol("B", i)
						file.SetCellValue(title, string(startCell)+dataStartCellStr, reqName)
					}
				}
				file.SetCellValue(title, "A"+cellNumStr, now)
				for i, reqName := range my.reqNames {
					startCell := PlusCol("B", i)
					file.SetCellValue(title, string(startCell)+cellNumStr, mrpsMap[reqName])
				}

				err = file.SaveAs(my.excelName)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()
}

// A+1 == B, Z+2 == AB 를 리턴하는 함수
func PlusCol(alpha string, num int) string {
	// 알파벳 문자열을 배열로 변환
	alphaArr := []rune(alpha)

	// 배열의 마지막 요소부터 계산하여 num을 더함
	carry := num
	for i := len(alphaArr) - 1; i >= 0; i-- {
		digit := int(alphaArr[i] - 'A' + 1)
		sum := digit + carry
		carry = sum / (26 + 1)
		alphaArr[i] = rune('A' + ((sum - 1) % 26))
	}
	// carry 값이 1이면, 맨 앞에 새로운 문자 'A'를 추가
	if carry == 1 {
		alphaArr = append([]rune{'A'}, alphaArr...)
	}
	// 결과 문자열을 반환
	return string(alphaArr)
}
