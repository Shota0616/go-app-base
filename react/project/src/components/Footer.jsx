import React from 'react';
import reactLogo from '/src/assets/react.svg'
import viteLogo from '/vite.svg'
import '/src/App.css'

const Footer = () => (
    <div>
        <a href="https://vitejs.dev" target="_blank">
        <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
        <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
        <footer className="bg-dark text-white p-3">
            <p>&copy; 2024 Shota Isoda. All rights reserved.</p>
        </footer>
    </div>
);

export default Footer;
