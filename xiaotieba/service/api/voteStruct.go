package api

import "sync"

var VOTING sync.Map

type Voting struct {
	topic   string
	starter string
	VoteID  int
	joiner  map[string]string
}
