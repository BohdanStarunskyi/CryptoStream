# CryptoStream Web Application

A simple, beautiful real-time cryptocurrency price tracker that connects to your gateway via WebSockets.

## Features

- 🚀 **Real-time Updates**: Live crypto price updates via WebSocket
- 💎 **Clean Design**: Simple, responsive UI with smooth gradients
- 🔍 **Search**: Quick search by name, symbol, or ID
- 📊 **Sort Options**: Sort by name, price, or symbol
- 📱 **Mobile Friendly**: Fully responsive design
- ⚡ **Connection Status**: Live connection indicator
- 📈 **Price Changes**: Visual price change indicators
- 🔄 **Auto Reconnect**: Automatic reconnection on disconnect

## File Structure

```
web/
├── index.html      # Main HTML file with semantic structure
├── styles.css      # Beautiful CSS with animations and responsive design
├── app.js          # JavaScript WebSocket client and application logic
└── README.md       # This documentation file
```

## How to Use

### Prerequisites

Make sure you have the backend services running:

1. **Fetcher Service**: Fetches crypto data and sends it to the gateway
   ```bash
   cd backend/fetcher
   go run .
   ```

2. **Gateway Service**: Receives data from fetcher and serves WebSocket connections
   ```bash
   cd backend/gateway
   go run .
   ```

### Running the Web Application

1. **Simple HTTP Server**: Open the web application using any HTTP server
   
   **Option 1: Python 3 (recommended)**
   ```bash
   cd web
   python3 -m http.server 3000
   # Then visit http://localhost:3000
   ```
   
   **Option 2: Node.js (if installed)**
   ```bash
   cd web
   npx serve -p 3000
   # Then visit http://localhost:3000
   ```
   
   **Option 3: VS Code Live Server**
   - Install the "Live Server" extension in VS Code
   - Right-click on `index.html` and select "Open with Live Server"

2. **Direct File Access**: You can also open `index.html` directly in your browser, but some browsers may block WebSocket connections from file:// URLs.

### WebSocket Connection

The application automatically connects to the gateway WebSocket endpoint at `ws://localhost:8080/ws`.

If you need to change the WebSocket URL, edit the connection string in `app.js`:

```javascript
this.ws = new WebSocket('ws://localhost:8080/ws');
```

## User Interface

### Header
- **Logo**: CryptoStream branding with Bitcoin icon
- **Connection Status**: Live indicator showing connection state (Connected/Disconnected/Reconnecting)

### Statistics Overview
- **Total Cryptocurrencies**: Count of available cryptocurrencies
- **Last Update**: Timestamp of most recent data update
- **Connection Time**: Duration of current WebSocket connection

### Search and Controls
- **Search Bar**: Filter cryptocurrencies by name, symbol, or ID
- **Sort Dropdown**: Sort by name (A-Z), price (high to low), or symbol (A-Z)

### Crypto Grid
- **Real-time Cards**: Each cryptocurrency displayed in an animated card
- **Price Information**: Current price with change indicators
- **Visual Feedback**: Cards pulse when updated with new data
- **Responsive Layout**: Automatically adjusts to screen size

## Features in Detail

### Real-time Updates
- Receives live data from the gateway via WebSocket
- Shows price changes with green (increase) and red (decrease) indicators
- Animates cards when new data arrives

### Search Functionality
- Type to filter cryptocurrencies instantly
- Searches across name, symbol, and ID fields
- Clear results with Escape key
- Keyboard shortcut: Ctrl/Cmd + F to focus search

### Connection Management
- Automatic reconnection with exponential backoff
- Visual connection status indicator
- Connection duration timer
- Graceful handling of network issues

### Responsive Design
- Mobile-first approach
- Breakpoints at 768px and 480px
- Touch-friendly interface
- Optimized for all screen sizes

## Browser Support

- Modern browsers with WebSocket support
- Chrome 76+
- Firefox 72+
- Safari 14+
- Edge 79+

## Troubleshooting

### WebSocket Connection Issues

1. **Check Backend Services**:
   - Ensure gateway service is running on port 8080
   - Ensure fetcher service is running and connected to gateway

2. **Browser Console**:
   - Open browser developer tools (F12)
   - Check console for error messages
   - Use `window.cryptoDebug` for debugging information

3. **Network Issues**:
   - Check if localhost:8080 is accessible
   - Verify no firewall blocking the connection
   - Try refreshing the page

### No Data Appearing

1. **Check Fetcher Service**: Make sure the fetcher is successfully fetching data from the API
2. **Gateway Logs**: Check gateway console for incoming data
3. **WebSocket Messages**: Check browser network tab for WebSocket messages

### Performance Issues

1. **Too Many Cryptocurrencies**: The interface handles large datasets well, but you can use search to filter
2. **Memory Usage**: Refresh the page occasionally for long-running sessions
3. **Browser Resources**: Close other tabs if experiencing slowdown

## Customization

### Styling
- Edit `styles.css` to customize colors, animations, and layout
- CSS custom properties make theming easy
- Supports CSS Grid and Flexbox for layout modifications

### Functionality
- Modify `app.js` to add new features
- WebSocket message handling in `handleCryptoUpdates()`
- Add new filtering or sorting options
- Customize reconnection behavior

### Data Display
- Modify `createCryptoCard()` to change card layout
- Add new statistics in the stats overview
- Customize price formatting in `formatPrice()`

## Development

For development and debugging, the application exposes a global `window.cryptoDebug` object with useful methods:

```javascript
// Check WebSocket connection status
window.cryptoDebug.getConnectionStatus()

// Get all crypto data
window.cryptoDebug.getCryptoData()

// Force reconnection
window.cryptoDebug.reconnect()

// Get currently filtered/displayed data
window.cryptoDebug.getFilteredData()
```

## Security Notes

- The WebSocket connection allows all origins for development
- For production, configure proper CORS settings in the gateway
- Consider adding authentication for production deployments
- Use HTTPS/WSS for production environments