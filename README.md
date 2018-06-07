# slack-api-lambda-example

> Go binary deployed as a lambda.
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

IN THIS REPO:
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

    local fl="-w -extldflags \"-static\"";
    export CGO_ENABLED=0;
    su-exec root go build --ldflags "$fl" -o $GOBIN/pupkin .

```

## DEPLOY

```bash
cd $__WD

cd tf

cp -a $GOBIN/pupkin .
zip pupkin.zip pupkin
rm pupkin

# export your AWS creds and then ...
terraform init
terraform plan -input=false
terraform apply -input=false -auto-approve
```
