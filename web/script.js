let ipHistory = [],
    filteredHistory = [],
    currentPage = 1,
    entriesPerPage = 10;

$(document).ready(function() {
    $.getJSON('/history', function(data) {
        ipHistory = data;
        filteredHistory = ipHistory;
        displayTable();
        updateTotalIPs();
    }).fail(function() {
        console.error('Error fetching IP history');
    });

    $('#searchBar').on('input', function() {
        const query = $(this).val().toLowerCase();
        filteredHistory = ipHistory.filter(entry =>
            entry.timestamp.toLowerCase().includes(query) ||
            (entry.ipv4 && entry.ipv4.toLowerCase().includes(query)) ||
            (entry.ipv6 && entry.ipv6.toLowerCase().includes(query))
        );
        currentPage = 1;
        displayTable();
        updateTotalIPs();
    });
});

function displayTable() {
    const tbody = $('#ipTable tbody').empty();
    const start = (currentPage - 1) * entriesPerPage;
    const end = start + entriesPerPage;
    const currentEntries = filteredHistory.slice(start, end);

    currentEntries.forEach(entry => {
        const tr = $('<tr>');
        tr.append($('<td>').text(entry.timestamp));
        tr.append($('<td>').text(entry.ipv4 || 'N/A').css('min-width', '120px')); // Set static width
        tr.append($('<td>').text(entry.ipv6 || 'N/A').css('min-width', '180px')); // Set static width
        tbody.append(tr);
    });
    updatePaginationButtons();
    updateTotalIPs();
}

function updateTotalIPs() {
    const totalIPs = filteredHistory.length;
    $('#searchBar').attr('placeholder', `Search IP history... (Total IPs: ${totalIPs})`);
}

function updatePaginationButtons() {
    $('#prevBtn').prop('disabled', currentPage === 1);
    $('#nextBtn').prop('disabled', currentPage === Math.ceil(filteredHistory.length / entriesPerPage));
}

function nextPage() {
    if (currentPage < Math.ceil(filteredHistory.length / entriesPerPage)) {
        currentPage++;
        displayTable();
    }
}

function prevPage() {
    if (currentPage > 1) {
        currentPage--;
        displayTable();
    }
}

function sortTable(field) {
    const isAscending = $(`#ipTable th:contains(${field.charAt(0).toUpperCase() + field.slice(1)})`).hasClass('sort-asc');
    const orderModifier = isAscending ? 1 : -1;

    filteredHistory.sort((a, b) => {
        if (a[field] && b[field]) {
            return (a[field].localeCompare(b[field])) * orderModifier;
        }
        return 0;
    });

    $('#ipTable th').removeClass('sort-asc sort-desc');
    $(`#ipTable th:contains(${field.charAt(0).toUpperCase() + field.slice(1)})`).addClass(isAscending ? 'sort-desc' : 'sort-asc');

    displayTable();
}

function exportToCSV() {
    const rows = [
        ['Timestamp', 'IPv4 Address', 'IPv6 Address'],
        ...filteredHistory.map(entry => [entry.timestamp, entry.ipv4 || '', entry.ipv6 || ''])
    ];
    const csvContent = "data:text/csv;charset=utf-8," + rows.map(e => e.join(",")).join("\n");
    const encodedUri = encodeURI(csvContent);
    const link = $('<a>').attr({
        href: encodedUri,
        download: 'ip_history.csv'
    });

    $('body').append(link);
    link[0].click();
    link.remove();
}
