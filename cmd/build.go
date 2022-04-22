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
	Headers struct {
		ContentLength string `json:"Content-Length"`
		ContentType   string `json:"Content-Type"`
		Host          string `json:"Host"`
	} `json:"headers"`
	JSON struct {
		Branch string `json:"branch"`
		Name   string `json:"name"`
	} `json:"json"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

func NewBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "build",
		Aliases: []string{"b"},
		Short:   "build",
		Long:    `This command can be used to create/get and show build(s)`,
	}
	cmd.AddCommand(
		NewBuildGetCmd(),
		NewBuildGetAllCmd(),
		NewBuildProgressCmd(),
		NewBuildLogsCmd(),
		//NewBuildCreateCmd(),
	)
	return cmd
}

func NewBuildGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get --code=[build-code]",
		Long:  `This command can be used to get build`,

		Run: func(cmd *cobra.Command, args []string) {

			code, _ := cmd.Flags().GetString("code")

			if code != "" && len(args) <= 0 {

				body := httpGet(client, SAP_CLOUD_API_URL+"/builds/"+code)
				var build Build

				if err := json.Unmarshal(body, &build); err != nil { // Parse []byte to go struct pointer
					log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
				} else {
					fmt.Println(utils.PrettyPrintJSON(build))
				}

			} else {
				fmt.Println(utils.UnknownCommandMsg("build get"))
				return
			}
		},
	}
	cmd.PersistentFlags().String("code", "", "To get build by its code")
	return cmd
}

func NewBuildGetAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-all",
		Short: "get-all",
		Long:  `This command can be used to get all builds`,

		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				fmt.Println(utils.UnknownCommandMsg("build get-all"))
				return
			}

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

		Run: func(cmd *cobra.Command, args []string) {

			code, _ := cmd.Flags().GetString("code")

			if code != "" && len(args) <= 0 {

				buildProgress := getBuildProgress(code)
				if buildProgress.BuildStatus != "" {
					fmt.Println("------------------------------------------------")
					fmt.Printf("progress: %d\ttasks: %d\tstatus: %s", buildProgress.Percentage, buildProgress.NumberOfTasks, buildProgress.BuildStatus)
				}

			} else {
				fmt.Println(utils.UnknownCommandMsg("build progress"))
				return
			}
		},
	}
	cmd.PersistentFlags().String("code", "", "To get build progress by its code")
	return cmd
}

func NewBuildLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "logs --code=[build-code]",
		Long:  `This command can be used to get build logs`,

		Run: func(cmd *cobra.Command, args []string) {

			code, _ := cmd.Flags().GetString("code")
			zipFileName := "build-" + code + ".zip"

			if code != "" {

				body := httpGet(client, SAP_CLOUD_API_URL+"/builds/"+code+"/logs")

				fmt.Println("[STARTING!...] download logs for build :" + code)

				err := utils.DownloadZipFile(LOGS_DIR, zipFileName, body)
				if err != nil {
					log.Fatalf("[FAILED!...] Failed downloading logs: %s", err)
				}

				fmt.Println("[FINISHED!...]. Logs saved into " + LOGS_DIR + "/" + zipFileName)

			} else {
				fmt.Println(utils.UnknownCommandMsg("build logs"))
				return
			}
		},
	}
	cmd.PersistentFlags().String("code", "", "To get build logs by its code")
	return cmd
}

// not complete yet !!! catch code returned
func NewBuildCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create --branch=[branch] --name=[name]",
		Long:  `This command can be used to create build`,

		Run: func(cmd *cobra.Command, args []string) {

			branch, _ := cmd.Flags().GetString("branch")
			name, _ := cmd.Flags().GetString("name")

			if branch != "" && name != "" && len(args) <= 0 {
				reqBody, err := json.Marshal(map[string]string{
					"branch": branch,
					"name":   name,
				})
				if err != nil {
					return
				}

				fmt.Println("[STARTING!...] Build branch " + branch)
				body := httpPost(client, SAP_CLOUD_API_URL+"/builds", reqBody)
				var buildCreateResp BuildCreateResp
				if err := json.Unmarshal(body, &buildCreateResp); err != nil {
					log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
				} else {
					fmt.Println(utils.PrettyPrintJSON(buildCreateResp))
				}

				isFinished := false

				for !isFinished {

					buildProgress := getBuildProgress(branch)

					fmt.Println("------------------------------------------------")
					fmt.Printf("progress: %d \ttasks: %d\tstatus: %s", buildProgress.Percentage, buildProgress.NumberOfTasks, buildProgress.BuildStatus)

					if buildProgress.BuildStatus == "SUCCESS" {
						isFinished = true
						_, err := utils.WriteFile(BUILDS_DIR, "last_success", buildProgress.BuildCode)
						if err != nil {
							log.Fatalf("[FAILED!...] Failed saving build: %s", err)
						}
						fmt.Println("\n[FINISHED!...] Build Success " + BUILDS_DIR + "/" + "last_success")

					} else if buildProgress.BuildStatus == "FAILED" {
						log.Fatalf("[FAILED!...] Build Failed :(")
					}

				}
			} else {
				fmt.Println(utils.UnknownCommandMsg("build create"))
				return
			}

		},
	}
	cmd.PersistentFlags().String("branch", "", "Branch to build")
	cmd.PersistentFlags().String("name", "", "Build's name")
	return cmd
}

func getBuildProgress(code string) (buildProgress BuildProgress) {

	body := httpGet(client, SAP_CLOUD_API_URL+"/builds/"+code+"/progress")
	if err := json.Unmarshal(body, &buildProgress); err != nil {
		log.Fatalf("[ERROR!...] Couldn't unmarshal JSON")
	}
	return
}
