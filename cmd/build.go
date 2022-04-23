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

//build model
type Build struct {
	ApplicationCode              string    `json:"applicationCode"`
	ApplicationDefinitionVersion string    `json:"applicationDefinitionVersion"`
	Branch                       string    `json:"branch"`
	BuildEndTimestamp            time.Time `json:"buildEndTimestamp"`
	BuildStartTimestamp          time.Time `json:"buildStartTimestamp"`
	BuildVersion                 string    `json:"buildVersion"`
	Code                         string    `json:"code"`
	CreatedBy                    string    `json:"createdBy"`
	Name                         string    `json:"name"`
	Properties                   []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"properties"`
	Status string `json:"status"`
}

//build model array
type Builds struct {
	Count int     `json:"count"`
	Value []Build `json:"value"`
}

//build progress model
type BuildProgress struct {
	ErrorMessage  string `json:"errorMessage"`
	NumberOfTasks int    `json:"numberOfTasks"`
	Percentage    int    `json:"percentage"`
	BuildStatus   string `json:"buildStatus"`
	BuildCode     string `json:"buildCode"`
}
type BuildCreateResp struct {
	Code string `json:"code"`
}

type buildOptions struct {
	code   string
	name   string
	branch string
}

var ob *buildOptions

func isValidBuildArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return errors.New("Requires at least one arg")
	}

	if args[0] == "get" || args[0] == "get-all" || args[0] == "logs" || args[0] == "progress" || args[0] == "create" {
		return nil
	} else {
		return errors.New("Invalid argument: " + args[0])
	}
}

func NewBuildCmd() *cobra.Command {
	ob = &buildOptions{}
	cmd := &cobra.Command{
		Use:     "build",
		Aliases: []string{"b"},
		Short:   "build",
		Long:    `This command can be used to create/get and show build(s)`,
		Args:    isValidBuildArgs,
		Run:     func(cmd *cobra.Command, args []string) {},
	}
	cmd.AddCommand(
		NewBuildGetCmd(),
		NewBuildGetAllCmd(),
		NewBuildProgressCmd(),
		NewBuildLogsCmd(),
		NewBuildCreateCmd(),
	)
	return cmd
}

func NewBuildGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get --code=[build-code]",
		Long:  `This command can be used to get build`,
		Args:  cobra.MaximumNArgs(0),

		Run: func(cmd *cobra.Command, args []string) {

			body := httpGet(client, SAP_CLOUD_API_URL+"/builds/"+ob.code)
			var build Build

			if err := json.Unmarshal(body, &build); err != nil { // Parse []byte to go struct pointer
				log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
			} else {
				fmt.Println(utils.PrettyPrintJSON(build))
			}

		},
	}

	cmd.PersistentFlags().StringVar(&ob.code, "code", "", "To get build by its code")
	cmd.MarkPersistentFlagRequired("code")

	return cmd
}

func NewBuildGetAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-all",
		Short: "get-all",
		Long:  `This command can be used to get all builds`,
		Args:  cobra.MaximumNArgs(0),

		Run: func(cmd *cobra.Command, args []string) {

			body := httpGet(client, SAP_CLOUD_API_URL+"/builds")
			var builds Builds
			if err := json.Unmarshal(body, &builds); err != nil { // Parse []byte to go struct pointer
				log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
			} else {
				fmt.Println(utils.PrettyPrintJSON(builds))
			}
		},
	}

	return cmd
}

func NewBuildProgressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "progress",
		Short: "progress --code=[build-code]",
		Long:  `This command can be used to get build progress`,
		Args:  cobra.MaximumNArgs(0),

		Run: func(cmd *cobra.Command, args []string) {

			buildProgress := getBuildProgress(ob.code)
			if buildProgress.BuildStatus != "" {
				fmt.Println("------------------------------------------------")
				fmt.Printf("progress: %d\ttasks: %d\tstatus: %s", buildProgress.Percentage, buildProgress.NumberOfTasks, buildProgress.BuildStatus)
			}

		},
	}

	cmd.PersistentFlags().StringVar(&ob.code, "code", "", "To get build progress by its code")
	cmd.MarkPersistentFlagRequired("code")
	return cmd
}

func NewBuildLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "logs --code=[build-code]",
		Long:  `This command can be used to get build logs`,
		Args:  cobra.MaximumNArgs(0),

		Run: func(cmd *cobra.Command, args []string) {

			//code, _ := cmd.Flags().GetString("code")
			zipFileName := "build-" + ob.code + ".zip"

			body := httpGet(client, SAP_CLOUD_API_URL+"/builds/"+ob.code+"/logs")
			fmt.Println("[STARTING!...] download logs for build :" + ob.code)

			err := utils.DownloadZipFile(LOGS_DIR, zipFileName, body)
			if err != nil {
				log.Fatalf("[FAILED!...] Failed downloading logs: %s", err)
			}

			fmt.Println("[FINISHED!...]. Logs saved into " + LOGS_DIR + zipFileName)

		},
	}

	cmd.PersistentFlags().StringVar(&ob.code, "code", "", "To get build logs by its code")
	cmd.MarkPersistentFlagRequired("code")
	return cmd
}

// not complete yet !!! catch code returned
func NewBuildCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create --branch=[branch] --name=[name]",
		Long:  `This command can be used to create build`,
		Args:  cobra.MaximumNArgs(0),

		Run: func(cmd *cobra.Command, args []string) {

			reqBody, err := json.Marshal(map[string]string{
				"branch": ob.branch,
				"name":   ob.name,
			})
			if err != nil {
				return
			}

			fmt.Println("[STARTING!...] Build branch " + ob.branch)
			body := httpPost(client, SAP_CLOUD_API_URL+"/builds", reqBody)
			var buildCreateResp BuildCreateResp
			if err := json.Unmarshal(body, &buildCreateResp); err != nil {
				log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
			} else {
				fmt.Println(utils.PrettyPrintJSON(buildCreateResp))
			}

			buildCode := buildCreateResp.Code
			isFinished := false

			for !isFinished {

				buildProgress := getBuildProgress(buildCode)

				fmt.Println("------------------------------------------------")
				fmt.Printf("progress: %d \ttasks: %d\tstatus: %s", buildProgress.Percentage, buildProgress.NumberOfTasks, buildProgress.BuildStatus)

				if buildProgress.BuildStatus == "SUCCESS" {
					isFinished = true
					_, err := utils.WriteFile(BUILDS_DIR, "last_success", buildProgress.BuildCode)
					if err != nil {
						log.Fatalf("[FAILED!...] Failed saving build: %s", err)
					}
					fmt.Println("\n[FINISHED!...] Build Success " + BUILDS_DIR + "last_success")

				} else if buildProgress.BuildStatus == "FAILED" {
					log.Fatalf("[FAILED!...] Build Failed :(")
				}

			}

		},
	}

	cmd.PersistentFlags().StringVar(&ob.branch, "branch", "", "")
	cmd.PersistentFlags().StringVar(&ob.name, "name", "", "")
	cmd.MarkPersistentFlagRequired("branch")
	cmd.MarkPersistentFlagRequired("name")
	return cmd
}

func getBuildProgress(code string) (buildProgress BuildProgress) {

	body := httpGet(client, SAP_CLOUD_API_URL+"/builds/"+code+"/progress")
	if err := json.Unmarshal(body, &buildProgress); err != nil {
		log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
	}
	return
}
