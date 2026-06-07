// WishlistController handles GET /wishlist - protected SSR page.
package controllers

type WishlistController struct {
	BaseController
}

// Get renders the wishlist page. Auth is guaranteed by the filter - no need to check here.
func (c *WishlistController) Get() {
	c.Data["ActiveNav"] = "wishlist"
	c.Data["WishlistItems"] = nil
	c.TplName = "wishlist.tpl"
	c.Layout = "layout.tpl"
}
