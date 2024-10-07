package jtpslimiter

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var limiterMap = sync.Map{}

// TpsLimit : TPS 제한
// key : 키
// tps : 1초에 몇번 요청을 받을 것인지 (20이면 1초에 20번 요청을 받을 수 있음)
// timeoutDuration : 타임아웃 처리할 시간
// return : error - 대기 완료, false - 요청 수가 종료 기준 수를 넘음
func TpsLimit(key string, tps int, timeoutDuration time.Duration) error {

	limiterAny, _ := limiterMap.LoadOrStore(key, rate.NewLimiter(rate.Limit(tps), tps))
	limiter := limiterAny.(*rate.Limiter)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	return limiter.Wait(ctx)
}
