# Tenable-Scan-Launcher [![Build Status](https://travis-ci.org/Invoca/tenable-scan-launcher.svg?branch=master)](https://travis-ci.org/Invoca/tenable-scan-launcher) [![Coverage Status](https://coveralls.io/repos/github/Invoca/tenable-scan-launcher/badge.svg?branch=master)](https://coveralls.io/github/Invoca/tenable-scan-launcher?branch=master)
This scan launcher collects the private IP addresses of Google Cloud and AWS instances and then launches a tenable scan
with the option to downlaod the scan as a pdf. 

## Installation
Installing this repo is as simple as cloning the repo into your $Go/src/github.com directory. 
```bash
git clone git@github.com:Invoca/tenable-scan-launcher.git
```

## Usage
### Running
There are two methods to run the scan launcher. With Docker or by running the executable file. 
Docker:
```shell script
    docker run invoca:SOMETHING   
```
Shell:
```shell script
  go build -mod=readonly -o $PWD/tenable-scan-launcher $PWD/cmd/tenable-scan-launcher
  ./tenable-scan-launcher $FLAGS
```

### Flags
The scanner will list private IPs from all regions of each cloud provider given. To enable AWS, include the 
`--include_aws` flag. It will use the shared aws configuration settings, so it will use the standard order of precedence
for AWS service accounts. To include Google Cloud, use the `--include_gcloud` flag and be sure to specify the service 
account file location with `--gcloud_json` and the desired project with `--gcloud_project`.

The following Tenable flags are needed to preform a scan:

* `--tenable-access-key` which is the access key generated from https://cloud.tenable.com/#/ . 
* `--tenable-secret-key` is the secret key generated from https://cloud.tenable.com/#/
* `--scanner_id` is the scanner to id of the scanner to use.

#### Reports

|Flag|Description|
|---|---|
|`--generate_report`|Generates a report|
|`--format`|Specifies the format of the report. Formats are Nessus, HTML, PDF, CSV, or DB. Defaults to empty string|
|`--report-file-location`|The file location to save the file|
|`--chapters`|Specify which chapters of the report to use. Supported chapters are vuln_hosts_summary, vuln_by_host, compliance_exec, remediations, vuln_by_plugin, compliance. Defaults to empty string.|
|`--summary-report`|Only includes the `vuln_hosts_summary` chapter|
|`--full-report`|Includes all chapters|

Note that `--summary-report` will override `--chapters` and `--full-report` overrides `--summary-report`

#### Filtering
In order to filter on the severity within the report, include the `--[low,medium,high,critical]_severity` flags. The
search filter can be modified with `--search_type`. The supported values are `and` and `or`. It is not recommended 
changing it to the `and` type since each vulnerability can only have a single severity level. 

#### Logging
Log level can be specified with `--log-level`. The levels are trace, info, fatal, panic, warn, and debug. Log format can
be specified with `--log-type`. The supported types are `json`, and `text`. 


## Contributions

Contributions to this project are always welcome!  Please read our [Contribution Guidelines](https://github.com/Invoca/tenable-scan-launcher/blob/master/CONTRIBUTING.md) before starting any work.
