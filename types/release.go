package types

const Description = "## Problem \n%s \n\n ### Description \n%s \n\n ### Additional info \n%s \n\n ## Definition of done \n%s \n\n"

type Release struct {
	Name       string
	Repository string
	Tag        string
}

type Project struct {
	Project     ProjectId     `json:"project"`
	Summary     string        `json:"summary"`
	Description string        `json:"description"`
	CustomField []interface{} `json:"customFields"`
}

type ProjectId struct {
	Id string `json:"id"`
}
