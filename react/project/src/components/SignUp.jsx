import React, { useState } from 'react';
import { useNavigate, Link as RouterLink } from 'react-router-dom';
import axios from 'axios';
import { useTranslation } from 'react-i18next';

import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';
import Link from '@mui/material/Link';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import { Alert, Card, CardContent, InputAdornment, IconButton } from '@mui/material';
import Snackbar from '@mui/material/Snackbar';
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';

export default function SignUp() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState('error');
  const [showPassword, setShowPassword] = useState(false);
  const [snackbarOpen, setSnackbarOpen] = useState(false);

  const validateEmail = (email) => {
    // シンプルなメールアドレス正規表現
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
  };

  const handleClickShowPassword = () => setShowPassword((show) => !show);
  const handleMouseDownPassword = (event) => {
    event.preventDefault();
  };

  const handleRegister = async (event) => {
    event.preventDefault();
    if (!email || !password) {
      setMessage(t('input_required'));
      setMessageType('error');
      setSnackbarOpen(true);
      return;
    }
    if (!validateEmail(email)) {
      setMessage(t('email_invalid'));
      setMessageType('error');
      setSnackbarOpen(true);
      return;
    }
    if (password.length < 8) {
      setMessage(t('password_min_length', { min: 8 }));
      setMessageType('error');
      setSnackbarOpen(true);
      return;
    }

    try {
      const response = await axios.post(`${import.meta.env.VITE_APP_API_URL}/api/register`, {
        email,
        password,
      });
      setMessage(response.data.message);
      setMessageType('success');
      setSnackbarOpen(true);
      setTimeout(() => navigate('/auth/verify', { state: { email } }), 1000);
    } catch (error) {
      setMessage(error.response?.data?.error || t('registration_failed'));
      setMessageType('error');
      setSnackbarOpen(true);
    }
  };

  const handleCloseSnackbar = (event, reason) => {
    if (reason === 'clickaway') return;
    setSnackbarOpen(false);
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
              {t('register')}
            </Typography>
            <Box component="form" noValidate onSubmit={handleRegister} sx={{ mt: 3 }}>
              <Grid container spacing={2}>
                <Grid item xs={12}>
                  <TextField
                    required
                    fullWidth
                    id="email"
                    label={t('email')}
                    name="email"
                    autoComplete="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                  />
                </Grid>
                <Grid item xs={12}>
                  <TextField
                    required
                    fullWidth
                    name="password"
                    label={t('password')}
                    type={showPassword ? 'text' : 'password'}
                    id="password"
                    autoComplete="new-password"
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
                </Grid>
              </Grid>
              <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
              >
                {t('register')}
              </Button>
              <Grid container justifyContent="flex-end">
                <Grid item>
                  <Link component={RouterLink} to="/auth/login" variant="body2">
                    {t('signin_link_text')}
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
