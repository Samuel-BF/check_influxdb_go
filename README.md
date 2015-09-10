check_influxdb_go
=================

Nagios/icinga/... plugin to query value from influxdb and check result. Usage:

	./check_influxdb_go --help
	Usage of /home/ruslan/go/bin/check_influxdb_go:
	  -H string
		Host to connect to (default "localhost")
	  -P int
		Port to connect to (default 8086)
	  -c string
		Critical range
	  -d string
		Database (default "metrics")
	  -p string
		Influxdb user password
	  -q string
		Database query
	  -r int
		Timeout (milliseconds) (default 1000)
	  -u string
		Influxdb user
	  -w string
		Warning range

Tested on influxdb 0.9.x.


Build
=====

Make shure that you have working go environment. Get check_influxdb_go and dependencies:

	go get github.com/neezgee/check_influxdb_go
	go get github.com/influxdb/influxdb/client
	go get github.com/olorin/nagiosplugin

Build and install it with:

	go install github.com/neezgee/check_influxdb_go
