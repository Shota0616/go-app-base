import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation, Link as RouterLink } from 'react-router-dom';
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
import { Alert, Card, CardContent } from '@mui/material';

export default function ResetPassword() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();
  
  const [newPassword, setNewPassword] = useState('');
  const [token, setToken] = useState('');
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState('error');

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const tokenParam = params.get('token');
    if (tokenParam) {
      setToken(tokenParam);
    } else {
        setMessage('Invalid password reset link.');
        setMessageType('error');
    }
  }, [location.search]);

  const handleReset = async (event) => {
    event.preventDefault();
    if (!newPassword) {
      setMessage(t('input_required'));
      setMessageType('error');
      return;
    }

    try {
      const response = await axios.post(`${import.meta.env.VITE_APP_API_URL}/api/reset-password`, { token, newPassword });
      setMessage(response.data.message);
      setMessageType('success');
      setTimeout(() => navigate('/auth/login'), 2000); // Redirect to login after 2 seconds
    } catch (error) {
      setMessage(error.response?.data?.error || t('password_reset_failed'));
      setMessageType('error');
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
              {t('reset_password')}
            </Typography>
            <Box component="form" onSubmit={handleReset} noValidate sx={{ mt: 1 }}>
              <TextField
                margin="normal"
                required
                fullWidth
                name="newPassword"
                label={t('new_password')}
                type="password"
                id="newPassword"
                autoFocus
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
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
                disabled={!token}
              >
                {t('reset_password')}
              </Button>
              <Grid container>
                <Grid item>
                  <Link component={RouterLink} to="/auth/login" variant="body2">
                    {t('back_to_login')}
                  </Link>
                </Grid>
              </Grid>
            </Box>
          </Box>
        </CardContent>
      </Card>
    </Container>
  );
}
