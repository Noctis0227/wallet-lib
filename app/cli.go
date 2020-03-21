package app

import (
	"flag"
	"fmt"
	"git.diabin.com/BlockChain/wallet-lib/conf"
	"git.diabin.com/BlockChain/wallet-lib/core"
	"git.diabin.com/BlockChain/wallet-lib/rpc"
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

type outCommand struct {
	*flag.FlagSet
	all           bool
	unspent       bool
	spent         bool
	locked        bool
	memoryUnspent bool
	usbale        bool
	addr          string
}

type addrCommand struct {
	*flag.FlagSet
	generate bool
	list     bool
	def      string
}

type tokenCommand struct {
	*flag.FlagSet
	get    bool
	addr   string
	amount float64
}

type statsCommand struct {
	*flag.FlagSet
	block bool
	index int
	size  int
	chain bool
	sAddr bool
	addr  string
}

type peerCommand struct {
	*flag.FlagSet
	record bool
	index  int
	size   int
}

type minerCommand struct {
	*flag.FlagSet
	record bool
	index  int
	size   int
	find   bool
	height int64
	addr   string
}

var cmdList []command
var startIdx int

func init() {
	cmdList = make([]command, 0)
	cmdList = append(cmdList, &txCommand{FlagSet: flag.NewFlagSet("tx", flag.ExitOnError)})
	//cmdList = append(cmdList, &outCommand{FlagSet: flag.NewFlagSet("out", flag.ExitOnError)})
	cmdList = append(cmdList, &addrCommand{FlagSet: flag.NewFlagSet("addr", flag.ExitOnError)})
	//cmdList = append(cmdList, &tokenCommand{FlagSet: flag.NewFlagSet("token", flag.ExitOnError)})
	//cmdList = append(cmdList, &statsCommand{FlagSet: flag.NewFlagSet("stats", flag.ExitOnError)})
	//cmdList = append(cmdList, &peerCommand{FlagSet: flag.NewFlagSet("peer", flag.ExitOnError)})
	//cmdList = append(cmdList, &minerCommand{FlagSet: flag.NewFlagSet("miner", flag.ExitOnError)})
}

func StartCli(cmd string) {
	//core.LoadKeystore()
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
	cmd.BoolVar(&cmd.record, "r", false, "query latest tx records")
	cmd.BoolVar(&cmd.send, "s", false, "send a tx to a specified address")
	cmd.BoolVar(&cmd.find, "f", false, "query a transaction by id")
	cmd.BoolVar(&cmd.trading, "t", false, "send the raw tx to the network")
	cmd.BoolVar(&cmd.out, "o", false, "query vout by tx id and index")
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
	case cmd.record:
		println("get latest tx recodes")
		tf := "%-68s%-15s%-10s%-38s%-10s%-15s%-15s\n"
		bf := "%-68s%-15d%-10d%-38s%-10d%-15d%-15d\n"
		fmt.Printf(tf, "Key", "Amount", "Oder", "From", "Stat", "Change", "Fee")
		rs := core.GetLatestTxRecords(cmd.index, cmd.size)
		for _, it := range rs {
			fmt.Printf(bf, it.TxId, it.Amount, it.BlockOrder, it.From, it.Stat, it.Change, it.Fee)
		}
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
		output := core.GetOutput(cmd.key)
		if output != nil {
			tf := "%-70s%-30s%-70s%-10s%-40s%-20s\n"
			bf := "%-70s%-30s%-70s%-10d%-40s%-20d\n"
			fmt.Printf(tf, "Key", "time", "hash", "height", "addr", "amount")
			fmt.Printf(bf, output.TxId, output.Timestamp, output.BlockHash, output.BlockOrder, output.Address, output.Amount)
		}
	}
}

/*end tx*/

