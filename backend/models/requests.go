package models

type SimulationParameters struct {
	MarketCap           int    `json:"marketCap" binding:"required"`
	OnePercentBuy       int    `json:"onePercentBuy" binding:"required"`
	TenPercentBuy       int    `json:"tenPercentBuy" binding:"required"`
	NegativeTrend       int    `json:"negativeTrend"`
	LastBuyLimit        int    `json:"lastBuyLimit"`
	PortfolioProportion int    `json:"portfolioProportion"`
	Version             string `json:"version,omitempty"`
}
