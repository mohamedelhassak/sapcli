/*
Copyright © 2022
This file is part of CLI application SAPCLI.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mohamedelhassak/sapcli/utils"
	"github.com/spf13/cobra"
)

//deployment model
type Deployment struct {
	Code                string    `json:"code"`
	CreatedBy           string    `json:"createdBy"`
	CreatedTimestamp    time.Time `json:"createdTimestamp"`
	BuildCode           string    `json:"buildCode"`
	EnvironmentCode     string    `json:"environmentCode"`
	DatabaseUpdateMode  string    `json:"databaseUpdateMode"`
	Strategy            string    `json:"strategy"`
	ScheduledTimestamp  time.Time `json:"scheduledTimestamp"`
	DeployedTimestamp   time.Time `json:"deployedTimestamp"`
	FailedTimestamp     time.Time `json:"failedTimestamp"`
	UndeployedTimestamp time.Time `json:"undeployedTimestamp"`
	Status              string    `json:"status"`
	Cancelation         struct {
		CanceledBy        string    `json:"canceledBy"`
		StartTimestamp    time.Time `json:"startTimestamp"`
		FinishedTimestamp time.Time `json:"finishedTimestamp"`
		Failed            bool      `json:"failed"`
		RollbackDatabase  bool      `json:"rollbackDatabase"`
	} `json:"cancelation"`
}

//deployment model array
type Deployments struct {
	Count int          `json:"count"`
	Value []Deployment `json:"value"`
}

//deployment progress model
type DeploymentProgress struct {
	DeploymentCode   string `json:"deploymentCode"`
	DeploymentStatus string `json:"deploymentStatus"`
	Percentage       int    `json:"percentage"`
}

type DeploymentCancelResp struct {
	Detail string `json:"detail"`
	Title  string `json:"title"`
}

type DeploymentCreateResp struct {
	Code string `json:"code"`
}

type deployOptions struct {
	code               string
	rollbackDatabase   bool
	buildCode          string
	env                string
	strategy           string
	databaseUpdateMode string
}

var od *deployOptions
var validDeployArgs = []string{"get", "getAll", "cancel", "progress", "create", "getCancelOpts"}

func NewDeployCmd() *cobra.Command {
	od = &deployOptions{}
	cmd := &cobra.Command{
		Use:                   "deploy [command]",
		Aliases:               []string{"d"},
		Short:                 "Trigger and manage deployment on SAP Cloud",
		Long:                  `This command can be used to create/get/cancel and show deployment(s)`,
		ValidArgs:             validDeployArgs,
		Args:                  utils.IsOneAndOnlyValidArgs,
		DisableFlagsInUseLine: true,
		Run:                   func(cmd *cobra.Command, args []string) {},
	}
	cmd.AddCommand(
		NewDeployGetCmd(),
		NewDeployGetAllCmd(),
		NewDeployProgressCmd(),
		NewDeployGetCancellationOptionsCmd(),
		NewDeployCreateCancellationCmd(),
		NewDeployCreateCmd(),
	)
	return cmd
}

func NewDeployGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get --code=[deploy-code]",
		Long:  `This command can be used to get deployment`,
		Args:  cobra.ExactArgs(0),

		Run: func(cmd *cobra.Command, args []string) {

			body := utils.HttpGet(client, SAP_CLOUD_API_URL+"/deployments/"+od.code, API_TOKEN)
			var deployment Deployment
			if err := json.Unmarshal(body, &deployment); err != nil { // Parse []byte to go struct pointer
				log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
			} else {
				fmt.Println(utils.PrettyPrintJSON(deployment))
			}
		},
	}
	cmd.PersistentFlags().StringVarP(&od.code, "code", "c", "", "To get deployment by its code")
	cmd.MarkPersistentFlagRequired("code")
	return cmd
}

func NewDeployGetAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getAll",
		Short: "getAll",
		Long:  `This command can be used to get all deployments`,
		Args:  cobra.ExactArgs(0),

		Run: func(cmd *cobra.Command, args []string) {

			body := utils.HttpGet(client, SAP_CLOUD_API_URL+"/deployments", API_TOKEN)

			var deployments Deployments
			if err := json.Unmarshal(body, &deployments); err != nil { // Parse []byte to go struct pointer
				log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
			} else {
				fmt.Println(utils.PrettyPrintJSON(deployments))
			}
		},
	}
	return cmd
}

func NewDeployProgressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "progress",
		Short: "progress --code=[deploy-code]",
		Long:  `This command can be used to get deploy progress`,
		Args:  cobra.ExactArgs(0),

		Run: func(cmd *cobra.Command, args []string) {

			deploymentProgress := getDeployProgress(od.code)
			if deploymentProgress.DeploymentStatus != "" {
				fmt.Println("------------------------------------------------")
				fmt.Printf("progress: %d\tstatus: %s", deploymentProgress.Percentage, deploymentProgress.DeploymentStatus)
			}

		},
	}

	cmd.PersistentFlags().StringVarP(&od.code, "code", "c", "", "To get deployment progress by its code")
	cmd.MarkPersistentFlagRequired("code")
	return cmd
}

func NewDeployGetCancellationOptionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "getCancelOpts",
		Aliases: []string{"getCO"},
		Short:   "getCancelOpts --code=[deploy-code]",
		Long:    `This command can be used to get deployment cancellation options`,
		Args:    cobra.ExactArgs(0),

		Run: func(cmd *cobra.Command, args []string) {

			body := utils.HttpGet(client, SAP_CLOUD_API_URL+"/deployments/"+od.code+"/cancellationoptions", API_TOKEN)

			fmt.Println(string(body))

		},
	}
	cmd.PersistentFlags().StringVarP(&od.code, "code", "c", "", "To get deployment cancel options by its code")
	cmd.MarkPersistentFlagRequired("code")
	return cmd
}

func NewDeployCreateCancellationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel",
		Short: "cancel --code=[deploy-code] --rollbackDatabase=[true | false | default: false]",
		Long:  `This command can be used to cancel a deployment`,
		Args:  cobra.ExactArgs(0),

		Run: func(cmd *cobra.Command, args []string) {

			reqBody, err := json.Marshal(map[string]bool{
				"rollbackDatabase": od.rollbackDatabase,
			})
			if err != nil {
				return
			}

			body := utils.HttpPost(client, SAP_CLOUD_API_URL+"/deployments/"+od.code+"/cancellation", API_TOKEN, reqBody)
			var deploymentCancelResp DeploymentCancelResp
			if err := json.Unmarshal(body, &deploymentCancelResp); err != nil { // Parse []byte to go struct pointer
				log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
			} else {
				fmt.Println(utils.PrettyPrintJSON(deploymentCancelResp))
			}

		},
	}
	cmd.PersistentFlags().StringVarP(&od.code, "code", "c", "", "To cancel deployment by its code")
	cmd.PersistentFlags().BoolVarP(&od.rollbackDatabase, "rollback-database", "r", false, "To cancel deployment by its code ,Values [true | false] default (false)")

	cmd.MarkPersistentFlagRequired("code")
	return cmd
}

var deployCreateDesc = `create --build-code=[build-code] 
		  	 --env=[environment-code]
		  	 --strategy=[strategy] 
		  	 --database-update-mode=[database-update-mode]`

// not complete yet !!! catch deploy code returned
func NewDeployCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: deployCreateDesc,
		Long:  `This command can be used to create a deployment`,
		Args:  cobra.ExactArgs(0),

		Run: func(cmd *cobra.Command, args []string) {

			reqBody, err := json.Marshal(map[string]string{
				"buildCode":          od.buildCode,
				"databaseUpdateMode": od.databaseUpdateMode,
				"environmentCode":    od.env,
				"strategy":           od.strategy,
			})
			if err != nil {
				return
			}

			fmt.Println("[STARTING!...] Deploying build " + od.buildCode)
			body := utils.HttpPost(client, SAP_CLOUD_API_URL+"/deployments", API_TOKEN, reqBody)
			var deploymentCreateResp DeploymentCreateResp
			if err := json.Unmarshal(body, &deploymentCreateResp); err != nil {
				log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
			} else {
				fmt.Println(utils.PrettyPrintJSON(deploymentCreateResp))
			}

			deployCode := deploymentCreateResp.Code
			isFinished := false

			for !isFinished {

				deploymentProgress := getDeployProgress(deployCode)

				fmt.Println("------------------------------------------------")
				fmt.Printf("status: %s\tprogress: %d", deploymentProgress.DeploymentStatus, deploymentProgress.Percentage)

				if deploymentProgress.DeploymentStatus == "DEPLOYED" {
					isFinished = true

					fmt.Println("\n[FINISHED!...] Deployment Success: " + deploymentProgress.DeploymentCode)

				} else if deploymentProgress.DeploymentStatus == "FAILED" {
					log.Fatalf("[FAILED!...] Deploy Failed :(")
				}

			}

		},
	}
	cmd.PersistentFlags().StringVarP(&od.buildCode, "build-code", "c", "", "build to deploy")
	cmd.PersistentFlags().StringVarP(&od.env, "env", "e", "", "target environment ")
	cmd.PersistentFlags().StringVarP(&od.strategy, "strategy", "s", "ROLLING_UPDATE", "deployment strategy, Values [ ROLLING_UPDATE | RECREATE ]")
	cmd.PersistentFlags().StringVarP(&od.databaseUpdateMode, "database-update-mode", "m", "NONE", "database update mode options, Values [ NONE | UPDATE | INITIALIZE ]")

	cmd.MarkPersistentFlagRequired("build-code")
	cmd.MarkPersistentFlagRequired("env")
	return cmd
}

func getDeployProgress(code string) (deploymentProgress DeploymentProgress) {

	body := utils.HttpGet(client, SAP_CLOUD_API_URL+"/deployments/"+code+"/progress", API_TOKEN)
	if err := json.Unmarshal(body, &deploymentProgress); err != nil {
		log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
	}
	return
}
