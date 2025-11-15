import { Route, Routes } from 'react-router-dom';
import Logout from '/src/components/Logout';
import MyPage from '/src/components/MyPage';
import Home from '/src/components/Home';
import SignIn from '/src/components/SignIn';
import SignUp from '/src/components/SignUp';
import Verify from '/src/components/Verify';
import RequestPasswordReset from '/src/components/RequestPasswordReset';
import ResetPassword from '/src/components/ResetPassword';
import Settings from '/src/components/Settings'; // Import Settings component

const RoutesConfig = () => {
    return (
        <Routes>
            {/* auth画面 */}
            <Route path="/auth/register" element={<SignUp />} />
            <Route path="/auth/login" element={<SignIn />} />
            <Route path="/auth/verify" element={<Verify />} />
            <Route path="/auth/request-password-reset" element={<RequestPasswordReset />} />
            <Route path="/auth/reset-password" element={<ResetPassword />} />
            
            {/* ログアウト */}
            <Route path="/logout" element={<Logout />} />
            
            {/* マイページ */}
            <Route path="/mypage" element={<MyPage />} />

            {/* 設定画面 */}
            <Route path="/settings" element={<Settings />} />

            <Route path="/" element={<Home />} />
        </Routes>
    );
};

export default RoutesConfig;
