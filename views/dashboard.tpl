<div class="page-body">
    <div class="page-header">
        <h1 class="page-title">Travel Dashboard</h1>
        <p class="page-subtitle">
            Your saved trips at a glance. Stats refresh automatically
            when your wishlist changes.
        </p>
    </div>

    <!-- AJAX replaces only this container — stat cards only -->
    <div id="dashboard-stats" class="stats-row">
        <div class="stat-card">
            <span class="stat-label">TOTAL SAVED</span>
            <span class="stat-number">{{.TotalSaved}}</span>
        </div>
        <div class="stat-card">
            <span class="stat-label">PLANNED</span>
            <span class="stat-number">{{.Planned}}</span>
        </div>
        <div class="stat-card">
            <span class="stat-label">VISITED</span>
            <span class="stat-number">{{.Visited}}</span>
        </div>
    </div>

    <!-- Saved destinations list — not part of AJAX refresh -->
    <div class="section">
        <h2 class="section-title">Saved destinations</h2>
        <div class="dest-list" id="dest-list">
            {{if .WishlistItems}}
                {{range .WishlistItems}}
                <div class="dest-row">
                    <span class="dest-row-name">{{.CountryName}}</span>
                    <span class="dest-row-sep">&mdash;</span>
                    <span class="dest-row-status {{.Status}}">{{.Status}}</span>
                    {{if .Note}}
                    <span class="dest-row-note">&middot; {{.Note}}</span>
                    {{end}}
                </div>
                {{end}}
            {{else}}
                <div class="empty-state">
                    <p>No saved destinations yet.
                       <a href="/countries">Browse countries</a> to get started.
                    </p>
                </div>
            {{end}}
        </div>
    </div>
</div>

<script src="/static/js/dashboard.js"></script>