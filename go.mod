module github.com/streamingfast/hivemapper-data-logger

go 1.20

replace github.com/daedaleanai/ublox => github.com/streamingfast/ublox v0.0.0-20230612141255-f8202e5f6890

//replace github.com/streamingfast/gnss-controller => ../gnss-controller
//replace github.com/streamingfast/imu-controller => ../imu-controller

require (
	github.com/bufbuild/connect-go v1.8.0
	github.com/google/uuid v1.3.0
	github.com/paulmach/go.geojson v1.4.0
	github.com/rosshemsley/kalman v0.0.0-20190615074247-f4b900823fd1
	github.com/rs/cors v1.9.0
	github.com/spf13/cobra v1.7.0
	github.com/streamingfast/gnss-controller v0.1.20-0.20230621182813-f4940c0a6996
	github.com/streamingfast/imu-controller v0.0.0-20230626153723-274a49c06e63
	github.com/stretchr/testify v1.8.4
	golang.org/x/net v0.5.0
	google.golang.org/protobuf v1.30.0
	modernc.org/sqlite v1.22.1
)

require (
	github.com/daedaleanai/ublox v0.0.0-00010101000000-000000000000 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/streamingfast/shutter v1.5.0 // indirect
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	golang.org/x/tools v0.1.12 // indirect
	gonum.org/v1/gonum v0.0.0-20190606121551-14af50e936aa // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	lukechampine.com/uint128 v1.2.0 // indirect
	modernc.org/cc/v3 v3.40.0 // indirect
	modernc.org/ccgo/v3 v3.16.13 // indirect
	modernc.org/libc v1.22.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.5.0 // indirect
	modernc.org/opt v0.1.3 // indirect
	modernc.org/strutil v1.1.3 // indirect
	modernc.org/token v1.0.1 // indirect
	periph.io/x/conn/v3 v3.7.0 // indirect
	periph.io/x/host/v3 v3.8.2 // indirect
)
