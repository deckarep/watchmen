package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/stathat/go"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	//"strconv"
	"strings"
	"syscall"
	"time"
)

// Matches the Output critical string
var cfg *config
var outputRegexp = regexp.MustCompile(`CRITICAL:\s(\d+)`)

func init() {
	loadConfig()
}

func main() {

	// Kick off uchiwa poller
	go pollEvents()

	// Poor-man wait: Just block until user presses enter
	log.Println(color.MagentaString(fmt.Sprintf("Started: Waiting for Uchiwa events to process with %d second interval.", cfg.Uchiwa.Interval)))

	blockUntilSignal()
}

func loadConfig() {
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(color.RedString("Couldn't load config.json data."))
	}

	c := &config{}
	err = json.Unmarshal(b, c)
	if err != nil {
		log.Fatal(color.RedString("Couldn't unmarshal config.json data."))
	}

	cfg = c
	log.Println(color.MagentaString("Loaded config.json file from disk."))
}

func blockUntilSignal() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		done <- true
	}()

	<-done
	log.Println(color.YellowString("Received signal to quit."))
}

func pollEvents() {
	errorCount := 0
	for {

		// Use this block for polling of uchiwa
		events, err := fetchEvents(cfg.Alerts)
		if err != nil {
			log.Println(color.YellowString("Could not pollEvents from Uchiwa"))
			errorCount++
		} else {
			errorCount = 0
		}

		eventCount := len(events)
		if eventCount > 0 {
			log.Println(color.MagentaString(fmt.Sprintf("Found %d Uchiwa events for processing...", eventCount)))
		} else {
			log.Println(color.MagentaString("Found 0 Uchiwa events, nothing to process currently."))
		}

		for i, ev := range events {
			//TODO: send to stathat
			//TODO: may have to aggregate stats manually to see per datacenter (not sure stathat can do that for us)

			key := ev.Check.Name + "-" + ev.Client.Name

			// Only sending 5 keys for now since we're on a demo plan
			if i < 5 {
				stathat.PostEZCount(key, cfg.StatsAccount, 1)
				fmt.Println("stathat -> " + key)
			}

			// results := outputRegexp.FindStringSubmatch(ev.Check.Output)
			// if len(results) == 2 {
			// 	if staleCount, err := strconv.Atoi(strings.TrimSpace(results[1])); err == nil {
			// 		log.Println(color.MagentaString(fmt.Sprintf("Host %s has %d stale files > 1 hour old.", ev.Client.Name, staleCount)))
			// 	}
			// }
		}

		if errorCount > cfg.Uchiwa.MaxErrors {
			log.Fatalln(color.MagentaString("Fatal: Reached max consecutive allowable errors with Uchiwa."))
		}

		time.Sleep(time.Second * time.Duration(cfg.Uchiwa.Interval))
	}
}

func fetchEvents(filters []string) ([]*Event, error) {
	resp, err := http.Get(cfg.Uchiwa.Host)
	if err != nil {
		return nil, fmt.Errorf("Couldn't GET unchiwa events data with error: %s", err.Error())
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read http body for uchiwa /events with error: %s", err.Error())
	}

	events := make([]*Event, 0)

	if err := json.Unmarshal(bytes, &events); err != nil {
		return nil, fmt.Errorf("Could unmarshal body for uchiwa /events with error: %s", err.Error())
	}

	// Loop through all events returned and only filter out the subset of alerts we're interested in
	filterdEvents := make([]*Event, 0)

	for _, ev := range events {
		for _, filter := range filters {
			if strings.Contains(ev.Check.Name, filter) {
				filterdEvents = append(filterdEvents, ev)
			}
		}
	}

	return filterdEvents, nil
}
