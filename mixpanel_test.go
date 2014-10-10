package mixpanel_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/nitrous-io/go-mixpanel"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

func TestContext(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "mixpanel")
}

var _ = Describe("mixpanel", func() {
	var server *ghttp.Server
	var baseURL string

	BeforeEach(func() {
		server = ghttp.NewServer()
		baseURL = server.URL()
	})

	AfterEach(func() {
		server.Close()
	})

	decodeBase64 := func(str string) string {
		data, err := base64.StdEncoding.DecodeString(str)
		Expect(err).To(BeNil())
		return bytes.NewBuffer(data).String()
	}

	verifyRequestResponse := func(server *ghttp.Server, method, uriRegexp, expectedData, responseData string) {
		var verifier http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			Expect(r.Method).To(Equal(method))
			Expect(r.RequestURI).To(MatchRegexp(uriRegexp))
			url, err := url.ParseRequestURI(r.RequestURI)
			Expect(err).To(BeNil())
			data := decodeBase64(url.Query().Get("data"))
			Expect(data).To(MatchJSON(expectedData))
			fmt.Fprint(w, responseData)
		}

		server.AppendHandlers(verifier)
	}

	Describe("NewMixpanelClient", func() {
		Context("with just a token", func() {
			It("should initialize a Mixpanel struct with the default base URL", func() {
				m := mixpanel.NewMixpanelClient("token")
				Expect(m.Token).To(Equal("token"))
				Expect(m.BaseURL).To(Equal(mixpanel.BASE_URL))
			})
		})

		Context("with token and base url", func() {
			It("should initialize a Mixpanel struct with the provided token and base url", func() {
				m := mixpanel.NewMixpanelClient("token", "http://localhost:3000")
				Expect(m.Token).To(Equal("token"))
				Expect(m.BaseURL).To(Equal("http://localhost:3000"))
			})
		})
	})

	Describe("Track", func() {
		Context("when mixpanel responds with a valid response", func() {
			BeforeEach(func() {
				verifyRequestResponse(server,
					"GET",
					`\A\/track\/\?data=.*?\z`,
					`{"event":"User Signed Up","properties":{"$distinct_id":"1","token":"token"}}`,
					"1",
				)
			})

			It("should send an base64 encoded version of the event along with the corresponding properties", func() {
				m := mixpanel.NewMixpanelClient("token", baseURL)
				err := m.Track("User Signed Up", map[string]interface{}{"$distinct_id": "1"})
				Expect(err).To(BeNil())
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("when mixpanel responds with an error", func() {
			BeforeEach(func() {
				verifyRequestResponse(server,
					"GET",
					`\A\/track\/\?data=.*?\z`,
					`{"event":"User Signed Up","properties":{"$distinct_id":"1","token":"token"}}`,
					"error",
				)
			})

			It("should send an base64 encoded version of the event along with the corresponding properties", func() {
				m := mixpanel.NewMixpanelClient("token", baseURL)
				err := m.Track("User Signed Up", map[string]interface{}{"$distinct_id": "1"})
				Expect(err).To(Equal(mixpanel.ErrUnexpectedTrackResponse))
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})
	})

	Describe("CreateProfile", func() {
		Context("when mixpanel responds with a valid response", func() {
			BeforeEach(func() {
				verifyRequestResponse(server,
					"GET",
					`\A\/engage\/\?data=.*?\z`,
					`{"$token":"token","$distinct_id":"1","$ip":"192.168.16.43","$set":{"full_name": "Mclovin", "Company": "Acme Organ Donation"}}`,
					"1",
				)
			})

			It("should send an base64 encoded version of the user profile along with the corresponding properties", func() {
				m := mixpanel.NewMixpanelClient("token", baseURL)
				err := m.CreateProfile("1", "192.168.16.43", map[string]interface{}{"full_name": "Mclovin", "Company": "Acme Organ Donation"})
				Expect(err).To(BeNil())
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("when mixpanel responds with an error", func() {
			BeforeEach(func() {
				verifyRequestResponse(server,
					"GET",
					`\A\/engage\/\?data=.*?\z`,
					`{"$token":"token","$distinct_id":"1","$ip":"192.168.16.43","$set":{"full_name": "Mclovin", "Company": "Acme Organ Donation"}}`,
					"error",
				)
			})

			It("should send an base64 encoded version of the event along with the corresponding properties", func() {
				m := mixpanel.NewMixpanelClient("token", baseURL)
				err := m.CreateProfile("1", "192.168.16.43", map[string]interface{}{"full_name": "Mclovin", "Company": "Acme Organ Donation"})
				Expect(err).To(Equal(mixpanel.ErrUnexpectedEngageResponse))
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})
	})

	Describe("SetPropertiesOnProfileOnce", func() {
		Context("when mixpanel responds with a valid response", func() {
			BeforeEach(func() {
				verifyRequestResponse(server,
					"GET",
					`\A\/engage\/\?data=.*?\z`,
					`{"$token":"token","$distinct_id":"1","$set_once":{"full_name": "Mclovin", "Company": "Acme Organ Donation"}}`,
					"1",
				)
			})

			It("should send an base64 encoded version of the user profile along with the corresponding properties", func() {
				m := mixpanel.NewMixpanelClient("token", baseURL)
				err := m.SetPropertiesOnProfileOnce("1", map[string]interface{}{"full_name": "Mclovin", "Company": "Acme Organ Donation"})
				Expect(err).To(BeNil())
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("when mixpanel responds with an error", func() {
			BeforeEach(func() {
				verifyRequestResponse(server,
					"GET",
					`\A\/engage\/\?data=.*?\z`,
					`{"$token":"token","$distinct_id":"1","$set_once":{"full_name": "Mclovin", "Company": "Acme Organ Donation"}}`,
					"error",
				)
			})

			It("should send an base64 encoded version of the event along with the corresponding properties", func() {
				m := mixpanel.NewMixpanelClient("token", baseURL)
				err := m.SetPropertiesOnProfileOnce("1", map[string]interface{}{"full_name": "Mclovin", "Company": "Acme Organ Donation"})
				Expect(err).To(Equal(mixpanel.ErrUnexpectedEngageResponse))
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})
	})

	Describe("IncrementPropertiesOnProfile", func() {
		Context("when mixpanel responds with a valid response", func() {
			BeforeEach(func() {
				verifyRequestResponse(server,
					"GET",
					`\A\/engage\/\?data=.*?\z`,
					`{"$token":"token","$distinct_id":"1","$add":{"items_created": 10, "invites_sent": -1}}`,
					"1",
				)
			})

			It("should send an base64 encoded version of the user profile along with the corresponding properties", func() {
				m := mixpanel.NewMixpanelClient("token", baseURL)
				err := m.IncrementPropertiesOnProfile("1", map[string]int{"items_created": 10, "invites_sent": -1})
				Expect(err).To(BeNil())
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("when mixpanel responds with an error", func() {
			BeforeEach(func() {
				verifyRequestResponse(server,
					"GET",
					`\A\/engage\/\?data=.*?\z`,
					`{"$token":"token","$distinct_id":"1","$add":{"items_created": 10, "invites_sent": -1}}`,
					"error",
				)
			})

			It("should send an base64 encoded version of the event along with the corresponding properties", func() {
				m := mixpanel.NewMixpanelClient("token", baseURL)
				err := m.IncrementPropertiesOnProfile("1", map[string]int{"items_created": 10, "invites_sent": -1})
				Expect(err).To(Equal(mixpanel.ErrUnexpectedEngageResponse))
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})
	})

	Describe("AppendPropertiesOnProfile", func() {
		Context("when mixpanel responds with a valid response", func() {
			BeforeEach(func() {
				verifyRequestResponse(server,
					"GET",
					`\A\/engage\/\?data=.*?\z`,
					`{"$token":"token","$distinct_id":"1","$append":{"level_ups": "sword obtained", "power_ups": "bubble lead"}}`,
					"1",
				)
			})

			It("should send an base64 encoded version of the user profile along with the corresponding properties", func() {
				m := mixpanel.NewMixpanelClient("token", baseURL)
				err := m.AppendPropertiesOnProfile("1", map[string]interface{}{"level_ups": "sword obtained", "power_ups": "bubble lead"})
				Expect(err).To(BeNil())
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("when mixpanel responds with an error", func() {
			BeforeEach(func() {
				verifyRequestResponse(server,
					"GET",
					`\A\/engage\/\?data=.*?\z`,
					`{"$token":"token","$distinct_id":"1","$append":{"level_ups": "sword obtained", "power_ups": "bubble lead"}}`,
					"error",
				)
			})

			It("should send an base64 encoded version of the event along with the corresponding properties", func() {
				m := mixpanel.NewMixpanelClient("token", baseURL)
				err := m.AppendPropertiesOnProfile("1", map[string]interface{}{"level_ups": "sword obtained", "power_ups": "bubble lead"})
				Expect(err).To(Equal(mixpanel.ErrUnexpectedEngageResponse))
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})
	})

	Describe("UnionPropertiesOnProfile", func() {
		Context("when mixpanel responds with a valid response", func() {
			BeforeEach(func() {
				verifyRequestResponse(server,
					"GET",
					`\A\/engage\/\?data=.*?\z`,
					`{"$token":"token","$distinct_id":"1","$union":{"items_purchased": ["socks", "shirts"]}}`,
					"1",
				)
			})

			It("should send an base64 encoded version of the user profile along with the corresponding properties", func() {
				m := mixpanel.NewMixpanelClient("token", baseURL)
				err := m.UnionPropertiesOnProfile("1", map[string]interface{}{"items_purchased": []string{"socks", "shirts"}})
				Expect(err).To(BeNil())
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("when mixpanel responds with an error", func() {
			BeforeEach(func() {
				verifyRequestResponse(server,
					"GET",
					`\A\/engage\/\?data=.*?\z`,
					`{"$token":"token","$distinct_id":"1","$union":{"items_purchased": ["socks", "shirts"]}}`,
					"error",
				)
			})

			It("should send an base64 encoded version of the event along with the corresponding properties", func() {
				m := mixpanel.NewMixpanelClient("token", baseURL)
				err := m.UnionPropertiesOnProfile("1", map[string]interface{}{"items_purchased": []string{"socks", "shirts"}})
				Expect(err).To(Equal(mixpanel.ErrUnexpectedEngageResponse))
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})
	})

	Describe("UnsetPropertiesOnProfile", func() {
		Context("when mixpanel responds with a valid response", func() {
			BeforeEach(func() {
				verifyRequestResponse(server,
					"GET",
					`\A\/engage\/\?data=.*?\z`,
					`{"$token":"token","$distinct_id":"1","$unset":["Days Purchased"]}`,
					"1",
				)
			})

			It("should send an base64 encoded version of the user profile along with the corresponding properties", func() {
				m := mixpanel.NewMixpanelClient("token", baseURL)
				err := m.UnsetPropertiesOnProfile("1", []string{"Days Purchased"})
				Expect(err).To(BeNil())
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("when mixpanel responds with an error", func() {
			BeforeEach(func() {
				verifyRequestResponse(server,
					"GET",
					`\A\/engage\/\?data=.*?\z`,
					`{"$token":"token","$distinct_id":"1","$unset":["Days Purchased"]}`,
					"error",
				)
			})

			It("should send an base64 encoded version of the event along with the corresponding properties", func() {
				m := mixpanel.NewMixpanelClient("token", baseURL)
				err := m.UnsetPropertiesOnProfile("1", []string{"Days Purchased"})
				Expect(err).To(Equal(mixpanel.ErrUnexpectedEngageResponse))
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})
	})

	Describe("DeleteProfile", func() {
		Context("when mixpanel responds with a valid response", func() {
			BeforeEach(func() {
				verifyRequestResponse(server,
					"GET",
					`\A\/engage\/\?data=.*?\z`,
					`{"$token":"token","$distinct_id":"1","$delete":""}`,
					"1",
				)
			})

			It("should send an base64 encoded version of the user profile along with the corresponding properties", func() {
				m := mixpanel.NewMixpanelClient("token", baseURL)
				err := m.DeleteProfile("1")
				Expect(err).To(BeNil())
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("when mixpanel responds with an error", func() {
			BeforeEach(func() {
				verifyRequestResponse(server,
					"GET",
					`\A\/engage\/\?data=.*?\z`,
					`{"$token":"token","$distinct_id":"1","$delete":""}`,
					"error",
				)
			})

			It("should send an base64 encoded version of the event along with the corresponding properties", func() {
				m := mixpanel.NewMixpanelClient("token", baseURL)
				err := m.DeleteProfile("1")
				Expect(err).To(Equal(mixpanel.ErrUnexpectedEngageResponse))
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})
	})

	Describe("Alias", func() {
		Context("when mixpanel responds with a valid response", func() {
			BeforeEach(func() {
				verifyRequestResponse(server,
					"GET",
					`\A\/track\/\?data=.*?\z`,
					`{"event":"$create_alias","properties":{"token": "token", "distinct_id":"deadbeef","alias":"1"}}`,
					"1",
				)
			})

			It("should send an base64 encoded version of the event along with the corresponding properties", func() {
				m := mixpanel.NewMixpanelClient("token", baseURL)
				err := m.Alias("deadbeef", "1")
				Expect(err).To(BeNil())
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("when mixpanel responds with an error", func() {
			BeforeEach(func() {
				verifyRequestResponse(server,
					"GET",
					`\A\/track\/\?data=.*?\z`,
					`{"event":"$create_alias","properties":{"token": "token", "distinct_id":"deadbeef","alias":"1"}}`,
					"error",
				)
			})

			It("should send an base64 encoded version of the event along with the corresponding properties", func() {
				m := mixpanel.NewMixpanelClient("token", baseURL)
				err := m.Alias("deadbeef", "1")
				Expect(err).To(Equal(mixpanel.ErrUnexpectedTrackResponse))
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})
	})
})
