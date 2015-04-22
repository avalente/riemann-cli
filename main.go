package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/amir/raidman"
	"gopkg.in/alecthomas/kingpin.v1"
)

func connect(protocol string, address string) *raidman.Client {
	conn, err := raidman.Dial(protocol, address)
	if err != nil {
		fmt.Printf("Can't connect to riemann: %v\n", err)
		os.Exit(1)
	}

	return conn
}

func main() {
	hostName, _ := os.Hostname()

	var (
		cli = kingpin.New("riemann-cli", "Command-line interface for Riemann")

		protocol = cli.Flag("protocol", "Network protocol").Short('p').Default("tcp").String()
		address  = cli.Flag("address", "Server address").Short('h').Default("localhost:5555").String()
		verbose  = cli.Flag("verbose", "Verbose").Short('v').Bool()

		send        = cli.Command("send", "Send an event to Riemann")
		host        = send.Flag("host", "Event Host").Short('h').Default(hostName).String()
		service     = send.Flag("riemann-cli", "Event service").Short('s').Default("riemann-cli").String()
		ttl         = send.Flag("ttl", "Event TTL").Float()
		description = send.Flag("description", "Event description").Short('d').String()
		time        = send.Flag("time", "Event timestamp").Int()
		tags        = send.Flag("tag", "Event tag (can be specified multiple times)").Short('t').Strings()
		state       = send.Flag("state", "Event state").Default("ok").String()
		metric      = send.Flag("metric", "Event metric").Short('m').Float()
		attributes  = send.Flag("attribute", "Event attributes").Short('a').StringMap()
		jsonFile    = send.Flag("json", "File with a JSON representation of the event").Short('j').ExistingFile()

		query    = cli.Command("query", "Query Riemann index")
		querystr = query.Flag("query", "Riemann Query").Default("true").String()
		jsonFmt  = query.Flag("json", "Output JSON").Short('j').Bool()
	)

	switch kingpin.MustParse(cli.Parse(os.Args[1:])) {

	case send.FullCommand():
		conn := connect(*protocol, *address)

		ev := raidman.Event{}

		if *jsonFile != "" {
			file, err := os.Open(*jsonFile)
			if err != nil {
				fmt.Printf("Can't open %v: %v\n", *jsonFile, err)
				os.Exit(3)
			}
			decoder := json.NewDecoder(file)
			err = decoder.Decode(&ev)
			if err != nil {
				fmt.Printf("Can't decode %v: %v\n", *jsonFile, err)
			}
		}

		if ev.Ttl == 0 {
			ev.Ttl = float32(*ttl)
		}

		if ev.Time == 0 {
			ev.Time = int64(*time)
		}

		if ev.Host == "" {
			ev.Host = *host
		}

		if ev.Service == "" {
			ev.Service = *service
		}

		if ev.Description == "" {
			ev.Description = *description
		}

		if len(ev.Tags) == 0 {
			ev.Tags = *tags
		}

		if ev.State == "" {
			ev.State = *state
		}

		if ev.Metric == nil || ev.Metric.(float64) == 0.0 {
			ev.Metric = *metric
		}

		if len(ev.Attributes) == 0 {
			ev.Attributes = *attributes
		}

		if *verbose {
			fmt.Printf("Sending to %v/%v\n%#v\n", *protocol, *address, ev)
		}

		err := conn.Send(&ev)
		if err != nil {
			fmt.Printf("Can't send event to riemann: %v\n", err)
			os.Exit(2)
		}

	case query.FullCommand():
		conn := connect(*protocol, *address)

		res, err := conn.Query(*querystr)
		if err != nil {
			fmt.Printf("Can't query riemann: %v\n", err)
		}

		if *jsonFmt {
			out, _ := json.Marshal(res)
			fmt.Printf("%s\n", out)
		} else {
			for _, ev := range res {
				fmt.Printf("%#v\n", ev)
			}
		}

	default:
		fmt.Printf("Please specify command: \"send\" or \"query\"\n")
		os.Exit(0)
	}
}
