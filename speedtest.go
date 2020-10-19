package main

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type SpeedTestResult struct {
	Download float64 `json:"download"`
	Upload   float64 `json:"upload"`
	Ping     float64 `json:"ping"`
	Server   struct {
		URL     string  `json:"url"`
		Lat     string  `json:"lat"`
		Lon     string  `json:"lon"`
		Name    string  `json:"name"`
		Country string  `json:"country"`
		Cc      string  `json:"cc"`
		Sponsor string  `json:"sponsor"`
		ID      string  `json:"id"`
		Host    string  `json:"host"`
		D       float64 `json:"d"`
		Latency float64 `json:"latency"`
	} `json:"server"`
	Timestamp     time.Time   `json:"timestamp"`
	BytesSent     int         `json:"bytes_sent"`
	BytesReceived int         `json:"bytes_received"`
	Share         interface{} `json:"share"`
	Client        struct {
		IP        string `json:"ip"`
		Lat       string `json:"lat"`
		Lon       string `json:"lon"`
		Isp       string `json:"isp"`
		Isprating string `json:"isprating"`
		Rating    string `json:"rating"`
		Ispdlavg  string `json:"ispdlavg"`
		Ispulavg  string `json:"ispulavg"`
		Loggedin  string `json:"loggedin"`
		Country   string `json:"country"`
	} `json:"client"`
}

type SpeedTest struct {
	cmd  string
	args []string

	mrw           sync.Mutex
	metricsResult SpeedTestResult
	doneChan      chan bool
}

func NewSpeedTest(cmd string, args []string) *SpeedTest {
	cmdArgs := make([]string, len(args)+1)
	copy(cmdArgs, args)
	cmdArgs[len(args)] = "--json"

	return &SpeedTest{
		cmd:  cmd,
		args: cmdArgs,
	}
}

// runSpeedTest measures the network bandwidth by doing a speedtest using speedtest-cli
// if multiple calls are made concurrnetly, only one test will run
func (st *SpeedTest) Run() *SpeedTestResult {
	st.mrw.Lock()
	if st.doneChan != nil {
		st.mrw.Unlock()
		<-st.doneChan
		return &st.metricsResult
	}

	st.doneChan = make(chan bool)
	st.mrw.Unlock()

	out, err := exec.Command(*cmd, st.args...).Output()
	if err != nil {
		log.Fatalf("failde to execute speedtest command: %v", err)
	}

	err = json.Unmarshal(out, &st.metricsResult)
	if err != nil {
		log.Fatalf("failed to parse speedtest output: %v", err)
	}

	st.mrw.Lock()
	close(st.doneChan)
	st.doneChan = nil
	st.mrw.Unlock()

	return &st.metricsResult
}

func checkSpeedTestVersion() {
	out, err := exec.Command(*cmd, "--version").Output()
	if err != nil {
		log.Fatalf("failed to verify speedtest version: %v", err)
	}

	ver := strings.Split(string(out), "\n")[0]
	ver = strings.Fields(ver)[1]
	log.Printf("speedtest version %s", ver)
}
