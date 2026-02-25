package graphql

import (
	"fmt"
	"log"
	"tracker/internal/service"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

var git *gitlab.Client

func init() {
	var err error
	git, err = gitlab.NewClient(service.Settings.GitlabToken, gitlab.WithBaseURL(service.Settings.GitlabUrl))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	log.Println("Init client")
}

func GetIssueInfo(issue string) IssueInfo {
	query := &gitlab.GraphQLQuery{Query: fmt.Sprintf(issueInfoQuery, issue)}
	var resp IssueInfoResponse
	_, err := git.GraphQL.Do(*query, &resp)
	if err != nil {
		log.Fatal(err)
	}
	if len(resp.Data.Project.WorkItems.Nodes) > 0 {
		return resp.Data.Project.WorkItems.Nodes[0]
	}
	return IssueInfo{}
}

func AddSpendTime(issue string, minutes int, date string) {
	info := GetIssueInfo(issue)
	if (IssueInfo{}) == info  {
		info = GetIssueInfo(service.Settings.OtherIssue)
		if (IssueInfo{}) == info  {
			return
		}
	}
	input := WorkItemInput{
		ID:                 info.Id,
		TimeTrackingWidget: TimeTrackingWidgetInput{
			Timelog: TimelogInput{
				SpentAt: date,
				TimeSpent: fmt.Sprintf("%dm", minutes),
			},
		},
	}
	query := &gitlab.GraphQLQuery{
		Query: addSpentTimeQuery,
		Variables: map[string]any{"input": input},
	}
	var resp ErrorsResponse
	_, err := git.GraphQL.Do(*query, &resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}

var s =
`
{
  "operationName": "workItemUpdate",
  "variables": {
    "input": {
      "id": "gid://gitlab/WorkItem/159385",
      "timeTrackingWidget": {
        "timelog": {
          "spentAt": "2026-02-16",
          "summary": "",
          "timeSpent": "3m"
        }
      }
    }
  },
  "query": "mutation workItemUpdate($input: WorkItemUpdateInput!) {\n  workItemUpdate(input: $input) {\n    workItem {\n      ...WorkItem\n      __typename\n    }\n    errors\n    __typename\n  }\n}\n\nfragment WorkItem on WorkItem {\n  id\n  iid\n  __typename\n}"
}

`