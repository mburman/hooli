A distributed chat client that displays messages based on your proximity to the location of the posts

To run an acceptor do

```go install github.com/mburman/hooli/arunner``` to install
```./bin/arunner -aport=9010```

To run a proposer do (type ```./bin/prunner --help``` for more info

```go install github.com/mburman/hooli/arunner``` to install
```./bin/prunner -pport=9009 -ports=9010```

To install and start a sample client do

```go install github.com/mburman/hooli/client``` to install
```./bin/client --host:port=localhost:9009``` - Change hostport appropriately.
