// Handles Save and Delete AJAX on the wishlist page.
// Only #wishlist-rows is replaced — never a full page reload.

(function () {
  const container = document.getElementById('wishlist-rows');
  if (!container) return;

  container.addEventListener('click', function (e) {
    const row = e.target.closest('tr[data-id]');
    if (!row) return;

    if (e.target.classList.contains('btn-save')) {
      saveRow(row);
    } else if (e.target.classList.contains('btn-delete')) {
      deleteRow(row);
    }
  });

  function saveRow(row) {
    const id     = row.dataset.id;
    const note   = row.querySelector('.wl-note-input').value.trim();
    const status = row.querySelector('.wl-status-select').value;

    fetch('/api/wishlist/' + encodeURIComponent(id), {
      method:  'PUT',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ note, status }),
    })
      .then(function (res) {
        return res.json().then(function (json) {
          if (!res.ok) throw new Error(json.message || 'Update failed');
          return json;
        });
      })
      .then(function () {
        // PUT returns only the updated item — fetch the full list to re-render.
        return fetchAndRender();
      })
      .then(function () {
        showSuccess('✓ Saved successfully!');
      })
      .catch(showError);
  }

  function deleteRow(row) {
    const id = row.dataset.id;

    fetch('/api/wishlist/' + encodeURIComponent(id), { method: 'DELETE' })
      .then(function (res) {
        if (res.status === 204) {
          // 204 No Content — success, fetch the updated list.
          return fetchAndRender();
        }
        return res.json().then(function (json) {
          throw new Error(json.message || 'Delete failed');
        });
      })
      .then(function () {
        showSuccess('✓ Deleted successfully!');
      })
      .catch(showError);
  }

  // fetchAndRender fetches the current wishlist and re-renders #wishlist-rows.
  function fetchAndRender() {
    return fetch('/api/wishlist')
      .then(function (res) { return res.json(); })
      .then(function (json) {
        renderRows(json.data || []);
        notifyDashboard();
      });
  }

  // renderRows rebuilds the entire tbody from the updated items array.
  function renderRows(items) {
    const wrap = document.getElementById('wishlist-rows');
    if (!items || items.length === 0) {
      wrap.innerHTML =
        '<div class="empty-state"><p>Your wishlist is empty. ' +
        '<a href="/countries">Browse destinations</a> to add one.</p></div>';
      return;
    }

    // Rebuild the full table with fresh data from the server.
    let html =
      '<table class="wishlist-table">' +
        '<thead><tr>' +
          '<th>COUNTRY</th><th>NOTE</th><th>STATUS</th><th>ACTIONS</th>' +
        '</tr></thead>' +
        '<tbody>';

    items.forEach(function (item) {
      const visitedClass = item.Status === 'Visited' ? 'status-visited' : '';
      html +=
        '<tr data-id="' + item.id + '">' +
          '<td class="wl-country">' + escHtml(item.CountryName) + '</td>' +
          '<td class="wl-note">' +
            '<input type="text" class="wl-note-input" ' +
              'value="' + escHtml(item.Note) + '" placeholder="Add a note...">' +
          '</td>' +
          '<td class="wl-status">' +
            '<select class="wl-status-select ' + visitedClass + '">' +
              '<option value="Planned"' + sel(item.Status, 'Planned') + '>Planned</option>' +
              '<option value="Visited"' + sel(item.Status, 'Visited') + '>Visited</option>' +
            '</select>' +
          '</td>' +
          '<td class="wl-actions">' +
            '<button class="btn-save">Save</button>' +
            '<button class="btn-delete">Delete</button>' +
          '</td>' +
        '</tr>';
    });

    html += '</tbody></table>';
    wrap.innerHTML = html;
  }

  function sel(current, value) {
    return current === value ? ' selected' : '';
  }

  function escHtml(str) {
    return String(str || '')
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;');
  }

  function notifyDashboard() {
    window.dispatchEvent(new CustomEvent('wishlist-updated'));
    refreshDestList();
  }

  function refreshDestList() {
    const destList = document.getElementById('dest-list');
    if (!destList) return;

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
              '<span class="dest-row-status ' + item.Status + '">' + item.Status + '</span>' +
              note +
            '</div>'
          );
        }).join('');
      })
      .catch(function () {});
  }

  function showError() {
    const wrap = document.getElementById('wishlist-rows');
    wrap.insertAdjacentHTML('afterbegin',
      '<div class="alert alert-error">Action failed. Please try again.</div>',
    );
  }

  function showSuccess(message) {
    const wrap = document.getElementById('wishlist-rows');
    wrap.insertAdjacentHTML('afterbegin',
      '<div class="alert alert-success">' + escHtml(message) + '</div>',
    );
    // Auto-remove after 2 seconds
    setTimeout(function () {
      const alert = wrap.querySelector('.alert-success');
      if (alert) alert.remove();
    }, 2000);
  }

}());