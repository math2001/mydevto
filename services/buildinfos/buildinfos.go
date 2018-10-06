package buildinfos

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// V is the version that will be set at build time, see the Makefile
var V string

// Major is the major version number
var Major int

// Minor is the minor version number
var Minor int

// Patch is the patch version number
var Patch int

// Testing is true if this is the testing version
var Testing bool

// we can affort to Fatal because this is done as soon as the program is run,
// not during the first request
func parseversion(v string) (major, minor, patch int, testing bool) {
	// remove the leading 'v'
	str := V[1:]
	testing = strings.HasSuffix(str, "-testing")

	all := strings.Split(str, "-")
	version := all[0]
	bits := strings.Split(version, ".")
	var err error
	major, err = strconv.Atoi(bits[0])
	if err != nil {
		log.Fatal(err)
	}
	minor, err = strconv.Atoi(bits[1])
	if err != nil {
		log.Fatal(err)
	}
	patch, err = strconv.Atoi(bits[2])
	if err != nil {
		log.Fatal(err)
	}
	return
}

func init() {
	if V == "" {
		log.Print("No build information. Use make to build")
		log.Print("$ go clean")
		log.Print("$ make")
		os.Exit(1)
	}
	Major, Minor, Patch, Testing = parseversion(V)
}
