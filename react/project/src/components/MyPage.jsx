// src/components/MyPage.jsx
import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { useTranslation } from 'react-i18next';
import { Container, Paper, Typography, Box, CircularProgress, Avatar } from '@mui/material';
import { useUser } from '/src/context/UserContext';

const MyPage = () => {
    const { t } = useTranslation();
    const { user, setUser } = useUser();
    const navigate = useNavigate();

    useEffect(() => {
        const checkAuth = async () => {
            const token = localStorage.getItem('token');
            if (!token) {
                navigate('/auth/login');
                return;
            }

            try {
                const response = await axios.get(`${import.meta.env.VITE_APP_API_URL}/api/getuser`, {
                    headers: {
                        'Authorization': token,
                    },
                });
                setUser(response.data);
            } catch (error) {
                localStorage.removeItem('token');
                window.dispatchEvent(new Event("storage"));
                navigate('/auth/login');
            }
        };

        checkAuth();
    }, [navigate]);

    return (
        <Container maxWidth="sm" sx={{ mt: 4 }}>
            <Paper elevation={3} sx={{ p: 4, borderRadius: 5, bgcolor: 'grey.800', color: 'white', boxShadow: 10 }}>
                <Typography variant="h5" component="h1" gutterBottom>
                    {t('mypage_title')}
                </Typography>
                {user ? (
                    <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                        {/* <Avatar sx={{ width: 80, height: 80, mb: 2 }}>{user.username.charAt(0)}</Avatar> */}
                        <Typography variant="h6">{t('username')}: {user.username}</Typography>
                    </Box>
                ) : (
                    <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100px' }}>
                        <CircularProgress />
                    </Box>
                )}
            </Paper>
        </Container>
    );
};

export default MyPage;
