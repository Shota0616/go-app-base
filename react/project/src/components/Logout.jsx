import React from 'react';
import { useState, useEffect } from 'react';
import { Button } from '@mui/material';
import { useTranslation } from 'react-i18next';

const Logout = () => {
    const { t } = useTranslation();
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
            setLogoutMessage(t('logout_successful'));
        } catch (error) {
            setLogoutMessage(t('logout_failed'));
        }
    };

    return (
        <>
            {logoutMessage ? (
                // logoutMessageが空じゃないときはメッセージ表示
                <p>{logoutMessage}</p>
            ) : (
                <Button onClick={handleLogout} type="submit" variant="contained" color="primary" sx={{ mt: 2, width: '150px', mx: 'auto', display: 'block', height: 50, borderRadius: 3 }}>
                    {t('logout')}
                </Button>
            )}
        </>
    );
};

export default Logout;
