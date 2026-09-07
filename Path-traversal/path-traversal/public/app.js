document.addEventListener('DOMContentLoaded', () => {
    const productGrid = document.getElementById('product-grid');
    const errorMessage = document.getElementById('error-message');

    fetchProducts('products.json');

    function fetchProducts(filePath) {

        fetch(`/api/content?file=${filePath}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error(`Failed to load product data. Status: ${response.status}`);
                }
                return response.json();
            })
            .then(products => {
                renderProducts(products);
            })
            .catch(error => {
                console.error("Error loading products:", error);
                showError("System Failure: Unable to connect to the inventory database.");
                productGrid.innerHTML = '';
            });
    }

    function renderProducts(products) {
        productGrid.innerHTML = '';
        products.forEach(product => {
            const card = document.createElement('div');
            card.className = 'product-card';

            card.innerHTML = `
                <img src="${product.image}" alt="${product.name}" class="product-image">
                <div class="product-info">
                    <h2 class="product-title">${product.name}</h2>
                    <p class="product-desc">${product.description}</p>
                    <div class="product-price">$${product.price.toFixed(2)}</div>
                    <button class="buy-button" onclick="alert('Order placed securely.')">Purchase Hardware</button>
                </div>
            `;

            productGrid.appendChild(card);
        });
    }

    function showError(message) {
        errorMessage.textContent = message;
        errorMessage.classList.remove('hidden');
    }
});
