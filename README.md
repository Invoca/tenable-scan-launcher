# Tenable-Scan-Launcher [![Build Status](https://travis-ci.org/Invoca/tenable-scan-launcher.svg?branch=master)](https://travis-ci.org/Invoca/tenable-scan-launcher) [![Coverage Status](https://coveralls.io/repos/github/Invoca/tenable-scan-launcher/badge.svg?branch=master)](https://coveralls.io/github/Invoca/tenable-scan-launcher?branch=master)
This scan launcher collects the private ip addresses of Google Cloud and AWS instances and then launches a tenable scan
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
The scanner will list private ips from all regions of each cloud provider given. To enable AWS, include the 
`--include_aws` flag. It will use the shared aws configuration settings, so it will use the standard order of precedence
for AWS service accounts. To include Google Cloud, use the `--include_gcloud` flag and be sure to specify the service 
account file location with `--gcloud_json` and the desired project with `--gcloud_project`.

The following Tenable flags are needed to preform a scan:
`--tenable-access-key` which is the access key generated from https://cloud.tenable.com/#/ . 
`--tenable-secret-key` is the secret key generated from https://cloud.tenable.com/#/
`--scanner_id` is the scanner to id of the scanner to use.

To generate a report, include the `--generate_report` flag. To specify a format of the report include the `--format` 
flag with the desired format. `--report-file-location` spcifies the file location to save the file. To specify which 
chapters of the report to use, use the `--chapters` flag with the desired chapters. `--summary-report` only includes 
the `vuln_hosts_summary` chapter while `--full-report` includes all of them. Note that `--summary-report` will override 
`--chapters` and `--full-report` overrides `--summary-report`.

In order to filter on the severity within the report, include the `--[low,medium,high,critical]_severity` flags. One can
also change the search type with `--search_type`. It is not recommended to change it to the `and` type. 


## Contributions

Contributions to this project are always welcome!  Please read our [Contribution Guidelines](https://github.com/Invoca/tenable-scan-launcher/blob/master/CONTRIBUTING.md) before starting any work.