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
import { Alert, Card, CardContent } from '@mui/material';

export default function SignIn() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { setUser } = useUser();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [rememberMe, setRememberMe] = useState(false); // Add rememberMe state
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState('error');

  const handleLogin = async (event) => {
    event.preventDefault();
    if (!email || !password) {
        setMessage(t('login_failed'));
        setMessageType('error');
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
        }
    }
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
                type="password"
                id="password"
                autoComplete="current-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
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
              {message && (
                <Alert severity={messageType} sx={{ mt: 2, width: '100%' }}>
                    {message}
                </Alert>
              )}
              <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
              >
                {t('login')}
              </Button>
              <Grid container>
                <Grid item xs>
                  <Link component={RouterLink} to="/auth/request-password-reset" variant="body2">
                    {t('forgot_password')}
                  </Link>
                </Grid>
                <Grid item>
                                <Link component={RouterLink} to="/auth/register" variant="body2">
                                  {t('signup_link_text')}
                                </Link>                </Grid>
              </Grid>
            </Box>
          </Box>
        </CardContent>
      </Card>
    </Container>
  );
}
