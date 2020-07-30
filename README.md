# cf-ondemand-asg
This app helps to bind ASGs to all spaces in a cloudfoundry org. It expects an ASG with the same name as the org already present.
This allows for on-demand binding of ASGs to new spaces.

## deployment
git clone this repo, cd into the directory, then run ```cf push <appname> --no-start```. Now set the CFUSER and CFPASSWORD environment variables using ```cf set-env``` and run ```cf start <app name>```

## usage
do a post request to /api/v1/synch with a json in this format:
```json
{
  "orgName": <org name>
}
```
