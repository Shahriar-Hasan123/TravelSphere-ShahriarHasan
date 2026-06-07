// WishlistController handles GET /wishlist - protected SSR page

package controllers

type WishlistController struct {
	BaseController
}

// Get renders the wishlist page for authenticated users
func (c *WishlistController) Get() {
	if !c.RequireLogin() {
		return
	}
	c.Data["ActiveNav"] = "wishlist"
	c.TplName = "wishlist.tpl"
	c.Layout = "layout.tpl"
}
