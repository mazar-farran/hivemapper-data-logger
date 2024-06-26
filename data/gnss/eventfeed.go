package gnss

import (
	"fmt"
	"time"

	"github.com/rosshemsley/kalman"
	"github.com/rosshemsley/kalman/models"
	"github.com/Hivemapper/gnss-controller/device/neom9n"
)

type GnssDataHandler func(data *neom9n.Data) error
type TimeHandler func(now time.Time) error

type GnssFilteredData struct {
	*neom9n.Data
	initialized bool
	lonModel    *models.SimpleModel
	lonFilter   *kalman.KalmanFilter
	latModel    *models.SimpleModel
	latFilter   *kalman.KalmanFilter
}

func NewGnssFilteredData() *GnssFilteredData {
	return &GnssFilteredData{Data: &neom9n.Data{}}
}

func (g *GnssFilteredData) init(d *neom9n.Data) {
	g.initialized = true
	g.lonModel = models.NewSimpleModel(d.Timestamp, 0.0, models.SimpleModelConfig{
		InitialVariance:     0.0,
		ProcessVariance:     2.0,
		ObservationVariance: 2.0,
	})
	g.lonFilter = kalman.NewKalmanFilter(g.lonModel)
	g.latModel = models.NewSimpleModel(d.Timestamp, 0.0, models.SimpleModelConfig{
		InitialVariance:     0.0,
		ProcessVariance:     2.0,
		ObservationVariance: 2.0,
	})
	g.latFilter = kalman.NewKalmanFilter(g.latModel)
	g.Data = d
}

type Option func(*GnssFeed)

type GnssFeed struct {
	dataHandlers     []GnssDataHandler
	timeHandlers     []TimeHandler
	gnssFilteredData *GnssFilteredData

	skipFiltering bool
}

func NewGnssFeed(dataHandlers []GnssDataHandler, timeHandlers []TimeHandler, opts ...Option) *GnssFeed {
	g := &GnssFeed{
		dataHandlers:     dataHandlers,
		timeHandlers:     timeHandlers,
		gnssFilteredData: NewGnssFilteredData(),
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

func WithSkipFiltering() func(*GnssFeed) {
	return func(f *GnssFeed) {
		f.skipFiltering = true
	}
}

func (f *GnssFeed) Run(gnssDevice *neom9n.Neom9n, timeValidThreshold string) error {
	//todo: datafeed is ugly
	dataFeed := neom9n.NewDataFeed(f.HandleData)
	err := gnssDevice.Run(dataFeed, timeValidThreshold, func(now time.Time) {
		dataFeed.SetStartTime(now)
		for _, handler := range f.timeHandlers {
			err := handler(now)
			if err != nil {
				fmt.Printf("handling gnss time: %s\n", err)
			}
		}
	})
	if err != nil {
		return fmt.Errorf("running gnss device: %w", err)
	}

	return nil
}

func (f *GnssFeed) HandleData(d *neom9n.Data) {
	if !f.gnssFilteredData.initialized {
		f.gnssFilteredData.init(d)
	}

	if !f.skipFiltering {
		filteredLon := d.Longitude
		filteredLat := d.Latitude

		err := f.gnssFilteredData.lonFilter.Update(d.Timestamp, f.gnssFilteredData.lonModel.NewMeasurement(d.Longitude))
		if err != nil {
			panic("updating lon filter: " + err.Error())
		}
		err = f.gnssFilteredData.latFilter.Update(d.Timestamp, f.gnssFilteredData.latModel.NewMeasurement(d.Latitude))
		if err != nil {
			panic("updating lat filter: " + err.Error())
		}

		filteredLon = f.gnssFilteredData.lonModel.Value(f.gnssFilteredData.lonFilter.State())
		filteredLat = f.gnssFilteredData.latModel.Value(f.gnssFilteredData.latFilter.State())

		d.Longitude = filteredLon
		d.Latitude = filteredLat
	}

	for _, handler := range f.dataHandlers {
		err := handler(d)
		if err != nil {
			fmt.Printf("handling gnss data: %s\n", err)
		}
	}
}
