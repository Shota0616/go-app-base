import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useLocation, useNavigate, Link } from 'react-router-dom';
import { Box, Button, TextField, Typography, Paper, IconButton, InputAdornment } from '@mui/material';
import { useTranslation } from 'react-i18next';
import { Visibility, VisibilityOff } from '@mui/icons-material';

// 認証状態を管理するカスタムフック
const useAuthState = () => {
    // ユーザー名の状態を管理
    const [username, setUsername] = useState('');
    // パスワードの状態を管理
    const [password, setPassword] = useState('');
    // メールアドレスの状態を管理
    const [email, setEmail] = useState('');
    // メッセージの状態を管理（エラーメッセージや成功メッセージ）
    const [message, setMessage] = useState('');
    // メッセージの種類を管理（エラーか成功か）
    const [messageType, setMessageType] = useState('');
    // 認証コードの状態を管理
    const [verificationCode, setVerificationCode] = useState('');
    // 新しいパスワードの状態を管理
    const [newPassword, setNewPassword] = useState('');
    // トークンの状態を管理
    const [token, setToken] = useState('');

    // 状態とその更新関数を返す
    return {
        username,
        setUsername,
        password,
        setPassword,
        email,
        setEmail,
        message,
        setMessage,
        messageType,
        setMessageType,
        verificationCode,
        setVerificationCode,
        newPassword,
        setNewPassword,
        token,
        setToken,
    };
};