/*begin out*/
func (cmd *outCommand) name() string {
	return cmd.Name()
}
func (cmd *outCommand) parse() bool {
	cmd.BoolVar(&cmd.all, "a", false, "query spent outputs")
	cmd.BoolVar(&cmd.spent, "s", false, "query spent outputs")
	cmd.BoolVar(&cmd.unspent, "b", false, "query unspent outputs")
	cmd.BoolVar(&cmd.memoryUnspent, "m", false, "query in memory unspent outputs")
	cmd.BoolVar(&cmd.usbale, "u", false, "query usable unspent outputs")
	cmd.BoolVar(&cmd.locked, "l", false, "query locked outputs (in latest 16 blocks)")
	cmd.StringVar(&cmd.addr, "addr", "", "the address of output")
	args := os.Args[startIdx:]
	err := cmd.Parse(args)
	if len(args) == 0 || err != nil {
		return false
	}
	return true
}
func (cmd *outCommand) usage() {
	cmd.Usage()
}
func (cmd *outCommand) handle() {
	if !cmd.parse() {
		cmd.usage()
		return
	}
	var outs []core.Output
	var status = ""
	total := float64(0)
	switch true {
	case cmd.all:
		outs, total = core.GetAllOuts(cmd.addr)
		status = `"all"`
	case cmd.spent:
		outs, total = core.GetSpentOuts(cmd.addr)
		status = `"spent"`
	case cmd.unspent:
		outs, total = core.GetUnspentOuts(cmd.addr)
		status = `"unspent"`
	case cmd.memoryUnspent:
		outs, total = core.GetMemoryUnspent(cmd.addr)
		status = `"memory unspent"`
	case cmd.locked:
		outs, total = core.GetLockedUnspent(cmd.addr)
		status = `locked unspent`
	case cmd.usbale:
		outs, total = core.GetUsableOuts(cmd.addr)
		status = "usable unspent"
	}
	tf := "%-70s%-30s%-70s%-10s%-40s%-20s%-10s\n"
	bf := "%-70s%-30s%-70s%-10d%-40s%-20d%-10s\n"
	fmt.Printf(tf, "Key", "time", "hash", "Order", "addr", "amount", "status")
	for _, it := range outs {
		fmt.Printf(bf, it.TxId, it.Timestamp, it.BlockHash, it.BlockOrder, it.Address, it.Amount, status)
	}
	fmt.Printf("total: %f\n", total)
}

/*end out*/

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

/*begin token*/
func (cmd *tokenCommand) name() string {
	return cmd.Name()
}
func (cmd *tokenCommand) parse() bool {
	cmd.BoolVar(&cmd.get, "g", false, "get tokens from miner")
	cmd.Float64Var(&cmd.amount, "amt", 0, "the amount to be got")
	cmd.StringVar(&cmd.addr, "addr", core.GetDefAddr(), "the addr to be received tokens")
	args := os.Args[startIdx:]
	err := cmd.Parse(args)
	if len(args) == 0 || err != nil {
		return false
	}
	return true
}
func (cmd *tokenCommand) usage() {
	cmd.Usage()
}
func (cmd *tokenCommand) handle() {
	if !cmd.parse() {
		cmd.usage()
		return
	}
	//fmt.Printf("r: %t\ns: %t\nfrom: %s\namount: %d\n", *latest, *send, *from, *amount)
	switch true {
	case cmd.get:
		println("get tokens from miner, amount:", cmd.amount)
		cfg := conf.Setting.Miner
		core.SendTransaction(cfg.Address, cfg.PrivateKey, cmd.addr, cmd.amount, core.Trading_Speed_Slow)
	}
}

/*end token*/

/*begin stats*/
func (cmd *statsCommand) name() string {
	return cmd.Name()
}
func (cmd *statsCommand) parse() bool {
	cmd.BoolVar(&cmd.block, "b", false, "gets statistics for the block records")
	cmd.IntVar(&cmd.index, "idx", 1, "start index of the query record")
	cmd.IntVar(&cmd.size, "size", 50, "number of records returned")
	cmd.BoolVar(&cmd.chain, "c", false, "gets statistics for chain")
	cmd.BoolVar(&cmd.sAddr, "a", false, "gets statistics for the address")
	cmd.StringVar(&cmd.addr, "addr", core.GetDefAddr(), "statistical address. if it is empty, list the statistics of all addresses ")
	args := os.Args[startIdx:]
	err := cmd.Parse(args)
	if len(args) == 0 || err != nil {
		return false
	}
	return true
}
func (cmd *statsCommand) usage() {
	cmd.Usage()
}
func (cmd *statsCommand) handle() {
	if !cmd.parse() {
		cmd.usage()
		return
	}
	switch true {
	case cmd.block:
		tf := "%-70s%-10s%-10s%-20s%-20s\n"
		fmt.Printf(tf, "Hash", "Order", "TxCount", "Size", "CreateTime")
		bf := "%-70s%-10d%-10d%-20d%-20s\n"
		for _, it := range core.GetLatestBlocks(cmd.index, cmd.size) {
			fmt.Printf(bf, it.Hash, it.BlockOrder, it.Transactions, it.Size, it.CreateTime)
		}
	case cmd.chain:
		tf := "%-20s%-20s%-20s%-20s%-20s\n"
		fmt.Printf(tf, "BlockNum", "TxNum", "UnconfirmedTxNum", "LatestHeight", "PeerCount")
		bf := "%-20s%-20s%-20s%-20s%-20s\n"
		rs := core.GetChainStats()
		if rs != nil {
			fmt.Printf(bf, rs[string(core.StatsKey.BlockNum)], rs[string(core.StatsKey.TxNum)],
				rs[string(core.StatsKey.UnconfirmedTxNum)], rs[string(core.StatsKey.LatestHeight)],
				rs[string(core.StatsKey.PeerInfoCount)])
		}
	case cmd.sAddr:
		//println("Gets statistics for the address", cmd.addr)
		if rs := core.GetAddressAmounts(cmd.addr); rs != nil {
			tf := "%-20s%-20s%-20s%-20s%-20s\n"
			fmt.Printf(tf, "Balance", "UseableAmount", "LockedAmount", "MemoryAmount", "SpentAmount")
			bf := "%-20f%-20f%-20f%-20f%-20f\n"
			fmt.Printf(bf, rs.Balance, rs.UseableAmount, rs.LockedAmount, rs.MemoryAmount, rs.SpentAmount)
		}
	}
}

