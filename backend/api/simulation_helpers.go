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

func (ctx *processingContext) fetchPriceData() error {
	var data []models.PricePoint
	result := ctx.db.Model(&models.Historical{}).
		Select("date, price").
		Where("symbol = ?", ctx.symbol).
		Order("date ASC").
		Find(&data)

	if result.Error != nil {
		return fmt.Errorf("error querying data for %s: %v", ctx.symbol, result.Error)
	}
	if len(data) == 0 {
		return fmt.Errorf("no price data found for %s", ctx.symbol)
	}

	ctx.priceData = &data // Store pointer to the slice
	return nil
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

func (ctx *processingContext) determineStartingPoint() error {
	startDate, err := findStartDate(*ctx.priceData, ctx.parameters.NegativeTrend)
	if err != nil {
		return err
	}

	ctx.startIndex, err = getStartIndex(*ctx.priceData, startDate)
	return err
}

func (ctx *processingContext) passFallInPrice(dataPoint *models.PricePoint) bool {
	return true
}

func (ctx *processingContext) passLastBuyCheck(dataPoint *models.PricePoint) bool {
	var investment models.Investment
	result := db.Where("symbol = ? AND version = ?", symbol, parameters.Version).First(&investment)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// no investment was found then there's not previous transaction
		return true
	}
	cutOffDate := dataPoint.Date.AddDate(0, 0, -parameters.LastBuyLimit)

	return investment.LastTransaction.Before(cutOffDate)
}

func (ctx *processingContext) passNegativeTrend(dataPoint *models.PricePoint) bool {
	return true
}

func (ctx *processingContext) passProportion(dataPoint *models.PricePoint) bool {
	return true
}

func (ctx *processingContext) shouldPurchase(dataPoint *models.PricePoint) bool {
	return ctx.passFallInPrice(dataPoint) &&
		ctx.passLastBuyCheck(dataPoint) &&
		ctx.passNegativeTrend(dataPoint) &&
		ctx.passProportion(dataPoint)
}

func (ctx *processingContext) makePurchase(data *models.PricePoint) error {
	fmt.Println("invest here")
	return nil
}
