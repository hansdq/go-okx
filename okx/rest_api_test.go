package okx

import (
	"fmt"
	"github.com/hansdq/go-okx/common/utils"
	"math"
	"strconv"
	"testing"
	"time"
)

var apiConfig = InitRestConfig("", "", "", "", true)

func TestGetTime(t *testing.T) {
	res, err := apiConfig.GetTime()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(res.Ts)
}

func TestGetPositions(t *testing.T) {
	positions, err := apiConfig.Positions("", "", "")
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range positions {
		t.Logf("产品id:%v 产品类型:%v 持仓id:%v 持仓方向:%v 持仓数量:%v 交易币余额:%v 计价币余额:%v 交易币已借:%v 交易币计息:%v 计价币已借:%v 计价币计息:%v", v.InstId, v.InstType, v.PosId, v.PosSide, v.Pos, v.BaseBal, v.QuoteBal, v.BaseBorrowed, v.BaseInterest, v.QuoteBorrowed, v.QuoteInterest)
		t.Logf("开仓均价:%v 预估强平价格:%v 开仓时间:%v 累计手续费:%v 累计资金费用:%v 利息:%v 可平仓数量:%v", v.AvgPx, v.LiqPx, v.UTime, v.Fee, v.FundingFee, v.Interest, v.AvailPos)
		t.Logf("负债额:%v 负债币种:%v 已实现收益:%v 未实现收益:%v 盈亏平衡价:%v", v.Liab, v.LiabCcy, v.RealizedPnl, v.Upl, v.BePx)
	}
}

