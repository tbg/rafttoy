// Copyright 2019 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package quorum

import (
	"fmt"
	"math"
	"strconv"
)

// IndexLookuper allows looking up a commit index for a given ID of a voter
// from a corresponding MajorityConfig.
type IndexLookuper interface {
	Index(voterID uint64) (idx uint64, found bool)
}

type mapLookuper map[uint64]uint64

func (m mapLookuper) Index(id uint64) (uint64, bool) {
	idx, ok := m[id]
	return idx, ok
}

// CommitRange is a tuple of commit indexes, Definitely, which is known to be
// committed, and Maybe, the highest index that could potentially be committed
// as more voters report back (while existing voters don't change their answer).
//
// For example, if the configuration consists of voters with IDs 1, 3, and 8, and
// - voter 3 has reported an index of 12,
// - voter 8 has reported an index of 11, and
// - voter 1 hasn't reported anything yet,
// then the CommitRange will have
// - Definitely=11 because voters 3 and 8 have acked index 11, but not 12, and
// - Maybe=12 (because if voter 1 reported in with index 12 or higher, the
//   committed index would be 12, but never higher).
//
// Maybe is mostly informational, though it is used by VoteResult() to detect
// definite outcomes (by observing that Definitely equals Maybe).
type CommitRange struct {
	Definitely uint64
	Maybe      uint64
}

// String implements fmt.Stringer. A CommitRange is printed as a uint64 if it
// contains only one possible commit index; otherwise it's printed as a-b with
// some special casing applied to the maximum unit64 value.
func (cr CommitRange) String() string {
	if cr.Maybe == math.MaxUint64 {
		if cr.Definitely == math.MaxUint64 {
			return "∞"
		}
		return fmt.Sprintf("%d-∞", cr.Definitely)
	}
	if cr.Definitely == cr.Maybe {
		return strconv.FormatUint(cr.Definitely, 10)
	}
	return fmt.Sprintf("%d-%d", cr.Definitely, cr.Maybe)
}

// VoteResult indicates the outcome of a vote.
//
//go:generate stringer -type=VoteResult
type VoteResult uint8

const (
	// VotePending indicates that the decision of the vote depends on future
	// votes, i.e. neither "yes" or "no" has reached quorum yet.
	VotePending VoteResult = 1 + iota
	// VoteLost indicates that the quorum has voted "no".
	VoteLost
	// VoteWon indicates that the quorum has voted "yes".
	VoteWon
)
