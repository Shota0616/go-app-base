import React from 'react'
import ReactDOM from 'react-dom/client'
import App from '/src/App.jsx'
import '/src/index.css'
import { BrowserRouter } from 'react-router-dom';
import { UserProvider } from './context/UserContext';

// index.htmlのid=rootに一致するもの
ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <BrowserRouter>
      <UserProvider>
        {/* App.jsxが呼び出される */}
        <App />
      </UserProvider>
    </BrowserRouter>
  </React.StrictMode>,
)
