package sideloop

import (
	"time"
)

type NoArgFunc func()

type Injector struct {
	WorkRequest chan struct{}
	WorkNotify  chan struct{}
}

func NewInjector() *Injector {
	return &Injector{
		WorkRequest: make(chan struct{}),
		WorkNotify:  make(chan struct{}),
	}
}

func (inj *Injector) Do(f NoArgFunc) {
	inj.WorkRequest <- struct{}{}
	f()
	inj.WorkNotify <- struct{}{}
}

type Repeater struct {
	DoneChan    chan bool
	ConfirmChan chan bool
	Ticker      *time.Ticker
}

func NewRepeater(f NoArgFunc, d time.Duration, inj *Injector) *Repeater {
	ticker := time.NewTicker(d)
	done := make(chan bool)
	confirm := make(chan bool)
	r := &Repeater{
		DoneChan:    done,
		ConfirmChan: confirm,
		Ticker:      ticker,
	}
	go func() {
		for {
			select {
			case <-r.DoneChan:
				return
			case <-r.Ticker.C:
				if inj == nil {
					f()
				} else {
					select {
					case <-r.DoneChan:
						return
					case inj.WorkRequest <- struct{}{}:
						f()
						inj.WorkNotify <- struct{}{}
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
	<-r.ConfirmChan
}
