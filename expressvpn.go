package expressvpn

import (
	"fmt"
	"math/rand"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	baseCommand = "expressvpn"
	all         = "all"
	list        = "list"
	connect     = "connect"
	disconnect  = "disconnect"
)

// Location describes the alias and location name of a given location
type Location struct {
	Alias    string
	Location string
}

func command(commands ...string) ([]byte, error) {
	cmd := exec.Command(commands[0], commands[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("cannot execute command %s", err.Error())
	}
	return out, nil
}

func formatResponse(rawLocation string) ([]string, error) {
	reg, err := regexp.Compile(`\s+(Smart Location)\s+|\s+(([A-Z][a-z]).+\s\([A-Z]+\))\s+|\s{2,}[Y]|\s{2,}`)
	if err != nil {
		return nil, err
	}
	tabbedString := reg.ReplaceAll([]byte(rawLocation), []byte{'\t'})
	return strings.SplitAfter(string(tabbedString), "\t"), nil
}

// ListAllLocations return a list of all the locations available for use
func ListAllLocations() ([]Location, error) {
	terminalOutput, err := command(baseCommand, list, all)
	if err != nil {
		return nil, fmt.Errorf("%s", terminalOutput)
	}
	locations := strings.SplitAfter(string(terminalOutput), "\n")
	formattedLocations := []Location{}
	for _, rawLocation := range locations[2 : len(locations)-1] {
		formattedLocation, err := formatResponse(rawLocation)
		if err != nil {
			return nil, err
		}
		formattedLocations = append(formattedLocations, Location{Alias: formattedLocation[0], Location: formattedLocation[1]})
	}
	return formattedLocations, nil
}

// ListRecommendedLocations returns a list of the recommended nodes available
func ListRecommendedLocations() ([]Location, error) {
	terminalOutput, err := command(baseCommand, list)
	if err != nil {
		return nil, fmt.Errorf("%s", terminalOutput)
	}
	locations := strings.SplitAfter(string(terminalOutput), "\n")
	formattedLocations := []Location{}
	for _, rawLocation := range locations[3 : len(locations)-3] {
		formattedLocation, err := formatResponse(rawLocation)
		if err != nil {
			return nil, err
		}
		formattedLocations = append(formattedLocations, Location{Alias: formattedLocation[0], Location: formattedLocation[1]})
	}
	return formattedLocations, nil
}

// RandomConnect connects to a location from the recommended location list
func RandomConnect() error {
	locationList, err := ListRecommendedLocations()
	if err != nil {
		return nil
	}
	rand.Seed(time.Now().Unix())
	index := rand.Int63n(int64(len(locationList)))
	location := locationList[index]

	cmd := exec.Command(baseCommand, connect, location.Alias)
	terminalOutput, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", terminalOutput)
	}
	return nil
}

// Connect connects to a given location and returns an error if the connection cannot be made
func Connect(location string) error {
	terminalOutput, err := command(baseCommand, connect, location)
	if err != nil {
		return fmt.Errorf("%s", terminalOutput)
	}
	return nil
}

// Disconnect disconnects from the the connected server
func Disconnect() error {
	terminalOutput, err := command(baseCommand, disconnect)
	if err != nil {
		return fmt.Errorf("%s", terminalOutput)
	}
	return nil
}
