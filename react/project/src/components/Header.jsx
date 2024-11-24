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
import SearchIcon from '@mui/icons-material/Search';
import NotificationsIcon from '@mui/icons-material/Notifications';
import Tooltip from '@mui/material/Tooltip';
import PersonAdd from '@mui/icons-material/PersonAdd';
import Settings from '@mui/icons-material/Settings';
import Logout from '@mui/icons-material/Logout';
import ListItemIcon from '@mui/material/ListItemIcon';
import Divider from '@mui/material/Divider';
import Avatar from '@mui/material/Avatar';
import Drawer from '@mui/material/Drawer';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemText from '@mui/material/ListItemText';
import CloseIcon from '@mui/icons-material/Close';
import { useTranslation } from 'react-i18next';
import { useUser } from '/src/context/UserContext';


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

function AccountMenu({ isLoggedIn }) {
    const { t } = useTranslation();
    const [drawerOpen, setDrawerOpen] = useState(false);
    const { user } = useUser();

    const toggleDrawer = (open) => (event) => {
        if (event.type === 'keydown' && (event.key === 'Tab' || event.key === 'Shift')) {
            return;
        }
        setDrawerOpen(open);
    };

    const menuItems = isLoggedIn ? [
        { text: t('mypage'), link: '/mypage', icon: <AccountCircle /> },
        { text: t('settings'), link: '/settings', icon: <Settings /> },
        { divider: true },
        { text: t('logout'), link: '/logout', icon: <Logout /> },
    ] : [
        { text: t('register'), link: '/auth/register', icon: <PersonAdd /> },
        { text: t('login'), link: '/auth/login', icon: <AccountCircle /> },
    ];

    return (
        <>
            <Tooltip title={t('account_settings')}>
                <IconButton
                    onClick={toggleDrawer(true)}
                    size="large"
                    color="inherit"
                >
                    <AccountCircle />
                </IconButton>
            </Tooltip>
            <Drawer anchor="right" open={drawerOpen} onClose={toggleDrawer(false)}>
                <Box
                    sx={{ width: 250, display: 'flex', flexDirection: 'column', height: '100%', bgcolor: '#242424', color: 'whitesmoke' }}
                    role="presentation"
                >
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mx: 2, my: 1 }}>
                        {isLoggedIn && (
                            <Typography variant="body1" sx={{ color: 'white' }}>
                                hi! {user?.username}
                            </Typography>
                        )}
                        <IconButton onClick={toggleDrawer(false)}>
                            <CloseIcon sx={{color: 'gray'}}/>
                        </IconButton>
                    </Box>
                    <Divider sx={{ bgcolor: '#666666', mx: 2 }}/>
                    <List>
                        {menuItems.map((item, index) => (
                            item.divider ? ( <Divider key={index} sx={{ bgcolor: '#666666', mx: 2, my: 1 }}/> ) : (
                            <ListItem button component={Link} to={item.link} key={index} dense sx={{ py: 0.25 }} onClick={toggleDrawer(false)}>
                                <ListItemIcon sx={{ color: '#999999', minWidth: 36 }} fontSize="small">{item.icon}</ListItemIcon>
                                <ListItemText primary={item.text} />
                            </ListItem>
                            )
                        ))}
                    </List>
                </Box>
            </Drawer>
        </>
    );
}

export default function PrimarySearchAppBar() {
    const { t } = useTranslation();
    const [isLoggedIn, setIsLoggedIn] = useState(Boolean(localStorage.getItem('token')));

    useEffect(() => {
        const handleStorageChange = () => {
            setIsLoggedIn(Boolean(localStorage.getItem('token')));
        };

        window.addEventListener('storage', handleStorageChange);
        window.dispatchEvent(new Event("storage"));

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
                        <AccountMenu isLoggedIn={isLoggedIn} />
                    </Box>
                </Toolbar>
            </AppBar>
        </Box>
    );
}