import React, { useState } from 'react';
import { useNavigate, Link as RouterLink } from 'react-router-dom';
import axios from 'axios';
import { useUser } from '/src/context/UserContext';
import { useTranslation } from 'react-i18next';

import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';
import FormControlLabel from '@mui/material/FormControlLabel';
import Checkbox from '@mui/material/Checkbox';
import Link from '@mui/material/Link';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import { Alert, Card, CardContent, InputAdornment, IconButton } from '@mui/material';
import Snackbar from '@mui/material/Snackbar';
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';

export default function SignIn() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { setUser } = useUser();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [rememberMe, setRememberMe] = useState(false);
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState('error');
  const [showAlert, setShowAlert] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [snackbarOpen, setSnackbarOpen] = useState(false);

  const handleClickShowPassword = () => setShowPassword((show) => !show);
  const handleMouseDownPassword = (event) => {
    event.preventDefault();
  };

  const handleLogin = async (event) => {
    event.preventDefault();
    if (!email || !password) {
        setMessage(t('login_failed'));
        setMessageType('error');
        setShowAlert(true);
        setSnackbarOpen(true);
        return;
    }

    try {
        const response = await axios.post(`${import.meta.env.VITE_APP_API_URL}/api/login`, {
            email,
            password,
        });

        const token = response.data.token;
        const refreshToken = response.data.refreshtoken;

        if (rememberMe) {
            localStorage.setItem('token', token);
            localStorage.setItem('refreshtoken', refreshToken);
            sessionStorage.removeItem('token');
            sessionStorage.removeItem('refreshtoken');
        } else {
            sessionStorage.setItem('token', token);
            sessionStorage.setItem('refreshtoken', refreshToken);
            localStorage.removeItem('token');
            localStorage.removeItem('refreshtoken');
        }
        
        window.dispatchEvent(new Event("storage"));
        setUser(response.data.user);
        navigate('/mypage');
    } catch (error) {
        if (error.response?.status === 303) {
            navigate("/auth/verify", { state: { message: error.response.data.error, messageType: 'error' } });
        } else {
            setMessage(error.response?.data?.error || t('login_failed'));
            setMessageType('error');
            setShowAlert(true);
            setSnackbarOpen(true);
        }
    }
  };

  const handleCloseSnackbar = (event, reason) => {
    if (reason === 'clickaway') return;
    setSnackbarOpen(false);
    setShowAlert(false);
    setMessage('');
  };

  return (
    <Container component="main" maxWidth="xs">
      <CssBaseline />
      <Card sx={{ marginTop: 8, borderRadius: 2 }}>
        <CardContent sx={{ p: 4 }}>
          <Box
            sx={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
            }}
          >
            <Typography component="h1" variant="h5">
              {t('login')}
            </Typography>
            <Box component="form" onSubmit={handleLogin} noValidate sx={{ mt: 1 }}>
              <TextField
                margin="normal"
                required
                fullWidth
                id="email"
                label={t('email')}
                name="email"
                autoComplete="email"
                autoFocus
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
              <TextField
                margin="normal"
                required
                fullWidth
                name="password"
                label={t('password')}
                type={showPassword ? 'text' : 'password'}
                id="password"
                autoComplete="current-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
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
                }}
              />
              <FormControlLabel
                control={
                  <Checkbox 
                    value="remember" 
                    color="primary" 
                    checked={rememberMe} 
                    onChange={(e) => setRememberMe(e.target.checked)} 
                  />
                }
                label={t('remember_me')}
              />
              <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
              >
                {t('login')}
              </Button>
              <Grid container justifyContent="space-between">
                <Grid item>
                  <Link component={RouterLink} to="/auth/request-password-reset" variant="body2">
                    {t('forgot_password')}
                  </Link>
                </Grid>
                <Grid item>
                  <Link component={RouterLink} to="/auth/register" variant="body2">
                    {t('signup_link_text')}
                  </Link>
                </Grid>
              </Grid>
            </Box>
          </Box>
        </CardContent>
      </Card>
      <Snackbar open={snackbarOpen} autoHideDuration={6000} onClose={handleCloseSnackbar} anchorOrigin={{ vertical: 'top', horizontal: 'center' }}>
        <Alert onClose={handleCloseSnackbar} severity={messageType} sx={{ width: '100%' }}>
          {message}
        </Alert>
      </Snackbar>
    </Container>
  );
}
