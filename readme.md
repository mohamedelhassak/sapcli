# SAPCLI
A sample SAP CLI implementation in Go

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