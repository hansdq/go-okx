package okx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ZYKJShadow/go-okx/common/utils"
	"github.com/pkg/errors"
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type RestConfig struct {
	ApiKey    string
	SecretKey string
	Password  string
	Simulate  bool // 模拟盘标识
	Proxy     string
	Host      string
	Timeout   int
}

type ResponseBean struct {
	Code string        `json:"code"`
	Msg  string        `json:"msg"`
	Data []interface{} `json:"data"`
}

type Params map[string]interface{}

type FundingRate struct {
	FundingRate     string `json:"fundingRate"`
	FundingTime     string `json:"fundingTime"`
	InstId          string `json:"instId"`
	InstType        string `json:"instType"`
	Method          string `json:"method"`
	MaxFundingRate  string `json:"maxFundingRate"`
	MinFundingRate  string `json:"minFundingRate"`
	NextFundingRate string `json:"nextFundingRate"`
	NextFundingTime string `json:"nextFundingTime"`
	SettFundingRate string `json:"settFundingRate"`
	SettState       string `json:"settState"`
	Ts              string `json:"ts"`
}

type Asset struct {
	Details struct {
		Classic string `json:"classic"`
		Earn    string `json:"earn"`
		Funding string `json:"funding"`
		Trading string `json:"trading"`
	} `json:"details"`
	TotalBal string `json:"totalBal"`
	Ts       string `json:"ts"`
}

type Balance struct {
	AvailBal  string `json:"availBal"`
	Bal       string `json:"bal"`
	Ccy       string `json:"ccy"`
	FrozenBal string `json:"frozenBal"`
}

type Account struct {
	AdjEq      string `json:"adjEq"`
	BorrowFroz string `json:"borrowFroz"`
	Details    []struct {
		AvailBal      string `json:"availBal"`
		AvailEq       string `json:"availEq"`
		CashBal       string `json:"cashBal"`
		Ccy           string `json:"ccy"`
		CrossLiab     string `json:"crossLiab"`
		DisEq         string `json:"disEq"`
		Eq            string `json:"eq"`
		EqUsd         string `json:"eqUsd"`
		FixedBal      string `json:"fixedBal"`
		FrozenBal     string `json:"frozenBal"`
		Interest      string `json:"interest"`
		IsoEq         string `json:"isoEq"`
		IsoLiab       string `json:"isoLiab"`
		IsoUpl        string `json:"isoUpl"`
		Liab          string `json:"liab"`
		MaxLoan       string `json:"maxLoan"`
		MgnRatio      string `json:"mgnRatio"`
		NotionalLever string `json:"notionalLever"`
		OrdFrozen     string `json:"ordFrozen"`
		Twap          string `json:"twap"`
		UTime         string `json:"uTime"`
		Upl           string `json:"upl"`
		UplLiab       string `json:"uplLiab"`
		StgyEq        string `json:"stgyEq"`
		SpotInUseAmt  string `json:"spotInUseAmt"`
		BorrowFroz    string `json:"borrowFroz"`
		SpotIsoBal    string `json:"spotIsoBal"`
	} `json:"details"`
	Imr         string `json:"imr"`
	IsoEq       string `json:"isoEq"`
	MgnRatio    string `json:"mgnRatio"`
	Mmr         string `json:"mmr"`
	NotionalUsd string `json:"notionalUsd"`
	OrdFroz     string `json:"ordFroz"`
	TotalEq     string `json:"totalEq"`
	UTime       string `json:"uTime"`
}

type Candles struct {
	Timestamp int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Confirm   int64 // k线是否完结 0未完结 1完结
}

type AccountConfig struct {
	AcctLv          string        `json:"acctLv"`
	AutoLoan        bool          `json:"autoLoan"`
	CtIsoMode       string        `json:"ctIsoMode"`
	GreeksType      string        `json:"greeksType"`
	Level           string        `json:"level"`
	LevelTmp        string        `json:"levelTmp"`
	MgnIsoMode      string        `json:"mgnIsoMode"`
	PosMode         string        `json:"posMode"`
	SpotOffsetType  string        `json:"spotOffsetType"`
	Uid             string        `json:"uid"`
	Label           string        `json:"label"`
	RoleType        string        `json:"roleType"`
	TraderInsts     []interface{} `json:"traderInsts"`
	SpotRoleType    string        `json:"spotRoleType"`
	SpotTraderInsts []interface{} `json:"spotTraderInsts"`
	OpAuth          string        `json:"opAuth"`
	KycLv           string        `json:"kycLv"`
	Ip              string        `json:"ip"`
	Perm            string        `json:"perm"`
	MainUid         string        `json:"mainUid"`
}

type Book struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
	Ts   string     `json:"ts"`
}

type BookOrder struct {
	Price         float64 // 价格
	Amount        float64 // 数量
	Value         float64 // 总价
	OrderQuantity int64   // 此价格的订单数量
}

type MaxSize struct {
	Ccy     string `json:"ccy"`
	InstId  string `json:"instId"`
	MaxBuy  string `json:"maxBuy"`
	MaxSell string `json:"maxSell"`
}

