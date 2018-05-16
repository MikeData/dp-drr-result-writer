# Dimension Recombination Reporter - Uploader

WORK IN PROGRESS - 10% project.

The DRR is a microservice based project to analyse all permutations from structural alterations (combinations of dimension items added/removed) that could be applied to a given source dataset. Presenting the findings as a human readable report.

The result writer checks each result, discarding those that don't meet the filter (for now hardcoded as those that don't reduce sparsity) and writing positive results to CSV before sending back to the controller for further permutations.

## Uploader

The uploader is a command line tool for starting the process.


## Configuration

| Environment variable        |  Description
| --------------------------- |  -----------
| AWS_REGION                  | an AWS credential
| AWS_SECRET_ACCESS_KEY       | an AWS credential
| AWS_ACCESS_KEY_ID           | an AWS credential
| SQS_SOURCE_QUEUE_URL        | the full url of the drr source queue
| SQS_TASK_QUEUE_URL          | the full url of the drr task queue
| DRR_IMPORT_BUCKET           | the name of the import bucket



## Usage

`dp-drr-uploader -filepath <path-to-csv>`
