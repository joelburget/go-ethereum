package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"net"
	"os"
)

func FromEnvironmentOrNil(badgerSock string, gethSock string) (net.Conn, net.Conn) {
	badgerSockPath := os.Getenv(badgerSock)
	gethSockPath := os.Getenv(gethSock)
	if badgerSockPath == "" || gethSockPath == "" {
		return nil, nil
	}

	myAddr, err := net.ResolveUnixAddr("unixgram", gethSockPath)
	if err != nil {
		panic(fmt.Sprintf("MustNew error: %v", err))
	}

	badgerAddr, err := net.ResolveUnixAddr("unixgram", badgerSockPath)
	if err != nil {
		panic(fmt.Sprintf("MustNew error: %v", err))
	}

	sender, err := net.DialUnix("unixgram", myAddr, badgerAddr)
	if err != nil {
		panic(fmt.Sprintf("MustNew error: %v", err))
	}

	// send an empty list to establish a connection with honeybadger
	encoding, err := rlp.EncodeToBytes(&types.Transactions{})
	if err != nil {
		panic(fmt.Sprintf("SendTxes error: %v", err))
	}
	sender.Write(encoding)

	return sender, sender
}

var Sender, Listener = FromEnvironmentOrNil("BADGER_SOCK", "GETH_SOCK")

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
