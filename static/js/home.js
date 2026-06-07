// Handles the home page destination search suggestion dropdown.
// Updates only #search-suggestions - no full page reload.

(function () {
  const input      = document.getElementById('home-search');
  const dropdown   = document.getElementById('search-suggestions');

  if (!input || !dropdown) return;

  let debounceTimer = null;

  input.addEventListener('input', function () {
    clearTimeout(debounceTimer);
    const query = input.value.trim();

    if (query.length < 1) {
      hideDropdown();
      return;
    }

    debounceTimer = setTimeout(function () {
      fetchSuggestions(query);
    }, 250);
  });

  // Close dropdown when clicking outside.
  document.addEventListener('click', function (e) {
    if (!input.contains(e.target) && !dropdown.contains(e.target)) {
      hideDropdown();
    }
  });

  function fetchSuggestions(query) {
    fetch('/api/countries/suggestions?q=' + encodeURIComponent(query))
      .then(function (res) { return res.json(); })
      .then(function (json) {
        renderSuggestions(json.data || []);
      })
      .catch(hideDropdown);
  }

  function renderSuggestions(items) {
    if (items.length === 0) {
      hideDropdown();
      return;
    }

    dropdown.innerHTML = items.map(function (item) {
      return (
        '<div class="suggestion-item" data-slug="' + item.Slug + '">' +
          item.Name + ' &mdash; ' + (item.Capital || '') +
        '</div>'
      );
    }).join('');

    // Navigate to destination page on suggestion click.
    dropdown.querySelectorAll('.suggestion-item').forEach(function (el) {
      el.addEventListener('click', function () {
        window.location.href = '/countries/' + el.dataset.slug;
      });
    });

    showDropdown();
  }

  function showDropdown() { dropdown.classList.remove('hidden'); }
  function hideDropdown()  { dropdown.classList.add('hidden'); }
}());