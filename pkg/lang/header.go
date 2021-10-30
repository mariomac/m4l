package lang

import (
	"bufio"
	"io"
	"regexp"
)

var ignoreLine = regexp.MustCompile(`^\s*(;.*)?\n?$`)
var headerProperty = regexp.MustCompile(`^\s*([\w\.]+)\s+([\w\.]+)\s*(;.*)?\n?$`)

// second argument: lines read
func parseHeader(reader io.ReadSeeker) (map[string]string, int, error) {
	props := map[string]string{}
	lineRead := bufio.NewReader(reader)
	lines := 0
	readBytes := 0
	for {
		line, err := lineRead.ReadString('\n')
		if err != nil {
			return nil, 0, err
		}
		lines++
		readBytes += len(line)
		if ignoreLine.MatchString(line) {
			continue
		}
		sm := headerProperty.FindStringSubmatch(line)
		if sm == nil {
			// no submatch, we assume end of the header zone so we rewind to the beginning of the line
			if _, err := reader.Seek(int64(readBytes-len(line)), io.SeekStart); err != nil {
				return nil, 0, err
			}
			return props, lines - 1, nil
		}
		props[sm[1]] = sm[2]
	}
}
