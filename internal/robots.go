package internal

import (
	"io"
	"strings"
)

type RobotRuleBlock struct {
	UserAgents []string
	DisallowedURLs []string
}

type RobotRules struct {
	RuleBlocks []RobotRuleBlock
	Sitemap string
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
		}
	}

	return
}
