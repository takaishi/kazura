package main

import (
	"context"
	"encoding/json"
	"github.com/fujiwara/tfstate-lookup/tfstate"
	"github.com/kayac/go-config"
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

func loadEventBridgeFile(opt Option) (*EventBridge, error) {
	ctx := context.Background()
	loader := config.New()
	if opt.TFState != "" {
		funcs, err := tfstate.FuncMap(ctx, opt.TFState)
		if err != nil {
			return nil, err
		}
		loader.Funcs(funcs)
	}
	b, err := loader.ReadWithEnv(opt.EventBridgeFilePath)
	if err != nil {
		return nil, err
	}

	var eb EventBridge
	err = json.Unmarshal(b, &eb)
	if err != nil {
		return nil, err
	}

	return &eb, nil
}
