package agent

import (
	"encoding/json"

	"github.com/shenghui0779/gochat/urls"
	"github.com/shenghui0779/gochat/wx"
)

type ResultAgentGet struct {
	AgentID            string          `json:"agentid"`
	Name               string          `json:"name"`
	SquareLogoURL      string          `json:"square_logo_url"`
	Description        string          `json:"description"`
	Close              int             `json:"close"`
	RedirectDomain     string          `json:"redirect_domain"`
	ReportLocationFlag int             `json:"report_location_flag"`
	ISReportenter      int             `json:"is_reportenter"`
	HomeURL            string          `json:"home_url"`
	AllowUserInfos     *AllowUserInfos `json:"allow_userinfos"`
	AllowPartys        *AllowPartys    `json:"allow_partys"`
	AllowTags          *AllowTags      `json:"allow_tags"`
}

type AllowUserInfos struct {
	User []*AllowUser `json:"user"`
}

type AllowUser struct {
	UserID string `json:"userid"`
}

type AllowPartys struct {
	PartID []int64 `json:"partid"`
}

type AllowTags struct {
	TagID []int64 `json:"tagid"`
}

func AgentGet(agentID string, result *ResultAgentGet) wx.Action {
	return wx.NewGetAction(urls.CorpAgentGet,
		wx.WithQuery("agent_id", agentID),
		wx.WithDecode(func(resp []byte) error {
			return json.Unmarshal(resp, result)
		}),
	)
}

type ResultAgentList struct {
	AgentList []*AgentListItem `json:"agentlist"`
}

type AgentListItem struct {
	AgentID       string `json:"agentid"`
	Name          string `json:"name"`
	SquareLogoURL string `json:"square_logo_url"`
}

func AgentList(result *ResultAgentList) wx.Action {
	return wx.NewGetAction(urls.CorpAgentList,
		wx.WithDecode(func(resp []byte) error {
			return json.Unmarshal(resp, result)
		}),
	)
}

type ParamsAgentSet struct {
	AgentID            string `json:"agentid"`
	Name               string `json:"name"`
	LogoMediaID        string `json:"logo_mediaid"`
	Description        string `json:"description"`
	RedirectDomain     string `json:"redirect_domain"`
	ReportLocationFlag int    `json:"report_location_flag"`
	ISReportenter      int    `json:"is_reportenter"`
	HomeURL            string `json:"home_url"`
}

func AgentSet(params *ParamsAgentSet) wx.Action {
	return wx.NewPostAction(urls.CorpAgentSet,
		wx.WithBody(func() ([]byte, error) {
			return json.Marshal(params)
		}),
	)
}
