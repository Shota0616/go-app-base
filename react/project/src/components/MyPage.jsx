// src/components/MyPage.jsx
import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { useTranslation } from 'react-i18next'; // 追加

// マイページ
const MyPage = () => {
    const { t } = useTranslation(); // 追加
    // apiを通じて取得したユーザ情報を格納
    const [user, setUser] = useState(null);
    // ページ遷移の関数をnavigateとする
    const navigate = useNavigate();

    useEffect(() => {
        const checkAuth = async () => {
            // localstorageに保存されたtokenを取得
            const token = localStorage.getItem('token');

            // もしtokenが存在しなかったらloginページに遷移
            if (!token) {
                navigate('/auth/login');
                return;
            }

            try {
                // ヘッダにtokenを載せてgetuserのapiを叩く
                const response = await axios.get(`${import.meta.env.VITE_APP_API_URL}/api/getuser`, {
                    headers: {
                        'Authorization': token,
                    },
                });
                // user変数に取得したユーザ情報を格納
                setUser(response.data);
            } catch (error) {
                // トークンが無効だったということなのでlocalstorageから削除
                localStorage.removeItem('token');
                // storageイベント発火
                window.dispatchEvent(new Event("storage"));
                navigate('/auth/login');
                console.log("/mypage : error")
            }
        };

        checkAuth();
    }, [navigate]);

    return (
        <div>
            <h2>{t('mypage_title')}</h2>
            {user ? (
                // userが存在するときユーザ情報を表示
                <div>
                    <p>{t('username')}: {user.username}</p>
                    <p>{t('email')}: {user.email}</p>
                </div>
            ) : (
                // ユーザ情報が取得できていないときはloadingを表示
                <p>{t('loading')}</p>
            )}
        </div>
    );
};

export default MyPage;
