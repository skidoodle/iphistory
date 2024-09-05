let ipHistory = [];
let filteredHistory = [];
let currentPage = 1;
let entriesPerPage = 10;

document.addEventListener('DOMContentLoaded', () => {
    fetch('/history')
        .then(response => response.json())
        .then(data => {
            ipHistory = data;
            filteredHistory = ipHistory;
            displayTable();
            updateTotalIPs();
        })
        .catch(error => console.error('Error fetching IP history:', error));

    const searchBar = document.getElementById('searchBar');
    searchBar.addEventListener('input', handleSearch);

    document.getElementById('prevBtn').addEventListener('click', prevPage);
    document.getElementById('nextBtn').addEventListener('click', nextPage);
    document.getElementById('exportBtn').addEventListener('click', exportToCSV);
    document.getElementById('sortTimestamp').addEventListener('click', () => sortTable('timestamp'));
    document.getElementById('sortIPv4').addEventListener('click', () => sortTable('ipv4'));
    document.getElementById('sortIPv6').addEventListener('click', () => sortTable('ipv6'));
});


function handleSearch() {
    const query = this.value.toLowerCase();
    filteredHistory = ipHistory.filter(entry =>
        entry.timestamp.toLowerCase().includes(query) ||
        (entry.ipv4 && entry.ipv4.toLowerCase().includes(query)) ||
        (entry.ipv6 && entry.ipv6.toLowerCase().includes(query))
    );
    currentPage = 1;
    displayTable();
    updateTotalIPs();
}

function displayTable() {
    const tbody = document.querySelector('#ipTable tbody');
    tbody.innerHTML = '';

    const start = (currentPage - 1) * entriesPerPage;
    const end = start + entriesPerPage;
    const currentEntries = filteredHistory.slice(start, end);

    currentEntries.forEach(entry => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${entry.timestamp}</td>
            <td>${entry.ipv4 || 'N/A'}</td>
            <td>${entry.ipv6 || 'N/A'}</td>
        `;
        tbody.appendChild(row);
    });
    updatePaginationButtons();
    updateTotalIPs();
}

function updateTotalIPs() {
    const totalIPs = filteredHistory.length;
    const searchBar = document.getElementById('searchBar');
    searchBar.placeholder = `Search IP history... (Total IPs: ${totalIPs})`;
}

function updatePaginationButtons() {
    const prevBtn = document.getElementById('prevBtn');
    const nextBtn = document.getElementById('nextBtn');
    prevBtn.disabled = currentPage === 1;
    nextBtn.disabled = currentPage === Math.ceil(filteredHistory.length / entriesPerPage);
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
    const isAscending = document.getElementById(`sort${capitalizeFirstLetter(field)}`).classList.contains('sort-asc');
    filteredHistory.sort((a, b) => (a[field] && b[field]) ? a[field].localeCompare(b[field]) * (isAscending ? 1 : -1) : 0);

    // Toggle the sort direction class
    document.querySelectorAll('th').forEach(th => th.classList.remove('sort-asc', 'sort-desc'));
    document.getElementById(`sort${capitalizeFirstLetter(field)}`).classList.add(isAscending ? 'sort-desc' : 'sort-asc');

    displayTable();
}

function exportToCSV() {
    const rows = [
        ['Timestamp', 'IPv4 Address', 'IPv6 Address'],
        ...filteredHistory.map(entry => [entry.timestamp, entry.ipv4 || '', entry.ipv6 || ''])
    ];
    const csvContent = rows.map(row => row.join(',')).join('\n');
    const blob = new Blob([csvContent], { type: 'text/csv' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = 'ip_history.csv';
    link.click();
}

function capitalizeFirstLetter(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
}
