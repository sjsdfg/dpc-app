package v2

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CMSgov/dpc/api/client"
	middleware2 "github.com/CMSgov/dpc/api/middleware"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SsasControllerTestSuite struct {
	suite.Suite
	sc  *SSASController
	msc *MockSsasClient
	mac *MockAttributionClient
}

func (suite *SsasControllerTestSuite) SetupTest() {
	mac := new(MockAttributionClient)
	msc := new(MockSsasClient)
	suite.mac = mac
	suite.msc = msc
	suite.sc = NewSSASController(suite.msc, suite.mac)

}

func TestSsasControllerTestSuite(t *testing.T) {
	suite.Run(t, new(SsasControllerTestSuite))
}

func (suite *SsasControllerTestSuite) TestCreateSystem() {

	req, _ := suite.SetupHappyPathMocks()

	//Do request
	w := httptest.NewRecorder()
	suite.sc.CreateSystem(w, req)
	res := w.Result()

	ja := jsonassert.New(suite.T())
	assert.Equal(suite.T(), http.StatusOK, res.StatusCode)
	resp, _ := ioutil.ReadAll(res.Body)
	ja.Assertf(string(resp), `
    {
        "client_id":"c001",
        "client_name":"Test Org",
        "ips":["ip-1", "ip-2"],
        "client_token":"client-token",
        "expires_at":"expiration"
    }`)
}

func (suite *SsasControllerTestSuite) TestCreateDuplicateSystem() {
	req, _ := suite.SetupHappyPathMocks()

	//Mock client calls
	managedOrg := client.ProviderOrg{
		OrgName:      "Test Org",
		OrgID:        "abc",
		Npi:          "npi-1",
		Status:       "Active",
		SsasSystemID: "42",
	}
	orgs := make([]client.ProviderOrg, 1)
	orgs[0] = managedOrg
	findExpectedCall(suite.mac.ExpectedCalls, "GetProviderOrgs").Return(orgs, nil)

	//Do request
	w := httptest.NewRecorder()
	suite.sc.CreateSystem(w, req)
	res := w.Result()

	assert.Equal(suite.T(), http.StatusConflict, res.StatusCode)
}

func (suite *SsasControllerTestSuite) TestCreateSystemForInactiveRelation() {

	req, _ := suite.SetupHappyPathMocks()

	//Mock client calls
	managedOrg := client.ProviderOrg{
		OrgName:      "Test Org",
		OrgID:        "abc",
		Npi:          "npi-1",
		Status:       "Inactive",
		SsasSystemID: "",
	}
	orgs := make([]client.ProviderOrg, 1)
	orgs[0] = managedOrg
	findExpectedCall(suite.mac.ExpectedCalls, "GetProviderOrgs").Return(orgs, nil)

	//Do request
	w := httptest.NewRecorder()
	suite.sc.CreateSystem(w, req)
	res := w.Result()
	resp, _ := ioutil.ReadAll(res.Body)

	assert.Equal(suite.T(), http.StatusBadRequest, res.StatusCode)
	assert.Contains(suite.T(), string(resp), "Implementor/Org relation is not active")
}

func (suite *SsasControllerTestSuite) SetupHappyPathMocks() (*http.Request, context.Context) {
	//Setup request
	reqBody := `{
        "client_name" : "Test Client",
        "public_key" : "public key",
        "ips" : ["ip-1","ip-2"]
    }`
	req := httptest.NewRequest("Post", "http://localohost/v2/Implementer/123/Org/abc/Systemfoo", strings.NewReader(reqBody))
	ctx := req.Context()
	ctx = context.WithValue(ctx, middleware2.ContextKeyImplementer, "123")
	req = req.WithContext(ctx)
	ctx = req.Context()
	ctx = context.WithValue(ctx, middleware2.ContextKeyOrganization, "abc")
	req = req.WithContext(ctx)

	//Mock client calls
	managedOrg := client.ProviderOrg{
		OrgName:      "Test Org",
		OrgID:        "abc",
		Npi:          "npi-1",
		Status:       "Active",
		SsasSystemID: "",
	}
	orgs := make([]client.ProviderOrg, 1)
	orgs[0] = managedOrg
	suite.mac.On("GetProviderOrgs", mock.Anything, mock.Anything).Return(orgs, nil)
	suite.mac.On("Get", mock.Anything, mock.Anything, mock.Anything).Return([]byte(`{"ssas_group_id":"0001"}`), nil)

	ips := make([]string, 2)
	ips[0] = "ip-1"
	ips[1] = "ip-2"

	ssasResp := client.CreateSystemResponse{
		SystemID:    "1",
		ClientName:  "Test Org",
		GroupID:     "0001",
		Scope:       "dpc",
		PublicKey:   "public-key",
		ClientID:    "c001",
		ClientToken: "client-token",
		ExpiresAt:   "expiration",
		XData:       "xdata",
		IPs:         ips,
	}
	suite.msc.On("CreateSystem", mock.Anything, mock.Anything).Return(ssasResp, nil)
	suite.mac.On("UpdateImplOrg", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(client.ImplementerOrg{}, nil)

	return req, ctx
}

func findExpectedCall(calls []*mock.Call, methodName string) *mock.Call {
	for _, c := range calls {
		if c.Method == methodName {
			return c
		}
	}
	return nil
}