type Instrument struct {
	Alias        string `json:"alias"`
	BaseCcy      string `json:"baseCcy"`
	Category     string `json:"category"`
	CtMult       string `json:"ctMult"`
	CtType       string `json:"ctType"`
	CtVal        string `json:"ctVal"`
	CtValCcy     string `json:"ctValCcy"`
	ExpTime      string `json:"expTime"`
	InstFamily   string `json:"instFamily"`
	InstId       string `json:"instId"`
	InstType     string `json:"instType"`
	Lever        string `json:"lever"`
	ListTime     string `json:"listTime"`
	LotSz        string `json:"lotSz"`
	MaxIcebergSz string `json:"maxIcebergSz"`
	MaxLmtSz     string `json:"maxLmtSz"`
	MaxMktSz     string `json:"maxMktSz"`
	MaxStopSz    string `json:"maxStopSz"`
	MaxTriggerSz string `json:"maxTriggerSz"`
	MaxTwapSz    string `json:"maxTwapSz"`
	MinSz        string `json:"minSz"`
	OptType      string `json:"optType"`
	QuoteCcy     string `json:"quoteCcy"`
	SettleCcy    string `json:"settleCcy"`
	State        string `json:"state"`
	Stk          string `json:"stk"`
	TickSz       string `json:"tickSz"`
	Uly          string `json:"uly"`
}

type SystemTime struct {
	Ts string `json:"ts"`
}

type Ticker struct {
	Last      string `json:"last" `    // 最新成交价
	LastSz    string `json:"lastSz" `  // 最新成交数量
	Open24h   string `json:"open24h" ` // 24小时开盘价
	AskSz     string `json:"askSz" `   // 卖一价对应的数量
	Low24h    string `json:"low24h" `  // 24小时最低价
	AskPx     string `json:"askPx" `   // 卖一价
	VolCcy24h string `json:"volCcy24h" `
	InstType  string `json:"instType" `
	InstId    string `json:"instId" `
	BidSz     string `json:"bidSz" ` // 买一价对应的数量
	BidPx     string `json:"bidPx" ` // 买一价
	High24h   string `json:"high24h" `
	SodUtc0   string `json:"sodUtc0" `
	Vol24h    string `json:"vol24h" `
	SodUtc8   string `json:"sodUtc8" `
	Ts        string `json:"ts" `
}

type Position struct {
	Adl           string `json:"adl"`
	AvailPos      string `json:"availPos"` // 可平仓数量
	AvgPx         string `json:"avgPx"`    // 开仓均价
	CTime         string `json:"cTime"`    // 持仓创建时间
	Ccy           string `json:"ccy"`      // 保证金币种
	DeltaBS       string `json:"deltaBS"`
	DeltaPA       string `json:"deltaPA"`
	GammaBS       string `json:"gammaBS"`
	GammaPA       string `json:"gammaPA"`
	Imr           string `json:"imr"`
	InstId        string `json:"instId"`
	InstType      string `json:"instType"`
	Interest      string `json:"interest"` // 利润
	Last          string `json:"last"`
	UsdPx         string `json:"usdPx"`
	Lever         string `json:"lever"`
	Liab          string `json:"liab"`     // 负债额
	LiabCcy       string `json:"liabCcy"`  // 负债币种
	LiqPx         string `json:"liqPx"`    // 预估强平价
	MarkPx        string `json:"markPx"`   // 标记价格
	Margin        string `json:"margin"`   // 保证金余额，可增减，仅适用于逐仓
	MgnMode       string `json:"mgnMode"`  // 保证金模式
	MgnRatio      string `json:"mgnRatio"` // 保证金率
	Mmr           string `json:"mmr"`
	NotionalUsd   string `json:"notionalUsd"`
	OptVal        string `json:"optVal"`
	PTime         string `json:"pTime"`
	Pos           string `json:"pos"`
	RealizedPnl   string `json:"realizedPnl"` // 已实现收益
	Pnl           string `json:"pnl"`         // 平仓累计订单收益额
	Fee           string `json:"fee"`         // 累计手续费金币，正数代表平台返佣，负数代表平台扣除
	FundingFee    string `json:"fundingFee"`  // 累计资金费用
	BePx          string `json:"bePx"`        // 盈亏平衡价
	PosCcy        string `json:"posCcy"`      // 仓位资产币种
	PosId         string `json:"posId"`
	PosSide       string `json:"posSide"`
	ThetaBS       string `json:"thetaBS"`
	ThetaPA       string `json:"thetaPA"`
	TradeId       string `json:"tradeId"`
	QuoteBal      string `json:"quoteBal"`      // 计价币余额
	BaseBal       string `json:"baseBal"`       // 交易币余额
	BaseBorrowed  string `json:"baseBorrowed"`  // 交易币已借数量
	BaseInterest  string `json:"baseInterest"`  // 交易币利息
	QuoteBorrowed string `json:"quoteBorrowed"` // 计价币已借数量
	QuoteInterest string `json:"quoteInterest"` // 计价币利息
	UTime         string `json:"uTime"`
	Upl           string `json:"upl"`      // 未实现收益
	UplRatio      string `json:"uplRatio"` // 未实现收益率
	VegaBS        string `json:"vegaBS"`
	VegaPA        string `json:"vegaPA"`
}

