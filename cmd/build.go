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
	Code   string
	Name   string
	Branch string
}

var ob *buildOptions
var validBuildArgs = []string{"get", "getAll", "logs", "progress", "create"}

func NewBuildCmd() *cobra.Command {
	ob = &buildOptions{}
	cmd := &cobra.Command{
		Use:                   "build [command]",
		Aliases:               []string{"b"},
		Short:                 "Trigger build on SAP Cloud",
		Long:                  `This command can be used to create/get and show build(s)`,
		ValidArgs:             validBuildArgs,
		Args:                  utils.IsOneAndOnlyValidArgs,
		DisableFlagsInUseLine: true,
		Run:                   func(cmd *cobra.Command, args []string) {},
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
		Args:  cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {

			body := utils.HttpGet(client, SAP_CLOUD_API_URL+"/builds/"+ob.Code, API_TOKEN)
			var build Build

			if err := json.Unmarshal(body, &build); err != nil { // Parse []byte to go struct pointer
				log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
			} else {
				fmt.Println(utils.PrettyFormatJSON(build))
			}

		},
	}

	cmd.PersistentFlags().StringVarP(&ob.Code, "code", "c", "", "To get build by its code")
	cmd.MarkPersistentFlagRequired("code")

	return cmd
}

func NewBuildGetAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getAll",
		Short: "getAll",
		Long:  `This command can be used to get all builds`,
		Args:  cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {

			body := utils.HttpGet(client, SAP_CLOUD_API_URL+"/builds", API_TOKEN)
			var builds Builds
			if err := json.Unmarshal(body, &builds); err != nil { // Parse []byte to go struct pointer
				log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
			} else {
				fmt.Println(utils.PrettyFormatJSON(builds))
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
		Args:  cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {

			buildProgress := getBuildProgress(ob.Code)
			if buildProgress.BuildStatus != "" {
				fmt.Println("------------------------------------------------")
				fmt.Printf("progress: %d\ttasks: %d\tstatus: %s", buildProgress.Percentage, buildProgress.NumberOfTasks, buildProgress.BuildStatus)
			}

		},
	}

	cmd.PersistentFlags().StringVarP(&ob.Code, "code", "c", "", "To get build progress by its code")
	cmd.MarkPersistentFlagRequired("code")
	return cmd
}

func NewBuildLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "logs --code=[build-code]",
		Long:  `This command can be used to get build logs`,
		Args:  cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {

			//code, _ := cmd.Flags().GetString("code")
			zipFileName := "build-" + ob.Code + ".zip"

			body := utils.HttpGet(client, SAP_CLOUD_API_URL+"/builds/"+ob.Code+"/logs", API_TOKEN)
			fmt.Println("[STARTING!...] download logs for build :" + ob.Code)

			err := utils.DownloadZipFile(LOGS_DIR, zipFileName, body)
			if err != nil {
				log.Fatalf("[FAILED!...] Failed downloading logs: %s", err)
			}

			fmt.Println("[FINISHED!...]. Logs saved into " + LOGS_DIR + zipFileName)

		},
	}

	cmd.PersistentFlags().StringVarP(&ob.Code, "code", "c", "", "To get build logs by its code")
	cmd.MarkPersistentFlagRequired("code")
	return cmd
}

// not complete yet !!! catch code returned
func NewBuildCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create --branch=[branch] --name=[name]",
		Long:  `This command can be used to create build`,
		Args:  cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {

			reqBody, err := json.Marshal(map[string]string{
				"branch": ob.Branch,
				"name":   ob.Name,
			})
			if err != nil {
				return
			}

			fmt.Println("[STARTING!...] Build branch " + ob.Branch)
			body := utils.HttpPost(client, SAP_CLOUD_API_URL+"/builds", API_TOKEN, reqBody)
			var buildCreateResp BuildCreateResp
			if err := json.Unmarshal(body, &buildCreateResp); err != nil {
				log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
			} else {
				fmt.Println(utils.PrettyFormatJSON(buildCreateResp))
			}

			buildCode := buildCreateResp.Code
			isFinished := false

			for !isFinished {

				buildProgress := getBuildProgress(buildCode)

				fmt.Println("------------------------------------------------")
				fmt.Printf("Status: %s\tProgress: %d%s\ttasks: %d", buildProgress.BuildStatus, buildProgress.Percentage, "%", buildProgress.NumberOfTasks)

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
				time.Sleep(time.Second * 3)

			}

		},
	}

	cmd.PersistentFlags().StringVarP(&ob.Branch, "branch", "b", "", "branch to build")
	cmd.PersistentFlags().StringVarP(&ob.Name, "name", "n", "", "build name")
	cmd.MarkPersistentFlagRequired("branch")
	cmd.MarkPersistentFlagRequired("name")
	return cmd
}

func getBuildProgress(code string) (buildProgress BuildProgress) {

	body := utils.HttpGet(client, SAP_CLOUD_API_URL+"/builds/"+code+"/progress", API_TOKEN)
	if err := json.Unmarshal(body, &buildProgress); err != nil {
		log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
	}
	return
}
