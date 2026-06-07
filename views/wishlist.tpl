<div class="page-body">
    <div class="page-header">
        <h1 class="page-title">Travel Wishlist</h1>
        <p class="page-subtitle">
            Edit notes, update trip status, or remove destinations.
            Changes save without reloading the page.
        </p>
    </div>

    <!-- AJAX replaces only this container -->
    <div id="wishlist-rows" class="wishlist-table-wrap">
        <table class="wishlist-table">
            <thead>
                <tr>
                    <th>COUNTRY</th>
                    <th>NOTE</th>
                    <th>STATUS</th>
                    <th>ACTIONS</th>
                </tr>
            </thead>
            <tbody>
                {{range .WishlistItems}}
                <tr data-id="{{.ID}}">
                    <td class="wl-country">{{.CountryName}}</td>
                    <td class="wl-note">
                        <input type="text" class="wl-note-input"
                               value="{{.Note}}" placeholder="Add a note...">
                    </td>
                    <td class="wl-status">
                        <select class="wl-status-select {{if eq .Status "Visited"}}status-visited{{end}}">
                            <option value="Planned" {{if eq .Status "Planned"}}selected{{end}}>Planned</option>
                            <option value="Visited" {{if eq .Status "Visited"}}selected{{end}}>Visited</option>
                        </select>
                    </td>
                    <td class="wl-actions">
                        <button class="btn-save">Save</button>
                        <button class="btn-delete">Delete</button>
                    </td>
                </tr>
                {{end}}
            </tbody>
        </table>
        {{if not .WishlistItems}}
        <div class="empty-state">
            <p>Your wishlist is empty.
               <a href="/countries">Browse destinations</a> to add one.
            </p>
        </div>
        {{end}}
    </div>
</div>

<script src="/static/js/wishlist.js"></script>