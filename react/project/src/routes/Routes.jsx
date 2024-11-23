import { Route, Routes } from 'react-router-dom';
import Auth from '/src/components/Auth';
import Logout from '/src/components/Logout';
import MyPage from '/src/components/MyPage'; // インポート
import Home from '/src/components/Home'; // インポート

const RoutesConfig = () => {
    return (
        <>
            <Routes>
                {/* auth画面 */}
                <Route path="/auth/register" element={<Auth />} />
                <Route path="/auth/login" element={<Auth />} />
                <Route path="/auth/verify" element={<Auth />} />
                <Route path="/auth/request-password-reset" element={<Auth />} />
                <Route path="/auth/reset-password" element={<Auth />} />
                {/* ログアウト */}
                <Route path="/logout" element={<Logout />} />
                {/* マイページ */}
                <Route path="/mypage" element={<MyPage />} />

                <Route path="/" element={<Home />} />
            </Routes>
        </>
    );
};

export default RoutesConfig;
