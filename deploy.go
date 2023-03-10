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
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/cobra"
	"log"
)

type deployCmdFlags struct {
	EventBridgeFilePath string
}

var deployFlags deployCmdFlags

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy or create trigger",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%+v\n", deployFlags)
		eb, err := loadEventBridgeFile(deployFlags.EventBridgeFilePath)
		if err != nil {
			return err
		}

		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}

		stsService := sts.NewFromConfig(cfg)
		identity, err := stsService.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
		if err != nil {
			return nil
		}
		accountID := *identity.Account

		if eb.Rule != nil {
			err = putRule(cfg, accountID, *eb)
			if err != nil {
				return err
			}

			err = putTarget(cfg, cfg.Region, accountID, *eb)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func putRule(cfg aws.Config, accountId string, eb EventBridge) error {
	eventBridgeSvc := eventbridge.NewFromConfig(cfg)

	if eb.Rule.EventPattern != nil {
		ep, err := json.Marshal(eb.Rule.EventPattern)
		if err != nil {
			return err
		}
		_, err = eventBridgeSvc.PutRule(context.TODO(), &eventbridge.PutRuleInput{
			Name:         aws.String(eb.Rule.Name),
			Description:  aws.String(eb.Rule.Description),
			EventBusName: aws.String(eb.Rule.EventBusName),
			EventPattern: aws.String(string(ep)),
		})
		return err
	} else if eb.Rule.ScheduleExpression != "" {
		_, err := eventBridgeSvc.PutRule(context.TODO(), &eventbridge.PutRuleInput{
			Name:               aws.String(eb.Rule.Name),
			Description:        aws.String(eb.Rule.Description),
			EventBusName:       aws.String(eb.Rule.EventBusName),
			ScheduleExpression: aws.String(eb.Rule.ScheduleExpression),
		})
		return err
	}
	return nil
}

func putTarget(cfg aws.Config, region string, accountId string, eb EventBridge) error {
	eventBridgeSvc := eventbridge.NewFromConfig(cfg)

	eventBridgeSvc.PutTargets(context.TODO(), &eventbridge.PutTargetsInput{
		Rule: aws.String(eb.Rule.Name),
		Targets: []ebTypes.Target{
			{
				Arn: aws.String(fmt.Sprintf("arn:aws:lambda:%s:%s:function:%s", region, accountId, eb.LambdaTarget.Name)),
				Id:  aws.String(fmt.Sprintf("%s-%s", eb.Rule.Name, eb.LambdaTarget.Name)),
			},
		},

		EventBusName: nil,
	})
	return nil
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	deployCmd.Flags().StringVarP(&deployFlags.EventBridgeFilePath, "eventbridge", "e", "eventbridge.json", "EventBridge file path")
}
