# slack-api-lambda-example

> Go binary deployed as an AWS lambda.
> Determines who is the funniest member of a channel, by totalling
> the number of amused reactions each member's comments received
> and publicising the victor back on the channel.

## FUTURE IMPROVEMENTS

* add this to a CI system - currently all automated but run from a local laptop.
   * triggered to build then run tests on commits to an open pull request
   * triggered to build->test->deploy on merges to master.
* add sufficient tests that we can reliably deploy to the environment.
   (i.e. *continuous deployment*)
* use a centralised secrets manager to store the slack api key e.g. AWS Secrets Manager or Vault or SSM parameter store (at the moment it is manually configured as an env var to the lambda, after deployment.

## BUILD

[![Run Status](https://api.shippable.com/projects/5b19bfbff9f9060700439319/badge?branch=master)](https://app.shippable.com/github/opsgang/slack-api-lambda-example)

(shippable does build the binary but does not deploy to AWS lambda yet)

To build manually:
```bash

# ... assuming you are set up for go1.x
# - if not consider using docker golang:alpine

    export __REPO=github.com/opsgang/slack-api-lambda-example
    export GOPATH=/go GOBIN=/usr/local/go/bin;
    export LGOBIN=$GOBIN PATH=$PATH:$GOBIN;
    export PATH=$GOBIN:$PATH;
    export __WD=$GOPATH/src/$__REPO

    go get $__REPO

    cd $__WD

    fl="-w -extldflags \"-static\"";
    export CGO_ENABLED=0;
    go build --ldflags "$fl" -o $GOBIN/pupkin .

```

## DEPLOY

```bash
cd $__WD/tf

cp -a $GOBIN/pupkin .
zip pupkin.zip pupkin
rm pupkin

# export your AWS creds and then ...
# ... export terraform vars - see tf/main.tf for variable descriptions.
export TF_VAR_channel_id=my-channel-id # replace val with your channel's id
export TF_VAR_results_posted_by="Rupert Pupkin Speaks! "
export TF_VAR_icon_url="http://blog.edtechie.net/wp-content/uploads/2015/07/kingofcomedy.jpg"

terraform init
terraform plan -input=false
terraform apply -input=false -auto-approve
```

## POST DEPLOY

Find the lambda and add the environment var `API_KEY`, with the value of an api token
you've created in your slack org. As this is a secret we don't hard code it in to the terraform.
