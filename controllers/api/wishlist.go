// WishlistAPIController handles all /api/wishlist JSON endpoints.

package apicontrollers

import (
	"TravelSphere/models"
	"TravelSphere/services"
	"TravelSphere/utils"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
)

type WishlistAPIController struct {
	beego.Controller
}

// sessionUsername reads the authenticated username from the Beego session.
func (c *WishlistAPIController) sessionUsername() string {
	val := c.GetSession("username")
	if val == nil {
		return ""
	}
	return val.(string)
}

// Get returns all wishlist entries for the authenticated user.
// GET /api/wishlist
func (c *WishlistAPIController) Get() {
	items := services.GetWishlistService().GetAll(c.sessionUsername())
	c.Data["json"] = utils.OKResponse(items)
	c.ServeJSON()
}

// Post creates a new wishlist entry.
func (c *WishlistAPIController) Post() {
	var body struct {
		CountryName string `json:"country_name"`
		Note        string `json:"note"`
		Status      string `json:"status"`
	}

	if err := json.NewDecoder(c.Ctx.Request.Body).Decode(&body); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = utils.ErrorResponse("invalid JSON body", 400)
		c.ServeJSON()
		return
	}

	if msg := utils.ValidateWishlistCreate(body.CountryName); msg != "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = utils.ErrorResponse(msg, 400)
		c.ServeJSON()
		return
	}

	// Reject duplicate entries — a country can only appear once per user.
	svc := services.GetWishlistService()
	if svc.Exists(c.sessionUsername(), body.CountryName) {
		c.Ctx.Output.SetStatus(409)
		c.Data["json"] = utils.ErrorResponse(body.CountryName+" is already in your wishlist", 409)
		c.ServeJSON()
		return
	}

	if body.Status == "" {
		body.Status = string(models.StatusPlanned)
	}

	if msg := utils.ValidateWishlistStatus(body.Status); msg != "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = utils.ErrorResponse(msg, 400)
		c.ServeJSON()
		return
	}

	item := svc.Create(c.sessionUsername(), body.CountryName, body.Note, body.Status)
	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = utils.CreatedResponse(item)
	c.ServeJSON()
}

// Update modifies note and status of a single wishlist entry.
// PUT /api/wishlist/:id — returns only the updated item, not the full list.
func (c *WishlistAPIController) Update() {
	id := c.Ctx.Input.Param(":id")

	var body struct {
		Note   string `json:"note"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(c.Ctx.Request.Body).Decode(&body); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = utils.ErrorResponse("invalid JSON body", 400)
		c.ServeJSON()
		return
	}

	if msg := utils.ValidateWishlistStatus(body.Status); msg != "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = utils.ErrorResponse(msg, 400)
		c.ServeJSON()
		return
	}

	updated, ok := services.GetWishlistService().Update(
		c.sessionUsername(), id, body.Note, body.Status,
	)
	if !ok {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = utils.ErrorResponse("wishlist entry not found", 404)
		c.ServeJSON()
		return
	}

	// Return only the updated item — REST convention for PUT on a single resource.
	c.Data["json"] = utils.OKResponse(updated)
	c.ServeJSON()
}

// Delete removes a wishlist entry by ID.
// DELETE /api/wishlist/:id — returns 204 No Content on success.
func (c *WishlistAPIController) Delete() {
	id := c.Ctx.Input.Param(":id")

	if !services.GetWishlistService().Delete(c.sessionUsername(), id) {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = utils.ErrorResponse("wishlist entry not found", 404)
		c.ServeJSON()
		return
	}

	// 204 No Content — deletion succeeded, nothing to return.
	c.Ctx.Output.SetStatus(204)
}
