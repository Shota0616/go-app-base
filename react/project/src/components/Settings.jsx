import React, { useState, useContext, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { useTranslation } from 'react-i18next';
import { 
    Container, 
    Typography, 
    Box, 
    TextField, 
    Button, 
    Card,
    CardContent,
    Divider,
    Link,
    Stack,
    Snackbar,
    FormControl, InputLabel, Select, MenuItem,
    InputAdornment, IconButton // Added InputAdornment, IconButton
} from '@mui/material';
import MuiAlert from '@mui/material/Alert';
import getToken from '/src/utils/getToken';
import { ThemeContext } from '/src/context/ThemeContext';
import Visibility from '@mui/icons-material/Visibility'; // Added Visibility icon
import VisibilityOff from '@mui/icons-material/VisibilityOff'; // Added VisibilityOff icon

// Custom Alert component for Snackbar
const Alert = React.forwardRef(function Alert(props, ref) {
  return <MuiAlert elevation={6} ref={ref} variant="filled" {...props} />;
});

function Copyright(props) {
  const { t } = useTranslation();
  return (
    <Typography variant="body2" color="text.secondary" align="center" {...props}>
      {'Copyright © '}
      <Link color="inherit" href="https://mui.com/">
        {t('site_name')}
      </Link>{' '}
      {new Date().getFullYear()}
      {'.'}
    </Typography>
  );
}

const Settings = () => {
    const { t, i18n } = useTranslation();
    const { toggleTheme, mode } = useContext(ThemeContext);
    const navigate = useNavigate();

    const [currentPassword, setCurrentPassword] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [confirmNewPassword, setConfirmNewPassword] = useState('');
    const [message, setMessage] = useState('');
    const [messageType, setMessageType] = useState('');
    const [snackbarOpen, setSnackbarOpen] = useState(false);
    const [currentLang, setCurrentLang] = useState(i18n.language);

    // States for password visibility
    const [showCurrentPassword, setShowCurrentPassword] = useState(false);
    const [showNewPassword, setShowNewPassword] = useState(false);
    const [showConfirmNewPassword, setShowConfirmNewPassword] = useState(false);

    const handleClickShowCurrentPassword = () => setShowCurrentPassword((show) => !show);
    const handleClickShowNewPassword = () => setShowNewPassword((show) => !show);
    const handleClickShowConfirmNewPassword = () => setShowConfirmNewPassword((show) => !show);

    const handleMouseDownPassword = (event) => {
        event.preventDefault();
    };

    useEffect(() => {
        setCurrentLang(i18n.language);
    }, [i18n.language]);

    const handleLangChange = (event) => {
        const newLang = event.target.value;
        i18n.changeLanguage(newLang);
        setCurrentLang(newLang);
    };

    const handlePasswordChange = async (e) => {
        e.preventDefault();
        if (newPassword !== confirmNewPassword) {
            setMessage(t('passwords_do_not_match'));
            setMessageType('error');
            setSnackbarOpen(true);
            return;
        }

        const token = getToken();
        if (!token) {
            navigate('/auth/login');
            return;
        }

        try {
            const response = await axios.put(
                `${import.meta.env.VITE_APP_API_URL}/api/user/password`,
                { currentPassword, newPassword },
                { headers: { 'Authorization': token } }
            );
            setMessage(response.data.message);
            setMessageType('success');
            setSnackbarOpen(true);
            setCurrentPassword('');
            setNewPassword('');
            setConfirmNewPassword('');
        } catch (error) {
            setMessage(error.response?.data?.error || t('password_update_failed'));
            setMessageType('error');
            setSnackbarOpen(true);
        }
    };

    const handleCloseSnackbar = (event, reason) => {
        if (reason === 'clickaway') {
            return;
        }
        setSnackbarOpen(false);
        setMessage('');
    };

    return (
        <Container component="main" maxWidth="sm" sx={{ mt: 4 }}>
            <Card elevation={3} sx={{ borderRadius: 2 }}>
                <CardContent sx={{ p: 4 }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
                        <Typography variant="h5" component="h1">
                            {t('settings')}
                        </Typography>
                    </Box>

                    <Divider sx={{ my: 2 }} />

                    <Stack spacing={3} sx={{ mb: 2 }}> {/* Increased spacing */}
                        {/* Language Settings */}
                        <Box>
                            <Typography variant="subtitle1" color="text.secondary" gutterBottom>
                                {t('language_settings')}
                            </Typography>
                            <FormControl variant="outlined" fullWidth size="small">
                                <InputLabel id="lang-select-label">{t('language')}</InputLabel>
                                <Select
                                    labelId="lang-select-label"
                                    value={currentLang}
                                    onChange={handleLangChange}
                                    label={t('language')}
                                >
                                    <MenuItem value="en">{t('english')}</MenuItem>
                                    <MenuItem value="ja">{t('japanese')}</MenuItem>
                                </Select>
                            </FormControl>
                        </Box>

                        {/* Theme Settings */}
                        <Box>
                            <Typography variant="subtitle1" color="text.secondary" gutterBottom>
                                {t('theme_settings')}
                            </Typography>
                            <FormControl variant="outlined" fullWidth size="small">
                                <InputLabel id="theme-select-label">{t('theme')}</InputLabel>
                                <Select
                                    labelId="theme-select-label"
                                    value={mode}
                                    onChange={toggleTheme}
                                    label={t('theme')}
                                >
                                    <MenuItem value="light">{t('light')}</MenuItem>
                                    <MenuItem value="dark">{t('dark')}</MenuItem>
                                </Select>
                            </FormControl>
                        </Box>

                        <Divider sx={{ my: 2 }} />

                        {/* Password Change */}
                        <Typography variant="subtitle1" color="text.secondary">
                            {t('change_password')}
                        </Typography>
                        <Box component="form" onSubmit={handlePasswordChange}>
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                name="currentPassword"
                                label={t('current_password')}
                                type={showCurrentPassword ? 'text' : 'password'} // Dynamic type
                                id="currentPassword"
                                autoComplete="current-password"
                                value={currentPassword}
                                onChange={(e) => setCurrentPassword(e.target.value)}
                                sx={{ mb: 2 }}
                                InputProps={{ // Added InputProps
                                    endAdornment: (
                                        <InputAdornment position="end">
                                            <IconButton
                                                aria-label="toggle current password visibility"
                                                onClick={handleClickShowCurrentPassword}
                                                onMouseDown={handleMouseDownPassword}
                                                edge="end"
                                            >
                                                {showCurrentPassword ? <VisibilityOff /> : <Visibility />}
                                            </IconButton>
                                        </InputAdornment>
                                    ),
                                }}
                            />
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                name="newPassword"
                                label={t('new_password')}
                                type={showNewPassword ? 'text' : 'password'} // Dynamic type
                                id="newPassword"
                                autoComplete="new-password"
                                value={newPassword}
                                onChange={(e) => setNewPassword(e.target.value)}
                                sx={{ mb: 2 }}
                                InputProps={{ // Added InputProps
                                    endAdornment: (
                                        <InputAdornment position="end">
                                            <IconButton
                                                aria-label="toggle new password visibility"
                                                onClick={handleClickShowNewPassword}
                                                onMouseDown={handleMouseDownPassword}
                                                edge="end"
                                            >
                                                {showNewPassword ? <VisibilityOff /> : <Visibility />}
                                            </IconButton>
                                        </InputAdornment>
                                    ),
                                }}
                            />
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                name="confirmNewPassword"
                                label={t('confirm_new_password')}
                                type={showConfirmNewPassword ? 'text' : 'password'} // Dynamic type
                                id="confirmNewPassword"
                                autoComplete="new-password"
                                value={confirmNewPassword}
                                onChange={(e) => setConfirmNewPassword(e.target.value)}
                                sx={{ mb: 2 }}
                                InputProps={{ // Added InputProps
                                    endAdornment: (
                                        <InputAdornment position="end">
                                            <IconButton
                                                aria-label="toggle confirm new password visibility"
                                                onClick={handleClickShowConfirmNewPassword}
                                                onMouseDown={handleMouseDownPassword}
                                                edge="end"
                                            >
                                                {showConfirmNewPassword ? <VisibilityOff /> : <Visibility />}
                                            </IconButton>
                                        </InputAdornment>
                                    ),
                                }}
                            />
                            <Button
                                type="submit"
                                fullWidth
                                variant="contained"
                                color="primary"
                                sx={{ mt: 2, mb: 2 }}
                            >
                                {t('change_password')}
                            </Button>
                        </Box>
                    </Stack>
                </CardContent>
            </Card>
            <Copyright sx={{ mt: 8, mb: 4 }} />

            <Snackbar open={snackbarOpen} autoHideDuration={6000} onClose={handleCloseSnackbar} anchorOrigin={{ vertical: 'top', horizontal: 'center' }}>
                <Alert onClose={handleCloseSnackbar} severity={messageType} sx={{ width: '100%' }}>
                    {message}
                </Alert>
            </Snackbar>
        </Container>
    );
};

export default Settings;