package graphql

import "fmt"

type IssueInfo struct {
	Id   string
	Iid  string
	Name string
}

func (info IssueInfo) Str() string {
	return fmt.Sprintf("%s %s", info.Iid, info.Name)
}

type IssueInfoResponse struct {
	Data struct {
		Project struct {
			WorkItems struct {
				Nodes []IssueInfo
			}
		}
	}
}

type TimelogInput struct {
	SpentAt   string `json:"spentAt"`
	Summary   string `json:"summary"`
	TimeSpent string `json:"timeSpent"`
}

type TimeTrackingWidgetInput struct {
	Timelog TimelogInput `json:"timelog"`
}

type WorkItemInput struct {
	ID                 string                  `json:"id"`
	TimeTrackingWidget TimeTrackingWidgetInput `json:"timeTrackingWidget"`
}

type ErrorsResponse struct {
	errors []error
}