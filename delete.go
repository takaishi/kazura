/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/spf13/cobra"
	"log"
)

type deleteCmdFlags struct {
	EventBridgeFilePath string
}

var deleteFlags deleteCmdFlags

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete trigger",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		eb, err := loadEventBridgeFile(deleteFlags.EventBridgeFilePath)
		if err != nil {
			return err
		}

		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}

		eventBridgeSvc := eventbridge.NewFromConfig(cfg)
		_, err = eventBridgeSvc.RemoveTargets(context.TODO(), &eventbridge.RemoveTargetsInput{
			Ids:  []string{fmt.Sprintf("%s-%s", eb.Rule.Name, eb.LambdaTarget.Name)},
			Rule: aws.String(eb.Rule.Name),
		})
		if err != nil {
			return err
		}
		_, err = eventBridgeSvc.DeleteRule(context.TODO(), &eventbridge.DeleteRuleInput{
			Name: aws.String(eb.Rule.Name),
		})
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	deleteCmd.Flags().StringVarP(&deleteFlags.EventBridgeFilePath, "eventbridge", "e", "eventbridge.json", "EventBridge file path")
}
