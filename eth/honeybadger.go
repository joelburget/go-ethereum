package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"net"
	"os"
)

func FromEnvironmentOrNil(name string) (net.Conn, net.Conn) {
	socketPath := os.Getenv(name)
	if socketPath == "" {
		return nil, nil
	}

	myAddr, err := net.ResolveUnixAddr("unixgram", "/tmp/gethsock")
	if err != nil {
		panic(fmt.Sprintf("MustNew error: %v", err))
	}

	listener, err := net.ListenUnixgram("unixgram", myAddr)
	if err != nil {
		panic(fmt.Sprintf("MustNew error: %v", err))
	}

	sender, err := net.Dial("unixgram", socketPath)
	if err != nil {
		panic(fmt.Sprintf("MustNew error: %v", err))
	}

	return sender, listener
}

var Sender, Listener = FromEnvironmentOrNil("SOCKET_PATH")

func SendTxes(txes []*types.Transaction) {
	log.Info("sending txes", "txes", txes)
	encoding, err := rlp.EncodeToBytes(txes)
	if err != nil {
		panic(fmt.Sprintf("SendTxes error: %v", err))
	}
	Sender.Write(encoding)
}

// block until honeybadger gives us a block
func ReceiveTxes() []*types.Transaction {
	log.Info("receiving txes")
	txes := make([]*types.Transaction, 0)
	err := rlp.Decode(Listener, &txes)
	if err != nil {
		panic(fmt.Sprintf("ReceiveTxes error: %v", err))
	}
	return txes
}
