// Copyright (c) 2018 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package rolldpos

import (
	"github.com/iotexproject/iotex-core/consensus/fsm"
	"github.com/iotexproject/iotex-core/delegate"
	"github.com/iotexproject/iotex-core/logger"
	"github.com/iotexproject/iotex-core/pkg/routine"
)

// proposerRotation is supposed to rotate the proposer per round of PBFT.
// However, use the first delegate as the proposer for now.
// We can propose based on the block height in the future.
type proposerRotation struct {
	*RollDPoS
}

// Handle handles transition to stateInitPropose
func (s *proposerRotation) Handle() {
	logger.Debug().Msg("determine if the node is the proposer")
	// If it's periodic proposer election on constant interval and the state is not ROUND_START, then returns
	if s.cfg.ProposerInterval != 0 && s.fsm.CurrentState() != stateRoundStart {
		return
	}
	height, err := s.bc.TipHeight()
	if err != nil {
		logger.Error().Err(err).Msg("failed to get blockchain height")
		return
	}
	if s.epochCtx == nil {
		logger.Error().Msg("epoch context is nil")
		return
	}
	pr, err := s.prCb(s.epochCtx.delegates, nil, 0, height+1)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get the proposer")
		return
	}
	// If proposer is not the current node, then returns
	if pr != s.self {
		logger.Info().
			Str("proposer", pr).
			Uint64("height", height+1).
			Msg("not the proposer")
		return
	}
	logger.Info().
		Str("proposer", pr).
		Uint64("height", height+1).
		Msg("propose new block height")

	s.enqueueEvent(&fsm.Event{
		State: stateInitPropose,
	})
}

// newProposerRotationNoDelay creates a ProposerRotation object
func newProposerRotationNoDelay(r *RollDPoS) *proposerRotation {
	return &proposerRotation{r}
}

// newProposerRotation creates a recurring task of proposer rotation.
func newProposerRotation(r *RollDPoS) *routine.RecurringTask {
	pr := &proposerRotation{r}
	return routine.NewRecurringTask(pr.Handle, r.cfg.ProposerInterval)
}

// FixedProposer will always choose the first in the delegate list as the proposer
func FixedProposer(delegates []string, _ []byte, _ uint64, _ uint64) (string, error) {
	if len(delegates) == 0 {
		return "", delegate.ErrZeroDelegate
	}
	return delegates[0], nil
}

// PseudoRotatedProposer will rotate among the delegates to choose the proposer
func PseudoRotatedProposer(delegates []string, _ []byte, _ uint64, height uint64) (string, error) {
	if len(delegates) == 0 {
		return "", delegate.ErrZeroDelegate
	}
	return delegates[height%uint64(len(delegates))], nil
}
