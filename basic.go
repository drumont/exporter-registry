package exporter_registry

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type NginxStats struct {
	ConnectionsActive float64
}

func ScanBasicStats(r io.Reader) ([]NginxStats, error) {
	s := bufio.NewScanner(r)

	var stats []NginxStats
	var nginxStats NginxStats
	for s.Scan() {
		fields := strings.Fields(string(s.Bytes()))
		fmt.Println(fields)

		if len(fields) == 3 && fields[0] == "Active" {
			c, err := strconv.ParseFloat(fields[2], 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing Active: %v", err)
			}
			nginxStats.ConnectionsActive = c
		}
	}

	stats = append(stats, nginxStats)

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("error scanning basic stats: %w", err)
	}

	return stats, nil
}
