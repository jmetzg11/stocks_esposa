package api

import (
	"fmt"
	"net/http"
	"stocks/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type processingContext struct {
	db         *gorm.DB
	parameters *models.SimulationParameters
	symbol     string
	priceData  *[]models.PricePoint
	startIndex int
}

func processStock(db *gorm.DB, parameters *models.SimulationParameters, symbol string) error {
	ctx := &processingContext{
		db:         db,
		parameters: parameters,
		symbol:     symbol,
	}

	// fetch price slice
	if err := ctx.fetchPriceData(); err != nil {
		return err
	}

	// find starting point in order to calculate negative trend
	if err := ctx.determineStartingPoint(); err != nil {
		return err
	}

	for i := ctx.startIndex; i < len(*ctx.priceData); i++ {
		dataPoint := &(*ctx.priceData)[i]
		// check 4 conditions
		if ctx.shouldPurchase(dataPoint) {
			// make investment
			if err := ctx.makePurchase(dataPoint); err != nil {
				return err
			}
		}
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
