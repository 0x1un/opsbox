package main

type TronMainnet struct {
	Tps struct {
		Data struct {
			BlockHeight int `json:"blockHeight"`
			CurrentTps  int `json:"currentTps"`
			MaxTps      int `json:"maxTps"`
		} `json:"data"`
		Type string `json:"type"`
	} `json:"tps"`
	Node struct {
		Total int `json:"total"`
		Code  int `json:"code"`
	} `json:"node"`
	StatsOverview struct {
		Success bool `json:"success"`
		Data    []struct {
			AccountWithTrx          int     `json:"accountWithTrx"`
			Date                    int64   `json:"date"`
			DateDayStr              string  `json:"dateDayStr"`
			TotalTransaction        int64   `json:"totalTransaction"`
			AvgBlockTime            int     `json:"avgBlockTime"`
			AvgBlockSize            int     `json:"avgBlockSize"`
			TotalBlockCount         int     `json:"totalBlockCount"`
			NewAddressSeen          int     `json:"newAddressSeen"`
			ActiveAccountNumber     int     `json:"active_account_number"`
			BlockchainSize          int64   `json:"blockchainSize"`
			TotalAddress            int     `json:"totalAddress"`
			NewBlockSeen            int     `json:"newBlockSeen"`
			NewTransactionSeen      int     `json:"newTransactionSeen"`
			NewTrigger              int     `json:"newTrigger"`
			NewTrc10                int     `json:"newTrc10"`
			NewTrc20                int     `json:"newTrc20"`
			TotalTrc10              int     `json:"totalTrc10"`
			TotalTrc20              int     `json:"totalTrc20"`
			Triggers                int     `json:"triggers"`
			TrxTransfer             int     `json:"trx_transfer"`
			Trc10Transfer           int     `json:"trc10_transfer"`
			FreezeTransaction       int     `json:"freeze_transaction"`
			UnfreezeTransaction     int     `json:"unfreeze_transaction"`
			VoteTransaction         int     `json:"vote_transaction"`
			ShieldedTransaction     int     `json:"shielded_transaction"`
			OtherTransaction        int     `json:"other_transaction"`
			EnergyUsage             int64   `json:"energy_usage"`
			NetUsage                int64   `json:"net_usage"`
			EnergyUsageChange24H    float64 `json:"energy_usage_change_24h"`
			NetUsageChange24H       float64 `json:"net_usage_change_24h"`
			ActiveAccountNumberRate float64 `json:"active_account_number_rate"`
			NewAddressSeenRate      float64 `json:"newAddressSeenRate"`
			TotalAddressRate        float64 `json:"totalAddressRate"`
		} `json:"data"`
	} `json:"statsOverview"`
	FreezeResource struct {
		Total int `json:"total"`
		Data  []struct {
			TotalNetCost               float64 `json:"total_net_cost"`
			NetCostChange24H           float64 `json:"net_cost_change_24h"`
			TotalEnergyWeight          int64   `json:"total_energy_weight"`
			TotalTurnOver              string  `json:"total_turn_over"`
			EnergyRate                 float64 `json:"energy_rate"`
			TotalFreezeWeightChange24H float64 `json:"total_freeze_weight_change_24h"`
			TotalFreezeWeight          int64   `json:"total_freeze_weight"`
			EnergyCostChange24H        float64 `json:"energy_cost_change_24h"`
			TotalTurnoverChange24H     float64 `json:"total_turnover_change_24h"`
			FreezingRateChange24H      float64 `json:"freezing_rate_change_24h"`
			FreezingRate               float64 `json:"freezing_rate"`
			NetRate                    float64 `json:"net_rate"`
			TotalNetWeight             int64   `json:"total_net_weight"`
			Day                        string  `json:"day"`
			TotalEnergyCost            float64 `json:"total_energy_cost"`
		} `json:"data"`
	} `json:"freezeResource"`
	TriggerStatistic struct {
		Total int `json:"total"`
		Min   struct {
			TriggersAmountRate float64 `json:"triggers_amount_rate"`
			TriggersAmount     int     `json:"triggers_amount"`
			Day                int64   `json:"day"`
		} `json:"min"`
		Data []struct {
			TriggersAmountRate float64 `json:"triggers_amount_rate"`
			TriggersAmount     int     `json:"triggers_amount"`
			Day                int64   `json:"day"`
		} `json:"data"`
		Max struct {
			TriggersAmountRate float64 `json:"triggers_amount_rate"`
			TriggersAmount     int     `json:"triggers_amount"`
			Day                int64   `json:"day"`
		} `json:"max"`
	} `json:"triggerStatistic"`
	CallerAddressStatistic struct {
		Total int `json:"total"`
		Min   struct {
			CallerAmountRate float64 `json:"caller_amount_rate"`
			Day              int64   `json:"day"`
			CallerAmount     int     `json:"caller_amount"`
		} `json:"min"`
		Data []struct {
			CallerAmountRate float64 `json:"caller_amount_rate"`
			Day              int64   `json:"day"`
			CallerAmount     int     `json:"caller_amount"`
		} `json:"data"`
		Max struct {
			CallerAmountRate float64 `json:"caller_amount_rate"`
			Day              int64   `json:"day"`
			CallerAmount     int     `json:"caller_amount"`
		} `json:"max"`
	} `json:"callerAddressStatistic"`
	AccountList struct {
		Total      int           `json:"total"`
		Data       []interface{} `json:"data"`
		RangeTotal int           `json:"rangeTotal"`
	} `json:"accountList"`
	Funds struct {
		GenesisBlockIssue  int64   `json:"genesisBlockIssue"`
		TotalDonateBalance int64   `json:"totalDonateBalance"`
		TotalFundBalance   int64   `json:"totalFundBalance"`
		TotalBlockPay      int     `json:"totalBlockPay"`
		TotalNodePay       int64   `json:"totalNodePay"`
		BurnPerDay         int     `json:"burnPerDay"`
		BurnUsddByCharge   float64 `json:"burnUsddByCharge"`
		BurnByCharge       float64 `json:"burnByCharge"`
		TotalTurnOver      float64 `json:"totalTurnOver"`
		FundSumBalance     int64   `json:"fundSumBalance"`
		DonateBalance      int64   `json:"donateBalance"`
		FundTrx            float64 `json:"fundTrx"`
		TurnOver           float64 `json:"turnOver"`
	} `json:"funds"`
	StableCoin []struct {
		Amount24H        float64 `json:"amount24h"`
		Icon             string  `json:"icon"`
		PriceInUsd       float64 `json:"priceInUsd"`
		ContractAddress  string  `json:"contractAddress"`
		Volumn24H        float64 `json:"volumn24h"`
		TransferCount    int     `json:"transferCount"`
		Supply           int64   `json:"supply"`
		TransferCount24H int     `json:"transferCount24h"`
		Holders          int     `json:"holders"`
		Transfers        []struct {
			Date          int64   `json:"date"`
			Amount        float64 `json:"amount"`
			TransferCount int     `json:"transfer_count"`
			Day           string  `json:"day"`
		} `json:"transfers"`
		Name       string  `json:"name"`
		Abbr       string  `json:"abbr"`
		PriceInTrx float64 `json:"priceInTrx"`
	} `json:"stableCoin"`
	TrxVolume24H string `json:"trx_volume_24h"`
	TrxPriceLine struct {
		Total int `json:"total"`
		Data  []struct {
			PriceUsd string `json:"priceUsd"`
			Time     string `json:"time"`
		} `json:"data"`
	} `json:"trxPriceLine"`
	Tvl struct {
		TvlLine []struct {
			T string  `json:"t"`
			V float64 `json:"v"`
		} `json:"tvlLine"`
		Total    float64 `json:"total"`
		Projects []struct {
			Project string  `json:"project"`
			Logo    string  `json:"logo"`
			Type    string  `json:"type"`
			Vip     string  `json:"vip,omitempty"`
			Locked  float64 `json:"locked"`
			URL     string  `json:"url"`
			Gain    float64 `json:"gain"`
		} `json:"projects"`
	} `json:"tvl"`
}

