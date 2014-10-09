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
	ErrUnexpectedTrackResponse  = fmt.Errorf("Unexpected Mixpanel Track Response")
	ErrUnexpectedEngageResponse = fmt.Errorf("Unexpected Mixpanel Engage Response")
)

type Mixpanel struct {
	Token   string
	BaseURL string
}

func NewMixpanelClient(args ...string) *Mixpanel {
	var m *Mixpanel

	if len(args) == 1 {
		m = &Mixpanel{Token: args[0], BaseURL: BASE_URL}
	} else if len(args) > 1 {
		m = &Mixpanel{Token: args[0], BaseURL: args[1]}
	}

	return m
}

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

func (m *Mixpanel) CreateProfile(distinctID string, properties map[string]interface{}) error {
	return m.engage(distinctID, "$set", properties)
}

func (m *Mixpanel) SetPropertiesOnProfileOnce(distinctID string, properties map[string]interface{}) error {
	return m.engage(distinctID, "$set_once", properties)
}

func (m *Mixpanel) IncrementPropertiesOnProfile(distinctID string, properties map[string]int) error {
	return m.engage(distinctID, "$add", properties)
}

func (m *Mixpanel) AppendPropertiesOnProfile(distinctID string, properties map[string]interface{}) error {
	return m.engage(distinctID, "$append", properties)
}

func (m *Mixpanel) UnionPropertiesOnProfile(distinctID string, properties map[string]interface{}) error {
	return m.engage(distinctID, "$union", properties)
}

func (m *Mixpanel) UnsetPropertiesOnProfile(distinctID string, properties []string) error {
	return m.engage(distinctID, "$unset", properties)
}

func (m *Mixpanel) DeleteProfile(distinctID string) error {
	return m.engage(distinctID, "$delete", "")
}

func (m *Mixpanel) Alias(oldID, newID string) error {
	return m.Track("$create_alias", map[string]interface{}{"distinct_id": oldID, "alias": newID})
}

func (m *Mixpanel) engage(distinctID string, op string, properties interface{}) error {
	var data map[string]interface{} = make(map[string]interface{})

	data["$token"] = m.Token
	data["$distinct_id"] = distinctID
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

	responseBody, err := ioutil.ReadAll(res.Body)

	return string(responseBody), err
}
