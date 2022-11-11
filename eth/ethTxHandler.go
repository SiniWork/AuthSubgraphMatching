package eth

import (
	"Corgi/smartcontract"
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wonderivan/logger"
	"io/ioutil"
	"math/big"
	"sync"
	//"fmt"
)

var initOnce sync.Once
// 这个是智能合约编译出来的golang文件中的，在编译时输入的-type的参数(一般为首字母大写)+Raw，建议自己去那个文件里搜一下，确定一下大小写
var ethTxRaw smartcontract.RwRaw
var client *ethclient.Client
const (
	//合约地址，在部署合约时，使用contract.address可以看到
	contractAddr = "0xadafdd9de51f9cc4026f7ffb2047d3bd43041ec7"
	//你的账户的密码
	passwd = "0000"
	// 下列两个参数组成的地址+文件是你的秘钥文件，地址需换成你所部署的以太坊/data/keystore地址，
	// filename为keystore文件夹下文件名的最后一段是你账户地址的文件。
	fileName = "UTC--2022-10-31T02-09-26.488472748Z--c16606a5bc2fbc4fbeb14d114fe17ae5257253ad"
	keystorePath = "/home/eth/go/src/ethereum-private/data/keystore/"
	//chainID = "ef872"
)

func init() {
	initOnce.Do(initClient)
}
//func main() {
//	CommitEth("abcde", "123")
//	result := QueryEth("abcde")
//	fmt.Println(result)
//}
 
func initClient() {
	var err error
	//自己的以太坊的ip及端口，端口在启动网络时设定
	client, err = ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("eth network connected")
	// 智能合约的地址
	address := common.HexToAddress(contractAddr)
	// 这两个函数需要去智能合约生成的golang文件中查看一下。
	// 一般情况下为"New+编译时输入的-type的参数"
	storeObj, err := smartcontract.NewRw(address, client)
	// 一般情况下为"编译时输入的-type的参数(一般为首字母大写)+Raw"
	ethTxRaw = smartcontract.RwRaw{storeObj}
}

func newAuth() *bind.TransactOpts {
	keyJson, err := ioutil.ReadFile(keystorePath + fileName)
	if err != nil {
		logger.Fatal(err)
	}
	key, err := keystore.DecryptKey(keyJson, passwd)
	if err != nil {
		logger.Fatal(err)
	}
	publicKey := key.PrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		logger.Fatal("error casting public key to ECDSA", err)
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		logger.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		logger.Fatal(err)
	}
	//auth := bind.NewKeyedTransactor(key.PrivateKey)
	// 搭建启动以太坊时确认的chainID
	chainID := big.NewInt(981106)
	//新版需要chainID，因此不再使用NewKeyedTransactor
	// 生成auth用于transact调用
	auth,err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, chainID)
	if err != nil {
		logger.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	// 发起事务如果出现错误可能是因为gas不足，设置大一点就好了
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	nonce1, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		logger.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce1))
	return auth
}


// 分别使用以下两个函数调用智能合约，其中只有调用Call、Transact的区别
// 这两个函数实际上不以是读还是写来进行区分，而是在编写智能合约.sol时，你在编写对应函数时，是否给函数加了view字段来区分
// 如果你添加了view，则你可以在生成的.go文件中，同名函数前发现 *** is a free data retrieval call binding the contract method注释，则你可以使用Call
// 如果没添加，可以发现 *** is a paid mutator transaction binding the contract method，则可以使用Transact

func QueryEth(keys string) interface{} {
	var result []interface{}
	err := ethTxRaw.Call(nil, &result, "retrieve", keys)
	if err != nil {
		logger.Error("call query of eth failed", err)
	}
	logger.Info("query done")
	return result
}

func CommitEth(keys string, values string) error {
	auth := newAuth()
	_, err := ethTxRaw.Transact(auth, "store", keys, values)
	if err != nil {
		logger.Fatal(err)
		return err
	}
	logger.Info("Write to blockChain success", keys)
	return nil
}
