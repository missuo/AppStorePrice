/*
 * @Author: Vincent Young
 * @Date: 2023-10-30 18:28:13
 * @LastEditors: Vincent Young
 * @LastEditTime: 2023-10-30 23:01:20
 * @FilePath: /AppStoreAPI/main.go
 * @Telegram: https://t.me/missuo
 *
 * Copyright © 2023 by Vincent, All Rights Reserved.
 */

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type Item struct {
	Name  string
	Price string
}

func fetchPrice(countryCode string, appId string) ([]map[string]string, error) {
	url := launcher.New().Headless(true).MustLaunch()

	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	pageUrl := fmt.Sprintf("https://apps.apple.com/%s/app/id%s", countryCode, appId)
	page := browser.MustPage(pageUrl)

	element := page.MustElementX("/html/body/div[3]/main/div[2]/section[8]/div[1]/dl/div[9]/dd/ol/div/button")
	element.MustClick()

	liElements := page.MustElementsX("/html/body/div[3]/main/div[2]/section[8]/div[1]/dl/div[9]/dd/ol/div/li")

	var result []map[string]string

	for _, liElement := range liElements {
		text, err := liElement.Text()
		if err != nil {
			log.Fatalf("Failed to get text from element: %v", err)
		}
		parts := strings.Split(text, "\n")
		if len(parts) >= 2 {
			item := make(map[string]string)
			item["name"] = parts[0]
			item["price"] = strings.Replace(parts[1], "\t", "", -1)
			result = append(result, item)
		}
	}
	return result, nil
}

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Creating a new context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Creating a new channel to receive the result
		ch := make(chan struct{}, 1)

		// Using a goroutine to execute the handler
		go func() {
			c.Next()
			ch <- struct{}{}
		}()

		select {
		case <-ch:
			// The request has completed, return directly
			return
		case <-ctx.Done():
			// The request has timed out
			c.AbortWithStatus(http.StatusGatewayTimeout)
			fmt.Println("Request Timeout")
		}
	}
}

func main() {
	fmt.Println("App Store Pricer is running...")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(TimeoutMiddleware(time.Second * 10))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Welcome to App Store Pricer.",
		})
	})
	r.GET("/as", func(c *gin.Context) {
		countryCode := c.Query("countrycode")
		appId := c.Query("appid")
		if countryCode == "" {
			countryCode = "US"
		}

		if appId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "Bad Parameter",
			})
			return
		}
		result, _ := fetchPrice(countryCode, appId)

		if result != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"data": result,
			})
		} else {
			c.JSON(http.StatusNoContent, gin.H{
				"code": http.StatusNoContent,
				"data": "No Content",
			})
		}

	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Path not found",
		})
	})
	r.Run(":7777")
}
