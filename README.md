# mvm-sint-predict

## Background
Every year Sinterklaas comes from Spain to Moeders Voor Moeders to give toys and candy to the nice children.
Since this is not an easy task to predict how many children will come we use the past data of their family to try to predict this.
The data is exported as CSV out of our Zoho setup then is processed to filter out "frequent clients" and get their family composition.
All data is processed using internal IDs, no names or other identifiable info is ever passed through this system.

You might find the setup weird to be using gRPC, but everything has its reason. Some are just not the right reasons ;)

## Usage
This app consists of two parts, a server and a client. The server for testing purposes can be run locally.
```
$ sint-server serve
```
This will by default listens on port 8080

To use the CLI client you need to use the CSV data exported from ZOHO, sample data is included in `fixtures`.
To get the frequency a famili vistis from the visit data you can use the `frequency` command which will output JSON data this can optionally be written to disk using `-f`.
```
$ sint-client frequency -f frequency.json -c ./fixtures/visits-large.csv 
```
Once extracted the frequency data you can use this to run the prediction system.
```
$ sint-client children-count -f frequency.json -c ./fixtures/family-composition.csv
```

### Fixtures
This repo comes with a few fixtures to present how we export our data from ZOHO, this is based on scrambeled data from the actual data, for this reason you might see weird things here.
