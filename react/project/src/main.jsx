import React from 'react'
import ReactDOM from 'react-dom/client'
import App from '/src/App.jsx'
import '/src/index.css'
import { BrowserRouter } from 'react-router-dom';


// index.htmlのid=rootに一致するもの
ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <BrowserRouter>
      {/* App.jsxが呼び出される */}
      <App />
    </BrowserRouter>
  </React.StrictMode>,
)
