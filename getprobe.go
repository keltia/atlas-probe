package main

import (
	"flag"
	"time"

	"fmt"
	"os"
	"path/filepath"
	"strings"
	"io/ioutil"
	"errors"

	"github.com/naoina/toml"
	"log"
	"github.com/bndr/gopencils"
)


const (
	TAG = "ripe-atlas"
	API = "http://atlas.ripe.net/api/v1"
)

var (
	fId string
	config Config
)

type Probe struct {
	// visible data
	Address_v4 string
	Address_v6 string
	Country_code string
	Is_anchor bool
	Is_public bool
	Latitude float32
	Longitude float32
	Status int32
	Status_name string
	Status_since time.Time

	// Filters
	Asn_v4 int32
	Asn_v6 int32
	Id int32 // Filter as Id or Id_In
	Prefix_v4 string
	Prefix_v6 string

	// hidden data
	Asn int32 // Can be filtered
	Location string
	Resource_uri string
}

type Metadata struct {
	Limit int
	Next string
	Offset int
	Previous string
	TotalCount int
	UseIsoTime bool
}

type RestAnswer struct {
	Objets []Probe
	Meta Metadata
}

type Config struct {
	Id string
}

// Load a file as a YAML document and return the structure
func LoadConfig(file string) (*Config, error) {
	var sFile string

	// Check for tag
	if !strings.HasSuffix(file, ".toml") {
		// file must be a tag so add a "."
		sFile = filepath.Join(os.Getenv("HOME"),
			fmt.Sprintf(".%s", file),
			"config.toml")
	} else {
		sFile = file
	}

	c := new(Config)
	buf, err := ioutil.ReadFile(sFile)
	if err != nil {
		return c, errors.New(fmt.Sprintf("Can not read %s file.", sFile))
	}

	err = toml.Unmarshal(buf, &c)
	if err != nil {
		return c, errors.New(fmt.Sprintf("Can not parse %s: %v", sFile, err))
	}

	return c, err
}

func init() {
	flag.StringVar(&fId, "i", config.Id, "Probe Id")
}

func main() {

	// Load defaultd if any
	config, err := LoadConfig(TAG)
	if err != nil {
		log.Printf("Warning: missing or unreadable config.toml as %s\n", TAG)
	}

	flag.Parse()

	// Check if we have anything on the command-line
	if fId == "" {
		if err != nil {
			log.Fatalf("Error: no default Id and nothing on command-line")
		}
		fId = config.Id
	}

	api := gopencils.Api(API)

	query := map[string]string{"id": fId}
	answer := new(RestAnswer)

	api.Res("probe/", answer).Get(query)

	fmt.Printf("Result:\n  %v", answer)
}
