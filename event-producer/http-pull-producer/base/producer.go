package base

import (
	"com.fs/event-service/utils"
	"errors"
	"time"
)

type Producer struct {
	Config *Config

	Pull func() error

	stopChan chan bool
}

func (prod *Producer) Start() error {
	if prod.Config == nil {
		utils.PrintErr("Producer.Start", "没有传递配置")
		return errors.New("没有传递配置")
	}

	if prod.Pull == nil {
		utils.PrintErr("Producer.Start", "没有创建OnPull")
		return errors.New("没有创建OnPull")
	}

	err := prod.pull()
	if err != nil {
		utils.PrintCallErr("Producer.Start", "prod.pull", err)
		return err
	}

	prod.stopChan = make(chan bool)

	ticker := time.NewTicker(time.Duration(prod.Config.PullPeriodSec) * time.Second)
	go func(p *Producer, ticker *time.Ticker) {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_ = p.pull()
			case stop := <-p.stopChan:
				if stop {
					p.stopChan <- true
					return
				}
			}
		}
	}(prod, ticker)

	return nil
}

func (prod *Producer) Stop() error {
	prod.stopChan <- true
	<-prod.stopChan

	return nil
}

func (prod *Producer) pull() error {
	return prod.Pull()
}
