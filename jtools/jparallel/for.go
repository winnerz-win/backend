package jparallel

import (
	"context"
	"fmt"
	"jtools/jlog"
	"sync"
	"time"
)

// 병렬 For. 아무것도 리턴하지 않는 버전 (제일빠름)
func ForNotReturn(maxLoopCount int, f func(int), runPerSec int) {
	dummyResultFunc := func(i int) (interface{}, error) {
		f(i)
		return nil, nil
	}
	forProcess(maxLoopCount, dummyResultFunc, runPerSec, nil, false)
}

// 병렬 For. 아무것도 리턴하지 않고 runPerSec 단위로 startTime ~ endTime 시간동안 돌아가는 버전
func ForTimeNotReturn(f func(int), startTime, endTime time.Time, runPerSec int) {
	dummyResultFunc := func(i int) (interface{}, error) {
		f(i)
		return nil, nil
	}
	now := time.Now()
	if startTime.Before(now) {
		jlog.Error("ForTimeNotReturn - startTime.Before(now)", "startTime:", startTime, "now:", now)
		return
	}
	if startTime.After(endTime) {
		jlog.Error("ForTimeNotReturn - startTime.After(endTime)", "startTime:", startTime, "endTime:", endTime)
		return
	}
	maxLoopCount := ((int)(endTime.Sub(startTime) / time.Second)) * runPerSec
	time.Sleep(time.Until(startTime))
	forProcess(maxLoopCount, dummyResultFunc, runPerSec, nil, false)
}

// 병렬 For. 에러들과 결과들을 리턴하는 버전
func For[TResult any](maxLoopCount int, f func(int) (TResult, error), runPerSec int) ([]TResult, PErrors, PTimeResults) {
	return forProcess(maxLoopCount, f, runPerSec, nil, true)
}

// 병렬 For. 점진적으로 RunPerSec 가 늘어나는 버전
func ForIncRunPerSec[TResult any](ctx context.Context, f func(int) (TResult, error), runSec, minRunPerSec, maxRunPerSec int) ([]TResult, PErrors, PTimeResults) {
	runPerSecTable := getIncRunPerSecTable(runSec, minRunPerSec, maxRunPerSec)
	return forProcess(len(runPerSecTable), f, 0, runPerSecTable, true)
}

// 병렬 For. runPerSec 단위로 startTime ~ endTime 시간동안 돌아가는 버전
func ForTime[TResult any](f func(int) (TResult, error), startTime, endTime time.Time, runPerSec int) ([]TResult, PErrors, PTimeResults) {
	if startTime.Before(time.Now()) {
		return nil, PErrors{PError{Index: 0, Err: fmt.Errorf("for time error - start time before time now")}}, nil
	}
	if startTime.After(endTime) {
		return nil, PErrors{PError{Index: 0, Err: fmt.Errorf("for time error - start time after end time")}}, nil
	}
	maxLoopCount := ((int)(endTime.Sub(startTime) / time.Second)) * runPerSec
	time.Sleep(time.Until(startTime))
	return forProcess(maxLoopCount, f, runPerSec, nil, true)
}

// 병렬 For 로직
func forProcess[TResult any](maxLoopCount int, f func(int) (TResult, error), runPerSec int, runPerSecTable []time.Duration, isErrorOrResult bool) ([]TResult, PErrors, PTimeResults) {
	results := make([]TResult, maxLoopCount)
	errors := make(PErrors, maxLoopCount)
	times := make(PTimeResults, maxLoopCount)

	wg := sync.WaitGroup{}
	wg.Add(maxLoopCount)

	for i := 0; i < maxLoopCount; i++ {
		go func(_i int) {
			startTime := time.Now()
			defer func() {
				wg.Done()
				r := recover()
				if r == nil {
					return
				}
				if err, ok := r.(error); ok {
					fmt.Println("Recovered from error:", err)
					if isErrorOrResult {
						results[_i] = interface{}(nil).(TResult)
						errors[_i] = PError{_i, err}
						times[_i] = PTimeResult{_i, startTime, time.Now()}
					}
				} else {
					fmt.Println("Recovered, but not from an error:", r)
					if isErrorOrResult {
						results[_i] = interface{}(nil).(TResult)
						errors[_i] = PError{_i, fmt.Errorf("recovered, but not from an error: %v", r)}
						times[_i] = PTimeResult{_i, startTime, time.Now()}
					}
				}
			}()
			if runPerSecTable == nil {
				skipI := _i
				//skipI := math.Max(float64(_i-700), 1)
				per_du := (time.Second / time.Duration(runPerSec))
				sleep_du := time.Duration(skipI) * per_du
				time.Sleep(sleep_du)
			} else {
				time.Sleep(runPerSecTable[_i])
			}
			result, err := f(_i)

			if isErrorOrResult {
				results[_i] = result
				errors[_i] = PError{_i, err}
				times[_i] = PTimeResult{_i, startTime, time.Now()}
			}
		}(i)
	}
	wg.Wait()

	return results, errors, times
}
