package okx

const (
	TimeMaxDiff = 2000 // 时间误差最大值，大于此值不能够下单
)

// 实盘
const (
	RestGlobalUrl      = "https://aws.okx.com"
	SocketPubUrl       = "wss://ws.okx.com:8443/ws/v5/public"
	SocketBusinessUrl  = "wss://ws.okx.com:8443/ws/v5/business"
	SocketPriGlobalUrl = "wss://wsaws.okx.com:8443/ws/v5/private"
)

// 模拟盘
const (
	RestSimulateUrl = "https://www.okx.com"
	SocketSimPubUrl = "wss://wspap.okx.com:8443/ws/v5/public?brokerId=9999"
	SocketSimPriUrl = "wss://wspap.okx.com:8443/ws/v5/private?brokerId=9999"
)

// public url
const (
	InstrumentsUrl = "/api/v5/public/instruments"
	FundingRateUrl = "/api/v5/public/funding-rate"
)

// asset url
const (
	AssetValuationUrl = "/api/v5/asset/asset-valuation"
	AssetBalancesUrl  = "/api/v5/asset/balances"
)

// account url
const (
	PositionsUrl        = "/api/v5/account/positions"
	PositionsHistoryUrl = "/api/v5/account/positions-history"
	SetPosModeUrl       = "/api/v5/account/set-position-mode"
	SetLeverageUrl      = "/api/v5/account/set-leverage"
	BalanceUrl          = "/api/v5/account/balance"
	MaxSizeUrl          = "/api/v5/account/max-size"
	AccountConfigUrl    = "/api/v5/account/config"
	TradeFeeUrl         = "/api/v5/account/trade-fee"
	InterestLimitsUrl   = "/api/v5/account/interest-limits"
	BillsUrl            = "/api/v5/account/bills"
	BillsArchiveUrl     = "/api/v5/account/bills-archive"
	InterestAccruedUrl  = "/api/v5/account/interest-accrued"
	InterestRateUrl     = "/api/v5/account/interest-rate"
)

// public url
const (
	TimeUrl = "/api/v5/public/time"
)

// market url
const (
	TickerUrl  = "/api/v5/market/ticker"
	TickersUrl = "/api/v5/market/tickers"

	CandlesUrl    = "/api/v5/market/candles"
	HisCandlesUrl = "/api/v5/market/history-candles"
	BooksUrl      = "/api/v5/market/books"
)

// trade url
const (
	OrderUrl             = "/api/v5/trade/order" // 直接下单
	BatchOrdersUrl       = "/api/v5/trade/batch-orders"
	ClosePositionUrl     = "/api/v5/trade/close-position"
	CancelOrderUrl       = "/api/v5/trade/cancel-order"
	OrdersPendingUrl     = "/api/v5/trade/orders-pending"
	PostOrderAlgo        = "/api/v5/trade/order-algo"   // 包含止盈止损的下单
	PostCancelOrderAlgos = "/api/v5/trade/cancel-algos" // 撤销策略订单
)

// 时间粒度
const (
	Min        = "1m"
	ThreeMin   = "3m"
	FiveMin    = "5m"
	FifteenMin = "15m"
	ThirtyMin  = "30m"
	H          = "1H"
	TwoH       = "2H"
	FourH      = "4H"
	Day        = "1D"
	TwoDay     = "2D"
	ThreeDay   = "3D"
	Week       = "1W"
	Mon        = "1M"
	Year       = "1Y"
)

// 毫秒
const (
	FiveMill    = 5 * 60 * 1000
	FifteenMill = 15 * 60 * 1000
	ThirtyMill  = 30 * 60 * 1000
	HMill       = 60 * 60 * 1000
	FourHMill   = 4 * HMill
)

// 产品类型
const (
	SPOT    = "SPOT"    // 币币
	MARGIN  = "MARGIN"  // 币币杠杆
	SWAP    = "SWAP"    // 永续合约
	FUTURES = "FUTURES" // 交割合约
	OPTION  = "OPTION"  // 期权
)

type Signal int

// 信号
const (
	Wait Signal = iota
	Long
	Short
)

type NetworkMode int

// 网络模式
const (
	HttpMode NetworkMode = iota + 1
	SocketMode
)

// 普通委托订单类型
const (
	Market   = "market"
	Limit    = "limit"
	PostOnly = "post_only"
)

// 条件委托订单类型
const (
	Condition     = "condition"       // 单向止盈止损
	Oco           = "oco"             // 双向止盈止损
	Plan          = "trigger"         // 计划委托
	MoveOrderStop = "move_order_stop" // 移动止盈止损
	Iceberg       = "iceberg"         // 冰山委托
	Twap          = "twap"            // 时间加权
)

// 交易模式
const (
	Isolated = "isolated" // 逐仓
	Cross    = "cross"    // 全仓
	Cash     = "cash"
)

// 方向
const (
	Buy       = "buy"
	Sell      = "sell"
	MakeLong  = "long"
	MakeShort = "short"
)

// 市价单委托数量的类型
const (
	BaseCcy  = "base_ccy"  // 交易货币
	QuoteCcy = "quote_ccy" // 计价货币
)

// 平仓价格类型
const (
	Last  = "last"  // 最新价格
	Index = "index" // 指数价格
	Mark  = "mark"  // 标记价格
)

// websocket
const (
	CandleChannel = "candle"
)

// 持仓模式
const (
	LongShortMode = "long_short_mode"
	NetMode       = "net_mode"
)
