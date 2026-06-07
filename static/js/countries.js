// Handles AJAX search and region filter on the Country Explorer page.

(function () {
  const searchInput  = document.getElementById('country-search');
  const regionSelect = document.getElementById('region-filter');
  const resultsEl    = document.getElementById('country-results');

  if (!searchInput || !resultsEl) return;

  let debounceTimer = null;

  // Attach listeners to both controls.
  searchInput.addEventListener('input', scheduleUpdate);
  regionSelect.addEventListener('change', fetchCountries);

  function scheduleUpdate() {
    // Debounce: wait 300ms after user stops typing before firing.
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(fetchCountries, 300);
  }

  function fetchCountries() {
    const search = searchInput.value.trim();
    const region = regionSelect.value;

    // Build query string.
    const params = new URLSearchParams();
    if (search) params.set('search', search);
    if (region) params.set('region', region);

    showSpinner();

    fetch('/api/countries?' + params.toString())
      .then(function (res) {
        if (!res.ok) throw new Error('Network error ' + res.status);
        return res.json();
      })
      .then(function (json) {
        renderCards(json.data || []);
      })
      .catch(function () {
        showError();
      });
  }

  // buildCard constructs one country card HTML string from a country object.
  function buildCard(country) {
    const langs = (country.Languages || []).join(', ');
    return (
      '<a href="/countries/' + country.Slug + '" class="country-card">' +
        '<div class="country-flag">' +
          '<img src="' + country.Flag + '" alt="' + country.Name + ' flag" loading="lazy">' +
        '</div>' +
        '<div class="country-info">' +
          '<h3 class="country-name">' + country.Name + '</h3>' +
          '<p><span class="info-label">Capital:</span> ' + (country.Capital || '—') + '</p>' +
          '<p><span class="info-label">Population:</span> ' + (country.FormattedPop || '—') + '</p>' +
          '<p><span class="info-label">Currency:</span> ' + (country.Currency || '—') + '</p>' +
          '<p><span class="info-label">Languages:</span> ' + (langs || '—') + '</p>' +
        '</div>' +
      '</a>'
    );
  }

  function renderCards(countries) {
    if (countries.length === 0) {
      resultsEl.innerHTML =
        '<div class="empty-state"><p>No countries match your search.</p></div>';
      return;
    }
    resultsEl.innerHTML = countries.map(buildCard).join('');
  }

  function showSpinner() {
    resultsEl.innerHTML = '<div class="spinner">Loading...</div>';
  }

  function showError() {
    resultsEl.innerHTML =
      '<div class="empty-state"><p>Failed to load countries. Please try again.</p></div>';
  }
}());