import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import reactLogo from '/src/assets/react.svg';
import { styled, alpha } from '@mui/material/styles';
import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import InputBase from '@mui/material/InputBase';
import Badge from '@mui/material/Badge';
import MenuItem from '@mui/material/MenuItem';
import Menu from '@mui/material/Menu';
import AccountCircle from '@mui/icons-material/AccountCircle';
import MailIcon from '@mui/icons-material/Mail';
import NotificationsIcon from '@mui/icons-material/Notifications';
import Tooltip from '@mui/material/Tooltip';
import PersonAdd from '@mui/icons-material/PersonAdd';
import Settings from '@mui/icons-material/Settings';
import Logout from '@mui/icons-material/Logout';
import ListItemIcon from '@mui/material/ListItemIcon';
import Divider from '@mui/material/Divider';
import Avatar from '@mui/material/Avatar';
import SearchIcon from '@mui/icons-material/Search';
import LoginIcon from '@mui/icons-material/Login';
import { useTranslation } from 'react-i18next';

const Search = styled('div')(({ theme }) => ({
    position: 'relative',
    borderRadius: theme.shape.borderRadius,
    backgroundColor: alpha(theme.palette.common.white, 0.15),
    '&:hover': {
        backgroundColor: alpha(theme.palette.common.white, 0.25),
    },
    marginRight: theme.spacing(2),
    marginLeft: 0,
    width: '100%',
    [theme.breakpoints.up('sm')]: {
        marginLeft: theme.spacing(3),
        width: 'auto',
    },
}));

const SearchIconWrapper = styled('div')(({ theme }) => ({
    padding: theme.spacing(0, 2),
    height: '100%',
    position: 'absolute',
    pointerEvents: 'none',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
}));

const StyledInputBase = styled(InputBase)(({ theme }) => ({
    color: 'inherit',
    '& .MuiInputBase-input': {
        padding: theme.spacing(1, 1, 1, 0),
        paddingLeft: `calc(1em + ${theme.spacing(4)})`,
        transition: theme.transitions.create('width'),
        width: '100%',
        [theme.breakpoints.up('md')]: {
            width: '20ch',
        },
    },
}));

