package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/dhamith93/SyMon/internal/api"
	"github.com/dhamith93/SyMon/internal/auth"
	"github.com/dhamith93/SyMon/internal/config"
	"github.com/dhamith93/SyMon/internal/logger"
	"github.com/dhamith93/SyMon/internal/monitor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	config := config.GetConfig("config.json")

	if config.LogFileEnabled {
		file, err := os.OpenFile(config.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		log.SetOutput(file)
	}

	var name, value, unit string

	initPtr := flag.Bool("init", false, "Initialize agent")
	customPtr := flag.Bool("custom", false, "Send custom metrics")
	flag.StringVar(&name, "name", "", "Name of the metric")
	flag.StringVar(&unit, "unit", "", "Unit of the metric")
	flag.StringVar(&value, "value", "", "Value of the metric")
	flag.Parse()

	if *initPtr {
		initAgent(&config)
		return
	} else if *customPtr {
		if len(name) > 0 && len(value) > 0 && len(unit) > 0 {
			sendCustomMetric(name, unit, value, &config)
		} else {
			fmt.Println("Metric name, unit, and value all required")
		}
		return
	}

	ticker := time.NewTicker(time.Minute)
	quit := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			select {
			case <-ticker.C:
				monitorData := monitor.MonitorAsJSON(&config)
				sendMonitorData(monitorData, &config)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	wg.Wait()
	fmt.Println("Exiting")
}

func initAgent(config *config.Config) {
	conn, c, ctx, cancel := createClient(config)
	defer conn.Close()
	defer cancel()
	response, err := c.InitAgent(ctx, &api.ServerInfo{
		ServerId: config.ServerId,
		Timezone: monitor.GetSystem().TimeZone,
	})
	if err != nil {
		logger.Log("error", "error adding agent: "+err.Error())
		os.Exit(1)
	}
	fmt.Printf("%s \n", response.Body)
}

func sendMonitorData(monitorData string, config *config.Config) {
	conn, c, ctx, cancel := createClient(config)
	defer conn.Close()
	defer cancel()
	_, err := c.HandleMonitorData(ctx, &api.MonitorData{MonitorData: monitorData})
	if err != nil {
		logger.Log("error", "error sending data: "+err.Error())
		os.Exit(1)
	}
}

func sendCustomMetric(name string, unit string, value string, config *config.Config) {
	customMetric := monitor.CustomMetric{
		Name:     name,
		Unit:     unit,
		Value:    value,
		Time:     strconv.FormatInt(time.Now().Unix(), 10),
		ServerId: config.ServerId,
	}
	jsonData, err := json.Marshal(&customMetric)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	conn, c, ctx, cancel := createClient(config)
	defer conn.Close()
	defer cancel()
	_, err = c.HandleCustomMonitorData(ctx, &api.MonitorData{MonitorData: string(jsonData)})
	if err != nil {
		logger.Log("error", "error sending custom data: "+err.Error())
		os.Exit(1)
	}
}

func generateToken() string {
	token, err := auth.GenerateJWT()
	if err != nil {
		logger.Log("error", "error generating token: "+err.Error())
		os.Exit(1)
	}
	return token
}

func createClient(config *config.Config) (*grpc.ClientConn, api.MonitorDataServiceClient, context.Context, context.CancelFunc) {
	conn, err := grpc.Dial(config.MonitorEndpoint, grpc.WithInsecure())
	if err != nil {
		logger.Log("error", "connection error: "+err.Error())
		os.Exit(1)
	}
	c := api.NewMonitorDataServiceClient(conn)
	token := generateToken()
	ctx, cancel := context.WithTimeout(metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{"jwt": token})), time.Second*1)
	return conn, c, ctx, cancel
}
