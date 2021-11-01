package main

import (
        "net/http"
//	"fmt"
	"github.com/ziutek/telnet"
	"log"
//	"os"
	"time"
        "strconv"
	"strings"
        "github.com/prometheus/client_golang/prometheus/promhttp"
        "github.com/prometheus/client_golang/prometheus"
)

const timeout = 10 * time.Second

type fooCollector struct {
        indoorTemp *prometheus.Desc
	outdoorTemp *prometheus.Desc
	boilerTemp *prometheus.Desc
	boilergazTemp *prometheus.Desc
	hotwaterTemp *prometheus.Desc
	power *prometheus.Desc
}

func newFooCollector() *fooCollector {
        return &fooCollector{
                indoorTemp: prometheus.NewDesc("indoor_temperature",
                        "Indoor Temperature measured by Heater",
                        nil, nil,
                ),
                outdoorTemp: prometheus.NewDesc("outdoor_temperature",
                        "Outdoor Temperature measured by Heater",
                        nil, nil,
                ),
                boilerTemp: prometheus.NewDesc("boiler_temperature",
                        "Boiler Temperature measured by Heater",
                        nil, nil,
                ),
                boilergazTemp: prometheus.NewDesc("boiler_gaz_temperature",
                        "Boiler Gaz Temperature measured by Heater",
                        nil, nil,
                ),
                hotwaterTemp: prometheus.NewDesc("hot_water_temperature",
                        "Hot Water Temperature measured by Heater",
                        nil, nil,
                ),
                power: prometheus.NewDesc("power",
                        "Power Temperature measured by Heater",
                        nil, nil,
                ),
        }
}

func (collector *fooCollector) Describe(ch chan<- *prometheus.Desc) {
        ch <- collector.indoorTemp
        ch <- collector.outdoorTemp
        ch <- collector.boilerTemp
        ch <- collector.boilergazTemp
        ch <- collector.hotwaterTemp
        ch <- collector.power
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln("Error:", err)
	}
}

func expect(t *telnet.Conn, d ...string) {
	checkErr(t.SetReadDeadline(time.Now().Add(timeout)))
	checkErr(t.SkipUntil(d...))
}

func sendln(t *telnet.Conn, s string) {
	checkErr(t.SetWriteDeadline(time.Now().Add(timeout)))
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = '\n'
	_, err := t.Write(buf)
	checkErr(err)
}

func getIndoorTemp(t *telnet.Conn)(float64,error) {
        t.SetUnixWriteMode(true)
        expect(t,"\n")
        sendln(t,"gvu indoor_temp")
        str, err := t.ReadString('C')
	checkErr(err)
	strbis := strings.Trim(str, " °C")
	res,err :=strconv.ParseFloat(strbis,64)
	return res, err
}

func getOutdoorTemp(t *telnet.Conn)(float64,error) {
        t.SetUnixWriteMode(true)
        expect(t,"\n")
        sendln(t,"gvu outdoor_temp")
        str, err := t.ReadString('C')
        checkErr(err)
        strbis := strings.Trim(str, " °C")
        res,err :=strconv.ParseFloat(strbis,64)
        return res, err
}

func getBoilerTemp(t *telnet.Conn)(float64,error) {
        t.SetUnixWriteMode(true)
        expect(t,"\n")
        sendln(t,"gvu boiler_temp")
        str, err := t.ReadString('C')
        checkErr(err)
        strbis := strings.Trim(str, " °C")
        res,err :=strconv.ParseFloat(strbis,64)
        return res, err
}

func getBoilerGazTemp(t *telnet.Conn)(float64,error) {
        t.SetUnixWriteMode(true)
        expect(t,"\n")
        sendln(t,"gvu boiler_gaz_temp")
        str, err := t.ReadString('C')
        checkErr(err)
        strbis := strings.Trim(str, " °C")
        res,err :=strconv.ParseFloat(strbis,64)
        return res, err
}

func getHotWaterTemp(t *telnet.Conn)(float64,error) {
        t.SetUnixWriteMode(true)
        expect(t,"\n")
        sendln(t,"gvu hot_water_temp")
        str, err := t.ReadString('C')
        checkErr(err)
        strbis := strings.Trim(str, " °C")
        res,err :=strconv.ParseFloat(strbis,64)
        return res, err
}

func getPower(t *telnet.Conn)(float64,error) {
        t.SetUnixWriteMode(true)
        expect(t,"\n")
        sendln(t,"gvu power")
        str, err := t.ReadString('%')
        checkErr(err)
        strbis := strings.Trim(str, " %")
        res,err :=strconv.ParseFloat(strbis,64)
        return res, err
}

func (collector *fooCollector) Collect(ch chan<- prometheus.Metric) {
        t, err := telnet.Dial("tcp", "127.0.0.1:3083")
        checkErr(err)
	a,err := getIndoorTemp(t)
        checkErr(err)
        b,err := getOutdoorTemp(t)
        checkErr(err)
        c,err := getBoilerTemp(t)
        checkErr(err)
        d,err := getBoilerGazTemp(t)
        checkErr(err)
        e,err := getHotWaterTemp(t)
        checkErr(err)
        f,err := getPower(t)
        checkErr(err)
        m1 := prometheus.MustNewConstMetric(collector.indoorTemp, prometheus.GaugeValue, a)
        m1 = prometheus.NewMetricWithTimestamp(time.Now(), m1)
        ch <- m1
        m2 := prometheus.MustNewConstMetric(collector.outdoorTemp, prometheus.GaugeValue, b)
        m2 = prometheus.NewMetricWithTimestamp(time.Now(), m2)
        ch <- m2
        m3 := prometheus.MustNewConstMetric(collector.boilerTemp, prometheus.GaugeValue, c)
        m3 = prometheus.NewMetricWithTimestamp(time.Now(), m3)
        ch <- m3
        m4 := prometheus.MustNewConstMetric(collector.boilergazTemp, prometheus.GaugeValue, d)
        m4 = prometheus.NewMetricWithTimestamp(time.Now(), m4)
        ch <- m4
        m5 := prometheus.MustNewConstMetric(collector.hotwaterTemp, prometheus.GaugeValue, e)
        m5 = prometheus.NewMetricWithTimestamp(time.Now(), m5)
        ch <- m5
        m6 := prometheus.MustNewConstMetric(collector.power, prometheus.GaugeValue, f)
        m6 = prometheus.NewMetricWithTimestamp(time.Now(), m6)
        ch <- m6
}


func main() {
        foo := newFooCollector()
        prometheus.MustRegister(foo)
        http.Handle("/metrics", promhttp.Handler())
        log.Fatal(http.ListenAndServe(":9101", nil))
}