type Order struct {
	InstType        string `json:"instType,omitempty"` // 产品类型
	InstId          string `json:"instId"`
	Ccy             string `json:"ccy,omitempty"` // 保证金币种，仅适用于单币种保证金模式下的全仓币币杠杆订单
	OrdId           string `json:"ordId"`
	ClOrdId         string `json:"clOrdId,omitempty"` //客户自定义订单ID
	Tag             string `json:"tag,omitempty"`
	Px              string `json:"px,omitempty"`  // 委托价格
	Sz              string `json:"sz"`            // 委托数量
	Pnl             string `json:"pnl,omitempty"` // 收益
	OrdType         string `json:"ordType"`
	Side            string `json:"side"`
	PosSide         string `json:"posSide,omitempty"`
	TdMode          string `json:"tdMode"`
	AccFillSz       string `json:"accFillSz,omitempty"`
	FillPx          string `json:"fillPx,omitempty"`
	TradeId         string `json:"tradeId,omitempty"`
	FillSz          string `json:"fillSz,omitempty"`
	FillTime        string `json:"fillTime,omitempty"`
	Source          string `json:"source,omitempty"`
	State           string `json:"state,omitempty"`           // 订单状态
	AvgPx           string `json:"avgPx,omitempty"`           // 成交均价，如果成交数量为0，该字段也为0
	Lever           string `json:"lever,omitempty"`           // 杠杆倍数，0.01到125之间的数值，仅适用于 币币杠杆/交割/永续
	TpTriggerPx     string `json:"tpTriggerPx,omitempty"`     // 止盈触发价
	TpTriggerPxType string `json:"tpTriggerPxType,omitempty"` // 止盈触发价类型
	TpOrdPx         string `json:"tpOrdPx,omitempty"`         // 止盈委托价
	SlTriggerPx     string `json:"slTriggerPx,omitempty"`     // 止损触发价
	SlTriggerPxType string `json:"slTriggerPxType,omitempty"`
	SlOrdPx         string `json:"slOrdPx,omitempty"`
	FeeCcy          string `json:"feeCcy,omitempty"` // 交易手续费币种
	Fee             string `json:"fee,omitempty"`    // 订单交易手续费，平台向用户收取的交易手续费，手续费扣除为负数。如： -0.01
	RebateCcy       string `json:"rebateCcy,omitempty"`
	Rebate          string `json:"rebate,omitempty"`
	TgtCcy          string `json:"tgtCcy,omitempty"`
	Category        string `json:"category,omitempty"` // 订单种类
	UTime           string `json:"uTime,omitempty"`    // 订单状态更新时间，Unix时间戳的毫秒数格式，如：1597026383085
	CTime           string `json:"cTime,omitempty"`    // 订单创建时间，Unix时间戳的毫秒数格式， 如 ：1597026383085
	SCode           string `json:"sCode,omitempty"`    // 错误码，仅在失败时返回
	SMsg            string `json:"sMsg,omitempty"`     // 错误信息，仅在失败时返回

	AttachAlgoOrds []*Trigger `json:"attachAlgoOrds,omitempty"` // 止盈止损
}

type Trigger struct {
	TpTriggerPx     string `json:"tpTriggerPx,omitempty"`     // 止盈触发价
	TpTriggerPxType string `json:"tpTriggerPxType,omitempty"` // 止盈触发价类型
	TpOrdPx         string `json:"tpOrdPx,omitempty"`         // 止盈委托价，为-1时执行市价止盈
	SlTriggerPx     string `json:"slTriggerPx,omitempty"`     // 止损触发价
	SlTriggerPxType string `json:"slTriggerPxType,omitempty"` // 止损触发类型
	SlOrdPx         string `json:"slOrdPx,omitempty"`         // 止损委托价，为-1时执行市价止损
}

type FundingRateArbitrage struct {
	Acc3DFundingRate string `json:"acc3dFundingRate"` // 三日累计费率
	Apy              string `json:"apy"`              // 参考年化
	ArbitrageId      string `json:"arbitrageId"`      // 策略id
	BuyInstId        string `json:"buyInstId"`        // 策略需要购买的产品id
	BuyInstType      string `json:"buyInstType"`      // 策略需要购买的产品类型
	Ccy              string `json:"ccy"`              // 货币
	FundingRate      string `json:"fundingRate"`      // 当期资金费率
	FundingTime      string `json:"fundingTime"`      // 结算时间
	NextFundingRate  string `json:"nextFundingRate"`  // 预期费率
	NotionalUsd      string `json:"notionalUsd"`      // 持仓价值
	SellInstId       string `json:"sellInstId"`       // 策略需要卖出的产品id
	SellInstType     string `json:"sellInstType"`     // 策略需要卖出的产品类型
	Spread           string `json:"spread"`           //
	State            string `json:"state"`            // 状态
	Ts               string `json:"ts"`
	Yield3DPer10K    string `json:"yield3dPer10K"` // 三日万份收益
}

type TradeFee struct {
	Category  string `json:"category"`
	Delivery  string `json:"delivery"`
	Exercise  string `json:"exercise"`
	InstType  string `json:"instType"`
	Level     string `json:"level"`
	Maker     string `json:"maker"`
	MakerU    string `json:"makerU"`
	MakerUSDC string `json:"makerUSDC"`
	Taker     string `json:"taker"`
	TakerU    string `json:"takerU"`
	TakerUSDC string `json:"takerUSDC"`
	Ts        string `json:"ts"`
	Fiat      []struct {
		Ccy   string `json:"ccy"`
		Taker string `json:"taker"`
		Maker string `json:"maker"`
	} `json:"fiat"`
}

