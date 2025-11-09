import React, { useState } from 'react';
import { Link as RouterLink } from 'react-router-dom';
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

export default function RequestPasswordReset() {
  const { t } = useTranslation();
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState('error');

  const handleRequest = async (event) => {
    event.preventDefault();
    if (!email) {
      setMessage(t('input_required'));
      setMessageType('error');
      return;
    }

    try {
      const response = await axios.post(`${import.meta.env.VITE_APP_API_URL}/api/request-password-reset`, { email });
      setMessage(response.data.message);
      setMessageType('success');
    } catch (error) {
      setMessage(error.response?.data?.error || t('password_reset_request_failed'));
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
              {t('send_password_reset_email')}
            </Typography>
            <Box component="form" onSubmit={handleRequest} noValidate sx={{ mt: 1 }}>
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
                {t('send')}
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
