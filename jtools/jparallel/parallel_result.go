package jparallel

import (
	"fmt"
	"math"
	"time"
)

// 발생한 에러
type PErrors []PError

func (my PErrors) Len() int {
	return len(my)
}
func (my PErrors) Swap(i, j int) {
	my[i].Index, my[j].Index = my[j].Index, my[i].Index
}
func (my PErrors) Less(i, j int) bool {
	return my[i].Index < my[j].Index
}

// 총 에러 개수
func (my PErrors) ErrorCount() int {
	errorCount := 0
	for _, forError := range my {
		if forError.Error() != nil {
			errorCount++
		}
	}
	return errorCount
}
func (my PErrors) Error() error {
	for _, forError := range my {
		if forError.Error() != nil {
			return forError.Error()
		}
	}
	return nil
}

type PError struct {
	Index int   // 실행한 순서
	Err   error // 발생한 에러
}

func (my PError) Error() error {
	if my.Err != nil {
		return fmt.Errorf(
			"parallel error - index[%d] error : %s",
			my.Index, my.Err)
	} else {
		return nil
	}
}

// 실행 결과
type PTimeResults []PTimeResult

func (my PTimeResults) Len() int {
	return len(my)
}
func (my PTimeResults) Swap(i, j int) {
	my[i].Index, my[j].Index = my[j].Index, my[i].Index
}
func (my PTimeResults) Less(i, j int) bool {
	return my[i].Index < my[j].Index
}

type PTimeResult struct {
	Index     int       // 실행한 순서
	StartTime time.Time // 작업 시작 시간
	EndTime   time.Time // 작업 종료 시간
}

// 하나의 작업에 대해 걸린 시간
func (my PTimeResult) TotalRunTime() time.Duration {
	return my.EndTime.Sub(my.StartTime)
}

// 점진적으로 늘어나는 RunPerSec 테이블 구하기
func getIncRunPerSecTable(runSec, minRunPerSec, maxRunPerSec int) []time.Duration {
	runSecF := float64(runSec)             // runSecond 초 동안
	minRunPerSecF := float64(minRunPerSec) // minRunPerSec 부터 시작하여
	maxRunPerSecF := float64(maxRunPerSec) // maxRunPerSec 까지 점진적으로 RunPerSec 늘리며 테스트 진행

	minSleepDelay := 1 * time.Second                 // sleep 최소단위를 피하기 위한 초기 딜레이
	runPerSecTable := []time.Duration{minSleepDelay} // 타임 테이블을 미리 짜고 그 타임 테이블대로 테스트 진행

	// 각 초의 진행률을 기반으로 현재 초의 RunPerSec 를 구해서 타임 테이블 짬
	count := 0
	for i := 1; i <= int(runSecF); i++ {
		runPer := float64(i) / runSecF
		curRunPerSec := math.Max(math.Round(maxRunPerSecF*runPer), minRunPerSecF)
		oneDuration := time.Duration(float64(time.Second) / curRunPerSec)
		for j := 1; j <= int(curRunPerSec); j++ {
			count++
			runPerSecTable = append(runPerSecTable, runPerSecTable[count-1]+oneDuration)
		}
	}

	return runPerSecTable
}
