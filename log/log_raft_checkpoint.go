package log

const (
	TxCreated      = "TX-CREATED"
	TxAccepted     = "TX-ACCEPTED"
	BecameMinter   = "BECAME-MINTER"
	BecameVerifier = "BECAME-VERIFIER"
)

var DoLogRaft = false

func LogRaftCheckpoint(checkpointName string, logValues ...interface{}) {
	if DoLogRaft {
		Info("QUORUM-CHECKPOINT\n", "name", checkpointName, "data", logValues)
	}
}
