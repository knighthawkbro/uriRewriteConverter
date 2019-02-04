package lib

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// Web.Config files typically have this structure
// configuration
//   system.webServer
//	   rewrite
//       rules
// 		   rule name, stopProcessing
//		     match url(rexex) ignoreCase
//			 conditions logicalGrouping, trackAllCaptures
//           action type, url

// Configuration struct for an URLRewrite XML file
type Configuration struct {
	XMLName         xml.Name        `xml:"configuration"`
	SystemWebServer SystemWebServer `xml:"system.webServer"`
}

// SystemWebServer struct
type SystemWebServer struct {
	Rewrite Rewrite `xml:"rewrite"`
}

// Rewrite struct
type Rewrite struct {
	Rules []Rule `xml:"rules>rule"`
}

// Rule struct
type Rule struct {
	Name           string     `xml:"name,attr"`
	StopProcessing *bool      `xml:"stopProcessing,attr"`
	Match          Match      `xml:"match"`
	Conditions     Conditions `xml:"conditions"`
	Action         Action     `xml:"action"`
}

// Match struct
type Match struct {
	URL        string `xml:"url,attr"`
	IgnoreCase *bool  `xml:"ignoreCase,attr"`
}

// Conditions struct
type Conditions struct {
	LogicalGrouping  string `xml:"logicalGrouping,attr"`
	TrackAllCaptures bool   `xml:"trackAllCaptures,attr"`
}

// Action struct The action of a Rule has two attributes, type and the URL
type Action struct {
	Type string `xml:"type,attr"`
	URL  string `xml:"url,attr"`
}

// Unmarshal function takes in data and unmarshals it using the XML library
func Unmarshal(data []byte) *Configuration {
	v := &Configuration{}
	err := xml.Unmarshal(data, v)
	CheckErr("Cannot Unmarshall data", err)
	return v
}

// Marshal function marshals configuration struct to string
func (x *Configuration) Marshal() string {
	output, err := xml.MarshalIndent(x, "", "    ")
	CheckErr("Could not marshal data", err)
	return fmt.Sprintf("%s%s", xml.Header, output)
}

// ToHTACL function takes a populated Configuration XML struct and transforms it to
// URL Rewrite Engine struct
func (x *Configuration) ToHTACL() *HTACL {
	res := &HTACL{}

	res.RewriteEngine = "on"
	res.RewriteBase = "/"
	for _, r := range x.SystemWebServer.Rewrite.Rules {
		rule := RewriteRule{}
		if r.Match.URL[0] == '^' {
			rule.Regex = r.Match.URL + ">"
		} else {
			rule.Regex = r.Match.URL
		}
		rule.URL = r.Action.URL
		if strings.Contains(r.Action.URL, "{R:1}") {
			rule.URL = strings.Replace(r.Action.URL, "{R:1}", "$1", 1)
		}
		tmp := []string{}
		if r.Match.IgnoreCase == nil || *r.Match.IgnoreCase {
			tmp = append(tmp, "NC")
		}
		if r.StopProcessing == nil || *r.StopProcessing {
			tmp = append(tmp, "L")
		}
		rule.Parameters = tmp
		res.RewriteRules = append(res.RewriteRules, rule)
	}

	return res
}
