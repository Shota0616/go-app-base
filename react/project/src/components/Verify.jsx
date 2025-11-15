import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import axios from 'axios';
import { useTranslation } from 'react-i18next';

import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';
import Link from '@mui/material/Link';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import { Alert, Card, CardContent } from '@mui/material';
import Snackbar from '@mui/material/Snackbar';

export default function Verify() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();
  
  const [verificationCode, setVerificationCode] = useState('');
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState(location.state?.message || '');
  const [messageType, setMessageType] = useState(location.state?.messageType || 'info');
  const [snackbarOpen, setSnackbarOpen] = useState(false);

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const urlEmail = params.get('email');
    const urlCode = params.get('code');

    if (urlEmail) {
      setEmail(urlEmail);
    } else if (location.state?.email) {
      setEmail(location.state.email);
    }

    if (urlCode) {
      setVerificationCode(urlCode);
    }
  }, [location.search, location.state]);

  const handleVerify = async (event) => {
    event.preventDefault();
    if (!email || !verificationCode) {
      setMessage(t('input_required'));
      setMessageType('error');
      setSnackbarOpen(true);
      return;
    }

    try {
      const response = await axios.post(`${import.meta.env.VITE_APP_API_URL}/api/verify`, { email, verificationCode }, { headers: { Authorization: '' } }); // Explicitly send no auth header
      setMessage(response.data.message);
      setMessageType('success');
      setSnackbarOpen(true);
      setTimeout(() => navigate('/auth/login'), 2000);
    } catch (error) {
      setMessage(error.response?.data?.error || t('verification_failed'));
      setMessageType('error');
      setSnackbarOpen(true);
    }
  };

  const handleResend = async () => {
    try {
        const response = await axios.post(`${import.meta.env.VITE_APP_API_URL}/api/resend-verification-code`, { email }, { headers: { Authorization: '' } }); // Explicitly send no auth header
        setMessage(response.data.message);
        setMessageType('success');
        setSnackbarOpen(true);
    } catch (error) {
        setMessage(error.response?.data?.error || t('resend_verification_code_failed'));
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
              {t('enter_verification_code')}
            </Typography>
            <Typography variant="body2" color="text.secondary" align="center" sx={{ mt: 1 }}>
              {t('verification_code_sent_to', { email: email })}
            </Typography>
            <Box component="form" onSubmit={handleVerify} noValidate sx={{ mt: 1 }}>
              <TextField
                margin="normal"
                required
                fullWidth
                id="verificationCode"
                label={t('verification_code')}
                name="verificationCode"
                autoFocus
                value={verificationCode}
                onChange={(e) => setVerificationCode(e.target.value)}
              />
              <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
              >
                {t('verify')}
              </Button>
              <Button
                fullWidth
                variant="text"
                onClick={handleResend}
              >
                {t('resend_verification_code')}
              </Button>
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
