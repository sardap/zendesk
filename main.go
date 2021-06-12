package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/namsral/flag"
	"github.com/sardap/zendesk/db"
)

type Args struct {
	// Files
	OrganizationsFile string
	UsersFile         string
	TicketsFile       string
	// Query
	Query db.Query
}

func ParseFlags() (Args, error) {
	var result Args

	flag.StringVar(&result.OrganizationsFile, "orgs_file", "", "path to organizations json file")
	flag.StringVar(&result.UsersFile, "users_file", "", "path to users json file")
	flag.StringVar(&result.TicketsFile, "tickets_file", "", "path to users json file")
	var queryStr string
	flag.StringVar(
		&queryStr, "query", "",
		fmt.Sprintf(
			"the query to be ran. should go \"RESOURCE FIELD TARGET VALUE\"\n"+
				"Example \"user name Cross Barlow\" will return the user along with any"+
				"tickets and organization associated with said user.\n"+
				"valid resoruce are %s %s and %s. Check the given json files for the field names\n",
			db.ResourceOrganization, db.ResourceUser, db.ResourceTicket,
		),
	)
	flag.Parse()

	if _, err := os.Stat(result.OrganizationsFile); err != nil {
		return result, fmt.Errorf("invalid or no organizations file given")
	}
	if _, err := os.Stat(result.UsersFile); err != nil {
		return result, fmt.Errorf("invalid or no users file given")
	}
	if _, err := os.Stat(result.TicketsFile); err != nil {
		return result, fmt.Errorf("invalid or no tickets file given")
	}

	// Parse query
	splits := strings.SplitN(queryStr, " ", 3)
	if len(splits) != 3 {
		return result, fmt.Errorf("invalid query string please check -h")
	}

	var resource db.ResourceType
	switch db.ResourceType(splits[0]) {
	case db.ResourceOrganization, db.ResourceUser, db.ResourceTicket:
		resource = db.ResourceType(splits[0])
	default:
		return result, fmt.Errorf("invalid resource given in query please check -h")
	}

	fullMatch := &db.FulLMatchCondition{
		Resource:  resource,
		Connector: db.ConnectorTypeUnion,
		Field:     splits[1],
		Match:     splits[2],
	}

	result.Query = db.Query{
		Conditions: []db.Condition{fullMatch},
	}

	return result, nil
}

func createDB(args Args) *db.DB {
	orgsF, err := os.Open(args.OrganizationsFile)
	if err != nil {
		panic(err)
	}
	defer orgsF.Close()

	usersF, err := os.Open(args.UsersFile)
	if err != nil {
		panic(err)
	}
	defer usersF.Close()

	ticketsF, err := os.Open(args.TicketsFile)
	if err != nil {
		panic(err)
	}
	defer ticketsF.Close()

	result, err := db.Create(orgsF, usersF, ticketsF)
	if err != nil {
		panic(err)
	}

	return result
}

func main() {
	args, err := ParseFlags()
	if err != nil {
		panic(err)
	}

	database := createDB(args)

	result, err := args.Query.Resolve(database)
	if err != nil {
		panic(err)
	}
	jsonBytes, _ := json.Marshal(result)
	fmt.Printf("%s\n", jsonBytes)
}
