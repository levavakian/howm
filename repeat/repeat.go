package repeat

import (
	"time"
)

type RepeatFunc func()

type Repeater struct {
	Func RepeatFunc
	DoneChan chan bool
	ConfirmChan chan bool
	Ticker *time.Ticker
}

func NewRepeater(f RepeatFunc, d time.Duration, workRequest chan struct{}{}, workNotify chan struct{}{}) *Repeater {
	ticker := time.NewTicker(d)
	done := make(chan bool)
	confirm := make(chan bool)
	r := &Repeater{
		DoneChan: done,
		ConfirmChan: confirm,
		Ticker: ticker,
	}
	go func(){
		for {
			select {
			case <-r.DoneChan:
				return
			case <-r.Ticker.C:
				if workRequest == nil || workNotify == nil {
					f()
				} else {
					slect {
					case <- r.DoneChan:
						return
					case workRequest <- struct{}{}:
						f()
						workNotify <- struct{}{}
					}
				}
			}
		}
		r.ConfirmChan <- true
	}()
	return r
}

func (r *Repeater) Stop() {
	r.Ticker.Stop()
	r.DoneChan <- true
	<- r.ConfirmChan
}