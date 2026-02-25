package graphql

var issueInfoQuery = `
	{ 
		project(fullPath: "greendata/greendata-core") { 
			workItems(iid: "%s") { 
				nodes { 
					id iid name
				} 
			} 
		} 
	}
 `

 var addSpentTimeQuery = `
 	mutation workItemUpdate($input: WorkItemUpdateInput!) {
		workItemUpdate(input: $input) {
			errors
		}
	}
 `
