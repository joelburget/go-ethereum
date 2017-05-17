package logger

import (
	"github.com/ethereum/go-ethereum/logger/glog"
)

const (
	TxCreated  = "TX-CREATED"
	TxAccepted = "TX-ACCEPTED"

	BecameMinter   = "BECAME-MINTER"
	BecameVerifier = "BECAME-VERIFIER"

	// TODO: it would be trivial to log the role in BlockVoting.run if we want
	BlockVotingStarted = "BLOCK-VOTING-STARTED"
)

var DoEmitCheckpoints = false

func EmitCheckpoint(checkpointName string, logValues ...interface{}) {
	if DoEmitCheckpoints {
		glog.V(Info).Infof("QUORUM-CHECKPOINT %s %v\n", checkpointName, logValues)
	}
}
