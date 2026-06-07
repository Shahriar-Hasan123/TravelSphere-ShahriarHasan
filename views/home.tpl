<section class="hero">
    <div class="hero-inner">
        <h1 class="hero-title">Discover your next destination</h1>
        <p class="hero-subtitle">
            Search countries, explore attractions, and curate your personal travel wishlist.
        </p>

        <div class="search-block">
            <label class="search-label">WHERE TO NEXT?</label>
            <div class="search-wrap">
                <input
                    type="text"
                    id="home-search"
                    class="search-input"
                    placeholder="Search a destination..."
                    autocomplete="off"
                >
                <div id="search-suggestions" class="suggestions-dropdown hidden"></div>
            </div>
        </div>
    </div>
</section>

<div class="page-body">
    <!-- Featured destinations -->
    <section class="section">
        <h2 class="section-title">Featured destinations</h2>
        <div class="featured-grid">
            {{range .FeaturedCountries}}
            <a href="/countries/{{.Slug}}" class="featured-card">
                <div class="featured-flag">
                    <img src="{{.Flag}}" alt="{{.Name}} flag">
                </div>
                <div class="featured-info">
                    <span class="featured-name">{{.Name}}</span>
                    <span class="featured-meta">{{.Capital}} · {{.Region}}</span>
                </div>
            </a>
            {{end}}
        </div>
    </section>

    <!-- Popular attractions -->
    <section class="section">
        <h2 class="section-title">Popular attractions</h2>
        <div class="attractions-list">
            {{range .PopularAttractions}}
            <div class="attraction-row">
                <span class="attraction-name">{{.Name}}</span>
                <span class="attraction-tags">
                    {{range $i, $tag := .Kinds}}{{if $i}},{{end}}{{$tag}}{{end}}
                </span>
            </div>
            {{end}}
        </div>
    </section>
</div>

<script src="/static/js/home.js"></script>