func TestSetLeverage(t *testing.T) {
	err := apiConfig.SetLeverage("BTC-USDT-SWAP", "30", Cross)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSetPosMode(t *testing.T) {
	err := apiConfig.SetPosMode(LongShortMode)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestMarginMarketBuyOrder(t *testing.T) {
	order, err := apiConfig.MarginMarketBuyOrder("BTC-USDT", Cross, "USDT", "1000")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(order.OrdId)
}

func TestMarginMarketSellOrder(t *testing.T) {
	order, err := apiConfig.MarginMarketSellOrder("BTC-USDT", Cross, "USDT", "0.1")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(order.OrdId)
}

func TestSpotMarketBuyOrder(t *testing.T) {
	order, err := apiConfig.SpotMarketBuyOrder("BTC-USDT", "1", BaseCcy)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(order.OrdId)
}

func TestSpotMarketSellOrder(t *testing.T) {
	order, err := apiConfig.SpotMarketSellOrder("BTC-USDT", "1", BaseCcy)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(order.OrdId)
}

func TestMakeMarketLongOrder(t *testing.T) {
	// 无止盈止损
	order, err := apiConfig.SwapMarketLongOrder("BTC-USDT-SWAP", Cross, "1", nil)
	if err != nil {
		t.Error(err)
		return
	}

	// 止盈止损
	order, err = apiConfig.SwapMarketLongOrder("ETH-USDT-SWAP", Cross, "10", []*Trigger{
		{
			TpTriggerPxType: Last,
			TpOrdPx:         "-1",
			TpTriggerPx:     "4500",
			SlTriggerPx:     "2000",
			SlOrdPx:         "-1",
			SlTriggerPxType: Last,
		},
	})

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(order.OrdId)
}

func TestBatchOrders(t *testing.T) {
	var balance float64 = 100

	ticker, err := apiConfig.Ticker("BTC-USDT-SWAP")
	if err != nil {
		t.Error(err)
		return
	}

	instruments, err := apiConfig.Instruments(SWAP, "", "", "BTC-USDT-SWAP")
	if err != nil {
		t.Error(err)
		return
	}

	instrument := instruments[0]
	ctVal := utils.MustParseFloat64(instrument.CtVal)
	fmt.Println(instrument.CtVal, instrument.CtMult, instrument.State, instrument.LotSz, instrument.MinSz)

	bidPx := utils.MustParseFloat64(ticker.BidPx)
	askPx := utils.MustParseFloat64(ticker.AskPx)

	bidSz := balance / bidPx
	askSz := balance / askPx

	sz := strconv.Itoa(int(bidSz / ctVal))

	fmt.Println(bidSz, askSz, sz)

	var data []*Order
	data = append(data, &Order{
		InstId:  "BTC-USDT-SWAP",
		TdMode:  Cross,
		OrdType: Limit,
		PosSide: MakeShort,
		Side:    Sell,
		Sz:      sz,
		Px:      ticker.BidPx,
	})

	ticker, err = apiConfig.Ticker("BTC-USDT")
	if err != nil {
		t.Error(err)
		return
	}

	bidPx = utils.MustParseFloat64(ticker.BidPx)
	askPx = utils.MustParseFloat64(ticker.AskPx)

	bidSz = balance / bidPx
	askSz = balance / askPx

	instruments, err = apiConfig.Instruments(SPOT, "", "", "BTC-USDT")
	if err != nil {
		t.Error(err)
		return
	}

	instrument = instruments[0]

	minSz := utils.MustParseFloat64(instrument.MinSz)
	fmt.Println(instrument.State, instrument.MinSz, minSz, instrument.LotSz, len(instrument.LotSz))

	// 计算精度因子
	precision := math.Log10(1 / utils.MustParseFloat64(instrument.LotSz))

	// 使用精度因子来格式化下单数量
	formattedQuantity := fmt.Sprintf("%.*f", int(precision), askSz)
	fmt.Println(formattedQuantity)

	data = append(data, &Order{
		InstId:  "BTC-USDT",
		TdMode:  Cross,
		OrdType: Limit,
		Side:    Buy,
		Ccy:     "USDT",
		Sz:      "1",
		Px:      ticker.AskPx,
	})

	//orders, err := apiConfig.BatchOrders(data)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//
	//for _, v := range orders {
	//	fmt.Println(v)
	//}
}

func TestCheckOrder(t *testing.T) {
	order, err := apiConfig.CheckOrder("ETH-USDT-SWAP", "663035144311250944", "")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(order)
}

func TestTicker(t *testing.T) {
	ticker, err := apiConfig.Ticker("BTC-USDT-SWAP")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(ticker)
}
func TestTickers(t *testing.T) {
	ticker, err := apiConfig.Tickers("SPOT")
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range ticker {
		fmt.Println(v.InstId)
	}

}
func TestCandles(t *testing.T) {
	candles, err := apiConfig.Candles("BTC-USDT-SWAP", "15m", "", "", "300")
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range candles {
		fmt.Println(v)
	}
}

func TestHistoryCandles(t *testing.T) {
	candles, err := apiConfig.HistoryCandles("BTC-USDT-SWAP", "15m", "", "", "300")
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range candles {
		fmt.Println(v)
	}
}

func TestBalance(t *testing.T) {
	balance, err := apiConfig.Balance([]string{"BTC", "ETH", "USDT"})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", balance)
	for _, item := range balance.Details {
		fmt.Println(item.Ccy, math.Round(utils.MustParseFloat64(item.AvailBal)*1000)/1000)
	}
}

func TestAssetValuation(t *testing.T) {
	assetValuation, err := apiConfig.AssetValuation("USDT")
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", assetValuation)
}

func TestAssetBalances(t *testing.T) {
	balances, err := apiConfig.AssetBalances([]string{"BTC", "ETH"})
	if err != nil {
		t.Error(err)
		return
	}

	for _, item := range balances {
		t.Logf("%+v", item)
	}
}

func TestInstruments(t *testing.T) {
	instruments, err := apiConfig.Instruments(SPOT, "", "", "PEOPLE-USDT")
	if err != nil {
		t.Error(err)
		return
	}

	for _, item := range instruments {
		t.Logf("%+v", item)
	}

	//marginInstruments, err := apiConfig.MarginInstruments("", "", "")
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//
	//for _, item := range marginInstruments {
	//	t.Logf("%+v", item)
	//}
}

func TestFundingRate(t *testing.T) {
	fundingRate, err := apiConfig.FundingRate("ARB-USDT")
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", fundingRate)
}

func TestMaxSize(t *testing.T) {
	maxSize, err := apiConfig.MaxSize("BTC-USDT", Cash, "", "", "", false)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(maxSize.Ccy, maxSize.MaxBuy, maxSize.MaxSell)
}

func TestAccountConfig(t *testing.T) {
	accountConfig, err := apiConfig.AccountConfig()
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", accountConfig)
}

func TestBooks(t *testing.T) {
	books, err := apiConfig.Books("PEOPLE-USDT", "400")
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", books)
}

func TestBookOrders(t *testing.T) {
	asks, bids, err := apiConfig.BookOrders("PEOPLE-USDT", "5000")
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", asks[0])
	t.Logf("%+v", bids[0])
}

func TestFundingRateArbitrage(t *testing.T) {
	fundingRateArbitrage, err := apiConfig.FundingRateArbitrage()
	if err != nil {
		t.Error(err)
		return
	}

	for _, item := range fundingRateArbitrage {
		t.Logf("%+v", item)
	}
}

func TestTradeFee(t *testing.T) {
	tradeFee, err := apiConfig.TradeFee("", SWAP, "", "")
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", tradeFee)
}

func TestInterestLimits(t *testing.T) {
	interestLimit, err := apiConfig.InterestLimits("", "ETH")
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", interestLimit)
	record := interestLimit.Records[0]

	nextDiscountTime := time.UnixMilli(utils.MustParseInt64(interestLimit.NextDiscountTime)).Format("2006-01-02 15:04:05")
	interestTime := time.UnixMilli(utils.MustParseInt64(interestLimit.NextInterestTime)).Format("2006-01-02 15:04:05")

	t.Logf("当前负债：%vUSDT 下次扣息时间:%v 下次计息时间:%v 日利率：%v%%", utils.MustParseFloat64(interestLimit.Debt), nextDiscountTime, interestTime, utils.MustParseFloat64(record.Rate)*100)
}

func TestBills(t *testing.T) {
	bills, err := apiConfig.Bills(MARGIN, "", "", "", "7", "", "", "", "", "", "")
	if err != nil {
		t.Error(err)
		return
	}

	for _, item := range bills {
		t.Logf("%+v", item)
	}

	bills, err = apiConfig.Bills(SWAP, "", "", "", "7", "", "", "", "", "", "")
	if err != nil {
		t.Error(err)
		return
	}

	for _, item := range bills {
		t.Logf("%+v", item)
	}

}

func TestInterestAccrued(t *testing.T) {
	interestAccrued, err := apiConfig.InterestAccrued("MAGIC-USDT", "USDT", "2", "", "", "1704648419000", "100")
	if err != nil {
		t.Error(err)
		return
	}

	for _, item := range interestAccrued {
		t.Logf("%+v", item)
	}
}

func TestInterestRate(t *testing.T) {
	interestRate, err := apiConfig.InterestRate("MAGIC")
	if err != nil {
		t.Error(err)
		return
	}

	for _, item := range interestRate {
		t.Logf("%+v", item)

		hourlyRate := utils.MustParseFloat64(item.InterestRate)
		hoursPerDay := 24

		// 计算连续复利的每天利率
		// dailyRate = (1 + hourlyRate)^(hoursPerDay) - 1
		dailyRate := math.Pow(1+hourlyRate, float64(hoursPerDay)) - 1

		// 将结果转换为百分比形式
		dailyRatePercent := dailyRate * 100

		fmt.Printf("每天的利率是: %.5f%%\n", dailyRatePercent)
	}
}

func TestClosePosition(t *testing.T) {
	err := apiConfig.ClosePosition("SHIB-USDT", "", Cross, "USDT", false)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
}
