import React, { useState, useEffect } from 'react'; // Reactとフックをインポート
import axios from 'axios'; // HTTPリクエスト用のaxiosをインポート
import { useLocation, useNavigate } from 'react-router-dom'; // ルーティングに必要なフックをインポート

// 認証状態に関するステートを管理するカスタムフック
const useAuthState = () => {
    const [username, setUsername] = useState(''); // ユーザー名のステート
    const [password, setPassword] = useState(''); // パスワードのステート
    const [email, setEmail] = useState(''); // メールアドレスのステート
    const [message, setMessage] = useState(''); // メッセージのステート
    const [pingResponse, setPingResponse] = useState(''); // Pingレスポンスのステート
    const [errorMessage, setErrorMessage] = useState(''); // エラーメッセージのステート

    return {
        username,
        setUsername,
        password,
        setPassword,
        email,
        setEmail,
        message,
        setMessage,
        pingResponse,
        setPingResponse,
        errorMessage,
        setErrorMessage,
    };
};

// /auth/activateパスに基づいて処理を行うカスタムフック
const useActivateEffect = (navigate, handleActivate) => {
    const location = useLocation(); // 現在のURLのパスを取得

    useEffect(() => {
        if (location.pathname === '/auth/activate') {
            const params = new URLSearchParams(window.location.search);
            const token = params.get('token'); // URLからトークンを取得
            const email = params.get('email'); // URLからメールアドレスを取得

            if (token && email) {
                handleActivate(token, email); // トークンとメールが存在する場合、アクティベーション処理を実行
            } else {
                handleInvalidActivationLink(navigate); // トークンやメールがない場合、無効なリンク処理を実行
            }
        }
    }, [location, navigate, handleActivate]); // locationが変わるたびに実行
};

// 無効なアクティベーションリンクの場合の処理
const handleInvalidActivationLink = (navigate) => {
    setTimeout(() => {
        navigate('/'); // 2秒後にホームページにリダイレクト
    }, 2000);
};

// ユーザー登録の処理
const handleRegister = async (username, password, email, setMessage) => {
    try {
        const response = await axios.post(`${import.meta.env.VITE_APP_API_URL}/api/register`, {
            username,
            password,
            email,
        });
        setMessage(response.data.message); // 成功メッセージを設定
    } catch (error) {
        handleError(error, setMessage); // エラーが発生した場合はエラーハンドリング
    }
};

// アクティベーションの処理
const handleActivate = async (token, email, setMessage, navigate) => {
    try {
        const response = await axios.post(`${import.meta.env.VITE_APP_API_URL}/api/activate`, {
            token,
            email,
        });
        setMessage(response.data.message); // 成功メッセージを設定
        handleInvalidActivationLink(navigate); // 無効なリンク処理を実行
    } catch (error) {
        handleError(error, setMessage); // エラーが発生した場合はエラーハンドリング
        handleInvalidActivationLink(navigate); // 無効なリンク処理を実行
    }
};

// ログインの処理
const handleLogin = async (email, password, setMessage) => {
    try {
        // 既存のトークンがlocalStorageに存在する場合、そのトークンの有効性を確認
        const existingToken = localStorage.getItem('token');
        if (existingToken) {
            const isTokenValid = await validateToken(existingToken);
            if (!isTokenValid) {
                localStorage.removeItem('token'); // 無効なトークンを削除
                window.dispatchEvent(new Event("storage")); // storageイベントを発火
            }
        }
        // emailとpasswordでログインし、トークンとリフレッシュトークンを取得
        const response = await axios.post(`${import.meta.env.VITE_APP_API_URL}/api/login`, {
            email,
            password,
        });
        localStorage.setItem('token', response.data.token); // アクセストークンをlocalStorageに保存
        localStorage.setItem('refrestoken', response.data.refreshtoken); // リフレッシュトークンも保存
        window.dispatchEvent(new Event("storage")); // storageイベントを発火
        setMessage('Login successful!'); // 成功メッセージを設定
    } catch (error) {
        handleError(error, setMessage); // エラーが発生した場合はエラーハンドリング
    }
};

// Ping APIからデータを取得する処理
const fetchPing = async (setPingResponse, setErrorMessage) => {
    try {
        const response = await axios.get(`${import.meta.env.VITE_APP_API_URL}/api/ping`);
        setPingResponse(response.data); // Pingレスポンスを設定
    } catch (error) {
        setErrorMessage('Error fetching ping response'); // エラーメッセージを設定
    }
};

// エラーハンドリングの処理
const handleError = (error, setMessage) => {
    if (error.response && error.response.data) {
        setMessage(error.response.data.error); // サーバーからのエラーメッセージを設定
    } else {
        setMessage("An unexpected error occurred."); // 一般的なエラーメッセージを設定
    }
};

// 登録画面のコンテンツをレンダリング
const renderRegister = (username, setUsername, password, setPassword, email, setEmail, handleRegister, setMessage) => (
    <>
        <h2>Register</h2>
        <input
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)} // ユーザー名の変更
        />
        <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)} // メールアドレスの変更
        />
        <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)} // パスワードの変更
        />
        <button onClick={() => handleRegister(username, password, email, setMessage)}>Register</button> {/* 登録処理の実行 */}
    </>
);

