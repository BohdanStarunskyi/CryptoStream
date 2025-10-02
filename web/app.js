class CryptoApp {
    constructor() {
        this.ws = null;
        this.cryptos = new Map();
        
        this.elements = {
            status: document.getElementById('status'),
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
            const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsHost = window.location.hostname === 'localhost' ? 'localhost:8080' : window.location.host;
            const wsUrl = `${wsProtocol}//${wsHost}/ws`;
            this.ws = new WebSocket(wsUrl);
            
            this.ws.onopen = () => {
                this.updateStatus('Connected', 'connected');
            };

            this.ws.onmessage = (event) => {
                const updates = JSON.parse(event.data);
                this.handleUpdates(updates);
            };

            this.ws.onclose = () => {
                this.updateStatus('Disconnected', 'disconnected');
                setTimeout(() => this.connect(), 3000);
            };

            this.ws.onerror = (error) => {
                this.updateStatus('Connection Error', 'disconnected');
            };

        } catch (error) {
            this.updateStatus('Failed to Connect', 'disconnected');
            setTimeout(() => this.connect(), 3000);
        }
    }

    handleUpdates(updates) {
        updates.forEach(crypto => {
            const prev = this.cryptos.get(crypto.id);
            crypto.priceChange24h = crypto.price_change_24h || 0;
            crypto.updated = true;
            crypto.priceDirection = prev && prev.current_price !== crypto.current_price 
                ? (crypto.current_price > prev.current_price ? 'up' : 'down') 
                : null;
            this.cryptos.set(crypto.id, crypto);
        });

        this.updateCount();
        this.render();
        
        setTimeout(() => {
            updates.forEach(crypto => {
                const stored = this.cryptos.get(crypto.id);
                if (stored) {
                    stored.updated = false;
                    stored.priceDirection = null;
                }
            });
            this.render();
        }, 1000);
    }



    render() {
        if (this.cryptos.size === 0) {
            this.elements.grid.innerHTML = `
                <div class="loading">
                    <div class="spinner"></div>
                    <p>Waiting for data...</p>
                </div>
            `;
            return;
        }

        const cryptoArray = Array.from(this.cryptos.values());
        cryptoArray.sort((a, b) => a.name.localeCompare(b.name));
        
        this.elements.grid.innerHTML = cryptoArray
            .map(crypto => this.createCard(crypto))
            .join('');
    }

    createCard(crypto) {
        const priceClass = crypto.priceChange24h > 0 ? 'positive' : 
                          crypto.priceChange24h < 0 ? 'negative' : 'neutral';
        
        const changeText = crypto.priceChange24h !== 0 ? 
            `${crypto.priceChange24h > 0 ? '+' : ''}${crypto.priceChange24h.toFixed(2)}%` : 
            '0.00%';

        const splashClass = crypto.priceDirection ? `price-splash-${crypto.priceDirection}` : '';

        return `
            <div class="crypto-card ${crypto.updated ? 'updated' : ''} ${splashClass}">
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
    }
}

document.addEventListener('DOMContentLoaded', () => {
    window.cryptoApp = new CryptoApp();
});