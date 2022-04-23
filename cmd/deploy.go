package cmd

import (
	"encoding/json"
	"errors"
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

func isValidDeployArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Requires at least one arg")
	}

	if args[0] == "get" || args[0] == "get-all" || args[0] == "cancel" || args[0] == "progress" || args[0] == "get-cancel-opts" {
		return nil
	} else {
		return errors.New("Invalid argument: " + args[0])
	}
}

func NewDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deploy",
		Aliases: []string{"d"},
		Short:   "deploy",
		Long:    `This command can be used to create/get/cancel and show deployment(s)`,
		Args:    isValidDeployArgs,
		Run:     func(cmd *cobra.Command, args []string) {},
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
		Args:  cobra.MinimumNArgs(1),

		Run: func(cmd *cobra.Command, args []string) {

			code, _ := cmd.Flags().GetString("code")

			if len(args) <= 0 {

				body := httpGet(client, SAP_CLOUD_API_URL+"/deployments/"+code)
				var deployment Deployment
				if err := json.Unmarshal(body, &deployment); err != nil { // Parse []byte to go struct pointer
					log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
				} else {
					fmt.Println(utils.PrettyPrintJSON(deployment))
				}

			} else {
				fmt.Println(utils.UnknownCommandMsg("deploy get"))
				return
			}
		},
	}
	cmd.PersistentFlags().String("code", "", "To get deployment by its code")
	cmd.MarkPersistentFlagRequired("code")
	return cmd
}

func NewDeployGetAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-all",
		Short: "get-all",
		Long:  `This command can be used to get all deployments`,

		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				fmt.Println(utils.UnknownCommandMsg("deploy get-all"))
				return
			}

			body := httpGet(client, SAP_CLOUD_API_URL+"/deployments")

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

		Run: func(cmd *cobra.Command, args []string) {

			code, _ := cmd.Flags().GetString("code")

			if len(args) <= 0 {

				deploymentProgress := getDeployProgress(code)
				if deploymentProgress.DeploymentStatus != "" {
					fmt.Println("------------------------------------------------")
					fmt.Printf("progress: %d\tstatus: %s", deploymentProgress.Percentage, deploymentProgress.DeploymentStatus)
				}

			} else {
				fmt.Println(utils.UnknownCommandMsg("deploy progress"))
				return
			}
		},
	}

	cmd.PersistentFlags().String("code", "", "To get deployment progress by its code")
	cmd.MarkPersistentFlagRequired("code")
	return cmd
}

func NewDeployGetCancellationOptionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-cancel-opts",
		Short: "get-cancel-opts --code=[deploy-code]",
		Long:  `This command can be used to get deployment cancellation options`,

		Run: func(cmd *cobra.Command, args []string) {

			code, _ := cmd.Flags().GetString("code")

			if len(args) <= 0 {
				body := httpGet(client, SAP_CLOUD_API_URL+"/deployments/"+code+"/cancellationoptions")

				fmt.Println(string(body))

			} else {
				fmt.Println(utils.UnknownCommandMsg("deploy get-cancel-opts"))
				return
			}
		},
	}
	cmd.PersistentFlags().String("code", "", "To get deployment cancel options by its code")
	cmd.MarkPersistentFlagRequired("code")
	return cmd
}

func NewDeployCreateCancellationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel",
		Short: "cancel --code=[deploy-code] --rollbackDatabase=[true | false | default: false]",
		Long:  `This command can be used to cancel a deployment`,

		Run: func(cmd *cobra.Command, args []string) {

			code, _ := cmd.Flags().GetString("code")
			rollbackDatabase, _ := cmd.Flags().GetBool("rollback-database")

			if len(args) <= 0 {
				reqBody, err := json.Marshal(map[string]bool{
					"rollbackDatabase": rollbackDatabase,
				})
				if err != nil {
					return
				}

				body := httpPost(client, SAP_CLOUD_API_URL+"/deployments/"+code+"/cancellation", reqBody)
				var deploymentCancelResp DeploymentCancelResp
				if err := json.Unmarshal(body, &deploymentCancelResp); err != nil { // Parse []byte to go struct pointer
					log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
				} else {
					fmt.Println(utils.PrettyPrintJSON(deploymentCancelResp))
				}

			} else {
				fmt.Println(utils.UnknownCommandMsg("deploy cancel"))
				return
			}
		},
	}
	cmd.PersistentFlags().String("code", "", "To cancel deployment by its code")
	cmd.PersistentFlags().Bool("rollback-database", false, "To cancel deployment by its code ,Values [true | false] default (false)")

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

		Run: func(cmd *cobra.Command, args []string) {
			buildCode, _ := cmd.Flags().GetString("build-code")
			env, _ := cmd.Flags().GetString("env")
			strategy, _ := cmd.Flags().GetString("strategy")
			databaseUpdateMode, _ := cmd.Flags().GetString("database-update-mode")

			if buildCode != "" && env != "" && len(args) <= 0 {
				reqBody, err := json.Marshal(map[string]string{
					"buildCode":          buildCode,
					"databaseUpdateMode": databaseUpdateMode,
					"environmentCode":    env,
					"strategy":           strategy,
				})
				if err != nil {
					return
				}

				fmt.Println("[STARTING!...] Deploying build " + buildCode)
				body := httpPost(client, SAP_CLOUD_API_URL+"/deployments", reqBody)
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
					fmt.Printf("progress: %d\tstatus: %s", deploymentProgress.Percentage, deploymentProgress.DeploymentStatus)

					if deploymentProgress.DeploymentStatus == "DEPLOYED" {
						isFinished = true

						fmt.Println("\n[FINISHED!...] Deployment Success: " + deploymentProgress.DeploymentCode)

					} else if deploymentProgress.DeploymentStatus == "FAILED" {
						log.Fatalf("[FAILED!...] Deploy Failed :(")
					}

				}
			} else {
				fmt.Println(utils.UnknownCommandMsg("deploy create"))
				return
			}
		},
	}
	cmd.PersistentFlags().String("build-code", "", "build to deploy")
	cmd.PersistentFlags().String("env", "", "target environment ")
	cmd.PersistentFlags().String("strategy", "ROLLING_UPDATE", "deployment strategy, Values [ ROLLING_UPDATE | RECREATE ]")
	cmd.PersistentFlags().String("database-update-mode", "NONE", "database update mode options, Values [ NONE | UPDATE | INITIALIZE ]")

	cmd.MarkPersistentFlagRequired("build-code")
	cmd.MarkPersistentFlagRequired("env")
	return cmd
}

func getDeployProgress(code string) (deploymentProgress DeploymentProgress) {

	body := httpGet(client, SAP_CLOUD_API_URL+"/deployments/"+code+"/progress")
	if err := json.Unmarshal(body, &deploymentProgress); err != nil {
		log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
	}
	return
}
