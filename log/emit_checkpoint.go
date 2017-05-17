package log

const (
	TxCreated      = "TX-CREATED"
	TxAccepted     = "TX-ACCEPTED"
	BecameMinter   = "BECAME-MINTER"
	BecameVerifier = "BECAME-VERIFIER"
)

var DoEmitCheckpoints = false

func EmitCheckpoint(checkpointName string, logValues ...interface{}) {
	if DoEmitCheckpoints {
		Info("QUORUM-CHECKPOINT\n", "name", checkpointName, "data", logValues)
	}
}
