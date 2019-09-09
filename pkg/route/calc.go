package route

import (
	"sort"

	"github.com/raddare/internal/osrm"
)

func (m *Manager) sortRoutes(responses []*osrm.Response) {
	sort.Slice(responses, func(i, j int) bool {

		// If durations are the same...
		if responses[i].Routes[0].Duration == responses[j].Routes[0].Duration {
			// Use distance choosing the shortes one
			return responses[i].Routes[0].Distance < responses[j].Routes[0].Distance
		}

		// Otherwise use duration.
		return responses[i].Routes[0].Duration < responses[j].Routes[0].Duration
	})
}

func (m *Manager) bestRoute(responses []*osrm.Response) *osrm.Response {
	sort.Slice(responses, func(i, j int) bool {

		// If durations are the same...
		if responses[i].Routes[0].Duration == responses[j].Routes[0].Duration {
			// Use distance choosing the shortes one
			return responses[i].Routes[0].Distance < responses[j].Routes[0].Distance
		}

		// Otherwise use duration.
		return responses[i].Routes[0].Duration < responses[j].Routes[0].Duration
	})

	return responses[0]
}
