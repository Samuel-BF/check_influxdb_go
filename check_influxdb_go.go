// vim: sts=4 sw=4 et

package main

import (
      "github.com/influxdb/influxdb/client"
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
var db = flag.String("d", "metrics", "Database")
var user = flag.String("u", "", "Influxdb user")
var password = flag.String("p", "", "Influxdb user password")
var query = flag.String("q", "", "Database query")

var warning = flag.String("w", "", "Warning range")
var critical = flag.String("c", "", "Critical range")
var timeout = flag.Int("r", 1000, "Timeout (milliseconds)")

func main() {

    flag.Parse()

    u, err := url.Parse(fmt.Sprintf("http://%s:%d", *host, *port))
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
        if err != nil {
            check.AddResult(nagiosplugin.UNKNOWN, "error parsing result")
        }

        check.AddPerfDatum("value", "", result)

        if *warning != "" {
            warnRange, err := nagiosplugin.ParseRange(*warning)
            if err != nil {
                check.AddResult(nagiosplugin.UNKNOWN, "error parsing warning range")
            }
            if warnRange.Check(result) {
                check.AddResult(nagiosplugin.WARNING, "Result in warning range")
            }
        } 

        if *critical != "" {
            criticalRange, err := nagiosplugin.ParseRange(*critical)
            if err != nil {
                check.AddResult(nagiosplugin.UNKNOWN, "error parsing critical range")
            }

            if criticalRange.Check(result) {
                check.AddResult(nagiosplugin.WARNING, "Result in critical range")
            }
        }
        check.AddResult(nagiosplugin.OK, "")
    }
}
