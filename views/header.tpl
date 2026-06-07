{{define "header.tpl"}}
<header class="navbar">
    <div class="navbar-inner">
        <div class="navbar-left">
            <a href="/" class="navbar-brand">TravelSphere</a>
            <nav class="navbar-links">
                <a href="/"
                   class="nav-link {{if eq .ActiveNav "home"}}active{{end}}">
                   Home
                </a>
                <a href="/countries"
                   class="nav-link {{if eq .ActiveNav "countries"}}active{{end}}">
                   Countries
                </a>
                <a href="/wishlist"
                   class="nav-link {{if eq .ActiveNav "wishlist"}}active{{end}}">
                   Wishlist
                </a>
                <a href="/dashboard"
                   class="nav-link {{if eq .ActiveNav "dashboard"}}active{{end}}">
                   Dashboard
                </a>
            </nav>
        </div>

        <div class="navbar-right">
            {{if .IsLoggedIn}}
                <span class="nav-greeting">Hi, {{.Username}}</span>
                <a href="/logout" class="nav-link">Logout</a>
            {{else}}
                <a href="/login" class="nav-link">Login</a>
            {{end}}
        </div>
    </div>
</header>
{{end}}