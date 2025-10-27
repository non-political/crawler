package internal

import (
	"regexp"
	"strings"
	"testing"
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

	yesNoRule, _ := regexp.Compile("/yes/no")
	bozoRule, _ := regexp.Compile("/bozo")
	noRule, _ := regexp.Compile("/no/")

	expected := RobotRules {
		RuleBlocks: []RobotRuleBlock{
			{
				UserAgents: []string{"Neng Li"},
				DisallowedURLs: []*regexp.Regexp{yesNoRule, bozoRule},
			},
			{
				UserAgents: []string{"Indeed"},
				DisallowedURLs: []*regexp.Regexp{noRule},
			},
		},
		Sitemap: "",
	}

	if !expected.IsEqual(&obtained) {
		t.Errorf("Expected:\n%v\n\nGot:\n%v\n", expected, obtained)
	}
}

func TestURLMatching(t *testing.T) {
	rule, _ := regexp.Compile("/*/world")
	if !MatchURLRule("/hello/world", rule) {
		t.Errorf("First assertion failed!")
	}

	rule, _ = regexp.Compile("/goodbye/world")
	if MatchURLRule("/hello/world", rule) {
		t.Errorf("Second assertion failed!")
	}

	rule, _ = regexp.Compile("/hello/world")
	if !MatchURLRule("/hello/world", rule) {
		t.Errorf("Third assertion failed!")
	}

	rule, _ = regexp.Compile("/hello/*")
	if !MatchURLRule("/hello/world", rule) {
		t.Errorf("Fourth assertion failed!")
	}
}
