import React from 'react';
// import axios from 'axios';
import { useState, useEffect } from 'react';


const Logout = () => {
    const [logoutMessage, setLogoutMessage] = useState(''); // 初期値を空文字に

    const handleLogout = async () => {
        try {
            // ローカルストレージにあるtokenというキーのものを削除
            localStorage.removeItem('token');
            localStorage.removeItem('refreshtoken');
            // storageイベント発火
            window.dispatchEvent(new Event("storage"));
            // トークンの削除はブラウザ側だけで実施可能なのでapiは使用しないでOK。トークン管理をサーバ側でも実施するときは削除も必要
            // await axios.post(`${import.meta.env.VITE_APP_API_URL}/api/logout`, {}, {
            //     headers: {
            //     'Authorization': token,
            //     },
            // });
            setLogoutMessage('Logout successful!');
        } catch (error) {
            setLogoutMessage('Logout failed!');
        }
    };

    return (
        <>
            <h2>Logout</h2>
            {logoutMessage ? (
                // logoutMessageが空じゃないときはメッセージ表示
                <p>{logoutMessage}</p>
            ) : (
                <button onClick={handleLogout}>Logout</button>
            )}
        </>
    );
};

export default Logout;
