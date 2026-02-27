package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jesusthecreator017/fswithgo/cmd/api/helpers"
	"github.com/jesusthecreator017/fswithgo/cmd/api/middleware"
	"github.com/jesusthecreator017/fswithgo/internal/store"
)

func (app *application) createIssueHandler(w http.ResponseWriter, req *http.Request) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := helpers.ReadJson(req, &input); err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	errs := make(map[string]string)

	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" {
		errs["title"] = "must not be blank"
	} else if len(input.Title) > 255 {
		errs["title"] = "must not be more than 255 characters"
	}

	if len(errs) > 0 {
		helpers.ValidationErrorJson(w, errs)
		return
	}

	userID := middleware.GetUserID(req)

	issue := &store.Issue{
		Title:       input.Title,
		UserID:      userID,
		Description: input.Description,
		CreatedAt:   time.Now(),
		Status:      store.Incomplete,
	}

	if err := app.store.Issues.Create(req.Context(), issue); err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to create issue")
		return
	}

	helpers.WriteJson(w, http.StatusCreated, helpers.Envelope{"issue": issue})
}

func (app *application) listIssueHandler(w http.ResponseWriter, req *http.Request) {
	userID := middleware.GetUserID(req)

	issueList, err := app.store.Issues.ListByUserID(req.Context(), userID)
	if err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to get issues")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"issues": issueList})
}

func (app *application) getIssueHandler(w http.ResponseWriter, req *http.Request) {
	// Get the id from path and validate it
	id := req.PathValue("id")

	// Check if the id is even passed in
	if id == "" {
		helpers.ErrorJson(w, http.StatusBadRequest, "id is required")
		return
	}

	// Check if the id is a number
	intId, err := strconv.Atoi(id)
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "id must be a number")
		return
	}

	// Get the issue of specific id
	issue, err := app.store.Issues.GetByID(req.Context(), int64(intId))
	if err != nil {
		helpers.ErrorJson(w, http.StatusNotFound, "issue not found")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"issue": issue})
}

func (app *application) updateIssueStatusHandler(w http.ResponseWriter, req *http.Request) {
	// Process the id path of the request
	id := req.PathValue("id")

	// Check if the id is even passed
	if id == "" {
		helpers.ErrorJson(w, http.StatusBadRequest, "id is required")
		return
	}

	// Check the id is a number
	intID, err := strconv.Atoi(id)
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "id has to be a number")
		return
	}

	// Look for the id in the database
	issue, err := app.store.Issues.GetByID(req.Context(), int64(intID))
	if err != nil {
		helpers.ErrorJson(w, http.StatusNotFound, "issue not found")
		return
	}

	// Now that we have the issue we need to update it
	var input struct {
		Status store.StatusType `json:"status"`
	}

	// Decode the input and handle any errors
	if err := helpers.ReadJson(req, &input); err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	errs := make(map[string]string)

	switch input.Status {
	case store.Incomplete, store.InProgress, store.Complete:
		// Valid Do nothing: pass
	default:
		errs["status"] = "must be one of: Incomplete, In-Progress, Complete"
	}

	if len(errs) > 0 {
		helpers.ValidationErrorJson(w, errs)
		return
	}

	// Update the status
	updatedIssue, err := app.store.Issues.UpdateStatus(req.Context(), issue.ID, input.Status)
	if err != nil {
		helpers.ErrorJson(w, http.StatusInternalServerError, "failed to update issue status")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"issue": updatedIssue})
}

func (app *application) deleteIssueHandler(w http.ResponseWriter, req *http.Request) {
	// Handle the id path
	id := req.PathValue("id")

	if id == "" {
		helpers.ErrorJson(w, http.StatusBadRequest, "id is required")
		return
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		helpers.ErrorJson(w, http.StatusBadRequest, "id must be a number")
		return
	}

	// Delete the issue
	err = app.store.Issues.Delete(req.Context(), int64(intID))
	if err != nil {
		helpers.ErrorJson(w, http.StatusNotFound, "Issue not found")
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"message": "issue deleted"})
}
