package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"release_youtracker/types"

	"github.com/rs/zerolog/log"
)

type Server struct {
	Config types.Config
}

func (as *Server) Init(configPath string) error {
	if err := as.Config.LoadConfig(configPath); err != nil {
		return err
	}
	return nil
}

func (as Server) AlertHandler(w http.ResponseWriter, r *http.Request) {
	var alert types.Alert
	var release types.Release
	if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Info().Msgf("Received alert: %+v\n", alert)

	for _, alerts := range alert.Alerts {
		for label, key := range alerts.Labels {
			switch label {
			case "name":
				release.Name = key
			case "repository":
				release.Repository = key
			case "tag":
				release.Tag = key
			}
		}
	}
	task := as.createTask(release)
	as.pushTask(task)
	w.WriteHeader(http.StatusOK)
}

func (as Server) createTask(release types.Release) types.Project {
	var task types.Project
	task.Project.Id = as.Config.Youtrack.ProjectID
	problem := fmt.Sprintf("There is a new version of %s", release.Name)
	description := fmt.Sprintf("Check the release notes, and if there is nothing critical, proceed with the update. Otherwise, put it on hold with a comment, notify the team, and schedule the update. %s", release.Name)
	additional := fmt.Sprintf("Check it at https://github.com/%s/releases/tag/%s", release.Repository, release.Tag)
	definition := fmt.Sprintf("Update %s or leave a comment why we don't have to update", release.Name)
	task.Description = fmt.Sprintf(types.Description, problem, description, additional, definition)
	task.Summary = fmt.Sprintf("%s update", release.Name)
	task.CustomField = append(task.CustomField, map[string]interface{}{ // just let it be hardcoded
		"value": map[string]interface{}{
			"name":  "Major",
			"id":    "99-2",
			"$type": "EnumBundleElement",
		},
		"name":  "Priority",
		"id":    "120-37",
		"$type": "SingleEnumIssueCustomField",
	})

	return task
}

func (as Server) pushTask(task types.Project) {
	var req *http.Request
	var resp *http.Response
	var payloadBytes []byte
	var err error

	if payloadBytes, err = json.Marshal(task); err != nil {
		log.Fatal().Err(err).Msg("Error marshaling payload")
	}

	url := fmt.Sprintf("%s/api/issues", as.Config.Youtrack.Host)
	if req, err = http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes)); err != nil {
		log.Fatal().Err(err).Msg("Error creating request")
	}

	req.Header.Set("Authorization", "Bearer "+as.Config.Youtrack.Key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	if resp, err = client.Do(req); err != nil {
		log.Fatal().Err(err).Msg("Error making request")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Info().Msgf("Task \"%s\" is created", task.Summary)
	} else {
		log.Info().Msgf("Failed to create \"%s\" task. Status: %s\n", task.Summary, resp.Status)
	}
}
