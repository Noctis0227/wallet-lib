package app

import (
	"flag"
	"fmt"
	"github.com/Noctis0227/wallet-lib/conf"
	"github.com/Noctis0227/wallet-lib/core"
	"github.com/Noctis0227/wallet-lib/rpc"
	"github.com/Noctis0227/wallet-lib/sign"
	"os"
)

type command interface {
	name() string
	handle()
	parse() bool
	usage()
}

type txCommand struct {
	*flag.FlagSet
	find    bool
	trading bool
	raw     string
	txId    string
	sign    bool
	code    string
	key     string
}

type addrCommand struct {
	*flag.FlagSet
	generate bool
	list     bool
	def      string
}

type syncCommand struct {
	*flag.FlagSet
}

var cmdList []command
var startIdx int

func init() {
	cmdList = make([]command, 0)
	cmdList = append(cmdList, &txCommand{FlagSet: flag.NewFlagSet("tx", flag.ExitOnError)})
	cmdList = append(cmdList, &addrCommand{FlagSet: flag.NewFlagSet("addr", flag.ExitOnError)})
}

func StartCli(cmd string) {
	var sel command
	for _, c := range cmdList {
		if c.name() == cmd {
			sel = c
			break
		}
	}
	if sel != nil {
		startIdx = 2
		println()
		sel.handle()
	} else {
		CliUsage(2)
	}
}

/*begin tx */
func (cmd *txCommand) name() string {
	return cmd.Name()
}
func (cmd *txCommand) parse() bool {
	cmd.BoolVar(&cmd.find, "f", false, "query a transaction by id")
	cmd.BoolVar(&cmd.trading, "t", false, "send the raw tx to the network")
	cmd.BoolVar(&cmd.sign, "s", false, "sign transaction")
	cmd.StringVar(&cmd.txId, "id", "", "the id of transaction")
	cmd.StringVar(&cmd.raw, "rawtx", "", "the raw tx of transaction")
	cmd.StringVar(&cmd.code, "code", "", "the transaction code")
	cmd.StringVar(&cmd.key, "key", "", "the private key")
	args := os.Args[startIdx:]
	err := cmd.Parse(args)
	if len(args) == 0 || err != nil {
		return false
	}

	//fmt.Printf("record: %t\nsend: %t\nfrom: %s\namount: %f\n", cmd.record, cmd.send, cmd.from, cmd.amount)

	return true
}
func (cmd *txCommand) usage() {
	cmd.Usage()
}

func (cmd *txCommand) handle() {
	if !cmd.parse() {
		cmd.usage()
		return
	}
	rpcClient := rpc.NewClient(&rpc.RpcConfig{
		Address: conf.Setting.Rpc.Host,
		User:    conf.Setting.Rpc.User,
		Pwd:     conf.Setting.Rpc.Pwd,
	})
	//fmt.Printf("r: %t\ns: %t\nfrom: %s\namount: %d\n", *latest, *send, *from, *amount)
	switch true {
	case cmd.trading:
		println("send raw tx ", cmd.raw)
		rs, err := rpcClient.SendTransaction(cmd.raw)
		if err == nil {
			println("send raw tx success, txid:", rs)
		} else {
			println("send raw tx failed!", rs)
		}
	case cmd.find:
		println("query a transaction")
		txs, err := rpcClient.GetTransaction(cmd.txId)
		if err != nil {
			println(err.Error())
			return
		}
		printTransaction(txs)
	case cmd.sign:
		rs, err := sign.TxSign(cmd.code, cmd.key, conf.Setting.Version)
		fmt.Println(rs, err)

	}
}

/*end tx*/

/*begin from*/
func (cmd *addrCommand) name() string {
	return cmd.Name()
}
func (cmd *addrCommand) parse() bool {
	cmd.BoolVar(&cmd.generate, "g", false, "generate a wallet address")
	cmd.BoolVar(&cmd.list, "l", false, "list the wallet address")
	cmd.StringVar(&cmd.def, "def", "", "get or set a default address, If it is empty, then get default address.")
	args := os.Args[startIdx:]
	err := cmd.Parse(args)
	if len(args) == 0 || err != nil {
		return false
	}
	return true
}
func (cmd *addrCommand) usage() {
	cmd.Usage()
}
func (cmd *addrCommand) handle() {
	if !cmd.parse() {
		cmd.usage()
		return
	}
	//fmt.Printf("r: %t\ns: %t\nfrom: %s\namount: %d\n", *latest, *send, *from, *amount)
	switch true {
	case cmd.generate:
		println("generate a address")

		println(core.GenerateAddr())
	case cmd.list:
		if rs := core.GetAddrs(); rs != nil {
			for k, v := range rs {
				println(k, ",", v)
			}
		}
	case cmd.def == "":
		println(core.GetDefAddr())
	case cmd.def != "":
		if err := core.SetDefAddr(cmd.def); err != nil {
			println(err)
		}
		println("set default address success")
	}
}

/*end from*/

func CliUsage(start int) {
	startIdx = start
	for _, c := range cmdList {
		c.parse()
		c.usage()
	}
}

func printTransaction(transaction *rpc.Transaction) {
	fmt.Println("result:{")
	fmt.Println(" hex:", transaction.Hex)
	fmt.Println(" hexnowit:", transaction.Hexnowit)
	fmt.Println(" hexwit:", transaction.Hexwit)
	fmt.Println(" txid:", transaction.Txid)
	fmt.Println(" txhash:", transaction.Txhash)
	fmt.Println(" version:", transaction.Version)
	fmt.Println(" locktime:", transaction.Locktime)
	fmt.Println(" expire:", transaction.Expire)
	fmt.Println(" blockheight:", transaction.Blockheight)
	fmt.Println(" confirmations:", transaction.Confirmations)
	fmt.Println(" vin:[")
	for _, vin := range transaction.Vin {
		fmt.Println("  {")
		fmt.Println("   txid:", vin.Txid)
		fmt.Println("   vout:", vin.Vout)
		fmt.Println("   sequence:", vin.Sequence)
		fmt.Println("   amountin:", vin.Amountin)
		fmt.Println("   blockheight:", vin.Blockheight)
		fmt.Println("   txindex:", vin.Txindex)
		fmt.Println("   scriptSig:{")
		fmt.Println("    asm:", vin.ScriptSig.Asm)
		fmt.Println("    hex:", vin.ScriptSig.Hex)
		fmt.Println("   }")
		fmt.Println("  }")
	}
	fmt.Println(" ]")
	fmt.Println(" vout:[")
	for _, vout := range transaction.Vout {
		fmt.Println("  {")
		fmt.Println("   amount:", vout.Amount)
		fmt.Println("   scriptPubKey:{")
		fmt.Println("    asm:", vout.ScriptPubKey.Asm)
		fmt.Println("    hex:", vout.ScriptPubKey.Hex)
		fmt.Println("    teqSigs:", vout.ScriptPubKey.ReqSigs)
		fmt.Println("    type:", vout.ScriptPubKey.Type)
		fmt.Println("    addresses:[")
		for _, addr := range vout.ScriptPubKey.Addresses {
			fmt.Printf("     \"%s\"\n", addr)
		}
		fmt.Println("    ]")
		fmt.Println("   }")
		fmt.Println("  }")
	}
	fmt.Println(" ]")
	fmt.Println("}")
}
