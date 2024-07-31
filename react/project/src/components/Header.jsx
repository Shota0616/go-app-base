import React from 'react';
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
import MenuIcon from '@mui/icons-material/Menu';
import SearchIcon from '@mui/icons-material/Search';
import AccountCircle from '@mui/icons-material/AccountCircle';
import MailIcon from '@mui/icons-material/Mail';
import NotificationsIcon from '@mui/icons-material/Notifications';
import MoreIcon from '@mui/icons-material/MoreVert';
import Tooltip from '@mui/material/Tooltip';
import PersonAdd from '@mui/icons-material/PersonAdd';
import Settings from '@mui/icons-material/Settings';
import Logout from '@mui/icons-material/Logout';
import ListItemIcon from '@mui/material/ListItemIcon';
import Divider from '@mui/material/Divider';
import Avatar from '@mui/material/Avatar';


// 検索ボックスのstyleを定義
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
    // vertical padding + font size from searchIcon
    paddingLeft: `calc(1em + ${theme.spacing(4)})`,
    transition: theme.transitions.create('width'),
    width: '100%',
    [theme.breakpoints.up('md')]: {
        width: '20ch',
    },
    },
}));

// アカウントメニュー
function AccountMenu() {
    const [anchorEl, setAnchorEl] = React.useState(null);
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
        <Typography sx={{ minWidth: 100 }}>Profile</Typography>
        <Tooltip title="Account settings">
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
        <Link to="/viteinit">
            <MenuItem onClick={handleClose}>
                <Avatar /> Profile
            </MenuItem>
        </Link>
        <Link to="/auth">
            <MenuItem onClick={handleClose}>
                <Avatar /> Register
            </MenuItem>
        </Link>
        <Link to="/logout">
            <MenuItem onClick={handleClose}>
                <Avatar /> Logout
            </MenuItem>
        </Link>
        <Divider />
        <MenuItem onClick={handleClose}>
            <ListItemIcon>
            <PersonAdd fontSize="small" />
            </ListItemIcon>
            Add another account
        </MenuItem>
        <MenuItem onClick={handleClose}>
            <ListItemIcon>
            <Settings fontSize="small" />
            </ListItemIcon>
            Settings
        </MenuItem>
        <MenuItem onClick={handleClose}>
            <ListItemIcon>
            <Logout fontSize="small" />
            </ListItemIcon>
            Logout
        </MenuItem>
        </Menu>
    </React.Fragment>
    );
}


export default function PrimarySearchAppBar() {
    // const [anchorEl, setAnchorEl] = React.useState(null);
    // const [mobileMoreAnchorEl, setMobileMoreAnchorEl] = React.useState(null);



    // const handleMobileMenuClose = () => {
    // setMobileMoreAnchorEl(null);
    // };


    const mobileMenuId = 'primary-search-account-menu-mobile';

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
                placeholder="Search…"
                inputProps={{ 'aria-label': 'search' }}
            />
            </Search>
            <Box sx={{ flexGrow: 1 }} />
            <Box sx={{ display: { xs: 'none', md: 'flex' } }}>
            <IconButton size="large" aria-label="show 4 new mails" color="inherit">
                <Badge badgeContent={4} color="error">
                <MailIcon />
                </Badge>
            </IconButton>
            <IconButton
                size="large"
                aria-label="show 17 new notifications"
                color="inherit"
            >
                <Badge badgeContent={17} color="error">
                <NotificationsIcon />
                </Badge>
            </IconButton>
            <AccountMenu />
            </Box>
        </Toolbar>
        </AppBar>
    </Box>
    );
}