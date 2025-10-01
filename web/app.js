class CryptoApp {
    constructor() {
        this.ws = null;
        this.cryptos = new Map();
        this.filteredCryptos = [];
        this.searchTerm = '';
        this.sortBy = 'name';
        
        this.elements = {
            status: document.getElementById('status'),
            search: document.getElementById('search'),
            sort: document.getElementById('sort'),
            grid: document.getElementById('cryptoGrid'),
            count: document.getElementById('count')
        };
        
        this.init();
    }

    init() {
        this.setupEvents();
        this.connect();
    }

    connect() {
        try {
            this.ws = new WebSocket('ws://localhost:8080/ws');
            
            this.ws.onopen = () => {
                console.log('Connected to WebSocket');
                this.updateStatus('Connected', 'connected');
            };

            this.ws.onmessage = (event) => {
                const updates = JSON.parse(event.data);
                this.handleUpdates(updates);
            };

            this.ws.onclose = () => {
                console.log('WebSocket closed');
                this.updateStatus('Disconnected', 'disconnected');
                setTimeout(() => this.connect(), 3000);
            };

            this.ws.onerror = (error) => {
                console.error('WebSocket error:', error);
                this.updateStatus('Connection Error', 'disconnected');
            };

        } catch (error) {
            console.error('Failed to connect:', error);
            this.updateStatus('Failed to Connect', 'disconnected');
            setTimeout(() => this.connect(), 3000);
        }
    }

    handleUpdates(updates) {
        updates.forEach(crypto => {
            const prev = this.cryptos.get(crypto.id);
            crypto.priceChange = prev ? crypto.current_price - prev.current_price : 0;
            crypto.updated = true;
            this.cryptos.set(crypto.id, crypto);
        });

        this.updateCount();
        this.filterAndRender();
        
        // Clear update flags after animation
        setTimeout(() => {
            updates.forEach(crypto => {
                const stored = this.cryptos.get(crypto.id);
                if (stored) stored.updated = false;
            });
        }, 500);
    }

    filterAndRender() {
        let data = Array.from(this.cryptos.values());
        
        // Filter by search term
        if (this.searchTerm) {
            data = data.filter(crypto => 
                crypto.name.toLowerCase().includes(this.searchTerm) ||
                crypto.symbol.toLowerCase().includes(this.searchTerm) ||
                crypto.id.toLowerCase().includes(this.searchTerm)
            );
        }
        
        // Sort data
        data.sort((a, b) => {
            switch (this.sortBy) {
                case 'price': return b.current_price - a.current_price;
                case 'symbol': return a.symbol.localeCompare(b.symbol);
                case 'name':
                default: return a.name.localeCompare(b.name);
            }
        });
        
        this.filteredCryptos = data;
        this.render();
    }

    render() {
        if (this.filteredCryptos.length === 0) {
            if (this.cryptos.size === 0) {
                this.elements.grid.innerHTML = `
                    <div class="loading">
                        <div class="spinner"></div>
                        <p>Waiting for data...</p>
                    </div>
                `;
            } else {
                this.elements.grid.innerHTML = `
                    <div class="empty">
                        <p>No cryptocurrencies found for "${this.searchTerm}"</p>
                    </div>
                `;
            }
            return;
        }

        this.elements.grid.innerHTML = this.filteredCryptos
            .map(crypto => this.createCard(crypto))
            .join('');
    }

    createCard(crypto) {
        const priceClass = crypto.priceChange > 0 ? 'positive' : 
                          crypto.priceChange < 0 ? 'negative' : 'neutral';
        
        const changeText = crypto.priceChange !== 0 ? 
            `${crypto.priceChange > 0 ? '+' : ''}$${Math.abs(crypto.priceChange).toFixed(4)}` : 
            'No change';

        return `
            <div class="crypto-card ${crypto.updated ? 'updated' : ''}">
                <div class="crypto-header">
                    <div class="crypto-info">
                        <img src="${crypto.image}" alt="${crypto.name}" class="crypto-icon" 
                             onerror="this.src='data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHZpZXdCb3g9IjAgMCA0MCA0MCIgZmlsbD0ibm9uZSI+PGNpcmNsZSBjeD0iMjAiIGN5PSIyMCIgcj0iMjAiIGZpbGw9IiNmMGYwZjAiLz48dGV4dCB4PSI1MCUiIHk9IjUwJSIgZm9udC1mYW1pbHk9IkFyaWFsIiBmb250LXNpemU9IjEyIiBmaWxsPSIjOTk5IiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBkb21pbmFudC1iYXNlbGluZT0ibWlkZGxlIj4/PC90ZXh0Pjwvc3ZnPg=='">
                        <div class="crypto-details">
                            <div class="crypto-name" title="${crypto.name}">${crypto.name}</div>
                            <div class="crypto-symbol" title="${crypto.symbol}">${crypto.symbol}</div>
                        </div>
                    </div>
                    <div class="crypto-price">
                        <div class="price-value">$${this.formatPrice(crypto.current_price)}</div>
                        <div class="price-change ${priceClass}">${changeText}</div>
                    </div>
                </div>
            </div>
        `;
    }

    formatPrice(price) {
        return price >= 1 ? 
            price.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) :
            price.toFixed(6);
    }

    updateStatus(text, className = '') {
        this.elements.status.textContent = text;
        this.elements.status.className = `status ${className}`;
    }

    updateCount() {
        this.elements.count.textContent = this.cryptos.size;
    }

    setupEvents() {
        // Search
        this.elements.search.addEventListener('input', (e) => {
            this.searchTerm = e.target.value.toLowerCase().trim();
            this.filterAndRender();
        });

        // Sort
        this.elements.sort.addEventListener('change', (e) => {
            this.sortBy = e.target.value;
            this.filterAndRender();
        });

        // Keyboard shortcuts
        document.addEventListener('keydown', (e) => {
            if ((e.ctrlKey || e.metaKey) && e.key === 'f') {
                e.preventDefault();
                this.elements.search.focus();
            }
            if (e.key === 'Escape' && document.activeElement === this.elements.search) {
                this.elements.search.value = '';
                this.searchTerm = '';
                this.filterAndRender();
            }
        });
    }
}

// Start the app
document.addEventListener('DOMContentLoaded', () => {
    console.log('Starting CryptoStream...');
    window.cryptoApp = new CryptoApp();
});