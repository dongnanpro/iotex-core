// Copyright (c) 2018 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package rolldpos

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/iotexproject/iotex-core/consensus/fsm"
	"github.com/iotexproject/iotex-core/proto"
)

func TestRuleProposeErrorVoteNil(t *testing.T) {
	t.Parallel()

	cb := rollDPoSCB{
		voteCb: func(msg proto.Message) error {
			vc, ok := (msg).(*iproto.ViewChangeMsg)
			assert.True(t, ok)
			assert.Nil(t, vc.Block)
			assert.Equal(t, vc.Vctype, iproto.ViewChangeMsg_PROPOSE)
			return nil
		},
	}
	h := rulePropose{
		RollDPoS: &RollDPoS{
			rollDPoSCB: cb,
		},
	}
	assert.True(t, h.Condition(&fsm.Event{Err: errors.New("err")}))
}
