// WishlistController serves GET /wishlist — protected by AuthFilter.
package controllers

import "TravelSphere/services"

type WishlistController struct {
	BaseController
}

// Get renders the wishlist page with the current user's saved destinations.
func (c *WishlistController) Get() {
	username := c.GetSession("username").(string)
	items := services.GetWishlistService().GetAll(username)
	c.Data["ActiveNav"]     = "wishlist"
	c.Data["WishlistItems"] = items
	c.TplName = "wishlist.tpl"
	c.Layout  = "layout.tpl"
}