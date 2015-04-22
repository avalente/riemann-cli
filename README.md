# riemann-cli

Command-line interface for Riemann. Can be used to query the Riemann index or send events from the command line

#### Build

    $ go build
    
#### Usage

    $ ./riemann-cli 
    Please specify command: "send" or "query"

##### Querying

Try 

    $ ./riemann-cli query --help

If you don't provide a `--query` parameter, all the indexed events will be returned (the query sent is just `'true'`)
    
##### Sending events

Try

    $ ./riemann-cli send --help
    
Without arguments, an event with some default properties will be sent (refer to the inline help for details).
You may use the `-v/--verbose` switch to display the event:

    $ ./riemann-cli -v send
    Sending to tcp/localhost:5555
    raidman.Event{Ttl:0, Time:0, Tags:[]string(nil), Host:"localhost", State:"ok", Service:"riemann-cli", Metric:0,         Description:"", Attributes:map[string]string{}}

The event can be read from a `json` file by using the `-j/--json` flag. You can override `json` attributes by passing them on the command line.

Binaries

[Ubuntu 12.04/14.04](https://github.com/avalente/riemann-cli/raw/master/binaries/ubuntu-12_14/riemann-cli)
