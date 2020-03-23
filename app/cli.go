package app

import (
	"flag"
	"fmt"
	"git.diabin.com/BlockChain/wallet-lib/core"
	"git.diabin.com/BlockChain/wallet-lib/database/sqlite"
	"git.diabin.com/BlockChain/wallet-lib/log"
	"git.diabin.com/BlockChain/wallet-lib/rpc"
	"os"
	"runtime"
)

type command interface {
	name() string
	handle()
	parse() bool
	usage()
}

type txCommand struct {
	*flag.FlagSet
	record  bool
	index   int
	size    int
	send    bool
	from    string
	to      string
	amount  float64
	speed   float64
	find    bool
	id      string
	trading bool
	rawtx   string
	out     bool
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
	cmdList = append(cmdList, &syncCommand{FlagSet: flag.NewFlagSet("sync", flag.ExitOnError)})
}

func initDB() error {
	db := &sqlite.Sqlite{}
	core.Storage = db
	core.List = db
	return db.Open("data.db")
}

func StartCli(cmd string) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := initDB(); err != nil {
		log.Errorf("open database failed! %s", err.Error())
		return
	}

	var sel command
	for _, c := range cmdList {
		if c.name() == cmd {
			sel = c
			break
		}
	}
	if sel != nil {
		startIdx = 1
		println()
		sel.handle()
	} else {
		CliUsage(1)
	}
}

/*begin tx */
func (cmd *txCommand) name() string {
	return cmd.Name()
}
func (cmd *txCommand) parse() bool {
	cmd.BoolVar(&cmd.send, "s", false, "send a tx to a specified address")
	cmd.BoolVar(&cmd.find, "f", false, "query a transaction by id")
	cmd.BoolVar(&cmd.trading, "t", false, "send the raw tx to the network")
	cmd.StringVar(&cmd.id, "id", "", "the id of transaction")
	cmd.StringVar(&cmd.from, "from", "", "the send address of transaction")
	cmd.StringVar(&cmd.to, "to", "", "the receive address of transaction")
	cmd.StringVar(&cmd.rawtx, "rawtx", "", "send raw tx of transaction")
	cmd.StringVar(&cmd.key, "key", "", "tx id and vout index key")
	cmd.Float64Var(&cmd.amount, "amt", 0, "the amount of transaction to be sent")
	cmd.Float64Var(&cmd.speed, "speed", core.Trading_Speed_Slow, "the speed of trade")
	cmd.IntVar(&cmd.index, "idx", 1, "start index of the query record")
	cmd.IntVar(&cmd.size, "size", 50, "number of records returned")
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
	//fmt.Printf("r: %t\ns: %t\nfrom: %s\namount: %d\n", *latest, *send, *from, *amount)
	switch true {
	case cmd.send:
		println("send tx", cmd.from, cmd.amount)
		key := core.GetPK(cmd.from)
		rs, ok := core.SendTransaction(cmd.from, key, cmd.to, cmd.amount, cmd.speed)
		if ok {
			println("send tx success, txid:", rs)
		} else {
			println("send tx failed!", rs)
		}
	case cmd.trading:
		println("send raw tx ", cmd.rawtx)
		rs, err := core.RpcClient.SendTransaction(cmd.rawtx)
		if err == nil {
			println("send raw tx success, txid:", rs)
		} else {
			println("send raw tx failed!", rs)
		}
	case cmd.find:
		println("query a transaction")
		txs, err := core.GetTransaction(cmd.id)
		if err != nil {
			println(err.Error())
			return
		}
		printTransaction(txs)
	case cmd.out:
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

/*begin stats*/
func (cmd *syncCommand) name() string {
	return cmd.Name()
}
func (cmd *syncCommand) parse() bool {
	args := os.Args[startIdx:]
	err := cmd.Parse(args)
	if len(args) == 0 || err != nil {
		return false
	}
	return true
}
func (cmd *syncCommand) usage() {
	cmd.Usage()
}
func (cmd *syncCommand) handle() {
	if !cmd.parse() {
		cmd.usage()
		return
	}
	core.StartSync()

}

/*end stats*/

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
