package jparallel

import (
	"fmt"
	"sync"
	"time"
)

// 병렬 Foreach. 아무것도 리턴하지 않는 버전 (제일빠름)
func ForeachNotReturn[TElement any](elements []TElement, f func(int, TElement), runPerSec int) {
	dummyResultFunc := func(i int, e TElement) (interface{}, error) {
		f(i, e)
		return nil, nil
	}
	foreachProcess(elements, dummyResultFunc, runPerSec, nil, false)
}

// 병렬 Foreach. 에러들과 결과들을 리턴하는 버전
func Foreach[TElement any, TResult any](elements []TElement, f func(int, TElement) (TResult, error), runPerSec int) ([]TResult, PErrors, PTimeResults) {
	return foreachProcess(elements, f, runPerSec, nil, true)
}

// 병렬 Foreach. 점진적으로 RunPerSec 가 늘어나는 버전
func ForeachIncRunPerSec[TElement any, TResult any](elements []TElement, f func(int, TElement) (TResult, error), runSec, minRunPerSec, maxRunPerSec int) ([]TResult, PErrors, PTimeResults) {
	runPerSecTable := getIncRunPerSecTable(runSec, minRunPerSec, maxRunPerSec)
	return foreachProcess(elements, f, 0, runPerSecTable, true)
}

// 병렬 Foreach 로직
func foreachProcess[TElement any, TResult any](elements []TElement, f func(int, TElement) (TResult, error), runPerSec int, runPerSecTable []time.Duration, isErrorOrResult bool) ([]TResult, PErrors, PTimeResults) {
	results := make([]TResult, len(elements))
	errors := make(PErrors, len(elements))
	times := make(PTimeResults, len(elements))

	wg := sync.WaitGroup{}
	wg.Add(len(elements))

	for i, e := range elements {
		go func(_i int, _e TElement) {
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
				time.Sleep(time.Duration(_i) * (time.Second / time.Duration(runPerSec)))
			} else {
				time.Sleep(runPerSecTable[_i])
			}
			defer wg.Done()
			result, err := f(_i, _e)
			if isErrorOrResult {
				results[_i] = result
				errors[_i] = PError{Index: _i, Err: err}
				times[_i] = PTimeResult{Index: _i, StartTime: startTime, EndTime: time.Now()}
			}
		}(i, e)
	}
	wg.Wait()

	return results, errors, times
}
