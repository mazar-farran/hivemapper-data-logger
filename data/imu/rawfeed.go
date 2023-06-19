package imu

import (
	"fmt"
	"time"

	"github.com/streamingfast/imu-controller/device/iim42652"
)

type RawFeed struct {
	imu      *iim42652.IIM42652
	handlers []RawFeedHandler
}

func NewRawFeed(imu *iim42652.IIM42652, handlers ...RawFeedHandler) *RawFeed {
	return &RawFeed{
		imu:      imu,
		handlers: handlers,
	}
}

type RawFeedHandler func(acceleration *Acceleration, angularRate *iim42652.AngularRate) error

func (f *RawFeed) Run() error {
	fmt.Println("Run imu raw feed")
	for {
		time.Sleep(10 * time.Millisecond)
		acceleration, err := f.imu.GetAcceleration()
		if err != nil {
			return fmt.Errorf("getting acceleration: %w", err)
		}
		angularRate, err := f.imu.GetGyroscopeData()
		if err != nil {
			return fmt.Errorf("getting angular rate: %w", err)
		}
		temperature, err := f.imu.GetTemperature()
		if err != nil {
			return fmt.Errorf("getting temperature: %w", err)
		}

		for _, handler := range f.handlers {
			err := handler(
				NewAcceleration(acceleration.CamX(), acceleration.CamY(), acceleration.CamZ(), acceleration.TotalMagnitude, *temperature, time.Now()),
				angularRate,
			)
			if err != nil {
				return fmt.Errorf("calling handler: %w", err)
			}
		}
	}

}
