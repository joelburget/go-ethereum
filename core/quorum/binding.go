// This file is an automatically generated Go binding. Do not modify as any
// change will likely be lost upon the next re-generation!

package quorum

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// BlockVotingABI is the input ABI used to generate the binding from.
const BlockVotingABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"setVoteThreshold\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"removeBlockMaker\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"voterCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"canCreateBlocks\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"voteThreshold\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"height\",\"type\":\"uint256\"}],\"name\":\"getCanonHash\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"height\",\"type\":\"uint256\"},{\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"vote\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"addBlockMaker\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"removeVoter\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"height\",\"type\":\"uint256\"},{\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getEntry\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isVoter\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"canVote\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"blockMakerCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getSize\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isBlockMaker\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"addVoter\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"name\":\"Vote\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"\",\"type\":\"address\"}],\"name\":\"AddVoter\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"\",\"type\":\"address\"}],\"name\":\"RemovedVoter\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"\",\"type\":\"address\"}],\"name\":\"AddBlockMaker\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"\",\"type\":\"address\"}],\"name\":\"RemovedBlockMaker\",\"type\":\"event\"}]"

// BlockVotingBin is the compiled bytecode used for deploying new contracts.
const BlockVotingBin = `0x6060604052341561000c57fe5b5b610a458061001c6000396000f300606060405236156100e35763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631290948581146100e5578063284d163c146100fa57806342169e4814610118578063488099a61461013a5780634fe437d51461016a578063559c390c1461018c57806368bb8bb6146101b157806372a571fc146101c957806386c1ff68146101e757806398ba676d14610205578063a7771ee31461022d578063adfaa72e1461025d578063cf5289851461028d578063de8fa431146102af578063e814d1c7146102d1578063f4ab9adf14610301575bfe5b34156100ed57fe5b6100f860043561031f565b005b341561010257fe5b6100f8600160a060020a0360043516610355565b005b341561012057fe5b610128610418565b60408051918252519081900360200190f35b341561014257fe5b610156600160a060020a036004351661041e565b604080519115158252519081900360200190f35b341561017257fe5b610128610433565b60408051918252519081900360200190f35b341561019457fe5b610128600435610439565b60408051918252519081900360200190f35b34156101b957fe5b6100f860043560243561054f565b005b34156101d157fe5b6100f8600160a060020a036004351661065f565b005b34156101ef57fe5b6100f8600160a060020a0360043516610717565b005b341561020d57fe5b6101286004356024356107da565b60408051918252519081900360200190f35b341561023557fe5b610156600160a060020a036004351661082f565b604080519115158252519081900360200190f35b341561026557fe5b610156600160a060020a0360043516610851565b604080519115158252519081900360200190f35b341561029557fe5b610128610866565b60408051918252519081900360200190f35b34156102b757fe5b61012861086c565b60408051918252519081900360200190f35b34156102d957fe5b610156600160a060020a0360043516610873565b604080519115158252519081900360200190f35b341561030957fe5b6100f8600160a060020a0360043516610895565b005b600160a060020a03331660009081526003602052604090205460ff161561034b5760018190555b610351565b60006000fd5b5b50565b600160a060020a03331660009081526005602052604090205460ff161561034b57600454600114156103875760006000fd5b600160a060020a03811660009081526005602052604090205460ff161561034657600160a060020a038116600081815260056020908152604091829020805460ff1916905560048054600019019055815192835290517f8cee3054364d6799f1c8962580ad61273d9d38ca1ff26516bd1ad23c099a60229281900390910190a15b5b610351565b60006000fd5b5b50565b60025481565b60056020526000908152604090205460ff1681565b60015481565b7f6a6f656c0000000000000000000000000000000000000000000000000000000060008080610547565b906000526020600020906002020160005b509250600090505b60018301548110156105435760018301805484916000918490811061049d57fe5b906000526020600020900160005b505481526020808201929092526040908101600090812054858252928690522054108015610512575060015483600001600085600101848154811015156104ee57fe5b906000526020600020900160005b5054815260208101919091526040016000205410155b1561053a576001830180548290811061052757fe5b906000526020600020900160005b505491505b5b60010161047c565b8193505b505050919050565b600160a060020a03331660009081526003602052604081205460ff161561034b57600054839010156105905760008054808503019061058e908261094d565b505b6000805460001985019081106105a257fe5b906000526020600020906002020160005b5060008381526020829052604090205490915015156105f6578060010180548060010182816105e2919061097f565b916000526020600020900160005b50839055505b600082815260208281526040918290208054600101905581514381529081018490528151600160a060020a033316927f3d03ba7f4b5227cdb385f2610906e5bcee147171603ec40005b30915ad20e258928290030190a25b610659565b60006000fd5b5b505050565b600160a060020a03331660009081526005602052604090205460ff161561034b57600160a060020a03811660009081526005602052604090205460ff16151561034657600160a060020a038116600081815260056020908152604091829020805460ff19166001908117909155600480549091019055815192835290517f1a4ce6942f7aa91856332e618fc90159f13a340611a308f5d7327ba0707e56859281900390910190a15b5b610351565b60006000fd5b5b50565b600160a060020a03331660009081526003602052604090205460ff161561034b57600254600114156107495760006000fd5b600160a060020a03811660009081526003602052604090205460ff161561034657600160a060020a038116600081815260036020908152604091829020805460ff1916905560028054600019019055815192835290517f183393fc5cffbfc7d03d623966b85f76b9430f42d3aada2ac3f3deabc78899e89281900390910190a15b5b610351565b60006000fd5b5b50565b600060006000600185038154811015156107f057fe5b906000526020600020906002020160005b509050806001018381548110151561081557fe5b906000526020600020900160005b505491505b5092915050565b600160a060020a03811660009081526003602052604090205460ff165b919050565b60036020526000908152604090205460ff1681565b60045481565b6000545b90565b600160a060020a03811660009081526005602052604090205460ff165b919050565b600160a060020a03331660009081526003602052604090205460ff161561034b57600160a060020a03811660009081526003602052604090205460ff16151561034657600160a060020a038116600081815260036020908152604091829020805460ff19166001908117909155600280549091019055815192835290517f0ad2eca75347acd5160276fe4b5dad46987e4ff4af9e574195e3e9bc15d7e0ff9281900390910190a15b5b610351565b60006000fd5b5b50565b8154818355818115116106595760020281600202836000526020600020918201910161065991906109a9565b5b505050565b815481835581811511610659576000838152602090206106599181019083016109d6565b5b505050565b61087091905b808211156109cf5760006109c660018301826109f7565b506002016109af565b5090565b90565b61087091905b808211156109cf57600081556001016109dc565b5090565b90565b508054600082559060005260206000209081019061035191906109d6565b5b505600a165627a7a72305820e88a43d2597ba7874e01e702e7dd35325871f33fe65353e6f2d379960cd9c3390029`

