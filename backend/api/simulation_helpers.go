package api

import (
	"errors"
	"fmt"
	"stocks/backend/models"
	"time"

	"gorm.io/gorm"
)

func generateVersion(p models.SimulationParameters) string {
	// Format: MC{marketCap}_1P{onePercentBuy}_10P{tenPercentBuy}_NT{negativeTrend}_LBL{lastBuyLimit}_PP{portfolioProportion}
	return fmt.Sprintf("MC%d_1P%d_10P%d_NT%d_LBL%d_PP%d",
		p.MarketCap,
		p.OnePercentBuy,
		p.TenPercentBuy,
		p.NegativeTrend,
		p.LastBuyLimit,
		p.PortfolioProportion)
}

func findStartDate(priceData []models.PricePoint, weeks int) (time.Time, error) {
	// the start date needs to be a spot where there is enough data to check for the negative trend
	minDate := priceData[0].Date
	targetDate := minDate.AddDate(0, 0, weeks*7)

	for _, point := range priceData {
		if point.Date.After(targetDate) || point.Date.Equal(targetDate) {
			return point.Date, nil
		}
	}
	return time.Time{}, fmt.Errorf("could not find startDate after %v weeks from %v", weeks, minDate)
}

func getStartIndex(priceData []models.PricePoint, startDate time.Time) (int, error) {
	for i, point := range priceData {
		if point.Date.Equal(startDate) || point.Date.After(startDate) {
			return i, nil
		}
	}
	return 0, fmt.Errorf("error finding starting index")
}

func passLastBuyCheck(db *gorm.DB, symbol string, dataPoint *models.PricePoint, parameters *models.SimulationParameters) bool {
	var investment models.Investment
	result := db.Where("symbol = ? AND version = ?", symbol, parameters.Version).First(&investment)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// no investment was found then there's not previous transaction
		return true
	}
	cutOffDate := dataPoint.Date.AddDate(0, 0, -parameters.LastBuyLimit)

	return investment.LastTransaction.Before(cutOffDate)
}