const Auth = ({ open, handleClose }) => {
    const { username, setUsername, password, setPassword, email, setEmail, message, setMessage, messageType, setMessageType, verificationCode, setVerificationCode, newPassword, setNewPassword, token, setToken } = useAuthState();
    const location = useLocation();
    const navigate = useNavigate();
    const { t } = useTranslation();
    const [showPassword, setShowPassword] = useState(false);
    const [emailError, setEmailError] = useState(false);
    const [passwordError, setPasswordError] = useState(false);
    const [usernameError, setUsernameError] = useState(false);
    const [verificationCodeError, setVerificationCodeError] = useState(false);
    const [newPasswordError, setNewPasswordError] = useState(false);

    const handleClickShowPassword = () => setShowPassword(!showPassword);
    const handleMouseDownPassword = (event) => event.preventDefault();

    // ページ遷移時にメッセージをリセット
    useEffect(() => {
        // location.stateからmessageを取得し、setMessageで状態を更新。messageがない場合は空文字を設定。
        setMessage(location.state?.message || '');
        // location.stateからmessageTypeを取得し、setMessageTypeで状態を更新。messageTypeがない場合は空文字を設定。
        setMessageType(location.state?.messageType || ''); // メッセージの種類をリセット
    // location.pathnameまたはlocation.stateが変更されるたびにこのuseEffectが実行される。
    }, [location.pathname, location.state]);

    useEffect(() => {
        // location.searchからクエリパラメータを取得し、URLSearchParamsオブジェクトを作成。
        const params = new URLSearchParams(location.search);
        // クエリパラメータから'token'を取得。
        const token = params.get('token');
        // tokenが存在する場合、setTokenで状態を更新。
        if (token) {
            setToken(token);
        }
    // location.searchが変更されるたびにこのuseEffectが実行される。
    }, [location.search]);

    // useEffect(() => {
    //     // location.searchからクエリパラメータを取得し、URLSearchParamsオブジェクトを作成。
    //     const params = new URLSearchParams(location.search);
    //     // クエリパラメータから'email'を取得。
    //     const emailParam = params.get('email');
    //     // emailParamが存在する場合、setEmailで状態を更新。
    //     if (emailParam) {
    //         setEmail(emailParam);
    //     }
    // // location.searchが変更されるたびにこのuseEffectが実行される。
    // }, [location.search]);

    // ページ遷移時にエラーメッセージをリセット
    useEffect(() => {
        setEmailError(false);
        setPasswordError(false);
        setUsernameError(false);
        setVerificationCodeError(false);
        setNewPasswordError(false);
    }, [location.pathname]);

    //////////////////////////////////////////
    ///////////// 各種処理の関数 ///////////////
    //////////////////////////////////////////

    // 新規登録の処理
    const handleRegister = async (username, password, email, setMessage, setMessageType, navigate) => {
        // ユーザー名、パスワード、メールアドレスのいずれかが空の場合、エラーメッセージを設定して処理を終了
        if (!username || !password || !email) {
            setUsernameError(!username);
            setEmailError(!email);
            setPasswordError(!password);
            return;
        }

        // ユーザー名、パスワード、メールアドレスをサーバーに送信して新規登録
        try {
            const response = await axios.post(`${import.meta.env.VITE_APP_API_URL}/api/register`, {
                username,
                password,
                email,
            });
            setMessage(response.data.message);
            setMessageType('success'); // 成功メッセージの種類を設定
            navigate('/auth/verify'); // 登録成功後にメールの検証画面にリダイレクト
        } catch (error) {
            // エラーレスポンスの処理
            if (error.response && error.response.data && error.response.data.error) {
                setMessage(error.response.data.error);
            } else {
                setMessage(t('registration_failed'));
            }
            setMessageType('error'); // エラーメッセージの種類を設定
        }
    };

    // ログインの処理
    const handleLogin = async (email, password, setMessage, setMessageType, navigate) => {
        // メールアドレス、パスワードのいずれかが空の場合、エラーメッセージを設定して処理を終了
        if (!email || !password) {
            setEmailError(!email);
            setPasswordError(!password);
            return;
        }

        try {
            const response = await axios.post(`${import.meta.env.VITE_APP_API_URL}/api/login`, {
                email,
                password,
            });

            localStorage.setItem('token', response.data.token); // アクセストークンをlocalStorageに保存
            localStorage.setItem('refrestoken', response.data.refreshtoken); // リフレッシュトークンも保存
            window.dispatchEvent(new Event("storage")); // storageイベントを発火
            setMessage(t('login_successful')); // 成功メッセージを設定
            setMessageType('success'); // 成功メッセージの種類を設定

        } catch (error) {
            // エラーレスポンスの処理

            if (error.response.status === 303) {
                navigate("/auth/verify", { state: { message: error.response.data.error, messageType: 'error' } });
            } else {
                setMessage(t('login_failed'));
            }
            setMessageType('error'); // エラーメッセージの種類を設定
        }
    };

    // 認証コードの処理
    const handleVerify = async (email, verificationCode, setMessage, setMessageType, navigate) => {
        // 認証コードが空の場合、エラーメッセージを設定して処理を終了
        if (!verificationCode) {
            setVerificationCodeError(!verificationCode);
            return;
        }

        try {
            const response = await axios.post(`${import.meta.env.VITE_APP_API_URL}/api/verify`, {
                email,
                verificationCode,
            });
            setMessage(response.data.message);
            setMessageType('success'); // 成功メッセージの種類を設定
            navigate('/auth/login'); // 認証成功後にログイン画面にリダイレクト
        } catch (error) {
            if (error.response && error.response.data && error.response.data.error) {
                setMessage(error.response.data.error);
            } else {
                setMessage(t('verification_failed'));
            }
            setMessageType('error'); // エラーメッセージの種類を設定
        }
    };

    // 認証コードの再送信の処理
    const handleResendVerificationCode = async (email, setMessage, setMessageType) => {
        // メールアドレスが空の場合、エラーメッセージを設定して処理を終了
        try {
            const response = await axios.post(`${import.meta.env.VITE_APP_API_URL}/api/resend-verification-code`, {
                email,
            });
            setMessage(response.data.message);
            setMessageType('success'); // 成功メッセージの種類を設定
        } catch (error) {
            if (error.response && error.response.data && error.response.data.error) {
                setMessage(error.response.data.error);
            }else {
                setMessage(t('resend_verification_code_failed'));
            }
            setMessageType('error');
        }
    };

    const handleRequestPasswordReset = async (email, setMessage, setMessageType) => {
        // メールアドレスが空の場合、エラーメッセージを設定して処理を終了
        if (!email) {
            setEmailError(!email);
            return;
        }
        try {
            const response = await axios.post(`${import.meta.env.VITE_APP_API_URL}/api/request-password-reset`, {
                email,
            });
            setMessage(response.data.message);
            setMessageType('success');
        } catch (error) {
            if (error.response && error.response.data && error.response.data.error) {
                setMessage(error.response.data.error);
            } else {
                setMessage(t('request_password_reset_failed'));
            }
            setMessageType('error');
        }
    };

    const handleResetPassword = async (token, newPassword, setMessage, setMessageType, navigate) => {
        // 新しいパスワードが空の場合、エラーメッセージを設定して処理を終了
        if (!newPassword) {
            setNewPasswordError(!newPassword);
            return;
        }
        try {
            const response = await axios.post(`${import.meta.env.VITE_APP_API_URL}/api/reset-password`, {
                token,
                newPassword,
            });
            setMessage(response.data.message);
            setMessageType('success');
            navigate('/auth/login');
        } catch (error) {
            if (error.response && error.response.data && error.response.data.error) {
                setMessage(error.response.data.error);
            } else {
                setMessage(t('reset_password_failed'));
            }
            setMessageType('error');
        }
    };



    //////////////////////////////////////////
    ///////////// 各種render処理 ///////////////
    //////////////////////////////////////////



    const renderRegister = () => {
        return (
            // <Box component="form" onSubmit={(e) => { e.preventDefault(); handleRegister(username, password, email, setMessage, setMessageType, navigate) }} sx={{ mt: 2 }}>
            <Box component="form" onSubmit={(e) => { e.preventDefault(); handleRegister(username, password, email, setMessage, setMessageType, navigate) }} sx={{ mt: 2 }}>

                <TextField
                    label={t('username')}
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    fullWidth
                    margin="normal"
                    error={usernameError}
                    helperText={usernameError ? `${t('username')}${t('input_required')}` : ''}
                    InputProps={{
                        style: {
                            backgroundColor: 'rgba(255, 255, 255, 0.1)',
                            color: 'white',
                        },
                    }}
                    InputLabelProps={{
                        style: { color: 'white' },
                    }}
                />
                <TextField
                    label={t('email')}
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    fullWidth
                    margin="normal"
                    error={emailError}
                    helperText={emailError ? `${t('email')}${t('input_required')}` : ''}
                    InputProps={{
                        style: {
                            backgroundColor: 'rgba(255, 255, 255, 0.1)',
                            color: 'white',
                        },
                    }}
                    InputLabelProps={{
                        style: { color: 'white' },
                    }}
                />
                <TextField
                    label={t('password')}
                    type={showPassword ? 'text' : 'password'}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    fullWidth
                    margin="normal"
                    error={passwordError}
                    helperText={passwordError ? `${t('password')}${t('input_required')}` : ''}
                    InputProps={{
                        endAdornment: (
                            <InputAdornment position="end">
                                <IconButton
                                    aria-label="toggle password visibility"
                                    onClick={handleClickShowPassword}
                                    onMouseDown={handleMouseDownPassword}
                                    edge="end"
                                >
                                    {showPassword ? <VisibilityOff /> : <Visibility />}
                                </IconButton>
                            </InputAdornment>
                        ),
                        style: {
                            backgroundColor: 'rgba(255, 255, 255, 0.1)',
                            color: 'white',
                        },                    }}
                    InputLabelProps={{
                        style: { color: 'white' },
                    }}
                />
                <Button type="submit" variant="contained" color="primary" sx={{ mt: 2, width: '50%', mx: 'auto', display: 'block', height: 50, borderRadius: 3 }}>
                    {t('register')}
                </Button>
            </Box>
        );
    };

    const renderLogin = () => {
        return (
            <Box component="form" onSubmit={(e) => {e.preventDefault(); handleLogin(email, password, setMessage, setMessageType, navigate)} } sx={{ mt: 2 }}>
                <TextField
                    label={t('email')}
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    fullWidth
                    margin="normal"
                    error={emailError}
                    helperText={emailError ? `${t('email')}${t('input_required')}` : ''}
                    InputProps={{
                        style: {
                            backgroundColor: 'rgba(255, 255, 255, 0.1)',
                            color: 'white',
                        },
                    }}
                    InputLabelProps={{
                        style: { color: 'white' },
                    }}
                />
                <TextField
                    label={t('password')}
                    type={showPassword ? 'text' : 'password'}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    fullWidth
                    margin="normal"
                    error={passwordError}
                    helperText={passwordError ? `${t('password')}${t('input_required')}` : ''}
                    InputProps={{
                        endAdornment: (
                            <InputAdornment position="end">
                                <IconButton
                                    aria-label="toggle password visibility"
                                    onClick={handleClickShowPassword}
                                    onMouseDown={handleMouseDownPassword}
                                    edge="end"
                                >
                                    {showPassword ? <VisibilityOff /> : <Visibility />}
                                </IconButton>
                            </InputAdornment>
                        ),
                        style: {
                            backgroundColor: 'rgba(255, 255, 255, 0.1)',
                            color: 'white',
                        },
                    }}
                    InputLabelProps={{
                        style: { color: 'white' },
                    }}
                />
                <Button type="submit" variant="contained" color="primary" sx={{ mt: 2, width: '50%', mx: 'auto', display: 'block', height: 50, borderRadius: 3 }}>
                    {t('login')}
                </Button>
                <Typography variant="body2" sx={{ mt: 2, textAlign: 'center' }}>
                    <Link to="/auth/request-password-reset" style={{ color: '#1976d2', textDecoration: 'none' }}>
                        {t('forgot_password')}
                    </Link>
                </Typography>
            </Box>
        );
    };

    const renderVerify = () => {
        return (
            <Box component="form" onSubmit={(e) => {e.preventDefault(); handleVerify(email, verificationCode, setMessage, setMessageType, navigate)}} sx={{ mt: 2 }}>
                <TextField
                    label={t('verification_code')}
                    value={verificationCode}
                    onChange={(e) => setVerificationCode(e.target.value)}
                    fullWidth
                    margin="normal"
                    error={verificationCodeError}
                    helperText={verificationCodeError ? `${t('verification_code')}${t('input_required')}` : ''}
                    InputProps={{
                        style: {
                            backgroundColor: 'rgba(255, 255, 255, 0.1)',
                            color: 'white',
                        },
                    }}
                    InputLabelProps={{
                        style: { color: 'white' },
                    }}
                />
                <Typography variant="body2" sx={{ mt: 2 }}>
                    {t('enter_verification_code')}
                </Typography>
                <Button type="submit" variant="contained" color="primary" sx={{ mt: 2, width: '50%', mx: 'auto', display: 'block', height: 50, borderRadius: 3 }}>
                    {t('verify')}
                </Button>
                {/* 認証コードの再送信ボタン */}
                <Button onClick={() => handleResendVerificationCode(email, setMessage, setMessageType)} variant="contained" color="secondary" sx={{ mt: 2, width: '50%', mx: 'auto', display: 'block', height: 50, borderRadius: 3 }}>
                    {t('resend_verification_code')}
                </Button>
            </Box>
        );
    };

    const renderRequestPasswordReset = () => {
        return (
            <Box component="form" onSubmit={(e) => { e.preventDefault(); handleRequestPasswordReset(email, setMessage, setMessageType); }} sx={{ mt: 2 }}>
            <TextField
                label={t('email')}
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                fullWidth
                margin="normal"
                error={emailError}
                helperText={emailError ? `${t('email')}${t('input_required')}` : ''}
                InputProps={{
                style: {
                    backgroundColor: 'rgba(255, 255, 255, 0.1)',
                    color: 'white',
                },
                }}
                InputLabelProps={{
                    style: { color: 'white' },
                }}
            />
            <Button type="submit" variant="contained" color="primary" sx={{ mt: 2, width: '50%', mx: 'auto', display: 'block', height: 50, borderRadius: 3 }}>
                {t('send')}
            </Button>
            </Box>
        );
    };

    const renderResetPassword = () => {
        return (
            <Box component="form" onSubmit={(e) => { e.preventDefault(); handleResetPassword(token, newPassword, setMessage, setMessageType, navigate); }} sx={{ mt: 2 }}>
            <TextField
                label={t('new_password')}
                type={showPassword ? 'text' : 'password'}
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                fullWidth
                margin="normal"
                error={newPasswordError}
                helperText={newPasswordError ? `${t('new_password')}${t('input_required')}` : ''}
                InputProps={{
                endAdornment: (
                    <InputAdornment position="end">
                    <IconButton
                        aria-label="toggle password visibility"
                        onClick={handleClickShowPassword}
                        onMouseDown={handleMouseDownPassword}
                        edge="end"
                    >
                        {showPassword ? <VisibilityOff /> : <Visibility />}
                    </IconButton>
                    </InputAdornment>
                ),
                style: {
                    backgroundColor: 'rgba(255, 255, 255, 0.1)',
                    color: 'white',
                },
                }}
                InputLabelProps={{
                style: { color: 'white' },
                }}
            />
            <Button type="submit" variant="contained" color="primary" sx={{ mt: 2, width: '50%', mx: 'auto', display: 'block', height: 50, borderRadius: 3 }}>
                {t('reset_password')}
            </Button>
            </Box>
        );
    };

    const renderContent = () => {
        switch (location.pathname) {
            case '/auth/register':
                return renderRegister();
            case '/auth/login':
                return renderLogin();
            case '/auth/verify':
                return renderVerify();
            case '/auth/request-password-reset':
                return renderRequestPasswordReset();
            case '/auth/reset-password':
                return renderResetPassword();
            default:
                return <Typography>Page not found</Typography>;
        }
    };

    const getTitle = () => {
        switch (location.pathname) {
            case '/auth/register':
                return t('register');
            case '/auth/login':
                return t('login');
            case '/auth/verify':
                return t('enter_verification_code');
            case '/auth/request-password-reset':
                return t('send_password_reset_email');
            case '/auth/reset-password':
                return t('reset_password');
            default:
                return '';
        }
    };

    return (
        <Paper elevation={3} sx={{ p: 4, maxWidth: 400, mx: 'auto', mt: 4, bgcolor: 'grey.800', color: 'white', borderRadius: 5, boxShadow: 10 }}>
            <Typography variant="h5" component="h1" gutterBottom>
                {getTitle()}
            </Typography>
            {renderContent()}
            {message && (
                <Typography
                    color={messageType === 'error' ? 'error' : 'success'}
                    sx={{ mt: 2 }}
                >
                    {message}
                </Typography>
            )}
        </Paper>
    );
};

export default Auth;