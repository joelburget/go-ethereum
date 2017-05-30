package eth

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/quorum"
	"github.com/ethereum/go-ethereum/rpc"
)

func (s *Ethereum) StartBlockVoting(client *rpc.Client, ks *keystore.KeyStore, voteAcct, makerAcct accounts.Account) error {
	blockMakerStrat := quorum.NewRandomDeadlineStrategy(s.eventMux, s.voteMinBlockTime, s.voteMaxBlockTime)
	quorum.Strategy = blockMakerStrat
	return s.blockVoting.Start(client, blockMakerStrat, ks, voteAcct, makerAcct)
}
