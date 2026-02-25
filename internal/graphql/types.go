package graphql

type IssueInfo struct {
	Id   string
	Iid  string
	Name string
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
    SpentAt     string `json:"spentAt"`
    Summary     string `json:"summary"`
    TimeSpent   string `json:"timeSpent"`
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