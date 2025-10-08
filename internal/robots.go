package internal

import (
	"io"
	"strings"
)

type RobotRuleBlock struct {
	UserAgents     []string
	DisallowedURLs []string
}

type RobotRules struct {
	RuleBlocks []RobotRuleBlock
	Sitemap    string
}

type RobotSyntaxError struct {
	Details string
}

func (err *RobotSyntaxError) Error() string {
	return err.Details
}

func (block *RobotRuleBlock) IsEmpty() bool {
	return len(block.DisallowedURLs) == 0
}

// This is only for testing purposes only - I don't see when we would
// ever use this in any actual code lol.
func (rules *RobotRules) IsEqual(b *RobotRules) bool {
	for i, block := range rules.RuleBlocks {
		for j, agent := range block.UserAgents {
			if b.RuleBlocks[i].UserAgents[j] != agent {
				return false
			}
		}

		for j, disallow := range block.DisallowedURLs {
			if b.RuleBlocks[i].DisallowedURLs[j] != disallow {
				return false
			}
		}
	}

	if rules.Sitemap != b.Sitemap {
		return false
	}

	return true
}

func ParseRobotTxt(reader io.Reader) (rules RobotRules, err error) {
	robotsFile, err := io.ReadAll(reader)
	if err != nil {
		return
	}

	currentBlock := RobotRuleBlock{}

	for line := range strings.Lines(string(robotsFile)) {
		// Comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Blank lines
		if strings.TrimSpace(line) == "" {
			if !currentBlock.IsEmpty() {
				rules.RuleBlocks = append(rules.RuleBlocks, currentBlock)
			}

			currentBlock = RobotRuleBlock{}
			continue
		}

		// TODO: Maybe do some more error checking?

		// Normal lines
		lineParts := strings.Split(line, ": ")
		if len(lineParts) != 2 {
			err = &RobotSyntaxError{"Unexpected ':' right here."}
			return
		}

		// Pretty sure the key is case insensitive, so we can just do this.
		key := strings.ToLower(lineParts[0])
		value := strings.TrimSpace(lineParts[1])

		switch key {
		case "user-agent":
			currentBlock.UserAgents = append(currentBlock.UserAgents, value)
		case "disallow":
			currentBlock.DisallowedURLs = append(currentBlock.DisallowedURLs, value)
		}

		// When we run into an unknown key we just ignore it lol
	}

	if !currentBlock.IsEmpty() {
		rules.RuleBlocks = append(rules.RuleBlocks, currentBlock)
	}

	return
}

// Since it seems like stripping the prefix and stuff would be too difficult,
// I decided that you are going to do that manually before calling this function.
// Therefore, do not worry if you don't end up dying from this shit hole
func MatchURLRule(url string, rule string) bool {
	urlParts := strings.Split(url, "/")
	ruleParts := strings.Split(rule, "/")

	// Obviously if they have different parts they won't match.
	if len(urlParts) != len(ruleParts) {
		return false
	}

	for i, urlPart := range(urlParts) {
		// '*' matches anything
		if ruleParts[i] == "*" {
			continue
		}

		if ruleParts[i] != urlPart {
			return false
		}
	}

	return true
}
