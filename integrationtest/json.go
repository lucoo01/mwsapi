package integrationtest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

// normalizeJSON normalizes a json string
func normalizeJSON(in string) (out string, err error) {
	var t interface{}

	// Unmarshal into a generic interface
	err = json.Unmarshal([]byte(in), &t)
	if err != nil {
		return
	}

	// Remarshal it
	outB, err := json.Marshal(&t)
	if err != nil {
		return
	}

	// and return
	out = string(outB)
	return
}

// Outputs json into {$asset}.out for debugging purposes
// returns the filename
func outputDebugJSON(t *testing.T, res interface{}, asset string) (filename string, err error) {
	filename = fmt.Sprintf("%s.out", asset)
	_, err = writeJSONFile(t, res, filename)
	return
}

// writeJSONFile writes a json version of res into filename
func writeJSONFile(t *testing.T, res interface{}, filename string) (bytes []byte, err error) {
	// Remarshal it
	bytes, err = json.MarshalIndent(res, "", "  ")
	if err != nil {
		return
	}

	// write the file
	return bytes, ioutil.WriteFile(filename, bytes, 0755)
}
