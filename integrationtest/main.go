package integrationtest

import (
	"flag"
	"net/http"
	"os"
	"path"
	"testing"
)

// Main is the entry point for an integration test
// service and urls are passed directly to StartDockerService and StopDockerService
func Main(m *testing.M, servicefile string, init func(client *http.Client) error, urls ...string) {
	// lock the tests -- only one integrationtest can run at the same time
	LockTests()
	defer UnLockTests()

	// resolve the test data path correctly
	servicefile = path.Join(TestDataPath, servicefile)

	// parse the flags manually to ensure that .Verbose() and .Short() are set
	flag.Parse()

	var code int
	defer func() {
		os.Exit(code)
	}()

	if !testing.Short() {
		// start the docker service
		client, err := StartDockerService(servicefile, urls...)
		if err != nil {
			code = 1
			panic(err)
		}
		defer StopDockerService(servicefile)

		// run init code
		if init != nil {
			if err := init(client); err != nil {
				code = 1
				panic(err)
			}
		}
	}

	code = m.Run()
}
