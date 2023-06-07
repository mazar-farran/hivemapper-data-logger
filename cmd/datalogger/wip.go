package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/streamingfast/gnss-controller/device/neom9n"
	"github.com/streamingfast/hivemapper-data-logger/data"
	"github.com/streamingfast/hivemapper-data-logger/data/gnss"
	"github.com/streamingfast/hivemapper-data-logger/data/imu"
	"github.com/streamingfast/hivemapper-data-logger/logger"
	"github.com/streamingfast/imu-controller/device/iim42652"
)

var WipCmd = &cobra.Command{
	Use:   "wip",
	Short: "Start the data logger",
	RunE:  wipRun,
}

func init() {
	// Imu
	WipCmd.Flags().String("imu-config-file", "imu-logger.json", "Imu logger config file. Default path is ./imu-logger.json")

	// GNSS
	WipCmd.Flags().String("gnss-config-file", "gnss-logger.json", "Neom9n logger config file. Default path is ./gnss-logger.json")
	WipCmd.Flags().String("gnss-json-destination-folder", "/mnt/data/gps", "json destination folder")
	WipCmd.Flags().Duration("gnss-json-save-interval", 15*time.Second, "json save interval")
	WipCmd.Flags().Int64("gnss-json-destination-folder-max-size", int64(30000*1024), "json destination folder maximum size") // 30MB
	WipCmd.Flags().String("gnss-serial-config-name", "/dev/ttyAMA1", "Config serial location")
	WipCmd.Flags().String("gnss-mga-offline-file-path", "/mnt/data/mgaoffline.ubx", "path to mga offline files")
	WipCmd.Flags().String("gnss-db-path", "/mnt/data/gnss.v1.0.3.db", "path to sqliteLogger database")
	WipCmd.Flags().Duration("gnss-db-log-ttl", 12*time.Hour, "ttl of logs in database")

	RootCmd.AddCommand(WipCmd)
}

func wipRun(cmd *cobra.Command, args []string) error {
	imuDevice := iim42652.NewSpi("/dev/spidev0.0", iim42652.AccelerationSensitivityG16, iim42652.GyroScalesG2000, true)
	err := imuDevice.Init()
	if err != nil {
		return fmt.Errorf("initializing IMU: %w", err)
	}
	err = imuDevice.UpdateRegister(iim42652.RegisterAccelConfig, func(currentValue byte) byte {
		return currentValue | 0x01
	})
	if err != nil {
		return fmt.Errorf("failed to update register: %w", err)
	}
	conf := imu.LoadConfig(mustGetString(cmd, "imu-config-file"))
	fmt.Println("Config: ", conf.String())

	rawImuEventFeed := imu.NewRawFeed(imuDevice)
	rawImuEventFeed.Start()

	correctedImuEventFeed := imu.NewCorrectedAccelerationFeed()
	correctedImuEventFeed.Start(rawImuEventFeed)

	serialConfigName := mustGetString(cmd, "gnss-serial-config-name")
	mgaOfflineFilePath := mustGetString(cmd, "gnss-mga-offline-file-path")
	gnssDevice := neom9n.NewNeom9n(serialConfigName, mgaOfflineFilePath)
	err = gnssDevice.Init(nil)
	if err != nil {
		return fmt.Errorf("initializing neom9n: %w", err)
	}

	gnssEventFeed := gnss.NewEventFeed()
	gnssEventFeed.Start(gnssDevice)

	jsonDestinationFolder := mustGetString(cmd, "gnss-json-destination-folder")
	jsonSaveInterval := mustGetDuration(cmd, "gnss-json-save-interval")
	jsonDestinationFolderMaxSize := mustGetInt64(cmd, "gnss-json-destination-folder-max-size")

	jsonLogger := logger.NewJsonFile(jsonDestinationFolder, jsonDestinationFolderMaxSize, jsonSaveInterval)
	gnssFileLoggerSubscription := gnssEventFeed.Subscribe("gnssFileLoggerSubscription")
	err = jsonLogger.Init(gnssFileLoggerSubscription)
	if err != nil {
		return fmt.Errorf("initializing json logger database: %w", err)
	}

	gnssEventSub := gnssEventFeed.Subscribe("merger")
	rawEventSub := rawImuEventFeed.Subscribe("merger")
	correctedImuEventSub := correctedImuEventFeed.Subscribe("merger")
	mergedEventFeed := data.NewEventFeedMerger(gnssEventSub, rawEventSub, correctedImuEventSub)
	mergedEventFeed.Start()

	mergedEventSub := mergedEventFeed.Subscribe("wip")

	fmt.Println("Starting to listen for events from mergedEventSub")
	var imuRawEvent *imu.RawImuEvent
	var correctedImuEvent *imu.CorrectedAccelerationEvent
	var gnssEvent *gnss.GnssEvent
	//fmt.Fprintf(os.Stderr, "Raw Acc X, Corrected Acc X, Tilt X, Raw Acc Y, Corrected Acc Y, Tilt Y, Raw Acc Z\n")
	for {
		select {
		case e := <-mergedEventSub.IncomingEvents:
			switch e := e.(type) {
			case *imu.RawImuEvent:
				imuRawEvent = e
			case *imu.CorrectedAccelerationEvent:
				correctedImuEvent = e
			case *gnss.GnssEvent:
				gnssEvent = e
			}
		}
		if imuRawEvent != nil && correctedImuEvent != nil {
			_ = gnssEvent
			_ = imuRawEvent
			_ = correctedImuEvent
			//printCVS(imuRawEvent, correctedImuEvent, gnssEvent)
			imuRawEvent = nil
			correctedImuEvent = nil
		}
	}
}

func printCVS(imuRawEvent *imu.RawImuEvent, correctedImuEvent *imu.CorrectedAccelerationEvent, gnss *gnss.GnssEvent) {
	fmt.Fprintf(os.Stderr, "%f,%f,%f,%f,%f,%f,%f\n",
		imuRawEvent.Acceleration.CamX(),
		correctedImuEvent.X,
		correctedImuEvent.XAngle,
		imuRawEvent.Acceleration.CamY(),
		correctedImuEvent.Y,
		correctedImuEvent.YAngle,
		imuRawEvent.Acceleration.CamZ(),
		gnss.Data.Fix,
		gnss.Data.Latitude,
		gnss.Data.Longitude,
		gnss.Data.Altitude,
		gnss.Data.Speed,
		gnss.Data.Heading,
		gnss.Data.Satellites.Seen,
		gnss.Data.Satellites.Used,
		gnss.Data.Eph,
	)
}