/*end stats*/

/*begin peer*/
func (cmd *peerCommand) name() string {
	return cmd.Name()
}
func (cmd *peerCommand) parse() bool {
	cmd.BoolVar(&cmd.record, "r", false, "query latest peer records")
	cmd.IntVar(&cmd.index, "idx", 1, "start index of the query record")
	cmd.IntVar(&cmd.size, "size", 50, "number of records returned")
	args := os.Args[startIdx:]
	err := cmd.Parse(args)
	if len(args) == 0 || err != nil {
		return false
	}
	return true
}
func (cmd *peerCommand) usage() {
	cmd.Usage()
}
func (cmd *peerCommand) handle() {
	if !cmd.parse() {
		cmd.usage()
		return
	}
	switch true {
	case cmd.record:
		tf := "%-20s%-30s%-35s%-40s%-20s%-20s%-20s\n"
		fmt.Printf(tf, "Id", "Addr", "ConnTime", "SubVer", "Mainorder", "Layer", "MainHeight")
		bf := "%-20d%-30s%-35s%-40s%-20d%-20d%-20d\n"
		for _, it := range core.GetLatestPeers(cmd.index, cmd.size) {
			fmt.Printf(bf, it.Id, it.Addr, it.ConnTime, it.SubVer, it.MainOrder, it.Layer, it.MainHeight)
		}
	}
}

/*end peer*/

/*begin miner*/
func (cmd *minerCommand) name() string {
	return cmd.Name()
}
func (cmd *minerCommand) parse() bool {
	cmd.BoolVar(&cmd.record, "r", false, "query latest miner records")
	cmd.IntVar(&cmd.index, "idx", 1, "start index of the query record")
	cmd.IntVar(&cmd.size, "size", 50, "number of records returned")
	cmd.BoolVar(&cmd.find, "f", false, "query miner by block height")
	cmd.Int64Var(&cmd.height, "height", -1, "the block height")
	cmd.StringVar(&cmd.addr, "addr", "", "the miner address")
	args := os.Args[startIdx:]
	err := cmd.Parse(args)
	if len(args) == 0 || err != nil {
		return false
	}
	return true
}
func (cmd *minerCommand) usage() {
	cmd.Usage()
}
func (cmd *minerCommand) handle() {
	if !cmd.parse() {
		cmd.usage()
		return
	}
	switch true {
	case cmd.record:
		tf := "%-40s%-20s%-20s\n"
		fmt.Printf(tf, "Addr", "CoinBase", "Block Number")
		bf := "%-40s%-20d%-20d\n"
		miners := core.GetLatestMiners(cmd.index, cmd.size)
		if miners != nil {
			for _, m := range miners {
				fmt.Printf(bf, m.Addr, m.CoinBase, len(m.Orders))
			}
		}
	case cmd.find && cmd.height != -1:
		tf := "%-40s%-20s%-20s\n"
		fmt.Printf(tf, "Addr", "CoinBase", "Block Number")
		bf := "%-40s%-20d%-20d\n"
		miner := core.GetMinerRecordByHeight(uint64(cmd.height))
		if miner != nil {
			fmt.Printf(bf, miner.Addr, miner.CoinBase, len(miner.Orders))
		}
	case cmd.find && cmd.addr != "":
		tf := "%-40s%-20s%-20s\n"
		fmt.Printf(tf, "Addr", "CoinBase", "Block Number")
		bf := "%-40s%-20d%-20d\n"
		miner := core.GetMinerRecordByAddr(cmd.addr)
		if miner != nil {
			fmt.Printf(bf, miner.Addr, miner.CoinBase, len(miner.Orders))
		}
	}
}

/*end miner*/

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
