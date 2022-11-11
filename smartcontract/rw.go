// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package smartcontract

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// RwMetaData contains all meta data concerning the Rw contract.
var RwMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_key\",\"type\":\"string\"}],\"name\":\"retrieve\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_key\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_value\",\"type\":\"string\"}],\"name\":\"store\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// RwABI is the input ABI used to generate the binding from.
// Deprecated: Use RwMetaData.ABI instead.
var RwABI = RwMetaData.ABI

// Rw is an auto generated Go binding around an Ethereum contract.
type Rw struct {
	RwCaller     // Read-only binding to the contract
	RwTransactor // Write-only binding to the contract
	RwFilterer   // Log filterer for contract events
}

// RwCaller is an auto generated read-only Go binding around an Ethereum contract.
type RwCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RwTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RwTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RwFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RwFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RwSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RwSession struct {
	Contract     *Rw               // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RwCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RwCallerSession struct {
	Contract *RwCaller     // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// RwTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RwTransactorSession struct {
	Contract     *RwTransactor     // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RwRaw is an auto generated low-level Go binding around an Ethereum contract.
type RwRaw struct {
	Contract *Rw // Generic contract binding to access the raw methods on
}

// RwCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RwCallerRaw struct {
	Contract *RwCaller // Generic read-only contract binding to access the raw methods on
}

// RwTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RwTransactorRaw struct {
	Contract *RwTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRw creates a new instance of Rw, bound to a specific deployed contract.
func NewRw(address common.Address, backend bind.ContractBackend) (*Rw, error) {
	contract, err := bindRw(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Rw{RwCaller: RwCaller{contract: contract}, RwTransactor: RwTransactor{contract: contract}, RwFilterer: RwFilterer{contract: contract}}, nil
}

// NewRwCaller creates a new read-only instance of Rw, bound to a specific deployed contract.
func NewRwCaller(address common.Address, caller bind.ContractCaller) (*RwCaller, error) {
	contract, err := bindRw(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RwCaller{contract: contract}, nil
}

// NewRwTransactor creates a new write-only instance of Rw, bound to a specific deployed contract.
func NewRwTransactor(address common.Address, transactor bind.ContractTransactor) (*RwTransactor, error) {
	contract, err := bindRw(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RwTransactor{contract: contract}, nil
}

// NewRwFilterer creates a new log filterer instance of Rw, bound to a specific deployed contract.
func NewRwFilterer(address common.Address, filterer bind.ContractFilterer) (*RwFilterer, error) {
	contract, err := bindRw(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RwFilterer{contract: contract}, nil
}

// bindRw binds a generic wrapper to an already deployed contract.
func bindRw(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RwMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rw *RwRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Rw.Contract.RwCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rw *RwRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rw.Contract.RwTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rw *RwRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rw.Contract.RwTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rw *RwCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Rw.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rw *RwTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rw.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rw *RwTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rw.Contract.contract.Transact(opts, method, params...)
}

// Retrieve is a free data retrieval call binding the contract method 0x64cc7327.
//
// Solidity: function retrieve(string _key) view returns(string, string)
func (_Rw *RwCaller) Retrieve(opts *bind.CallOpts, _key string) (string, string, error) {
	var out []interface{}
	err := _Rw.contract.Call(opts, &out, "retrieve", _key)

	if err != nil {
		return *new(string), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)

	return out0, out1, err

}

// Retrieve is a free data retrieval call binding the contract method 0x64cc7327.
//
// Solidity: function retrieve(string _key) view returns(string, string)
func (_Rw *RwSession) Retrieve(_key string) (string, string, error) {
	return _Rw.Contract.Retrieve(&_Rw.CallOpts, _key)
}

// Retrieve is a free data retrieval call binding the contract method 0x64cc7327.
//
// Solidity: function retrieve(string _key) view returns(string, string)
func (_Rw *RwCallerSession) Retrieve(_key string) (string, string, error) {
	return _Rw.Contract.Retrieve(&_Rw.CallOpts, _key)
}

// Store is a paid mutator transaction binding the contract method 0xf641090c.
//
// Solidity: function store(string _key, string _value) payable returns()
func (_Rw *RwTransactor) Store(opts *bind.TransactOpts, _key string, _value string) (*types.Transaction, error) {
	return _Rw.contract.Transact(opts, "store", _key, _value)
}

// Store is a paid mutator transaction binding the contract method 0xf641090c.
//
// Solidity: function store(string _key, string _value) payable returns()
func (_Rw *RwSession) Store(_key string, _value string) (*types.Transaction, error) {
	return _Rw.Contract.Store(&_Rw.TransactOpts, _key, _value)
}

// Store is a paid mutator transaction binding the contract method 0xf641090c.
//
// Solidity: function store(string _key, string _value) payable returns()
func (_Rw *RwTransactorSession) Store(_key string, _value string) (*types.Transaction, error) {
	return _Rw.Contract.Store(&_Rw.TransactOpts, _key, _value)
}
