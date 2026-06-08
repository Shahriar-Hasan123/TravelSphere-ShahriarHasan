// Refreshes the dashboard destination list when the wishlist changes,
// including when the user returns via browser back/forward navigation.

(function () {
  const destList = document.getElementById('dest-list');
  if (!destList) return;

  window.addEventListener('wishlist-updated', refreshDestList);
  window.addEventListener('pageshow', function (event) {
    if (event.persisted || sessionStorage.getItem('wishlist-updated')) {
      refreshDestList();
      sessionStorage.removeItem('wishlist-updated');
    }
  });

  function refreshDestList() {
    fetch('/api/wishlist')
      .then(function (res) { return res.json(); })
      .then(function (json) {
        const items = json.data || [];
        if (items.length === 0) {
          destList.innerHTML =
            '<div class="empty-state"><p>No saved destinations yet. ' +
            '<a href="/countries">Browse countries</a> to get started.</p></div>';
          return;
        }

        destList.innerHTML = items.map(function (item) {
          const note = item.Note
            ? '<span class="dest-row-note">&middot; ' + escHtml(item.Note) + '</span>'
            : '';
          return (
            '<div class="dest-row">' +
              '<span class="dest-row-name">' + escHtml(item.CountryName) + '</span>' +
              '<span class="dest-row-sep">&mdash;</span>' +
              '<span class="dest-row-status ' + escHtml(item.Status) + '">' + escHtml(item.Status) + '</span>' +
              note +
            '</div>'
          );
        }).join('');
      })
      .catch(function () {});
  }

  function escHtml(str) {
    return String(str || '')
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;');
  }
}());