type InterestLimit struct {
	Debt             string `json:"debt"`
	Interest         string `json:"interest"`
	NextDiscountTime string `json:"nextDiscountTime"` // 下次扣息时间
	NextInterestTime string `json:"nextInterestTime"` // 下次计息时间
	LoanAlloc        string `json:"loanAlloc"`
	Records          []struct {
		AvailLoan         string `json:"availLoan"`
		Ccy               string `json:"ccy"`
		Interest          string `json:"interest"`
		LoanQuota         string `json:"loanQuota"`
		PosLoan           string `json:"posLoan"`
		Rate              string `json:"rate"` // 日利率
		AvgRate           string `json:"avgRate"`
		SurplusLmt        string `json:"surplusLmt"`
		SurplusLmtDetails struct {
			AllAcctRemainingQuota string `json:"allAcctRemainingQuota"`
			CurAcctRemainingQuota string `json:"curAcctRemainingQuota"`
			PlatRemainingQuota    string `json:"platRemainingQuota"`
		} `json:"surplusLmtDetails"`
		UsedLmt  string `json:"usedLmt"`
		UsedLoan string `json:"usedLoan"`
	} `json:"records"`
}

type Bill struct {
	Bal         string `json:"bal"`
	BalChg      string `json:"balChg"`
	BillId      string `json:"billId"`
	Ccy         string `json:"ccy"`
	ClOrdId     string `json:"clOrdId"`
	ExecType    string `json:"execType"`
	Fee         string `json:"fee"`
	FillFwdPx   string `json:"fillFwdPx"`
	FillIdxPx   string `json:"fillIdxPx"`
	FillMarkPx  string `json:"fillMarkPx"`
	FillMarkVol string `json:"fillMarkVol"`
	FillPxUsd   string `json:"fillPxUsd"`
	FillPxVol   string `json:"fillPxVol"`
	FillTime    string `json:"fillTime"`
	From        string `json:"from"`
	InstId      string `json:"instId"`
	InstType    string `json:"instType"`
	Interest    string `json:"interest"`
	MgnMode     string `json:"mgnMode"`
	Notes       string `json:"notes"`
	OrdId       string `json:"ordId"`
	Pnl         string `json:"pnl"`
	PosBal      string `json:"posBal"`
	PosBalChg   string `json:"posBalChg"`
	Px          string `json:"px"`
	SubType     string `json:"subType"`
	Sz          string `json:"sz"`
	Tag         string `json:"tag"`
	To          string `json:"to"`
	TradeId     string `json:"tradeId"`
	Ts          string `json:"ts"`
	Type        string `json:"type"`
}

type InterestAccrued struct {
	Ccy          string `json:"ccy"`
	InstId       string `json:"instId"`
	Interest     string `json:"interest"`
	InterestRate string `json:"interestRate"`
	Liab         string `json:"liab"`
	MgnMode      string `json:"mgnMode"`
	Ts           string `json:"ts"`
	Type         string `json:"type"`
}

type InterestRate struct {
	Ccy          string `json:"ccy"`
	InterestRate string `json:"interestRate"`
}

func InitRestConfig(apiKey, secretKey, password, proxy string, simulate bool) *RestConfig {
	return &RestConfig{
		ApiKey:    apiKey,
		SecretKey: secretKey,
		Password:  password,
		Simulate:  simulate,
		Proxy:     proxy,
	}
}

func (c *RestConfig) CheckLocalTime() error {
	t, err := c.GetTime()
	if err != nil {
		return err
	}

	systemTime := utils.MustParseInt64(t.Ts)

	if math.Abs(float64(time.Now().UnixMilli()-systemTime)) > TimeMaxDiff {
		return errors.New("pls check your local time")
	}

	return nil
}

// GetTime 获取交易所时间
func (c *RestConfig) GetTime() (*SystemTime, error) {
	var systemTimes []*SystemTime
	_, err := c.request(nil, &systemTimes, http.MethodGet, TimeUrl, "", true)
	if err != nil {
		return nil, err
	}

	return systemTimes[0], nil
}

// Positions 获取持仓信息
func (c *RestConfig) Positions(instType, instId, posId string) ([]*Position, error) {
	data := url.Values{
		"instType": {instType},
		"instId":   {instId},
		"posId":    {posId},
	}

	var positions []*Position
	_, err := c.request(nil, &positions, http.MethodGet, fmt.Sprintf("%s?%s", PositionsUrl, data.Encode()), "", false)
	if err != nil {
		return nil, err
	}

	return positions, nil
}

// PositionsHistory 获取历史持仓信息
func (c *RestConfig) PositionsHistory(instType, instId, posId, mgnMode, tp, after, before, limit string) ([]*Position, error) {
	data := url.Values{
		"instType": {instType},
		"instId":   {instId},
		"posId":    {posId},
		"mgnMode":  {mgnMode},
		"type":     {tp},
		"after":    {after},
		"before":   {before},
		"limit":    {limit},
	}

	var positions []*Position
	_, err := c.request(nil, &positions, http.MethodGet, fmt.Sprintf("%s?%s", PositionsHistoryUrl, data.Encode()), "", false)
	if err != nil {
		return nil, err
	}

	return positions, nil
}

// ClosePosition 市价全平
func (c *RestConfig) ClosePosition(instId string, posSide string, mgnMode string, ccy string, autoCxl bool) error {
	data := Params{
		"instId":  instId,
		"posSide": posSide,
		"mgnMode": mgnMode,
		"ccy":     ccy,
		"autoCxl": autoCxl,
	}

	_, err := c.request(data, nil, http.MethodPost, ClosePositionUrl, "", false)
	if err != nil {
		return err
	}

	return nil
}

