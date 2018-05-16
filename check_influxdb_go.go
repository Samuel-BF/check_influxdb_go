// vim: sts=4 sw=4 et

/*

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.

*/

package main

import (
      "github.com/influxdata/influxdb/client"
      "github.com/olorin/nagiosplugin"
      "net/url"
      "fmt"
      "log"
      "time"
      "flag"
      "encoding/json"
)

var host = flag.String("H", "localhost", "Host to connect to")
var port = flag.Int("P", 8086, "Port to connect to")
var ssl = flag.Bool("ssl", false, "Use SSL/TLS for connection")
var db = flag.String("d", "metrics", "Database")
var user = flag.String("u", "", "Influxdb user")
var password = flag.String("p", "", "Influxdb user password")
var query = flag.String("q", "", "Database query")

var warning = flag.String("w", "", "Warning range")
var critical = flag.String("c", "", "Critical range")
var timeout = flag.Int("r", 1000, "Timeout (milliseconds)")

func main() {

    flag.Parse()

    start := time.Now()

    var protocol = "http"
    if *ssl {
        protocol = "https"
    }
    u, err := url.Parse(fmt.Sprintf("%s://%s:%d", protocol, *host, *port))
    if err != nil {
        log.Fatal(err)
    }

    conf := client.Config{
        URL:      *u,
        Username: *user,
        Password: *password,
        Timeout:  time.Duration(*timeout) * time.Millisecond,
    }

    check := nagiosplugin.NewCheck()
    defer check.Finish()

    con, err := client.NewClient(conf)
    if err != nil {
        check.AddResult(nagiosplugin.UNKNOWN, "Can't connect to database")
        log.Fatal(err)
    }

    q := client.Query{
        Command: *query,
        Database: *db,
    }

    if response, err := con.Query(q); err == nil {
        if response.Error() != nil {
            check.AddResult(nagiosplugin.UNKNOWN, "Can't execute query")
            log.Fatal(err)
        }

        result, err := (response.Results[0].Series[0].Values[0][1].(json.Number)).Float64()

        duration := time.Since(start)

        if err != nil {
            check.AddResult(nagiosplugin.UNKNOWN, "Error parsing result")
        }

        message := fmt.Sprintf("Got %v in %v", result, duration)

        check.AddPerfDatum("value", "", result)

        if *warning != "" {
            warnRange, err := nagiosplugin.ParseRange(*warning)
            if err != nil {
                check.AddResult(nagiosplugin.UNKNOWN, "Error parsing warning range")
            }
            if warnRange.Check(result) {
                check.AddResult(nagiosplugin.WARNING, message)
            }
        }

        if *critical != "" {
            criticalRange, err := nagiosplugin.ParseRange(*critical)
            if err != nil {
                check.AddResult(nagiosplugin.UNKNOWN, "Error parsing critical range")
            }

            if criticalRange.Check(result) {
                check.AddResult(nagiosplugin.CRITICAL, message)
            }
        }
        check.AddResult(nagiosplugin.OK, message)
    }
}
