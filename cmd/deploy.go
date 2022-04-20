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

var deployCmd = &cobra.Command{
	Use:     "deploy",
	Aliases: []string{"d"},
	Short:   "deploy",
	Long:    `This command can be used to create/get/cancel and show deployment(s)`,
}

var deployGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get --code=[deploy-code]",
	Long:  `This command can be used to get deployment`,

	Run: func(cmd *cobra.Command, args []string) {

		code, _ := cmd.Flags().GetString("code")

		if code != "" && len(args) <= 0 {

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

var deployGetAllCmd = &cobra.Command{
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

var deployProgressCmd = &cobra.Command{
	Use:   "progress",
	Short: "progress --code=[deploy-code]",
	Long:  `This command can be used to get deploy progress`,

	Run: func(cmd *cobra.Command, args []string) {

		code, _ := cmd.Flags().GetString("code")

		if code != "" && len(args) <= 0 {
			body := httpGet(client, SAP_CLOUD_API_URL+"/deployments/"+code+"/progress")

			deploymentProgress := getDeployProgress(body)
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

var deployGetCancellationOptionsCmd = &cobra.Command{
	Use:   "get-cancel-opts",
	Short: "get-cancel-opts --code=[deploy-code]",
	Long:  `This command can be used to get deployment cancellation options`,

	Run: func(cmd *cobra.Command, args []string) {

		code, _ := cmd.Flags().GetString("code")

		if code != "" && len(args) <= 0 {
			body := httpGet(client, SAP_CLOUD_API_URL+"/deployments/"+code+"/cancellationoptions")

			fmt.Println(string(body))

		} else {
			fmt.Println(utils.UnknownCommandMsg("deploy get-cancel-opts"))
			return
		}
	},
}

var deployCreateCancellationCmd = &cobra.Command{
	Use:   "cancel",
	Short: "cancel --code=[deploy-code] --rollbackDatabase=[true | false | default: false]",
	Long:  `This command can be used to cancel a deployment`,

	Run: func(cmd *cobra.Command, args []string) {

		code, _ := cmd.Flags().GetString("code")
		rollbackDatabase, _ := cmd.Flags().GetString("rollbackDatabase")

		if code != "" && len(args) <= 0 {
			reqBody, err := json.Marshal(map[string]string{
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

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.AddCommand(deployGetCmd)
	deployCmd.AddCommand(deployGetAllCmd)
	deployCmd.AddCommand(deployProgressCmd)
	deployCmd.AddCommand(deployGetCancellationOptionsCmd)
	deployCmd.AddCommand(deployCreateCancellationCmd)

	//deploy get flags
	deployGetCmd.PersistentFlags().String("code", "", "To get deployment by its code")
	deployProgressCmd.PersistentFlags().String("code", "", "To get deployment progress by its code")
	deployGetCancellationOptionsCmd.PersistentFlags().String("code", "", "To get deployment cancel options by its code")
	deployCreateCancellationCmd.PersistentFlags().String("code", "", "To cancel deployment by its code")
	deployCreateCancellationCmd.PersistentFlags().Bool("rollbackDatabase", false, "To cancel deployment by its code")
}

func getDeployProgress(body []byte) (deploymentProgress DeploymentProgress) {

	if err := json.Unmarshal(body, &deploymentProgress); err != nil {
		log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
	}
	return
}
