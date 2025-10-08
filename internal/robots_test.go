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

func TestURLMatching(t *testing.T) {
	if !MatchURLRule("/hello/world", "/*/world") {
		t.Errorf("First assertion failed!")
	}

	if MatchURLRule("/hello/world", "/goodbye/world") {
		t.Errorf("Second assertion failed!")
	}

	if !MatchURLRule("/hello/world", "/hello/world") {
		t.Errorf("Third assertion failed!")
	}

	if !MatchURLRule("/hello/world", "/hello/*") {
		t.Errorf("Fourth assertion failed!")
	}
}
