import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Box, Button, Typography, Paper } from '@mui/material';

const Logout = () => {
    const { t } = useTranslation();
    const [logoutMessage, setLogoutMessage] = useState('');

    const handleLogout = async () => {
        try {
            localStorage.removeItem('token');
            localStorage.removeItem('refreshtoken');
            window.dispatchEvent(new Event("storage"));
            setLogoutMessage(t('logout_successful'));
        } catch (error) {
            setLogoutMessage(t('logout_failed'));
        }
    };

    return (
        <Paper elevation={3} sx={{ p: 4, maxWidth: 400, mx: 'auto', mt: 4, bgcolor: 'grey.800', color: 'white', borderRadius: 5, boxShadow: 10 }}>
            {logoutMessage ? (
                <Typography>{logoutMessage}</Typography>
            ) : (
                <Box component="form" onSubmit={(e) => { e.preventDefault(); handleLogout(); }}>
                    <Typography variant="p" sx={{ mb: 2 }}>
                        {t('confirm_logout')}
                    </Typography>
                    <Button type="submit" variant="contained" color="primary" sx={{ mt: 2, width: '50%', height: 50, borderRadius: 3 }}>
                        {t('logout')}
                    </Button>
                </Box>
            )}
        </Paper>
    );
};

export default Logout;
