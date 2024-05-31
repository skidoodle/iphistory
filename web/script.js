let ipHistory = [], filteredHistory = [], currentPage = 1, entriesPerPage = 10;

$(document).ready(function() {
    $.getJSON('/history', function(data) {
        ipHistory = data.map(entry => entry.split(" - "));
        filteredHistory = ipHistory;
        displayTable();
    }).fail(function() {
        console.error('Error fetching IP history');
    });

    $('#searchBar').on('input', function() {
        const query = $(this).val().toLowerCase();
        filteredHistory = ipHistory.filter(row => row.some(cell => cell.toLowerCase().includes(query)));
        currentPage = 1;
        displayTable();
    });
});

function displayTable() {
    const tbody = $('#ipTable tbody').empty();
    const start = (currentPage - 1) * entriesPerPage;
    const end = start + entriesPerPage;
    const currentEntries = filteredHistory.slice(start, end);

    currentEntries.forEach(row => {
        const tr = $('<tr>');
        row.forEach(cell => {
            tr.append($('<td>').text(cell));
        });
        tbody.append(tr);
    });
    updatePaginationButtons();
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
    const th = $('#ipTable th').eq(0);
    const isAscending = !th.hasClass('sort-asc');
    const orderModifier = isAscending ? 1 : -1;

    filteredHistory.sort((a, b) => a[0].localeCompare(b[0]) * orderModifier);
    th.toggleClass('sort-asc', isAscending).toggleClass('sort-desc', !isAscending);

    displayTable();
}

function exportToCSV() {
    const rows = [['Timestamp', 'IP Address'], ...filteredHistory];
    const csvContent = "data:text/csv;charset=utf-8," + rows.map(e => e.join(",")).join("\n");
    const encodedUri = encodeURI(csvContent);
    const link = $('<a>').attr({href: encodedUri, download: 'ip_history.csv'});

    $('body').append(link);
    link[0].click();
    link.remove();
}