// ログイン画面のコンテンツをレンダリング
const renderLogin = (email, setEmail, password, setPassword, handleLogin, setMessage) => (
    <>
        <h2>Login</h2>
        <input
            type="text"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)} // メールアドレスの変更
        />
        <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)} // パスワードの変更
        />
        <button onClick={() => handleLogin(email, password, setMessage)}>Login</button> {/* ログイン処理の実行 */}
    </>
);

// Pingテスト画面のコンテンツをレンダリング
const renderPing = (fetchPing, pingResponse, errorMessage, setPingResponse, setErrorMessage) => (
    <>
        <h2>Ping</h2>
        <button onClick={() => fetchPing(setPingResponse, setErrorMessage)}>Ping</button> {/* Pingテスト実行 */}
        {errorMessage ? (
            <p style={{ color: 'red' }}>{errorMessage}</p> // エラーメッセージ表示
        ) : (
            pingResponse && <p>Message: {pingResponse.message}</p> // Pingレスポンス表示
        )}
    </>
);

// // リフレッシュトークンを使って新しいアクセストークンを取得する関数
// const refreshToken = async () => {
//     try {
//         // localStorageからリフレッシュトークンを取得
//         const refreshToken = localStorage.getItem('refreshToken');
//         if (!refreshToken) {
//             throw new Error('No refresh token found'); // リフレッシュトークンがない場合はエラー
//         }
//         // リフレッシュトークンを使って新しいアクセストークンを取得するAPIエンドポイントにリクエスト
//         const response = await axios.post(`${import.meta.env.VITE_APP_API_URL}/api/refresh`, {
//             refreshToken,
//         });
//         // 取得した新しいアクセストークンとリフレッシュトークンをlocalStorageに保存
//         localStorage.setItem('token', response.data.token);
//         localStorage.setItem('refreshToken', response.data.refreshToken);
//         window.dispatchEvent(new Event("storage")); // storageイベントを発火して、他のタブに通知
//         return response.data.token; // 新しいアクセストークンを返す
//     } catch (error) {
//         console.error('Failed to refresh token', error);
//         // リフレッシュトークンが無効な場合、トークンを削除
//         localStorage.removeItem('token');
//         localStorage.removeItem('refreshToken');
//         window.dispatchEvent(new Event("storage")); // storageイベントを発火して、他のタブに通知
//         throw error; // エラーを再スローして処理を停止
//     }
// };

// // Axiosのリクエストインターセプター
// axios.interceptors.request.use(async (config) => {
//     // localStorageからアクセストークンを取得
//     let token = localStorage.getItem('token');
//     if (token) {
//         try {
//             // トークンの有効性を確認
//             const isTokenValid = await validateToken(token);
//             if (!isTokenValid) {
//                 // トークンが無効な場合、リフレッシュトークンを使って新しいアクセストークンを取得
//                 token = await refreshToken();
//             }
//         } catch (error) {
//             console.error('Token refresh failed', error);
//         }
//     }
//     // 有効なアクセストークンがあれば、リクエストヘッダーに設定
//     if (token) {
//         config.headers.Authorization = `Bearer ${token}`;
//     }
//     return config; // リクエスト設定を返す
// }, (error) => {
//     return Promise.reject(error); // リクエストエラーをそのまま返す
// });


// メインのAuthコンポーネント
const Auth = () => {
    const location = useLocation(); // 現在のURLのパスを取得
    const navigate = useNavigate(); // ページ遷移を行うためのnavigateフックを取得
    const {
        username,
        setUsername,
        password,
        setPassword,
        email,
        setEmail,
        message,
        setMessage,
        pingResponse,
        setPingResponse,
        errorMessage,
        setErrorMessage,
    } = useAuthState(); // 認証状態に関するステートを取得

    // useActivateEffectを使用して、/auth/activateパスに基づくアクティベーション処理を実行
    useActivateEffect(navigate, (token, email) => handleActivate(token, email, setMessage, navigate));

    // 現在のパスに基づいて適切なコンテンツをレンダリング
    const renderContent = () => {
        switch (location.pathname) {
            case '/auth/register':
                return renderRegister(username, setUsername, password, setPassword, email, setEmail, handleRegister, setMessage);
            case '/auth/login':
                return renderLogin(email, setEmail, password, setPassword, handleLogin, setMessage);
            case '/auth/ping':
                console.log("/auth/ping : case")
                return renderPing(fetchPing, pingResponse, errorMessage, setPingResponse, setErrorMessage);
            case '/auth/activate':
                return <h2>Activating your account...</h2>;
            default:
                return <h2>Page not found</h2>;
        }
    };

    return (
        <div>
            {renderContent()} {/* 現在のパスに基づいたコンテンツをレンダリング */}
            {message && <p>{message}</p>} {/* メッセージがある場合、表示 */}
        </div>
    );
};

export default Auth; // Authコンポーネントをエクスポート
