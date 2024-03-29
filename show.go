/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	ebTypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/spf13/cobra"
)

type showCmdFlags struct {
}

var showFlags showCmdFlags

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show EventBridge rule and target",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		eb, err := loadEventBridgeFile(opt)
		if err != nil {
			return err
		}

		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return fmt.Errorf("unable to load SDK config, %v", err)
		}

		rule, err := getRule(cfg, eb.Rule.Name)
		if err != nil {
			return err
		}
		targets, err := getTargets(cfg, eb.Rule.Name)
		output := &ShowEventBridgeOutput{
			Rule:    rule,
			Targets: targets,
		}
		if err != nil {
			return err
		}

		j, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(j))

		return nil
	},
}

func getTargets(cfg aws.Config, ruleName string) ([]ebTypes.Target, error) {
	eventBridgeSvc := eventbridge.NewFromConfig(cfg)
	targetsOutput, err := eventBridgeSvc.ListTargetsByRule(context.TODO(), &eventbridge.ListTargetsByRuleInput{
		Rule: aws.String(ruleName),
	})
	if err != nil {
		return nil, err
	}

	return targetsOutput.Targets, nil
}

func getRule(cfg aws.Config, ruleName string) (*eventbridge.DescribeRuleOutput, error) {
	eventBridgeSvc := eventbridge.NewFromConfig(cfg)
	ruleOutput, err := eventBridgeSvc.DescribeRule(context.TODO(), &eventbridge.DescribeRuleInput{
		Name: aws.String(ruleName),
	})
	if err != nil {
		return nil, err
	}
	return ruleOutput, nil
}

type ShowEventBridgeOutput struct {
	Rule    *eventbridge.DescribeRuleOutput
	Targets []ebTypes.Target
}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	showCmd.Flags().StringVarP(&opt.EventBridgeFilePath, "eventbridge", "e", "eventbridge.json", "EventBridge file path")
	showCmd.Flags().StringVarP(&opt.TFState, "tfstate", "", "", "tfstate url")
}
