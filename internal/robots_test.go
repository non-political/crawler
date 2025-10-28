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

# Comment
User-agent: Indeed
User-agent: Interesting
# Comment
Disallow: /no/`

	obtained, err := ParseRobotTxt(strings.NewReader(input))
	if err != nil {
		t.Fatalf("It encountered an error while parsing: %v\n", err)
	}

	yesNoRule, _ := regexp.Compile("/yes/no")
	bozoRule, _ := regexp.Compile("/bozo")
	noRule, _ := regexp.Compile("/no/")

	expected := RobotRules{
		RuleBlocks: []RobotRuleBlock{
			{
				UserAgents:     []string{"Neng Li"},
				DisallowedURLs: []*regexp.Regexp{yesNoRule, bozoRule},
			},
			{
				UserAgents:     []string{"Indeed", "Interesting"},
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
	tests := map[int]struct {
		Rule     string
		URL      string
		Expected bool
	}{
		1:  {Rule: "/*/world", URL: "/hello/world", Expected: true},
		2:  {Rule: "/goodbye/world", URL: "/hello/world", Expected: false},
		3:  {Rule: "/hello/world", URL: "/hello/world", Expected: true},
		4:  {Rule: "/hello/*", URL: "/hello/world", Expected: true},
		5:  {Rule: "/", URL: "/hello/world", Expected: true},
		6:  {Rule: "/hello/", URL: "/hello/world", Expected: true},
		7:  {Rule: "/Hello/World", URL: "/hello/world", Expected: false}, // case sensitivity
		8:  {Rule: "/hello/$", URL: "/hello/world/hi", Expected: false},  // anchored end
		9:  {Rule: "/hello$", URL: "/hello", Expected: true},             // end anchor match
		10: {Rule: "/*.php$", URL: "/index.php", Expected: true},
		11: {Rule: "/*.php$", URL: "/index.php?parameter", Expected: false},
		12: {Rule: "/$", URL: "/", Expected: true},
		13: {Rule: "/$", URL: "/hello", Expected: false},
	}

	for id, test := range tests {
		rule, err := regexp.Compile(test.Rule)
		if err != nil {
			t.Errorf("Test %d: failed to compile regex %q: %v", id, test.Rule, err)
			continue
		}
		result := MatchURLRule(test.URL, rule)
		if result != test.Expected {
			t.Errorf("Test %d: rule=%q, url=%q, expected %t, got %t",
				id, test.Rule, test.URL, test.Expected, result)
		}
	}
}
