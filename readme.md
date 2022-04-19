# SAPCLI
A sample SAP CLI implementation in Go

[![GitHub release](https://img.shields.io/github/release/moul/banner.svg)](https://github.com/mohamedelhassak/sapcli/releases)

[![License](https://img.shields.io/badge/license-Apache--2.0%20%2F%20MIT-%2397ca00.svg)](https://github.com/mohamedelhassak/sapcli/blob/main/LICENSE)
![Made by Mohamed El hassak](https://img.shields.io/badge/made%20by-Mohamed%20El%20hassak-blue.svg?style=flat)

## One Step Before Go 🙂 <a name="requis"></a>
* having Go 1.14 or higher installed
* Having SAP Cloud Subscription ID
* Having SAP Cloud API Token
  

## Easy To Setup 😉 <a name="setup"></a>
1.  clone git repo 

```shell
git clone https://github.com/mohamedelhassak/sapcli.git
```
2.  Build project

```go
cd sapcli
#linux
env GOOS=linux GOARCH=amd64 go build -o build/sapcli
#windows
env GOOS=windows GOARCH=amd64 go build -o build/sapcli.exe
```
3. setup env var ``SAPCLI_WORK_DIR``
* ##### For windows, Run :
```powershell
set SAPCLI_WORK_DIR=/your/path
```
* ##### For linux, Run :
```bash
export SAPCLI_WORK_DIR=/your/path
```
4. Configure your default credentials config
* ##### create `.config.yaml` file in ``SAPCLI_WORK_DIR``
 ```bash
touch $SAPCLI_WORK_DIR/.config.yaml
```
* ##### Your`.config.yaml` should look like folow :
```yaml
# Auth configurations
creds:
	api-token: "YOUR CLOUD API TOKEN HERE"
	subscription-id: "YOUR CLOUD SUBSCRIPTION ID HERE"
```
  
## Easy To Use 😀 <a name="usage"></a>

`./sapcli help`

```bash
#show general info of tool
./sapcli info

#get build by its code
./sapcli build get --code=[BUILD_CODE]

#get build progress by its code
./sapcli build progress --code=[BUILD_CODE]

#create new build
./sapcli build create --branch=[BRANCH_NAME] --name=[BUILD_NAME]

#create new deploy
./sapcli deploy create --build-code=[BUILD_CODE] --database-update-mode=[DB_UPDAT_MODE] --strategy=[STRATEGY] --env=[ENV]
```

### It's possible also to use a custom credentials config file (json or yaml)
- Examples :
```yaml
#json confige file example
{
"creds": {
	"api-token": "YOUR CLOUD API TOKEN HERE",
	"subscription-id": "YOUR CLOUD SUBSCRIPTION ID HERE"
	}
}
```

```yaml
#yaml confige file example
creds:
	api-token: "YOUR CLOUD API TOKEN HERE"
	subscription-id: "YOUR CLOUD SUBSCRIPTION ID HERE"
```
- Usage examples:
```bash
#use custom json file for get build
./sapcli build get --config=/your/custom/path/config.json --code=[BUILD_CODE] 

#use custom yaml file for create build
./sapcli build create --config=/your/custom/path/config.yaml --branch=[BRANCH_NAME] --name=[BUILD_NAME]
```

## Releases

See https://github.com/mohamedelhassak/sapcli/releases

## License

© 2022 [Mohamed El hassak]()

Licensed under the [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0) ([`LICENSE-APACHE`](LICENSE-APACHE)) or the [MIT license](https://opensource.org/licenses/MIT) ([`LICENSE-MIT`](LICENSE-MIT)), at your option. See the [`LICENSE`](LICENSE) file for more details.

`SPDX-License-Identifier: (Apache-2.0 OR MIT)`