// ユーザのボタンを押下したときのモーダル
function AccountMenu({ isLoggedIn }) { // isLoggedInをプロパティとして受け取る
    const { t } = useTranslation();
    const [anchorEl, setAnchorEl] = useState(null);
    const open = Boolean(anchorEl);

    const handleClick = (event) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    return (
        <React.Fragment>
            <Box sx={{ display: 'flex', alignItems: 'center', textAlign: 'center' }}>
                {/* <Typography sx={{ minWidth: 100 }}>Profile</Typography> */}
                <Tooltip title={t('account_settings')}>
                    <IconButton
                        onClick={handleClick}
                        size="large"
                        aria-controls={open ? 'account-menu' : undefined}
                        aria-haspopup="true"
                        aria-expanded={open ? 'true' : undefined}
                        color="inherit"
                    >
                        <AccountCircle />
                    </IconButton>
                </Tooltip>
            </Box>
            <Menu
                anchorEl={anchorEl}
                id="account-menu"
                open={open}
                onClose={handleClose}
                onClick={handleClose}
                PaperProps={{
                    elevation: 0,
                    sx: {
                        overflow: 'visible',
                        filter: 'drop-shadow(0px 2px 8px rgba(0,0,0,0.32))',
                        mt: 1.5,
                        '& .MuiAvatar-root': {
                            width: 32,
                            height: 32,
                            ml: -0.5,
                            mr: 1,
                        },
                        '&::before': {
                            content: '""',
                            display: 'block',
                            position: 'absolute',
                            top: 0,
                            right: 14,
                            width: 10,
                            height: 10,
                            bgcolor: 'background.paper',
                            transform: 'translateY(-50%) rotate(45deg)',
                            zIndex: 0,
                        },
                    },
                }}
                transformOrigin={{ horizontal: 'right', vertical: 'top' }}
                anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
            >
                {isLoggedIn ? (
                    // ログインされている時の処理
                    [
                        <Link to="/mypage" key="mypage" style={{ textDecoration: 'none', color: 'inherit' }}>
                            <MenuItem onClick={handleClose}>
                                <Avatar /> {t('mypage')}
                            </MenuItem>
                        </Link>,
                        <Divider key="divider" />,
                        <Link to="/settings" key="settings" style={{ textDecoration: 'none', color: 'inherit' }}>
                            <MenuItem onClick={handleClose}>
                                <ListItemIcon>
                                    <Settings fontSize="small" />
                                </ListItemIcon>
                                {t('settings')}
                            </MenuItem>
                        </Link>,
                        <Link to="/logout" key="logout-option" style={{ textDecoration: 'none', color: 'inherit' }}>
                            <MenuItem onClick={handleClose}>
                                <ListItemIcon>
                                    <Logout fontSize="small" />
                                </ListItemIcon>
                                {t('logout')}
                            </MenuItem>
                        </Link>
                    ]
                                    ) : (
                    // ログインされていない時の処理
                    [
                        <Link to="/auth/register" key="register" style={{ textDecoration: 'none', color: 'inherit' }}>
                            <MenuItem onClick={handleClose}>
                                <ListItemIcon>
                                    <PersonAdd fontSize="small" />
                                </ListItemIcon>
                                {t('register')}
                            </MenuItem>
                        </Link>,
                        <Link to="/auth/login" key="login" style={{ textDecoration: 'none', color: 'inherit' }}>
                            <MenuItem onClick={handleClose}>
                                <ListItemIcon>
                                    <LoginIcon fontSize="small" />
                                </ListItemIcon>
                                {t('login')}
                            </MenuItem>
                        </Link>,
                        <Link to="/settings" key="settings" style={{ textDecoration: 'none', color: 'inherit' }}>
                            <MenuItem onClick={handleClose}>
                                <ListItemIcon>
                                    <Settings fontSize="small" />
                                </ListItemIcon>
                                {t('settings')}
                            </MenuItem>
                        </Link>,
                    ]
                )}
            </Menu>
        </React.Fragment>
    );
}

export default function PrimarySearchAppBar() {

    // ログイン状態の有無を判断するためにlocalstorageにtokenが存在��るか見る
    const { t } = useTranslation(); // 追加
    const [isLoggedIn, setIsLoggedIn] = useState(Boolean(localStorage.getItem('token')));
    console.log('isLoggedIn:', isLoggedIn);

    useEffect(() => {
        // storageイベントリスナーの設定
        const handleStorageChange = () => {
            setIsLoggedIn(Boolean(localStorage.getItem('token')));
        };

        // イベントリスナーを追加
        window.addEventListener('storage', handleStorageChange);
        // storageイベント発火
        window.dispatchEvent(new Event("storage"));

        // クリーンアップ関数でイベントリスナーを削除
        return () => {
            window.removeEventListener('storage', handleStorageChange);
        };
    }, [isLoggedIn]);

    return (
        <Box sx={{ flexGrow: 1 }}>
            <AppBar position="static" sx={{ bgcolor: 'black' }}>
                <Toolbar>
                    <Box sx={{ display: 'flex', alignItems: 'center', mr: 2 }}>
                        <img src={reactLogo} alt="React Logo" style={{ height: 40, width: 40 }} />
                    </Box>
                    <Search>
                        <SearchIconWrapper>
                            <SearchIcon />
                        </SearchIconWrapper>
                        <StyledInputBase
                            placeholder={t('search')}
                            inputProps={{ 'aria-label': 'search' }}
                        />
                    </Search>
                    <Box sx={{ flexGrow: 1 }} />
                    <Box sx={{ display: { xs: 'flex', md: 'flex' } }}>
                        <IconButton
                            size="large"
                            aria-label={t('show_new_notifications')}
                            color="inherit"
                        >
                            <Badge badgeContent={17} color="error">
                                <NotificationsIcon />
                            </Badge>
                        </IconButton>
                        <AccountMenu isLoggedIn={isLoggedIn}/>
                    </Box>
                </Toolbar>
            </AppBar>
        </Box>
    );
}