// AccountConfig 获取账户配置
func (c *RestConfig) AccountConfig() (*AccountConfig, error) {
	var accountConfig []*AccountConfig
	_, err := c.request(nil, &accountConfig, http.MethodGet, AccountConfigUrl, "", false)
	if err != nil {
		return nil, err
	}

	return accountConfig[0], nil
}

// Ticker 产品当前行情数据
func (c *RestConfig) Ticker(instId string) (*Ticker, error) {
	data := url.Values{
		"instId": {instId},
	}

	var tickers []*Ticker
	_, err := c.request(nil, &tickers, http.MethodGet, fmt.Sprintf("%s?%s", TickersUrl, data.Encode()), "", true)
	if err != nil {
		return nil, err
	}

	return tickers[0], nil
}

// MaxSize
// 币币：返回最大可买的交易币数量和最大可卖的计价币数量，例如：BTC-USDT，返回的是BTC的最大可买数量和USDT的最大可卖数量
// 合约：返回最大可开多的合约张数和最大可开空的合约张数
func (c *RestConfig) MaxSize(instId, tdMode, ccy, px, leverages string, unSpotOffset bool) (*MaxSize, error) {
	data := url.Values{
		"instId":       {instId},
		"tdMode":       {tdMode},
		"ccy":          {ccy},
		"px":           {px},
		"leverages":    {leverages},
		"unSpotOffset": {fmt.Sprintf("%t", unSpotOffset)},
	}

	var maxSize []*MaxSize
	_, err := c.request(nil, &maxSize, http.MethodGet, fmt.Sprintf("%s?%s", MaxSizeUrl, data.Encode()), "", false)
	if err != nil {
		return nil, err
	}

	return maxSize[0], nil
}

// Candles 近期k线数据
func (c *RestConfig) Candles(instId, bar, after, before, limit string) ([]*Candles, error) {
	data := url.Values{
		"instId": {instId},
		"bar":    {bar},
		"after":  {after},
		"before": {before},
		"limit":  {limit},
	}

	var candles [][]string
	_, err := c.request(nil, &candles, http.MethodGet, fmt.Sprintf("%s?%s", CandlesUrl, data.Encode()), "", true)
	if err != nil {
		return nil, err
	}

	var ret []*Candles
	for _, item := range candles {
		ret = append(ret, &Candles{
			Timestamp: utils.MustParseInt64(item[0]),
			Open:      utils.MustParseFloat64(item[1]),
			High:      utils.MustParseFloat64(item[2]),
			Low:       utils.MustParseFloat64(item[3]),
			Close:     utils.MustParseFloat64(item[4]),
			Confirm:   utils.MustParseInt64(item[5]),
		})
	}

	return ret, nil
}

// HistoryCandles 历史k线数据
func (c *RestConfig) HistoryCandles(instId, bar, after, before, limit string) ([]*Candles, error) {
	data := url.Values{
		"instId": {instId},
		"bar":    {bar},
		"after":  {after},
		"before": {before},
		"limit":  {limit},
	}

	var candles [][]string
	_, err := c.request(nil, &candles, http.MethodGet, fmt.Sprintf("%s?%s", HisCandlesUrl, data.Encode()), "", true)
	if err != nil {
		return nil, err
	}

	var ret []*Candles
	for _, item := range candles {
		ret = append(ret, &Candles{
			Timestamp: utils.MustParseInt64(item[0]),
			Open:      utils.MustParseFloat64(item[1]),
			High:      utils.MustParseFloat64(item[2]),
			Low:       utils.MustParseFloat64(item[3]),
			Close:     utils.MustParseFloat64(item[4]),
			Confirm:   utils.MustParseInt64(item[5]),
		})
	}

	return ret, nil
}

// Balance 指定币种账户余额
func (c *RestConfig) Balance(ccy []string) (*Account, error) {
	data := url.Values{
		"ccy": {strings.Join(ccy, ",")},
	}

	var account []*Account
	_, err := c.request(nil, &account, http.MethodGet, fmt.Sprintf("%s?%s", BalanceUrl, data.Encode()), "", false)
	if err != nil {
		return nil, err
	}

	return account[0], nil
}

// FundingRate 资金费率
func (c *RestConfig) FundingRate(instId string) (*FundingRate, error) {
	data := url.Values{
		"instId": {instId},
	}

	var fundingRate []*FundingRate
	_, err := c.request(nil, &fundingRate, http.MethodGet, fmt.Sprintf("%s?%s", FundingRateUrl, data.Encode()), "", true)
	if err != nil {
		return nil, err
	}

	return fundingRate[0], nil
}

// SpotInstruments 产品列表-币币
func (c *RestConfig) SpotInstruments(uly, instFamily, instId string) ([]*Instrument, error) {
	return c.Instruments(SPOT, uly, instFamily, instId)
}

// MarginInstruments 产品列表-币币杠杆
func (c *RestConfig) MarginInstruments(uly, instFamily, instId string) ([]*Instrument, error) {
	return c.Instruments(MARGIN, uly, instFamily, instId)
}

// SwapInstruments 产品列表-永续合约
func (c *RestConfig) SwapInstruments(uly, instFamily, instId string) ([]*Instrument, error) {
	return c.Instruments(SWAP, uly, instFamily, instId)
}

// FuturesInstruments 产品列表-交割合约
func (c *RestConfig) FuturesInstruments(uly, instFamily, instId string) ([]*Instrument, error) {
	return c.Instruments(FUTURES, uly, instFamily, instId)
}

