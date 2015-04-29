// This package implements the Mixpanel API as referenced here: https://mixpanel.com/help/reference/http
package mixpanel

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const BASE_URL = "https://api.mixpanel.com"

var (
	// This error is returned when Mixpanel returns a non-success message when tracking an event
	ErrUnexpectedTrackResponse = fmt.Errorf("Unexpected Mixpanel Track Response")
	// This error is returned when Mixpanel returns a non-success message when using an engage event
	ErrUnexpectedEngageResponse = fmt.Errorf("Unexpected Mixpanel Engage Response")
)

type Mixpanel struct {
	Token             string
	BaseURL           string
	OverrideIPAddress string
}

// NewMixpanelClient returns a Mixpanel struct with which you can perform other Mixpanel operations
// e.g. `m := mixpanel.NewMixpanelClient("your_mixpanel_token")`
func NewMixpanelClient(args ...string) *Mixpanel {
	var m *Mixpanel

	if len(args) == 1 {
		m = &Mixpanel{Token: args[0], BaseURL: BASE_URL}
	} else if len(args) > 1 {
		m = &Mixpanel{Token: args[0], BaseURL: args[1]}
	}

	return m
}

// Track creates a Mixpanel event for the "event" string along with other properties
// that are added to the event as meta-data
// e.g. `err := mc.Track("User Signed Up", map[string]interface{}{"$distinct_id": "1"})`
func (m *Mixpanel) Track(event string, properties map[string]interface{}) error {
	var data map[string]interface{} = make(map[string]interface{})

	data["event"] = event
	properties["token"] = m.Token
	data["properties"] = properties

	response, err := m.get(fmt.Sprintf("%s/track/", m.BaseURL), data)
	if err != nil {
		return err
	}

	if response != "1" {
		return ErrUnexpectedTrackResponse
	}

	return nil
}

// ProfileSet creates a "People" profile in Mixpanel with a distinctID (which is the primary key)
// along with properties that are added as meta-data to the profile
// e.g. `err := m.ProfileSet("1", map[string]interface{}{"full_name": "Mclovin", "Company": "Acme Organ Donation"})`
func (m *Mixpanel) ProfileSet(distinctID string, properties map[string]interface{}) error {
	return m.engage(distinctID, "$set", properties)
}

// ProfileSetOnce sets properties that are not already set in the profile
// that is referenced by the distinctID (which is the primary key)
// ip is optional
// e.g. `err := m.ProfileSetOnce("1", map[string]interface{}{"full_name": "Mclovin", "Company": "Acme Organ Donation"})`
func (m *Mixpanel) ProfileSetOnce(distinctID string, properties map[string]interface{}) error {
	return m.engage(distinctID, "$set_once", properties)
}

// ProfileAdd increments properties by the given amount for the profile
// that is referenced by the distinctID (which is the primary key)
// If you need to decrement a property, provide a negative value
// ip is optional
// e.g. `err := m.ProfileAdd("1", map[string]int{"items_created": 10, "invites_sent": -1})`
func (m *Mixpanel) ProfileAdd(distinctID string, properties map[string]int) error {
	return m.engage(distinctID, "$add", properties)
}

// ProfileAppend appends values to the given properties of the profile
// that is referenced by the distinctID (which is the primary key)
// ip is optional
// e.g. `err := m.ProfileAppend("1", map[string]interface{}{"level_ups": "sword obtained", "power_ups": "bubble lead"})`
func (m *Mixpanel) ProfileAppend(distinctID string, properties map[string]interface{}) error {
	return m.engage(distinctID, "$append", properties)
}

// ProfileUnion unions values to the given properties of the profile
// that is referenced by the distinctID (which is the primary key)
// ip is optional
// e.g. `err := m.ProfileUnion("1", map[string]interface{}{"items_purchased": []string{"socks", "shirts"}})`
func (m *Mixpanel) ProfileUnion(distinctID string, properties map[string]interface{}) error {
	return m.engage(distinctID, "$union", properties)
}

// ProfileUnset unions values to the given properties of the profile
// that is referenced by the distinctID (which is the primary key)
// ip is optional
// e.g. `err := m.ProfileUnset("1", []string{"Days Purchased"})`
func (m *Mixpanel) ProfileUnset(distinctID string, properties []string) error {
	return m.engage(distinctID, "$unset", properties)
}

// ProfileDelete deletes the profile that is referenced by the distinctID
// e.g. `err := m.ProfileDelete("1")`
func (m *Mixpanel) ProfileDelete(distinctID string) error {
	return m.engage(distinctID, "$delete", "")
}

// Alias alias'es an old distinct ID with the new distinct ID
// e.g. `err := m.ProfileCreateAliasDistinctIdToAlias("deadbeef", "1")`
func (m *Mixpanel) ProfileCreateAliasDistinctIdToAlias(oldID, newID string) error {
	return m.Track("$create_alias", map[string]interface{}{"distinct_id": oldID, "alias": newID})
}

func (m *Mixpanel) engage(distinctID string, op string, properties interface{}) error {
	var data map[string]interface{} = make(map[string]interface{})

	data["$token"] = m.Token
	data["$distinct_id"] = distinctID
	if len(m.OverrideIPAddress) > 0 {
		data["$ip"] = m.OverrideIPAddress
	}
	data[op] = properties

	response, err := m.get(fmt.Sprintf("%s/engage/", m.BaseURL), data)
	if err != nil {
		return err
	}

	if response != "1" {
		return ErrUnexpectedEngageResponse
	}

	return nil
}

func (m *Mixpanel) get(url string, data map[string]interface{}) (string, error) {
	jsonedData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	base64JSONData := base64.StdEncoding.EncodeToString(jsonedData)

	res, err := http.Get(fmt.Sprintf("%s?data=%s", url, base64JSONData))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)

	return string(responseBody), err
}
