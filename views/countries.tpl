<div class="page-body">
    <div class="page-header">
        <h1 class="page-title">Country Explorer</h1>
        <p class="page-subtitle">
            Browse every destination on first load. Search and filter update only the
            results below &mdash; no full page reload.
        </p>
    </div>

    <!-- Search and filter bar -->
    <div class="filter-bar">
        <div class="filter-group">
            <label class="filter-label">SEARCH</label>
            <input
                type="text"
                id="country-search"
                class="filter-input"
                placeholder="Country or capital..."
                value="{{.SearchQuery}}"
            >
        </div>
        <div class="filter-group">
            <label class="filter-label">REGION</label>
            <select id="region-filter" class="filter-select">
                <option value="">All regions</option>
                <option value="Africa"   {{if eq .RegionFilter "Africa"}}selected{{end}}>Africa</option>
                <option value="Americas" {{if eq .RegionFilter "Americas"}}selected{{end}}>Americas</option>
                <option value="Asia"     {{if eq .RegionFilter "Asia"}}selected{{end}}>Asia</option>
                <option value="Europe"   {{if eq .RegionFilter "Europe"}}selected{{end}}>Europe</option>
                <option value="Oceania"  {{if eq .RegionFilter "Oceania"}}selected{{end}}>Oceania</option>
            </select>
        </div>
    </div>

    <!-- Country results — AJAX replaces only this container -->
    <div id="country-results" class="country-grid">
        {{range .Countries}}
        <a href="/countries/{{.Slug}}" class="country-card">
            <div class="country-flag">
                <img src="{{.Flag}}" alt="{{.Name}} flag" loading="lazy">
            </div>
            <div class="country-info">
                <h3 class="country-name">{{.Name}}</h3>
                <p><span class="info-label">Capital:</span> {{.Capital}}</p>
                <p><span class="info-label">Population:</span> {{.Population}}</p>
                <p><span class="info-label">Currency:</span> {{.Currency}}</p>
                <p>
                    <span class="info-label">Languages:</span>
                    {{range $i, $l := .Languages}}{{if $i}}, {{end}}{{$l}}{{end}}
                </p>
            </div>
        </a>
        {{end}}
    </div>
</div>

<script src="/static/js/countries.js"></script>