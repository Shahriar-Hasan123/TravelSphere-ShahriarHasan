// Handles Add to Wishlist AJAX on the destination detail page.
// Only #wishlist-feedback is updated — no page reload.

(function () {
  const btn      = document.getElementById('add-wishlist-btn');
  const feedback = document.getElementById('wishlist-feedback');

  if (!btn || !feedback) return;

  btn.addEventListener('click', function () {
    const countryName = btn.dataset.country;
    const isLoggedIn  = btn.dataset.isLoggedIn === 'true';

    if (!isLoggedIn) {
      window.location.href = '/login?redirect_to=' + encodeURIComponent(window.location.pathname);
      return;
    }

    btn.disabled       = true;
    btn.textContent    = 'Adding...';
    feedback.textContent = '';
    feedback.className   = 'wishlist-feedback';

    fetch('/api/wishlist', {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ country_name: countryName, status: 'Planned' }),
    })
      .then(function (res) {
        return res.json().then(function (json) {
          return { status: res.status, json };
        });
      })
      .then(function ({ status, json }) {
        if (status === 201) {
          // Successfully added.
          feedback.textContent = '✓ Added to wishlist!';
          feedback.classList.add('success');
          btn.textContent = 'Added to Wishlist';
          // Button stays disabled — no point adding again.
          return;
        }

        if (status === 409) {
          // Already exists — treat as informational, not an error.
          feedback.textContent = '✓ Already in your wishlist.';
          feedback.classList.add('success');
          btn.textContent = 'Already in Wishlist';
          // Button stays disabled.
          return;
        }

        // All other errors — re-enable the button.
        throw new Error(json.message || 'Failed to add');
      })
      .catch(function (err) {
        feedback.textContent = err.message;
        feedback.classList.add('error');
        btn.disabled    = false;
        btn.textContent = 'Add to Wishlist';
      });
  });
}());