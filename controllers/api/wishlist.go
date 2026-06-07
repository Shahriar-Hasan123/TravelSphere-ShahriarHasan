// WishlistAPIController handles JSON CRUD API for the travel wishlist.

package apicontrollers

import beego "github.com/beego/beego/v2/server/web"

type WishlistAPIController struct {
	beego.Controller
}

// Get returns all wishlist entries
func (c *WishlistAPIController) Get() {
	c.Data["json"] = map[string]string{"status": "ok", "message": "wishlist stub"}
	c.ServeJSON()
}

// Post creates a new wishlist entry (stub).
func (c *WishlistAPIController) Post() {
	c.Data["json"] = map[string]string{"status": "ok", "message": "create stub"}
	c.ServeJSON()
}

// Update modifies an existing wishlist entry (stub).
func (c *WishlistAPIController) Update() {
	c.Data["json"] = map[string]string{"status": "ok", "message": "update stub"}
	c.ServeJSON()
}

// Delete removes a wishlist entry (stub).
func (c *WishlistAPIController) Delete() {
	c.Data["json"] = map[string]string{"status": "ok", "message": "delete stub"}
	c.ServeJSON()
}
