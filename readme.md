# SAPCLI
A sample SAP CLI implementation in Go

[![GitHub release](https://img.shields.io/github/release/moul/banner.svg)](https://github.com/mohamedelhassak/sapcli/releases)
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

- For linux 64 bits arch :
```go
cd sapcli
env GOOS=linux GOARCH=amd64 go build -o build/sapcli
```
- For windows 64 bits arch :
```go
cd sapcli
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
4. Configure your default credentials config within the work directory ``SAPCLI_WORK_DIR``, you can use either `.json` or `.yaml` file
* ##### Create `.config.yaml`
 ```bash
touch $SAPCLI_WORK_DIR/.config.yaml
```
* ##### Or
 ```bash
touch $SAPCLI_WORK_DIR/.config.json
```

* ##### Your`.config.yaml` should look like folow :
```yaml
creds:
	api-token: "YOUR CLOUD API TOKEN HERE"
	subscription-id: "YOUR CLOUD SUBSCRIPTION ID HERE"
```
* ##### Your`.config.json` should look like folow :
```json
{
"creds": {
	"api-token": "YOUR CLOUD API TOKEN HERE",
	"subscription-id": "YOUR CLOUD SUBSCRIPTION ID HERE"
	}
}
```
  
## Easy To Use 😀 <a name="usage"></a>



```bash
#show help
./sapcli help (or -h)

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

### It's possible also to use a custom credentials config file (json or yaml) passed with `--config` flag


- Examples :
```bash
#use custom json file for get build
./sapcli build get --config=/your/custom/path/config.json --code=[BUILD_CODE] 

#use custom yaml file for create build
./sapcli build create --config=/your/custom/path/config.yaml --branch=[BRANCH_NAME] --name=[BUILD_NAME]
```

## Releases

See https://github.com/mohamedelhassak/sapcli/releases
