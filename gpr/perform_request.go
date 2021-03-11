package gpr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/mtfelian/utils"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

const (
	authMethodBasic = "basic"
	authMethodToken = "token"
	authMethodNone  = "none"
)

const (
	DefaultBasicAuthLogin    = "test@example.com"
	DefaultBasicAuthPassword = "never mind"
)

// BasePath for server API
const BasePath = ""

// GPR abstracts ginkgo perform request
type GPR struct {
	authMethod        string
	logPerformRequest bool
	accessToken       strfmt.UUID // for token auth
	login, password   string      // for basic auth
	Logger            *logrus.Logger
	TestServer        **httptest.Server
}

// New returns a pointer to a new Ginkgo Perform Request (GPR) object with default settings
func New(ts *httptest.Server) *GPR {
	logrus.SetLevel(logrus.InfoLevel)
	return &GPR{
		authMethod:        authMethodBasic,
		logPerformRequest: true,
		login:             DefaultBasicAuthLogin,
		password:          DefaultBasicAuthPassword,
		Logger:            logrus.StandardLogger(),
		TestServer:        &ts,
	}
}

// SetLogPerformRequest to value v
func (g *GPR) SetLogPerformRequest(v bool) { g.logPerformRequest = v }

// SetAuthMethodNone for requests
func (g *GPR) SetAuthMethodNone() { g.authMethod = authMethodNone }

// SetAuthMethodBasic for requests
func (g *GPR) SetAuthMethodBasic() { g.authMethod = authMethodBasic }

// SetAuthMethodToken for requests
func (g *GPR) SetAuthMethodToken() { g.authMethod = authMethodToken }

// SetAccessToken for AuthModeToken authorization mode
func (g *GPR) SetAccessToken(token strfmt.UUID) { g.accessToken = token }

// AccessToken returns an access token set previously
func (g *GPR) AccessToken() strfmt.UUID { return g.accessToken }

// SetLogin for AuthModeBasic authorization mode
func (g *GPR) SetLogin(login string) { g.login = login }

// Login returns a login set previously
func (g *GPR) Login() string { return g.login }

// SetPassword for AuthModeBasic authorization mode
func (g *GPR) SetPassword(password string) { g.password = password }

// Password returns a password set previously
func (g *GPR) Password() string { return g.password }

// doAuthorizedRequest with HTTP basic auth with given method, url and content for body
func (g *GPR) doAuthorizedRequest(method, url string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url, body)
	ExpectWithOffset(3, err).NotTo(HaveOccurred(), "err from http.NewRequest in doAuthorizedRequest")
	request.Header.Add("content-type", "application/json")
	switch g.authMethod {
	case authMethodNone: // do nothing
	case authMethodToken:
		request.Header.Add("access-token", g.AccessToken().String())
	default: // default basic auth basic
		request.SetBasicAuth(g.Login(), g.Password())
	}
	return http.DefaultClient.Do(request)
}

// doRequest with given method, url and content for body
func (g *GPR) doRequest(method, url string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url, body)
	ExpectWithOffset(3, err).NotTo(HaveOccurred(), "err from http.NewRequest in doRequest")
	request.Header.Add("content-type", "application/json")
	return http.DefaultClient.Do(request)
}

// HTTP server API request helpers

// PerformRequest is a helper func to perform HTTP API requests to the dbcore
func (g *GPR) PerformRequest(addr, method string, JSON []byte, expectedStatusCode int, target interface{}) {
	a := fmt.Sprintf("%s%s%s", (*g.TestServer).URL, BasePath, addr)

	var bytesReader io.Reader
	if JSON != nil {
		bytesReader = bytes.NewReader(JSON)
	}

	response, err := g.doAuthorizedRequest(method, a, bytesReader)
	ExpectWithOffset(2, err).NotTo(HaveOccurred(), "err from doAuthorizedRequest in PerformRequest")
	ExpectWithOffset(2, response).NotTo(BeNil(), "in PerformRequest")
	responseBody, err := ioutil.ReadAll(response.Body)

	if g.logPerformRequest {
		g.Logger.Warnf(">>> PerformRequest() via %s on %s: <<< %s\n", method, addr, string(responseBody))
	}

	ExpectWithOffset(2, response.StatusCode).To(Equal(expectedStatusCode), "in PerformRequest")
	ExpectWithOffset(2, err).NotTo(HaveOccurred(), "err from ioutil.ReadAll in PerformRequest")

	if strings.Contains(addr, "/auth/") { // stop checking if it's auth request
		return
	}

	if target != nil {
		ExpectWithOffset(2, json.Unmarshal(responseBody, &target)).To(Succeed())
	}
}

// MustIndentJSON indents the given JSON bytes b, panics on error
func (g *GPR) MustIndentJSON(b []byte) []byte {
	var indentedJSON bytes.Buffer
	if err := json.Indent(&indentedJSON, b, "", "  "); err != nil {
		panic(err)
	}
	return indentedJSON.Bytes()
}

// MustMarshalJSONToString v into a JSON string, panics on error
func (g *GPR) MustMarshalJSONToString(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// MustMarshalJSON is just like MustMarshalJSONToString but returns []byte
func (g *GPR) MustMarshalJSON(v interface{}) []byte { return []byte(g.MustMarshalJSONToString(v)) }

// Marshal data to JSON and unmarshal it to out.
// Useful to convert interface{} to something more concrete via JSON
func (g *GPR) MarshalUnmarshalJSON(data, out interface{}) error {
	return utils.MarshalUnmarshalJSON(data, out)
}