// DeployBlockVoting deploys a new Ethereum contract, binding an instance of BlockVoting to it.
func DeployBlockVoting(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BlockVoting, error) {
	parsed, err := abi.JSON(strings.NewReader(BlockVotingABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BlockVotingBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BlockVoting{BlockVotingCaller: BlockVotingCaller{contract: contract}, BlockVotingTransactor: BlockVotingTransactor{contract: contract}}, nil
}

// BlockVoting is an auto generated Go binding around an Ethereum contract.
type BlockVoting struct {
	BlockVotingCaller     // Read-only binding to the contract
	BlockVotingTransactor // Write-only binding to the contract
}

// BlockVotingCaller is an auto generated read-only Go binding around an Ethereum contract.
type BlockVotingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockVotingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BlockVotingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockVotingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BlockVotingSession struct {
	Contract     *BlockVoting      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BlockVotingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BlockVotingCallerSession struct {
	Contract *BlockVotingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// BlockVotingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BlockVotingTransactorSession struct {
	Contract     *BlockVotingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// BlockVotingRaw is an auto generated low-level Go binding around an Ethereum contract.
type BlockVotingRaw struct {
	Contract *BlockVoting // Generic contract binding to access the raw methods on
}

// BlockVotingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BlockVotingCallerRaw struct {
	Contract *BlockVotingCaller // Generic read-only contract binding to access the raw methods on
}

// BlockVotingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BlockVotingTransactorRaw struct {
	Contract *BlockVotingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBlockVoting creates a new instance of BlockVoting, bound to a specific deployed contract.
func NewBlockVoting(address common.Address, backend bind.ContractBackend) (*BlockVoting, error) {
	contract, err := bindBlockVoting(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BlockVoting{BlockVotingCaller: BlockVotingCaller{contract: contract}, BlockVotingTransactor: BlockVotingTransactor{contract: contract}}, nil
}

// NewBlockVotingCaller creates a new read-only instance of BlockVoting, bound to a specific deployed contract.
func NewBlockVotingCaller(address common.Address, caller bind.ContractCaller) (*BlockVotingCaller, error) {
	contract, err := bindBlockVoting(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &BlockVotingCaller{contract: contract}, nil
}

// NewBlockVotingTransactor creates a new write-only instance of BlockVoting, bound to a specific deployed contract.
func NewBlockVotingTransactor(address common.Address, transactor bind.ContractTransactor) (*BlockVotingTransactor, error) {
	contract, err := bindBlockVoting(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &BlockVotingTransactor{contract: contract}, nil
}

// bindBlockVoting binds a generic wrapper to an already deployed contract.
func bindBlockVoting(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BlockVotingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BlockVoting *BlockVotingRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BlockVoting.Contract.BlockVotingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BlockVoting *BlockVotingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BlockVoting.Contract.BlockVotingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BlockVoting *BlockVotingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BlockVoting.Contract.BlockVotingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BlockVoting *BlockVotingCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BlockVoting.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BlockVoting *BlockVotingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BlockVoting.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BlockVoting *BlockVotingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BlockVoting.Contract.contract.Transact(opts, method, params...)
}

// BlockMakerCount is a free data retrieval call binding the contract method 0xcf528985.
//
// Solidity: function blockMakerCount() constant returns(uint256)
func (_BlockVoting *BlockVotingCaller) BlockMakerCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BlockVoting.contract.Call(opts, out, "blockMakerCount")
	return *ret0, err
}

// BlockMakerCount is a free data retrieval call binding the contract method 0xcf528985.
//
// Solidity: function blockMakerCount() constant returns(uint256)
func (_BlockVoting *BlockVotingSession) BlockMakerCount() (*big.Int, error) {
	return _BlockVoting.Contract.BlockMakerCount(&_BlockVoting.CallOpts)
}

// BlockMakerCount is a free data retrieval call binding the contract method 0xcf528985.
//
// Solidity: function blockMakerCount() constant returns(uint256)
func (_BlockVoting *BlockVotingCallerSession) BlockMakerCount() (*big.Int, error) {
	return _BlockVoting.Contract.BlockMakerCount(&_BlockVoting.CallOpts)
}

// CanCreateBlocks is a free data retrieval call binding the contract method 0x488099a6.
//
// Solidity: function canCreateBlocks( address) constant returns(bool)
func (_BlockVoting *BlockVotingCaller) CanCreateBlocks(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _BlockVoting.contract.Call(opts, out, "canCreateBlocks", arg0)
	return *ret0, err
}

// CanCreateBlocks is a free data retrieval call binding the contract method 0x488099a6.
//
// Solidity: function canCreateBlocks( address) constant returns(bool)
func (_BlockVoting *BlockVotingSession) CanCreateBlocks(arg0 common.Address) (bool, error) {
	return _BlockVoting.Contract.CanCreateBlocks(&_BlockVoting.CallOpts, arg0)
}

// CanCreateBlocks is a free data retrieval call binding the contract method 0x488099a6.
//
// Solidity: function canCreateBlocks( address) constant returns(bool)
func (_BlockVoting *BlockVotingCallerSession) CanCreateBlocks(arg0 common.Address) (bool, error) {
	return _BlockVoting.Contract.CanCreateBlocks(&_BlockVoting.CallOpts, arg0)
}

// CanVote is a free data retrieval call binding the contract method 0xadfaa72e.
//
// Solidity: function canVote( address) constant returns(bool)
func (_BlockVoting *BlockVotingCaller) CanVote(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _BlockVoting.contract.Call(opts, out, "canVote", arg0)
	return *ret0, err
}

// CanVote is a free data retrieval call binding the contract method 0xadfaa72e.
//
// Solidity: function canVote( address) constant returns(bool)
func (_BlockVoting *BlockVotingSession) CanVote(arg0 common.Address) (bool, error) {
	return _BlockVoting.Contract.CanVote(&_BlockVoting.CallOpts, arg0)
}

// CanVote is a free data retrieval call binding the contract method 0xadfaa72e.
//
// Solidity: function canVote( address) constant returns(bool)
func (_BlockVoting *BlockVotingCallerSession) CanVote(arg0 common.Address) (bool, error) {
	return _BlockVoting.Contract.CanVote(&_BlockVoting.CallOpts, arg0)
}

// GetCanonHash is a free data retrieval call binding the contract method 0x559c390c.
//
// Solidity: function getCanonHash(height uint256) constant returns(bytes32)
func (_BlockVoting *BlockVotingCaller) GetCanonHash(opts *bind.CallOpts, height *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _BlockVoting.contract.Call(opts, out, "getCanonHash", height)
	return *ret0, err
}

// GetCanonHash is a free data retrieval call binding the contract method 0x559c390c.
//
// Solidity: function getCanonHash(height uint256) constant returns(bytes32)
func (_BlockVoting *BlockVotingSession) GetCanonHash(height *big.Int) ([32]byte, error) {
	return _BlockVoting.Contract.GetCanonHash(&_BlockVoting.CallOpts, height)
}

// GetCanonHash is a free data retrieval call binding the contract method 0x559c390c.
//
// Solidity: function getCanonHash(height uint256) constant returns(bytes32)
func (_BlockVoting *BlockVotingCallerSession) GetCanonHash(height *big.Int) ([32]byte, error) {
	return _BlockVoting.Contract.GetCanonHash(&_BlockVoting.CallOpts, height)
}

// GetEntry is a free data retrieval call binding the contract method 0x98ba676d.
//
// Solidity: function getEntry(height uint256, n uint256) constant returns(bytes32)
func (_BlockVoting *BlockVotingCaller) GetEntry(opts *bind.CallOpts, height *big.Int, n *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _BlockVoting.contract.Call(opts, out, "getEntry", height, n)
	return *ret0, err
}

// GetEntry is a free data retrieval call binding the contract method 0x98ba676d.
//
// Solidity: function getEntry(height uint256, n uint256) constant returns(bytes32)
func (_BlockVoting *BlockVotingSession) GetEntry(height *big.Int, n *big.Int) ([32]byte, error) {
	return _BlockVoting.Contract.GetEntry(&_BlockVoting.CallOpts, height, n)
}

// GetEntry is a free data retrieval call binding the contract method 0x98ba676d.
//
// Solidity: function getEntry(height uint256, n uint256) constant returns(bytes32)
func (_BlockVoting *BlockVotingCallerSession) GetEntry(height *big.Int, n *big.Int) ([32]byte, error) {
	return _BlockVoting.Contract.GetEntry(&_BlockVoting.CallOpts, height, n)
}

// GetSize is a free data retrieval call binding the contract method 0xde8fa431.
//
// Solidity: function getSize() constant returns(uint256)
func (_BlockVoting *BlockVotingCaller) GetSize(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BlockVoting.contract.Call(opts, out, "getSize")
	return *ret0, err
}

// GetSize is a free data retrieval call binding the contract method 0xde8fa431.
//
// Solidity: function getSize() constant returns(uint256)
func (_BlockVoting *BlockVotingSession) GetSize() (*big.Int, error) {
	return _BlockVoting.Contract.GetSize(&_BlockVoting.CallOpts)
}

// GetSize is a free data retrieval call binding the contract method 0xde8fa431.
//
// Solidity: function getSize() constant returns(uint256)
func (_BlockVoting *BlockVotingCallerSession) GetSize() (*big.Int, error) {
	return _BlockVoting.Contract.GetSize(&_BlockVoting.CallOpts)
}

// IsBlockMaker is a free data retrieval call binding the contract method 0xe814d1c7.
//
// Solidity: function isBlockMaker(addr address) constant returns(bool)
func (_BlockVoting *BlockVotingCaller) IsBlockMaker(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _BlockVoting.contract.Call(opts, out, "isBlockMaker", addr)
	return *ret0, err
}

// IsBlockMaker is a free data retrieval call binding the contract method 0xe814d1c7.
//
// Solidity: function isBlockMaker(addr address) constant returns(bool)
func (_BlockVoting *BlockVotingSession) IsBlockMaker(addr common.Address) (bool, error) {
	return _BlockVoting.Contract.IsBlockMaker(&_BlockVoting.CallOpts, addr)
}

// IsBlockMaker is a free data retrieval call binding the contract method 0xe814d1c7.
//
// Solidity: function isBlockMaker(addr address) constant returns(bool)
func (_BlockVoting *BlockVotingCallerSession) IsBlockMaker(addr common.Address) (bool, error) {
	return _BlockVoting.Contract.IsBlockMaker(&_BlockVoting.CallOpts, addr)
}

// IsVoter is a free data retrieval call binding the contract method 0xa7771ee3.
//
// Solidity: function isVoter(addr address) constant returns(bool)
func (_BlockVoting *BlockVotingCaller) IsVoter(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _BlockVoting.contract.Call(opts, out, "isVoter", addr)
	return *ret0, err
}

// IsVoter is a free data retrieval call binding the contract method 0xa7771ee3.
//
// Solidity: function isVoter(addr address) constant returns(bool)
func (_BlockVoting *BlockVotingSession) IsVoter(addr common.Address) (bool, error) {
	return _BlockVoting.Contract.IsVoter(&_BlockVoting.CallOpts, addr)
}

// IsVoter is a free data retrieval call binding the contract method 0xa7771ee3.
//
// Solidity: function isVoter(addr address) constant returns(bool)
func (_BlockVoting *BlockVotingCallerSession) IsVoter(addr common.Address) (bool, error) {
	return _BlockVoting.Contract.IsVoter(&_BlockVoting.CallOpts, addr)
}

// VoteThreshold is a free data retrieval call binding the contract method 0x4fe437d5.
//
// Solidity: function voteThreshold() constant returns(uint256)
func (_BlockVoting *BlockVotingCaller) VoteThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BlockVoting.contract.Call(opts, out, "voteThreshold")
	return *ret0, err
}

// VoteThreshold is a free data retrieval call binding the contract method 0x4fe437d5.
//
// Solidity: function voteThreshold() constant returns(uint256)
func (_BlockVoting *BlockVotingSession) VoteThreshold() (*big.Int, error) {
	return _BlockVoting.Contract.VoteThreshold(&_BlockVoting.CallOpts)
}

// VoteThreshold is a free data retrieval call binding the contract method 0x4fe437d5.
//
// Solidity: function voteThreshold() constant returns(uint256)
func (_BlockVoting *BlockVotingCallerSession) VoteThreshold() (*big.Int, error) {
	return _BlockVoting.Contract.VoteThreshold(&_BlockVoting.CallOpts)
}

// VoterCount is a free data retrieval call binding the contract method 0x42169e48.
//
// Solidity: function voterCount() constant returns(uint256)
func (_BlockVoting *BlockVotingCaller) VoterCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BlockVoting.contract.Call(opts, out, "voterCount")
	return *ret0, err
}

// VoterCount is a free data retrieval call binding the contract method 0x42169e48.
//
// Solidity: function voterCount() constant returns(uint256)
func (_BlockVoting *BlockVotingSession) VoterCount() (*big.Int, error) {
	return _BlockVoting.Contract.VoterCount(&_BlockVoting.CallOpts)
}

// VoterCount is a free data retrieval call binding the contract method 0x42169e48.
//
// Solidity: function voterCount() constant returns(uint256)
func (_BlockVoting *BlockVotingCallerSession) VoterCount() (*big.Int, error) {
	return _BlockVoting.Contract.VoterCount(&_BlockVoting.CallOpts)
}

// AddBlockMaker is a paid mutator transaction binding the contract method 0x72a571fc.
//
// Solidity: function addBlockMaker(addr address) returns()
func (_BlockVoting *BlockVotingTransactor) AddBlockMaker(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _BlockVoting.contract.Transact(opts, "addBlockMaker", addr)
}

// AddBlockMaker is a paid mutator transaction binding the contract method 0x72a571fc.
//
// Solidity: function addBlockMaker(addr address) returns()
func (_BlockVoting *BlockVotingSession) AddBlockMaker(addr common.Address) (*types.Transaction, error) {
	return _BlockVoting.Contract.AddBlockMaker(&_BlockVoting.TransactOpts, addr)
}

// AddBlockMaker is a paid mutator transaction binding the contract method 0x72a571fc.
//
// Solidity: function addBlockMaker(addr address) returns()
func (_BlockVoting *BlockVotingTransactorSession) AddBlockMaker(addr common.Address) (*types.Transaction, error) {
	return _BlockVoting.Contract.AddBlockMaker(&_BlockVoting.TransactOpts, addr)
}

// AddVoter is a paid mutator transaction binding the contract method 0xf4ab9adf.
//
// Solidity: function addVoter(addr address) returns()
func (_BlockVoting *BlockVotingTransactor) AddVoter(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _BlockVoting.contract.Transact(opts, "addVoter", addr)
}

// AddVoter is a paid mutator transaction binding the contract method 0xf4ab9adf.
//
// Solidity: function addVoter(addr address) returns()
func (_BlockVoting *BlockVotingSession) AddVoter(addr common.Address) (*types.Transaction, error) {
	return _BlockVoting.Contract.AddVoter(&_BlockVoting.TransactOpts, addr)
}

// AddVoter is a paid mutator transaction binding the contract method 0xf4ab9adf.
//
// Solidity: function addVoter(addr address) returns()
func (_BlockVoting *BlockVotingTransactorSession) AddVoter(addr common.Address) (*types.Transaction, error) {
	return _BlockVoting.Contract.AddVoter(&_BlockVoting.TransactOpts, addr)
}

// RemoveBlockMaker is a paid mutator transaction binding the contract method 0x284d163c.
//
// Solidity: function removeBlockMaker(addr address) returns()
func (_BlockVoting *BlockVotingTransactor) RemoveBlockMaker(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _BlockVoting.contract.Transact(opts, "removeBlockMaker", addr)
}

// RemoveBlockMaker is a paid mutator transaction binding the contract method 0x284d163c.
//
// Solidity: function removeBlockMaker(addr address) returns()
func (_BlockVoting *BlockVotingSession) RemoveBlockMaker(addr common.Address) (*types.Transaction, error) {
	return _BlockVoting.Contract.RemoveBlockMaker(&_BlockVoting.TransactOpts, addr)
}

// RemoveBlockMaker is a paid mutator transaction binding the contract method 0x284d163c.
//
// Solidity: function removeBlockMaker(addr address) returns()
func (_BlockVoting *BlockVotingTransactorSession) RemoveBlockMaker(addr common.Address) (*types.Transaction, error) {
	return _BlockVoting.Contract.RemoveBlockMaker(&_BlockVoting.TransactOpts, addr)
}

// RemoveVoter is a paid mutator transaction binding the contract method 0x86c1ff68.
//
// Solidity: function removeVoter(addr address) returns()
func (_BlockVoting *BlockVotingTransactor) RemoveVoter(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _BlockVoting.contract.Transact(opts, "removeVoter", addr)
}

// RemoveVoter is a paid mutator transaction binding the contract method 0x86c1ff68.
//
// Solidity: function removeVoter(addr address) returns()
func (_BlockVoting *BlockVotingSession) RemoveVoter(addr common.Address) (*types.Transaction, error) {
	return _BlockVoting.Contract.RemoveVoter(&_BlockVoting.TransactOpts, addr)
}

// RemoveVoter is a paid mutator transaction binding the contract method 0x86c1ff68.
//
// Solidity: function removeVoter(addr address) returns()
func (_BlockVoting *BlockVotingTransactorSession) RemoveVoter(addr common.Address) (*types.Transaction, error) {
	return _BlockVoting.Contract.RemoveVoter(&_BlockVoting.TransactOpts, addr)
}

// SetVoteThreshold is a paid mutator transaction binding the contract method 0x12909485.
//
// Solidity: function setVoteThreshold(threshold uint256) returns()
func (_BlockVoting *BlockVotingTransactor) SetVoteThreshold(opts *bind.TransactOpts, threshold *big.Int) (*types.Transaction, error) {
	return _BlockVoting.contract.Transact(opts, "setVoteThreshold", threshold)
}

// SetVoteThreshold is a paid mutator transaction binding the contract method 0x12909485.
//
// Solidity: function setVoteThreshold(threshold uint256) returns()
func (_BlockVoting *BlockVotingSession) SetVoteThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _BlockVoting.Contract.SetVoteThreshold(&_BlockVoting.TransactOpts, threshold)
}

// SetVoteThreshold is a paid mutator transaction binding the contract method 0x12909485.
//
// Solidity: function setVoteThreshold(threshold uint256) returns()
func (_BlockVoting *BlockVotingTransactorSession) SetVoteThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _BlockVoting.Contract.SetVoteThreshold(&_BlockVoting.TransactOpts, threshold)
}

// Vote is a paid mutator transaction binding the contract method 0x68bb8bb6.
//
// Solidity: function vote(height uint256, hash bytes32) returns()
func (_BlockVoting *BlockVotingTransactor) Vote(opts *bind.TransactOpts, height *big.Int, hash [32]byte) (*types.Transaction, error) {
	return _BlockVoting.contract.Transact(opts, "vote", height, hash)
}

// Vote is a paid mutator transaction binding the contract method 0x68bb8bb6.
//
// Solidity: function vote(height uint256, hash bytes32) returns()
func (_BlockVoting *BlockVotingSession) Vote(height *big.Int, hash [32]byte) (*types.Transaction, error) {
	return _BlockVoting.Contract.Vote(&_BlockVoting.TransactOpts, height, hash)
}

// Vote is a paid mutator transaction binding the contract method 0x68bb8bb6.
//
// Solidity: function vote(height uint256, hash bytes32) returns()
func (_BlockVoting *BlockVotingTransactorSession) Vote(height *big.Int, hash [32]byte) (*types.Transaction, error) {
	return _BlockVoting.Contract.Vote(&_BlockVoting.TransactOpts, height, hash)
}