// MainnetYiTaiFang yitaifang.com
type MainnetYiTaiFang struct {
	Status bool   `json:"status"`
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Data   struct {
		Summary struct {
			Difficulty   int64   `json:"difficulty"`
			Hashrate     float64 `json:"hashrate"`
			Lastblock    int     `json:"lastblock"`
			Avgtime      float64 `json:"avgtime"`
			Transactions int     `json:"transactions"`
			Tps          float64 `json:"tps"`
		} `json:"summary"`
		Headers []struct {
			Hash             string        `json:"hash"`
			ParentHash       string        `json:"parentHash"`
			Sha3Uncles       string        `json:"sha3Uncles"`
			Miner            string        `json:"miner"`
			StateRoot        string        `json:"stateRoot"`
			TransactionsRoot string        `json:"transactionsRoot"`
			ReceiptsRoot     string        `json:"receiptsRoot"`
			TotalDifficulty  string        `json:"totalDifficulty"`
			Difficulty       string        `json:"difficulty"`
			Number           int           `json:"number"`
			GasLimit         string        `json:"gasLimit"`
			GasUsed          string        `json:"gasUsed"`
			GasPrice         string        `json:"gasPrice"`
			Timestamp        string        `json:"timestamp"`
			BlockTime        int           `json:"blockTime"`
			Nonce            string        `json:"nonce"`
			Size             int           `json:"size"`
			ExtraData        string        `json:"extraData"`
			Hashrate         float64       `json:"hashrate"`
			Txs              int           `json:"txs"`
			TxFees           interface{}   `json:"txFees"`
			Reward           int           `json:"reward"`
			Uncles           int           `json:"uncles"`
			UncleToReward    float64       `json:"uncleToReward"`
			UncleReward      float64       `json:"uncleReward"`
			Position         []interface{} `json:"position"`
			MinerScore       interface{}   `json:"minerScore"`
			MinerComment     string        `json:"minerComment"`
		} `json:"headers"`
		Transactions []struct {
			Block               string        `json:"block"`
			Tx                  string        `json:"tx"`
			Timestamp           int           `json:"timestamp"`
			From                string        `json:"from"`
			To                  string        `json:"to"`
			Nonce               string        `json:"nonce"`
			GasPrice            string        `json:"gasPrice"`
			GasLimit            string        `json:"gasLimit"`
			Value               string        `json:"value"`
			Input               string        `json:"input"`
			CumulativeGasUsed   string        `json:"cumulativeGasUsed"`
			GasUsed             string        `json:"gasUsed"`
			Status              int           `json:"status"`
			TransactionIndex    string        `json:"transactionIndex"`
			TxFee               string        `json:"txFee"`
			BlockNumber         int           `json:"blockNumber"`
			ContractAddress     string        `json:"contractAddress"`
			InternalTransaction []interface{} `json:"internalTransaction"`
			InputArr            []interface{} `json:"inputArr"`
			ContractSymbol      string        `json:"contractSymbol"`
			FromComment         string        `json:"fromComment"`
			ToComment           string        `json:"toComment"`
			FromScore           interface{}   `json:"fromScore"`
			ToScore             interface{}   `json:"toScore"`
		} `json:"transactions"`
		TransactionsCount15 struct {
			Num1659715200 int `json:"1659715200"`
			Num1659801600 int `json:"1659801600"`
			Num1659888000 int `json:"1659888000"`
			Num1659974400 int `json:"1659974400"`
			Num1660060800 int `json:"1660060800"`
			Num1660147200 int `json:"1660147200"`
			Num1660233600 int `json:"1660233600"`
			Num1660320000 int `json:"1660320000"`
			Num1660406400 int `json:"1660406400"`
			Num1660492800 int `json:"1660492800"`
			Num1660579200 int `json:"1660579200"`
			Num1660665600 int `json:"1660665600"`
			Num1660752000 int `json:"1660752000"`
			Num1660838400 int `json:"1660838400"`
		} `json:"transactionsCount15"`
		TransactionLastTime int `json:"transactionLastTime"`
	} `json:"data"`
}
type BlockChairState struct {
	Data struct {
		Blocks                   int         `json:"blocks"`
		Transactions             int         `json:"transactions"`
		Outputs                  int         `json:"outputs"`
		Circulation              int64       `json:"circulation"`
		Blocks24H                int         `json:"blocks_24h"`
		Transactions24H          int         `json:"transactions_24h"`
		Difficulty               int64       `json:"difficulty"`
		Volume24H                int64       `json:"volume_24h"`
		MempoolTransactions      int         `json:"mempool_transactions"`
		MempoolSize              int         `json:"mempool_size"`
		MempoolTps               float64     `json:"mempool_tps"`
		MempoolTotalFeeUsd       float64     `json:"mempool_total_fee_usd"`
		BestBlockHeight          int64       `json:"best_block_height"`
		BestBlockHash            string      `json:"best_block_hash"`
		BestBlockTime            string      `json:"best_block_time"`
		BlockchainSize           int64       `json:"blockchain_size"`
		AverageTransactionFee24H interface{} `json:"average_transaction_fee_24h"`
		Inflation24H             interface{} `json:"inflation_24h"`
		MedianTransactionFee24H  interface{} `json:"median_transaction_fee_24h"`
		Cdd24H                   interface{} `json:"cdd_24h"`
		MempoolOutputs           int         `json:"mempool_outputs"`
		LargestTransaction24H    struct {
			Hash     string      `json:"hash"`
			ValueUsd interface{} `json:"value_usd"`
		} `json:"largest_transaction_24h"`
		Nodes                             int           `json:"nodes"`
		Hashrate24H                       interface{}   `json:"hashrate_24h"`
		InflationUsd24H                   interface{}   `json:"inflation_usd_24h"`
		AverageTransactionFeeUsd24H       interface{}   `json:"average_transaction_fee_usd_24h"`
		MedianTransactionFeeUsd24H        interface{}   `json:"median_transaction_fee_usd_24h"`
		MarketPriceUsd                    interface{}   `json:"market_price_usd"`
		MarketPriceBtc                    interface{}   `json:"market_price_btc"`
		MarketPriceUsdChange24HPercentage interface{}   `json:"market_price_usd_change_24h_percentage"`
		MarketCapUsd                      interface{}   `json:"market_cap_usd"`
		MarketDominancePercentage         float64       `json:"market_dominance_percentage"`
		NextRetargetTimeEstimate          string        `json:"next_retarget_time_estimate"`
		NextDifficultyEstimate            int64         `json:"next_difficulty_estimate"`
		Countdowns                        []interface{} `json:"countdowns"`
		SuggestedTransactionFeePerByteSat int           `json:"suggested_transaction_fee_per_byte_sat"`
		HodlingAddresses                  int           `json:"hodling_addresses"`
	} `json:"data"`
	Context struct {
		Code           int         `json:"code"`
		Source         string      `json:"source"`
		State          int         `json:"state"`
		MarketPriceUsd interface{} `json:"market_price_usd"`
		Cache          struct {
			Live     bool        `json:"live"`
			Duration interface{} `json:"duration"`
			Since    string      `json:"since"`
			Until    string      `json:"until"`
			Time     float64     `json:"time"`
		} `json:"cache"`
		API struct {
			Version         string      `json:"version"`
			LastMajorUpdate string      `json:"last_major_update"`
			NextMajorUpdate interface{} `json:"next_major_update"`
			Documentation   string      `json:"documentation"`
			Notice          string      `json:"notice"`
		} `json:"api"`
		Servers     string  `json:"servers"`
		Time        float64 `json:"time"`
		RenderTime  float64 `json:"render_time"`
		FullTime    float64 `json:"full_time"`
		RequestCost int     `json:"request_cost"`
	} `json:"context"`
}
