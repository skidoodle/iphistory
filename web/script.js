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
        filteredHistory = ipHistory.filter(entry => entry.timestamp.toLowerCase().includes(query) || entry.ip_address.toLowerCase().includes(query));
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
        tr.append($('<td>').text(entry.ip_address));
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

function sortTable() {
    const isAscending = $('#ipTable th').eq(0).hasClass('sort-asc');
    const orderModifier = isAscending ? 1 : -1;

    filteredHistory.sort((a, b) => (a.timestamp.localeCompare(b.timestamp)) * orderModifier);

    $('#ipTable th').removeClass('sort-asc sort-desc');
    $('#ipTable th').eq(0).addClass(isAscending ? 'sort-desc' : 'sort-asc');

    displayTable();
}

function exportToCSV() {
    const rows = [
        ['Timestamp', 'IP Address'], ...filteredHistory.map(entry => [entry.timestamp, entry.ip_address])
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
