package main

import (
	"encoding/json"
	"os"
)

type EventBridge struct {
	Rule         *Rule         `json:"rule"`
	LambdaTarget *LambdaTarget `json:"lambdaTarget"`
}

type Rule struct {
	Name               string      `json:"name"`
	ScheduleExpression string      `json:"scheduleExpression"`
	Description        string      `json:"description"`
	EventBusName       string      `json:"eventBusName"`
	EventPattern       interface{} `json:"eventPattern"`
}

type LambdaTarget struct {
	Name string `json:"name"`
}

func loadEventBridgeFile(ebFilePath string) (*EventBridge, error) {
	f, err := os.Open(ebFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var eb EventBridge
	err = json.NewDecoder(f).Decode(&eb)
	if err != nil {
		return nil, err
	}
	return &eb, nil
}