// OptionInstruments 产品列表-期权
func (c *RestConfig) OptionInstruments(uly, instFamily, instId string) ([]*Instrument, error) {
	return c.Instruments(OPTION, uly, instFamily, instId)
}

// Instruments 产品列表
func (c *RestConfig) Instruments(instType, uly, instFamily, instId string) ([]*Instrument, error) {
	data := url.Values{
		"instType":   {instType},
		"uly":        {uly},
		"instFamily": {instFamily},
		"instId":     {instId},
	}

	var instruments []*Instrument
	_, err := c.request(nil, &instruments, http.MethodGet, fmt.Sprintf("%s?%s", InstrumentsUrl, data.Encode()), "", true)
	if err != nil {
		return nil, err
	}

	if len(instruments) == 0 {
		return nil, fmt.Errorf("instrument not found, instType: %s, uly: %s, instFamily: %s, instId: %s", instType, uly, instFamily, instId)
	}

	return instruments, nil
}

// AssetBalances 资金账户余额
func (c *RestConfig) AssetBalances(ccy []string) ([]*Balance, error) {
	data := url.Values{
		"ccy": {strings.Join(ccy, ",")},
	}

	var balances []*Balance
	_, err := c.request(nil, &balances, http.MethodGet, fmt.Sprintf("%s?%s", AssetBalancesUrl, data.Encode()), "", false)
	if err != nil {
		return nil, err
	}

	return balances, nil
}

// AssetValuation 资产估值
func (c *RestConfig) AssetValuation(ccy string) (*Asset, error) {
	data := url.Values{
		"ccy": {ccy},
	}

	var assets []*Asset
	_, err := c.request(nil, &assets, http.MethodGet, fmt.Sprintf("%s?%s", AssetValuationUrl, data.Encode()), "", false)
	if err != nil {
		return nil, err
	}

	return assets[0], nil
}

// MarginMarketBuyOrder 币币杠杆-市价买入
// 注意：市价买入时买入所使用的货币和数量都是ccy
func (c *RestConfig) MarginMarketBuyOrder(instId, tdMode, ccy, sz string) (*Order, error) {
	return c.MakeOrder(instId, tdMode, ccy, "", Buy, Market, "", sz, false, "", "", false, nil)
}

// MarginMarketSellOrder 币币杠杆-市价卖出
// 注意：市价卖出时卖出的货币和数量都是instId
func (c *RestConfig) MarginMarketSellOrder(instId, tdMode, ccy, sz string) (*Order, error) {
	return c.MakeOrder(instId, tdMode, ccy, "", Sell, Market, "", sz, false, "", "", false, nil)
}

// SpotMarketBuyOrder 币币-市价买入
func (c *RestConfig) SpotMarketBuyOrder(instId, sz, tgtCcy string) (*Order, error) {
	return c.MakeOrder(instId, Cash, "", "", Buy, Market, "", sz, false, "", tgtCcy, false, nil)
}

// SpotMarketSellOrder 币币-市价卖出
func (c *RestConfig) SpotMarketSellOrder(instId, sz, tgtCcy string) (*Order, error) {
	return c.MakeOrder(instId, Cash, "", "", Sell, Market, "", sz, false, "", tgtCcy, false, nil)
}

// SwapMarketShortOrder 合约市价做空
func (c *RestConfig) SwapMarketShortOrder(instId, tdMode, sz string, triggers []*Trigger) (*Order, error) {
	return c.MakeOrder(instId, tdMode, "", "", Sell, Market, "", sz, false, MakeShort, "", false, triggers)
}

// SwapMarketLongOrder 合约市价做多
func (c *RestConfig) SwapMarketLongOrder(instId, tdMode, sz string, triggers []*Trigger) (*Order, error) {
	return c.MakeOrder(instId, tdMode, "", "", Buy, Market, "", sz, false, MakeLong, "", false, triggers)
}

