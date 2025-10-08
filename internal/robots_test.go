package internal

import (
	"testing"
	"strings"
)

func TestRobotParsing(t *testing.T) {
	input := `User-agent: Neng Li
Disallow: /yes/no
Disallow: /bozo

User-agent: Indeed
Disallow: /no/`

	obtained, err := ParseRobotTxt(strings.NewReader(input))
	if err != nil {
		t.Fatalf("It encountered an error while parsing: %v\n", err)
	}

	expected := RobotRules {
		RuleBlocks: []RobotRuleBlock{
			{
				UserAgents: []string{"Neng Li"},
				DisallowedURLs: []string{"/yes/no", "/bozo"},
			},
			{
				UserAgents: []string{"Indeed"},
				DisallowedURLs: []string{"/no/"},
			},
		},
		Sitemap: "",
	}

	if !expected.IsEqual(&obtained) {
		t.Errorf("Expected:\n%v\n\nGot:\n%v\n", expected, obtained)
	}
}
