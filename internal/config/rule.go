package config

import (
	"github.com/nobe4/gh-not/internal/jq"
	"github.com/nobe4/gh-not/internal/notifications"
)

// Rule is a struct to filter and act on notifications.
//
//	rules:
//	  - name: showcasing conditionals
//	    action: debug
//	    filters:
//	      - .author.login == "dependabot[bot]"
//	      - >
//	        (.subject.title | contains("something unimportant")) or
//	        (.subject.title | contains("something already done"))
//
//	  - name: ignore ci failures for the current repo
//	    action: done
//	    filters:
//	      - .repository.full_name == "nobe4/gh-not"
//	      - .reason == "ci_activity"
type Rule struct {
	Name string `yaml:"name"`

	// Filters is a list of jq filters to filter the notifications.
	// The filters are applied in order, like they are joined by 'and'.
	// Having 'or' can be done via '(cond1) or (cond2) or ...'.
	//
	// E.g.:
	// filters: ["A", "B or C"]
	// Will filter `A and (B or C)`.
	Filters []string `yaml:"filters"`

	// Action is the action to take on the filtered notifications.
	// See github.com/nobe4/internal/actors for list of available actions.
	Action string `yaml:"action"`
}

// FilterIds filters the notifications with the jq filters and returns the IDs.
func (r Rule) FilterIds(n notifications.Notifications) ([]string, error) {
	var err error

	for _, filter := range r.Filters {
		n, err = jq.Filter(filter, n)
		if err != nil {
			return nil, err
		}
	}

	return n.IDList(), nil
}
