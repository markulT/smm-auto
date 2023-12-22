package scrapper

import (
	"regexp"
)

type ChatSubmissionParsingError struct{}

func (c ChatSubmissionParsingError) Error() string {
	return "Error parsing admin chat from channel description (no matches)"
}

func getAdminChatFromDescription(hc string) (string, error) {
	re, err := regexp.Compile(`Admin - (\w+)`)
	if err != nil {
		return "", err
	}

	// Find the match
	match := re.FindStringSubmatch(hc)
	if len(match) < 2 {
		return "", ChatSubmissionParsingError{}
	}
	return match[1], err
}
