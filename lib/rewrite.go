package lib

import (
	"fmt"
	"os"
	"strings"
)

// HTACL struct is the 3 types of lines that are contained inside of a ht.acl file
type HTACL struct {
	RewriteEngine string
	RewriteBase   string
	RewriteRules  []RewriteRule
}

// RewriteRule struct a Rule typically some regex to match a URI, a URL to go to, and some parameters like NoCase and Stop Execution
type RewriteRule struct {
	Regex      string
	URL        string
	Parameters []string
}

// Marshal function
func (a *HTACL) Marshal() string {
	result := ""
	if a.RewriteEngine != "" {
		result += "RewriteEngine " + a.RewriteEngine
	}
	if a.RewriteBase != "" {
		result += "\nRewriteBase " + a.RewriteBase
	}
	if len(a.RewriteRules) > 0 {
		for _, rule := range a.RewriteRules {
			result += fmt.Sprintf("\nRewriteRule %s %s [%s]", rule.Regex, rule.URL, strings.Join(rule.Parameters, ", "))
		}
	}
	return result
}

// Unmarshal function
func (a *HTACL) Unmarshal(input []string) {
	if len(input) < 2 || strings.TrimSpace(input[0]) == "" {
		return
	}
	if input[0] == "RewriteEngine" {
		a.RewriteEngine = input[1]
	} else if input[0] == "RewriteBase" {
		a.RewriteBase = input[1]
	} else if input[0] == "RewriteRule" {
		rule := RewriteRule{}
		rule.Regex = input[1]
		rule.URL = input[2]
		tmp := []string{}
		for n, i := range input[3:] {
			if n == len(input[3:])-1 {
				tmp = append(tmp, i[:len(i)-1])
				break
			} else if n == 0 {
				tmp = append(tmp, i[1:len(i)-1])
				continue
			}
			tmp = append(tmp, i[:len(i)])
		}
		rule.Parameters = tmp
		if a.exists(rule) {
			fmt.Println("Warn: duplicate line found")
			return
		}
		a.RewriteRules = append(a.RewriteRules, rule)
	} else if input[0] == "Header" {
		return
	} else {
		fmt.Printf("Unexpect HTACL argument: " + input[0])
		os.Exit(1)
	}
}

func (a *HTACL) exists(rule RewriteRule) bool {
	for _, r := range a.RewriteRules {
		if r.URL == rule.URL {
			if r.Regex == rule.Regex {
				flag := true
				for _, parameter := range rule.Parameters {
					if !Contains(r.Parameters, parameter) {
						flag = false
					}
				}
				if flag {
					return true
				}
			}
		}
	}
	return false
}

// ToWebConfig function transforms a HTACL file struct to a Configuration XML struct
func (a *HTACL) ToWebConfig() *Configuration {
	res := &Configuration{}
	for _, r := range a.RewriteRules {
		rule := Rule{}
		if Contains(r.Parameters, "L") {
			// TODO: Feel like there is a better way to do this.
			b := true
			rule.StopProcessing = &b
		}
		if Contains(r.Parameters, "NC") {
			b := true
			rule.Match.IgnoreCase = &b
		}
		if r.Regex[0] == '^' {
			rule.Name = r.Regex[:len(r.Regex)-2] + " to " + r.URL
			rule.Match.URL = r.Regex[:len(r.Regex)-2]
		} else {
			rule.Name = r.Regex + " to " + r.URL
			rule.Match.URL = r.Regex
		}
		rule.Conditions.LogicalGrouping = "MatchAll"
		rule.Conditions.TrackAllCaptures = false
		rule.Action.Type = "Rewrite"
		rule.Action.URL = r.URL
		// TODO: Rewrite this to use regex instead of using static strings
		if strings.Contains(r.URL, "$1") {
			rule.Action.URL = strings.Replace(r.URL, "$1", "{R:1}", 1)
		}
		res.SystemWebServer.Rewrite.Rules = append(res.SystemWebServer.Rewrite.Rules, rule)
	}
	return res
}
