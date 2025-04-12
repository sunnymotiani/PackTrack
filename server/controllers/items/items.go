package items

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sunnymotiani/PackTrack/server/models/items"
	"github.com/sunnymotiani/PackTrack/server/utils"
)

type ItemsController struct {
	IS *items.ItemsService
}

func (ic *ItemStatusController) AddItem(w http.ResponseWriter, r *http.Request) {
	var input items.Item
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, utils.JSONError{Msg: "err bad request"})
		return
	}
	err := ic.IS.AddItem(&input)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError,
			utils.JSONError{Msg: fmt.Sprintf("err adding item : %s", err.Error())})
		return
	}
	utils.RespondJSON(w, http.StatusOK, nil)
}
func (ic *ItemStatusController) GetItemByCategory(w http.ResponseWriter, r *http.Request) {
	catID := chi.URLParam(r, "catID")
	if catID == "" {
		utils.ResponseError(w, http.StatusBadRequest, utils.JSONError{Msg: "missing categoryID in URL"})
		return
	}
	items, err := ic.IS.GetItemByCategory(r.Context(), catID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError,
			utils.JSONError{Msg: fmt.Sprintf("err fetching items %s", err.Error())})
		return
	}
	utils.RespondJSON(w, http.StatusOK, items)
}

func (ic *ItemStatusController) UpdateItemStatus(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ItemID    string `json:"item_id"`
		UserID    string `json:"user_id"`
		NewStatus string `json:"new_status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, utils.JSONError{Msg: "err bad request"})
		return
	}
	err := ic.IS.UpdateItemStatus(input.ItemID, input.NewStatus, input.UserID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError,
			utils.JSONError{Msg: fmt.Sprintf("err updating item status %s", err.Error())})
		return
	}
	utils.RespondJSON(w, http.StatusOK, nil)
}
func (ic *ItemStatusController) AssignItem(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ItemID string `json:"item_id"`
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.ResponseError(w, http.StatusBadGateway, utils.JSONError{Msg: "err bad request"})
		return
	}
	err := ic.IS.AssignItem(input.ItemID, input.UserID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError,
			utils.JSONError{Msg: fmt.Sprintf("err assigning item %s", err.Error())})
		return
	}
	utils.RespondJSON(w, http.StatusOK, nil)
}
func (ic *ItemStatusController) UnAssignItem(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ItemID string `json:"item_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.ResponseError(w, http.StatusBadGateway, utils.JSONError{Msg: "err bad request"})
		return
	}
	err := ic.IS.UnassignItem(input.ItemID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError,
			utils.JSONError{Msg: fmt.Sprintf("err un-assigning item %s", err.Error())})
		return
	}
	utils.RespondJSON(w, http.StatusOK, nil)
}
func (ic *ItemsController) EditItem(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "itemID")
	if itemID == "" {
		utils.ResponseError(w, http.StatusBadRequest, utils.JSONError{Msg: "missing itemID in URL"})
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, utils.JSONError{Msg: "invalid JSON body"})
		return
	}

	err := ic.IS.EditItem(r.Context(), itemID, updates)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError,
			utils.JSONError{Msg: fmt.Sprintf("failed to update item: %s", err.Error())})
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"msg": "item updated successfully"})
}

func (ic *ItemsController) DeleteItem(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "itemID")
	if itemID == "" {
		utils.ResponseError(w, http.StatusBadRequest, utils.JSONError{Msg: "missing itemID in URL"})
		return
	}

	err := ic.IS.DeleteItem(itemID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError,
			utils.JSONError{Msg: fmt.Sprintf("failed to delete item: %s", err.Error())})
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"msg": "item deleted successfully"})
}
