# Zendesk coding challenge

## Assumptions
* no Joining queries for example find this user AND this user
* That we can just present the full found data and don't need to filter data.
* That search should care about case.

## Design notes
* I'm not happy with my foreign key implementation. Where the object knows who has refs to it. However I think it's easier to understand and use then making foreign key objects which could be queried.
* I Personally prefer large test functions with the message giving more info rather then lot's of test functions all testing one thing. However I Always do the single test functions in my code at work.
* For resolving what matches what. I started to go down the reflection road (using the type data at runtime) but I thought it was becoming fairly hard to read due to it becoming a lot of dense rarely used reflection functions stuff. So I opted to go with just writing match methods for the data objects. However If I needed to add 3 more data objects I would switch to the reflection route. 
* Support for running multiple queries exists in a limited capacity because I wanted to make sure my design would allow it. You can look at `db/query_test.go` for examples.
* Yes searching happens linearly when not searching by ID. I don't know enough about making DB's from scratch to create a indexing system.

## Arguments

`-h` Output
```
  -orgs_file="": path to organizations json file
  -query="": the query to be ran. should go "RESOURCE FIELD TARGET VALUE" Example "user name Cross Barlow" will return the user along with any tickets and organization associated with said user. valid resoruce are organization user and ticket. Check the given json files for the field names
  -tickets_file="": path to users json file
  -users_file="": path to users json file
```

full example 
```
	-orgs_file "db/db_testdata/organizations.json" \
	-users_file "db/db_testdata/users.json" \
	-tickets_file "db/db_testdata/tickets.json" \
	-query "user id 74"
```

Query Examples
* `user name 

## Using Docker

### Building
Run `docker build . -t zendesk-coding:paul-sarda`

### Running
Run `docker run -rm zendesk-coding:paul-sarda ARGUMENTS HERE`

## Local

### Building
run `go build .`

### Running
Windows `./zendesk.exe`

Linux / Mac `./zendesk ARGUMENTS HERE`
