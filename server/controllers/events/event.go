package events

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sunnymotiani/PackTrack/server/models/events"
	"github.com/sunnymotiani/PackTrack/server/utils"
)

type EventController struct {
	ES *events.EventService
}

func (ec *EventController) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		OwnerID     string `json:"owner_id" binding:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, utils.JSONError{Msg: "invalid request"})
		return
	}
	event, err := ec.ES.CreateEvent(r.Context(), input.Name, input.Description, input.OwnerID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError,
			utils.JSONError{Msg: fmt.Sprintf("err creating event %s", err.Error())})
	}
	utils.RespondJSON(w, http.StatusOK, event)
}

func (ec *EventController) AddMember(w http.ResponseWriter, r *http.Request) {
	var input struct {
		EventID string `json:"event_id" binding:"required"`
		UserID  string `json:"user_id" binding:"required"`
		Role    string `json:"role" binding:"required"` // owner/admin/member/viewer
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, utils.JSONError{Msg: "invalid request"})
		return
	}

	err := ec.ES.AddMember(r.Context(), input.EventID, input.UserID, input.Role)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError,
			utils.JSONError{Msg: fmt.Sprintf("err adding member %s", err.Error())})
		return
	}
	utils.RespondJSON(w, http.StatusOK, nil)
}
func (ec *EventController) GetEventMembers(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "eventID")
	if eventID == "" {
		utils.ResponseError(w, http.StatusBadRequest, utils.JSONError{Msg: "missing eventID in URL"})
		return
	}

	members, err := ec.ES.GetEventMembers(r.Context(), eventID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, utils.JSONError{Msg: fmt.Sprintf("err fetching members %s", err.Error())})
		return
	}

	utils.RespondJSON(w, http.StatusOK, members)
}
func (ec *EventController) GetEventsForUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		utils.ResponseError(w, http.StatusBadRequest, utils.JSONError{Msg: "missing userID in URL"})
		return
	}

	eventsList, err := ec.ES.GetEventsForUser(r.Context(), userID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, utils.JSONError{Msg: fmt.Sprintf("err fetching events %s", err.Error())})
		return
	}

	utils.RespondJSON(w, http.StatusOK, eventsList)
}

func (ec *EventController) RemoveMember(w http.ResponseWriter, r *http.Request) {
	var input struct {
		EventID string `json:"event_id"`
		UserID  string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, utils.JSONError{Msg: "invalid request"})
		return
	}

	if err := ec.ES.RemoveMember(r.Context(), input.EventID, input.UserID); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, utils.JSONError{Msg: fmt.Sprintf("err removing member %s", err.Error())})
		return
	}

	utils.RespondJSON(w, http.StatusOK, nil)
}
func (ec *EventController) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	var input struct {
		EventID string `json:"event_id"`
		UserID  string `json:"user_id"`
		NewRole string `json:"new_role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, utils.JSONError{Msg: "invalid request"})
		return
	}

	if err := ec.ES.UpdateMemberRole(r.Context(), input.EventID, input.UserID, input.NewRole); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, utils.JSONError{Msg: fmt.Sprintf("err updating role %s", err.Error())})
		return
	}

	utils.RespondJSON(w, http.StatusOK, nil)
}
