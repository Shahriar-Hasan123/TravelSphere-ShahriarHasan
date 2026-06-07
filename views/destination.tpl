<div class="page-body">
    <!-- Country info card -->
    <div class="dest-card">
        <div class="dest-flag">
            <img src="{{.Country.Flag}}" alt="{{.Country.Name}} flag">
        </div>
        <div class="dest-info">
            <span class="dest-region-badge">{{.Country.Region}}</span>
            <h1 class="dest-name">{{.Country.Name}}</h1>
            <p class="dest-official">{{.Country.OfficialName}}</p>
            <div class="dest-meta-row">
                <div class="dest-meta-item">
                    <span class="meta-label">CAPITAL</span>
                    <span class="meta-value">{{.Country.Capital}}</span>
                </div>
                <div class="dest-meta-item">
                    <span class="meta-label">POPULATION</span>
                    <span class="meta-value">{{.Country.FormattedPop}}</span>
                </div>
                <div class="dest-meta-item">
                    <span class="meta-label">REGION</span>
                    <span class="meta-value">
                        {{.Country.Region}}<br>{{.Country.Subregion}}
                    </span>
                </div>
                <div class="dest-meta-item">
                    <span class="meta-label">CURRENCY</span>
                    <span class="meta-value">{{.Country.Currency}}</span>
                </div>
                <div class="dest-meta-item">
                    <span class="meta-label">LANGUAGES</span>
                    <span class="meta-value">
                        {{range $i, $l := .Country.Languages}}{{if $i}}, {{end}}{{$l}}{{end}}
                    </span>
                </div>
            </div>
        </div>
    </div>

    <!-- Add to Wishlist — AJAX updates #wishlist-feedback only -->
    {{if .IsLoggedIn}}
    <div class="wishlist-action">
        <button id="add-wishlist-btn" class="btn-wishlist"
                data-country="{{.Country.Name}}">
            Add to Wishlist
        </button>
        <div id="wishlist-feedback" class="wishlist-feedback"></div>
    </div>
    {{end}}

    <!-- Weather + Attractions row -->
    <div class="dest-panels">
        <!-- Travel weather panel -->
        <div class="panel panel-weather">
            <h2 class="panel-title">Travel weather</h2>
            {{if .Weather}}
                <!-- Current conditions -->
                <div class="weather-current">
                    {{if .Weather.Icon}}
                    <img src="{{.Weather.Icon}}" alt="{{.Weather.Condition}}"
                         class="weather-icon">
                    {{end}}
                    <p class="weather-temp">
                        {{printf "%.1f" .Weather.TempC}}&deg;C
                        / {{printf "%.1f" .Weather.TempF}}&deg;F
                    </p>
                    <p class="weather-condition">{{.Weather.Condition}}</p>
                    <p class="weather-city">{{.Weather.City}}</p>
                </div>

                <!-- Detailed stats -->
                <div class="weather-stats">
                    <div class="weather-stat">
                        <span class="wstat-label">Feels like</span>
                        <span class="wstat-value">{{printf "%.1f" .Weather.FeelsLikeC}}&deg;C</span>
                    </div>
                    <div class="weather-stat">
                        <span class="wstat-label">Humidity</span>
                        <span class="wstat-value">{{.Weather.Humidity}}%</span>
                    </div>
                    <div class="weather-stat">
                        <span class="wstat-label">Wind</span>
                        <span class="wstat-value">{{printf "%.0f" .Weather.WindKph}} km/h</span>
                    </div>
                </div>

                <!-- Tomorrow's forecast -->
                {{if .Weather.Forecast}}
                <div class="weather-forecast">
                    <p class="forecast-label">Tomorrow</p>
                    <div class="forecast-row">
                        {{if .Weather.Forecast.Icon}}
                        <img src="{{.Weather.Forecast.Icon}}"
                             alt="{{.Weather.Forecast.Condition}}"
                             class="forecast-icon">
                        {{end}}
                        <div>
                            <p class="forecast-condition">{{.Weather.Forecast.Condition}}</p>
                            <p class="forecast-temps">
                                {{printf "%.1f" .Weather.Forecast.MaxTempC}}&deg;C
                                / {{printf "%.1f" .Weather.Forecast.MinTempC}}&deg;C
                            </p>
                        </div>
                    </div>
                </div>
                {{end}}

                <!-- Travel conditions advice -->
                <div class="travel-advice">
                    <p class="advice-label">Travel conditions</p>
                    <p class="advice-text">{{.Weather.TravelAdvice}}</p>
                </div>

            {{else}}
                <div class="weather-placeholder">
                    <p>
                        Weather data is optional. Add
                        <code>WEATHERAPI_KEY</code> to your
                        <code>.env</code> file to enable live conditions.
                    </p>
                </div>
            {{end}}
        </div>

        <!-- Attractions & landmarks panel -->
        <div class="panel panel-attractions">
            <h2 class="panel-title">Attractions &amp; landmarks</h2>
            {{if .Attractions}}
                <div class="attractions-list">
                    {{range .Attractions}}
                    <div class="attraction-row">
                        <span class="attraction-name">{{.Name}}</span>
                        <span class="attraction-tags">
                            {{range $i, $tag := .Kinds}}{{if $i}},{{end}}{{$tag}}{{end}}
                        </span>
                    </div>
                    {{end}}
                </div>
            {{else}}
                <p class="no-data">No attractions found for this destination.</p>
            {{end}}
        </div>
    </div>
</div>

<script src="/static/js/destination.js"></script>