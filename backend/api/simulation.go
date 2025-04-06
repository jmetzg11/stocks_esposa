package api

import (
	"fmt"
	"net/http"
	"stocks/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func processStock(db *gorm.DB, parameters *models.SimulationParameters, symbol string) error {
	var priceData []models.PricePoint
	result := db.Model(&models.Historical{}).
		Select("date, price").
		Where("symbol = ?", symbol).
		Order("date ASC").
		Find(&priceData)

	if result.Error != nil {
		return fmt.Errorf("error querying data for %s: %v", symbol, result.Error)
	}
	if len(priceData) == 0 {
		return fmt.Errorf("no price data found for %s", symbol)
	}

	startDate, err := findStartDate(priceData, parameters.NegativeTrend)
	if err != nil {
		return err
	}

	startIndex, err := getStartIndex(priceData, startDate)
	if err != nil {
		return err
	}

	for i := startIndex; i < len(priceData); i++ {
		dataPoint := &priceData[i]

		if !passLastBuyCheck(db, symbol, dataPoint, parameters) {
			continue
		}

		// check negative trend

		// check proportion of portfolio

		// make purchase if all conditions pass
	}

	return nil
}

func (h *Handler) StartSimulation(c *gin.Context) {
	var parameters models.SimulationParameters
	if err := c.ShouldBindJSON(&parameters); err != nil {
		fmt.Printf("Bad payload was passed %v\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parameters.Version = generateVersion(parameters)

	marketCapFilter := parameters.MarketCap * 1000
	fmt.Println("marketCapFilter", marketCapFilter)
	var symbols []string
	if err := h.db.Model(&models.MarketCap{}).
		Where("market_cap >= ?", marketCapFilter).
		Pluck("symbol", &symbols).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, symbol := range symbols {
		if err := processStock(h.db, &parameters, symbol); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	fmt.Println("Number of symbols", len(symbols))

	// then calculate last price and how many share I have and how much it will cost

	c.JSON(http.StatusOK, gin.H{"Message": "Simulation started"})
}