// OrdersPending 获取未成交订单列表
func (c *RestConfig) OrdersPending(instType, uly, instFamily, instId, ordType, state, after, before, limit string) ([]*Order, error) {
	data := url.Values{
		"instType":   {instType},
		"uly":        {uly},
		"instFamily": {instFamily},
		"instId":     {instId},
		"ordType":    {ordType},
		"state":      {state},
		"after":      {after},
		"before":     {before},
		"limit":      {limit},
	}

	var orders []*Order
	_, err := c.request(nil, &orders, http.MethodGet, fmt.Sprintf("%s?%s", OrdersPendingUrl, data.Encode()), "", false)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (c *RestConfig) BatchOrders(data []*Order) ([]*Order, error) {
	var ret []*Order
	_, err := c.request(data, &ret, http.MethodPost, BatchOrdersUrl, "", false)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// MakeOrder 下单
func (c *RestConfig) MakeOrder(instId string, tdMode string, ccy string, clOrdId string, side string, ordType string, px string, sz string, reduceOnly bool, posSide string, tgtCcy string, banAmend bool, triggers []*Trigger) (*Order, error) {
	data := Params{
		"instId":         instId,
		"tdMode":         tdMode,
		"ccy":            ccy,
		"clOrdId":        clOrdId,
		"side":           side,
		"ordType":        ordType,
		"px":             px,
		"sz":             sz,
		"reduceOnly":     reduceOnly,
		"posSide":        posSide,
		"tgtCcy":         tgtCcy,
		"banAmend":       banAmend,
		"attachAlgoOrds": triggers,
	}

	var order []*Order
	_, err := c.request(data, &order, http.MethodPost, OrderUrl, "", false)
	if err != nil {
		return nil, err
	}

	return order[0], nil
}

// CheckOrder 查询订单
func (c *RestConfig) CheckOrder(instId, ordId, clOrdId string) (*Order, error) {
	data := url.Values{
		"instId":  {instId},
		"ordId":   {ordId},
		"clOrdId": {clOrdId},
	}

	var orders []*Order
	_, err := c.request(nil, &orders, http.MethodGet, fmt.Sprintf("%s?%s", OrderUrl, data.Encode()), "", false)
	if err != nil {
		return nil, err
	}

	return orders[0], nil
}

// CancelOrder 撤单
func (c *RestConfig) CancelOrder(instId, ordId, clOrdId string) (*Order, error) {
	data := Params{
		"instId":  instId,
		"ordId":   ordId,
		"clOrdId": clOrdId,
	}

	var orders []*Order
	_, err := c.request(data, &orders, http.MethodPost, CancelOrderUrl, "", false)
	if err != nil {
		return nil, err
	}

	return orders[0], nil
}

// SetPosMode 设置持仓模式
func (c *RestConfig) SetPosMode(mode string) error {
	data := Params{"posMode": mode}
	_, err := c.request(data, nil, http.MethodPost, SetPosModeUrl, "", false)
	return err
}

// SetLeverage 设置杠杆倍数
func (c *RestConfig) SetLeverage(instId string, lever string, mgnMode string) error {
	data := Params{
		"instId":  instId,
		"lever":   lever,
		"mgnMode": mgnMode,
	}

	_, err := c.request(data, nil, http.MethodPost, SetLeverageUrl, "", false)
	return err
}

// Books 获取产品深度数据
func (c *RestConfig) Books(instId, sz string) (*Book, error) {
	data := url.Values{
		"instId": {instId},
		"sz":     {sz},
	}

	var books []*Book
	_, err := c.request(nil, &books, http.MethodGet, fmt.Sprintf("%s?%s", BooksUrl, data.Encode()), "", true)
	if err != nil {
		return nil, err
	}

	return books[0], nil
}

// BookOrders 获取解析为订单的形式的产品深度数据
func (c *RestConfig) BookOrders(instId, sz string) ([]*BookOrder, []*BookOrder, error) {
	books, err := c.Books(instId, sz)
	if err != nil {
		return nil, nil, err
	}

	var asks []*BookOrder
	for _, item := range books.Asks {
		asks = append(asks, &BookOrder{
			Price:         utils.MustParseFloat64(item[0]),
			Amount:        utils.MustParseFloat64(item[1]),
			OrderQuantity: utils.MustParseInt64(item[3]),
			Value:         math.Round(utils.MustParseFloat64(item[1]) * utils.MustParseFloat64(item[0])),
		})
	}

	var bids []*BookOrder
	for _, item := range books.Bids {
		bids = append(bids, &BookOrder{
			Price:         utils.MustParseFloat64(item[0]),
			Amount:        utils.MustParseFloat64(item[1]),
			OrderQuantity: utils.MustParseInt64(item[3]),
			Value:         math.Round(utils.MustParseFloat64(item[1]) * utils.MustParseFloat64(item[0])),
		})
	}

	sort.Slice(asks, func(i, j int) bool {
		return asks[i].Value > asks[j].Value
	})

	sort.Slice(bids, func(i, j int) bool {
		return bids[i].Value > bids[j].Value
	})

	return asks, bids, nil
}

func (c *RestConfig) FundingRateArbitrage() ([]*FundingRateArbitrage, error) {
	data := url.Values{
		"ctType":        {"linear"},
		"ccyType":       {"USDT"},
		"arbitrageType": {"futures_spot"},
		"countryFilter": {"1"},
		"t":             {fmt.Sprintf("%d", time.Now().UnixMilli())},
	}

	var fundingRateArbitrages []*FundingRateArbitrage
	_, err := c.request(nil, &fundingRateArbitrages, http.MethodGet, "", fmt.Sprintf("%s?%s", "https://www.okx.com/priapi/v5/rubik/web/public/funding-rate-arbitrage", data.Encode()), true)
	if err != nil {
		return nil, err
	}

	return fundingRateArbitrages, nil
}

func (c *RestConfig) TradeFee(instId, instType, uly, instFamily string) (*TradeFee, error) {
	data := url.Values{
		"instId":     {instId},
		"instType":   {instType},
		"uly":        {uly},
		"instFamily": {instFamily},
	}

	var tradeFee []*TradeFee
	_, err := c.request(nil, &tradeFee, http.MethodGet, fmt.Sprintf("%s?%s", TradeFeeUrl, data.Encode()), "", false)
	if err != nil {
		return nil, err
	}

	if len(tradeFee) == 0 {
		return nil, fmt.Errorf("trade fee not found, instId: %s, instType: %s, uly: %s, instFamily: %s", instId, instType, uly, instFamily)
	}

	return tradeFee[0], nil
}

// InterestLimits 获取借币利率与限额
func (c *RestConfig) InterestLimits(tp string, ccy string) (*InterestLimit, error) {
	data := url.Values{
		"ccy":  {ccy},
		"type": {tp},
	}

	var interestLimits []*InterestLimit
	_, err := c.request(nil, &interestLimits, http.MethodGet, fmt.Sprintf("%s?%s", InterestLimitsUrl, data.Encode()), "", false)
	if err != nil {
		return nil, err
	}

	if len(interestLimits) == 0 {
		return nil, fmt.Errorf("interest limit not found, ccy: %s, type: %s", ccy, tp)
	}

	return interestLimits[0], nil
}

// Bills 账单流水查询（最近三天）
func (c *RestConfig) Bills(instType string, ccy string, mgnMode string, ctType string, tp string, subType string, after string, before string, begin string, end string, limit string) ([]*Bill, error) {
	data := url.Values{
		"instType": {instType},
		"ccy":      {ccy},
		"mgnMode":  {mgnMode},
		"ctType":   {ctType},
		"type":     {tp},
		"after":    {after},
		"before":   {before},
		"limit":    {limit},
		"subType":  {subType},
		"begin":    {begin},
		"end":      {end},
	}

	var bills []*Bill
	_, err := c.request(nil, &bills, http.MethodGet, fmt.Sprintf("%s?%s", BillsUrl, data.Encode()), "", false)
	if err != nil {
		return nil, err
	}

	return bills, nil
}

// BillsArchive 账单流水查询（近三个月）
func (c *RestConfig) BillsArchive(instType string, ccy string, mgnMode string, ctType string, tp string, subType string, after string, before string, begin string, end string, limit string) ([]*Bill, error) {
	data := url.Values{
		"instType": {instType},
		"ccy":      {ccy},
		"mgnMode":  {mgnMode},
		"ctType":   {ctType},
		"type":     {tp},
		"after":    {after},
		"before":   {before},
		"limit":    {limit},
		"subType":  {subType},
		"begin":    {begin},
		"end":      {end},
	}

	var bills []*Bill
	_, err := c.request(nil, &bills, http.MethodGet, fmt.Sprintf("%s?%s", BillsArchiveUrl, data.Encode()), "", false)
	if err != nil {
		return nil, err
	}

	return bills, nil
}

func (c *RestConfig) InterestAccrued(instId, ccy, tp, mgnMode, after, before, limit string) ([]*InterestAccrued, error) {
	data := url.Values{
		"instId":  {instId},
		"ccy":     {ccy},
		"mgnMode": {mgnMode},
		"type":    {tp},
		"after":   {after},
		"before":  {before},
		"limit":   {limit},
	}

	var interestAccrued []*InterestAccrued
	_, err := c.request(nil, &interestAccrued, http.MethodGet, fmt.Sprintf("%s?%s", InterestAccruedUrl, data.Encode()), "", false)
	if err != nil {
		return nil, err
	}

	return interestAccrued, nil
}

func (c *RestConfig) InterestRate(ccy string) ([]*InterestRate, error) {
	data := url.Values{
		"ccy": {ccy},
	}

	var interestRate []*InterestRate
	_, err := c.request(nil, &interestRate, http.MethodGet, fmt.Sprintf("%s?%s", InterestRateUrl, data.Encode()), "", false)
	if err != nil {
		return nil, err
	}

	return interestRate, nil
}

func (c *RestConfig) request(msg interface{}, res interface{}, method string, addr string, customize string, public bool) (*ResponseBean, error) {
	host := RestGlobalUrl
	if c.Simulate {
		host = RestSimulateUrl
	}

	client := &http.Client{
		Timeout: time.Duration(1) * time.Second,
	}

	if c.Proxy != "" {
		u, _ := url.Parse(c.Proxy)
		t := &http.Transport{
			MaxIdleConns:    10,
			MaxConnsPerHost: 10,
			IdleConnTimeout: time.Duration(3) * time.Second,
			Proxy:           http.ProxyURL(u),
		}
		client.Transport = t
	}

	var body []byte
	var err error
	if msg != nil {
		if body, err = json.Marshal(msg); err != nil {
			return nil, err
		}
	}

	u := host + addr
	if customize != "" {
		u = customize
	}

	r, err := http.NewRequest(method, u, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	if method == http.MethodPost {
		r.Header.Set("content-type", "application/json")
	}

	if !public {
		timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
		r.Header.Set("OK-ACCESS-KEY", c.ApiKey)
		r.Header.Set("OK-ACCESS-SIGN", c.getAccessSign(method, addr, string(body), timestamp))
		r.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
		r.Header.Set("OK-ACCESS-PASSPHRASE", c.Password)
	}

	if c.Simulate {
		r.Header.Set("x-simulated-trading", "1")
	}

	rsp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	var ret []byte
	ret, err = io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	var resp ResponseBean
	if err = json.Unmarshal(ret, &resp); err != nil {
		return nil, err
	}

	if resp.Code != "0" {
		return nil, fmt.Errorf("msg: %s, data: %+v", resp.Msg, resp.Data)
	}

	if res != nil && len(resp.Data) > 0 {
		bytes, err := json.Marshal(resp.Data)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(bytes, res); err != nil {
			return nil, err
		}
	}

	return &resp, nil
}

func (c *RestConfig) getAccessSign(method, requestPath, body, timestamp string) string {
	// OK-ACCESS-SIGN的请求头是对timestamp + method + requestPath + body字符串（+表示字符串连接），以及SecretKey，使用HMAC SHA256方法加密，通过Base-64编码输出而得到的。
	return base64.StdEncoding.EncodeToString(hmacSha256(c.SecretKey, timestamp+method+requestPath+body))
}

func hmacSha256(key, data string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return h.Sum(nil)
}
