const API_BASE_URL = '/api/v1';
const ORDERS_PER_PAGE = 10;

const ordersContainer = document.getElementById('orders-container');
const paginationContainer = document.getElementById('pagination');
const searchInput = document.getElementById('search-input');
const searchButton = document.getElementById('search-button');
const clearSearchButton = document.getElementById('clear-search');
const searchInfo = document.getElementById('search-info');
const orderModal = document.getElementById('order-modal');
const orderDetails = document.getElementById('order-details');
const closeModal = document.querySelector('.close');
const notification = document.getElementById('notification');

let currentPage = 0;
let totalOrders = 0;
let currentSearchTerm = '';

async function loadOrders(page = 0, searchTerm = '') {
    try {
        ordersContainer.innerHTML = '<div class="loader">Загрузка...</div>';

        let url = `${API_BASE_URL}/orders?page=${page}&size=${ORDERS_PER_PAGE}`;
        if (searchTerm) url += `&search=${encodeURIComponent(searchTerm)}`;

        const response = await fetch(url);
        if (!response.ok) throw new Error(`Ошибка: ${response.status}`);

        const data = await response.json();
        if (!data.orders_ids || data.orders_ids.length === 0) {
            ordersContainer.innerHTML = '<div class="error">Заказы не найдены</div>';
            return;
        }

        totalOrders = data.total_orders;
        displayOrders(data.orders_ids);
        updatePagination(page, searchTerm);
        currentPage = page;
        currentSearchTerm = searchTerm;
        updateSearchInfo(searchTerm);

    } catch (error) {
        ordersContainer.innerHTML = '<div class="error">Ошибка загрузки</div>';
    }
}

function updateSearchInfo(searchTerm) {
    searchInfo.style.display = searchTerm ? 'block' : 'none';
    searchInfo.innerHTML = searchTerm ? `Поиск: <strong>"${searchTerm}"</strong>` : '';
}

function displayOrders(orders) {
    const list = document.createElement('ul');
    list.className = 'orders-list';

    orders.forEach(orderId => {
        const listItem = document.createElement('li');
        listItem.className = 'order-item';
        listItem.innerHTML = `
            <span class="order-id">Заказ ${orderId}</span>
            <button class="copy-btn">Копировать номер</button>
        `;

        listItem.querySelector('.copy-btn').addEventListener('click', (e) => {
            e.stopPropagation();
            copyToClipboard(orderId);
        });

        listItem.addEventListener('click', () => showOrderDetails(orderId));
        list.appendChild(listItem);
    });

    ordersContainer.innerHTML = '';
    ordersContainer.appendChild(list);
}

function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
        notification.classList.add('show');
        setTimeout(() => notification.classList.remove('show'), 2000);
    });
}

function updatePagination(currentPage, searchTerm = '') {
    paginationContainer.innerHTML = '';
    const totalPages = Math.ceil(totalOrders / ORDERS_PER_PAGE);

    if (currentPage > 0) {
        const prevButton = document.createElement('button');
        prevButton.textContent = '←';
        prevButton.addEventListener('click', () => loadOrders(currentPage - 1, searchTerm));
        paginationContainer.appendChild(prevButton);
    }

    const startPage = Math.max(0, currentPage - 2);
    const endPage = Math.min(totalPages - 1, currentPage + 2);

    for (let i = startPage; i <= endPage; i++) {
        const pageButton = document.createElement('button');
        pageButton.textContent = i + 1;
        if (i === currentPage) pageButton.classList.add('active');
        pageButton.addEventListener('click', () => loadOrders(i, searchTerm));
        paginationContainer.appendChild(pageButton);
    }

    if (currentPage < totalPages - 1) {
        const nextButton = document.createElement('button');
        nextButton.textContent = '→';
        nextButton.addEventListener('click', () => loadOrders(currentPage + 1, searchTerm));
        paginationContainer.appendChild(nextButton);
    }
}

async function showOrderDetails(orderId) {
    try {
        orderDetails.innerHTML = '<div class="loader">Загрузка...</div>';
        orderModal.style.display = 'flex';

        const response = await fetch(`${API_BASE_URL}/orders/${orderId}`);
        if (!response.ok) throw new Error(`Ошибка: ${response.status}`);

        const data = await response.json();
        if (!data.order) throw new Error('Данные не найдены');

        displayOrderDetails(data.order);

    } catch (error) {
        orderDetails.innerHTML = '<div class="error">Ошибка загрузки</div>';
    }
}

function displayOrderDetails(order) {
    orderDetails.innerHTML = `
        <h2>Заказ #${order.order_uid}</h2>
        <div class="order-details">
            <div class="detail-section">
                <h3>Основная информация</h3>
                <div class="detail-item"><span class="detail-label">Трек-номер:</span> ${order.track_number}</div>
                <div class="detail-item"><span class="detail-label">Дата:</span> ${new Date(order.date_created).toLocaleString()}</div>
                <div class="detail-item"><span class="detail-label">Служба доставки:</span> ${order.delivery_service}</div>
            </div>
            
            <div class="detail-section">
                <h3>Доставка</h3>
                <div class="detail-item"><span class="detail-label">Получатель:</span> ${order.delivery.name}</div>
                <div class="detail-item"><span class="detail-label">Адрес:</span> ${order.delivery.city}, ${order.delivery.address}</div>
                <div class="detail-item"><span class="detail-label">Регион:</span> ${order.delivery.region}</div>
            </div>
            
            <div class="detail-section">
                <h3>Оплата</h3>
                <div class="detail-item"><span class="detail-label">Сумма:</span> ${order.payment.amount} ${order.payment.currency}</div>
                <div class="detail-item"><span class="detail-label">Провайдер:</span> ${order.payment.provider}</div>
            </div>
            
            <div class="detail-section">
                <h3>Товары</h3>
                ${order.items.map(item => `
                    <div class="detail-item"><span class="detail-label">Наименование:</span> ${item.name}</div>
                    <div class="detail-item"><span class="detail-label">Цена:</span> ${item.price}</div>
                    <div class="detail-item"><span class="detail-label">Итоговая цена:</span> ${item.total_price}</div>
                `).join('')}
            </div>
        </div>
    `;
}

function searchOrders() {
    const searchTerm = searchInput.value.trim();
    loadOrders(0, searchTerm);
}

function clearSearch() {
    searchInput.value = '';
    loadOrders(0);
}

function init() {
    loadOrders(0);

    searchButton.addEventListener('click', searchOrders);
    clearSearchButton.addEventListener('click', clearSearch);

    searchInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') searchOrders();
    });

    closeModal.addEventListener('click', () => {
        orderModal.style.display = 'none';
    });

    window.addEventListener('click', (e) => {
        if (e.target === orderModal) orderModal.style.display = 'none';
    });
}

init();