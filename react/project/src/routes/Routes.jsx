import { Route, Routes } from 'react-router-dom';
import Auth from '/src/components/Auth';
import Logout from '/src/components/Logout';
import MyPage from '/src/components/MyPage'; // インポート

const RoutesConfig = () => {
    return (
        <>
            <Routes>
                {/* auth画面 */}
                <Route path="/auth/register" element={<Auth />} />
                <Route path="/auth/login" element={<Auth />} />
                <Route path="/auth/activate" element={<Auth />} />
                <Route path="/auth/ping" element={<Auth />} />
                {/* ログアウト */}
                <Route path="/logout" element={<Logout />} />
                {/* マイページ */}
                <Route path="/mypage" element={<MyPage />} />

                <Route path="/" element={<MyPage />} />
            </Routes>
        </>
    );
};

export default RoutesConfig;
