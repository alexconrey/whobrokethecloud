import Moment from 'moment';
import Typography from '@mui/material/Typography';
import Link from '@mui/material/Link';

function Copyright() {
    const now = Moment();
    return (
      <Typography variant="body2" color="text.secondary" align="center">
        {'Copyright © '}
        <Link color="inherit" href="https://whobrokethe.cloud/">
          whobrokethe.cloud
        </Link>{' '}
        {now.year()}
        {/* {new Date().getFullYear()} */}
        {'.'}
      </Typography>
    );
  }

function Footer() {
    Moment.locale('en')
    return (
        <footer>
            <Copyright />
        </footer>
    )
}

export default Footer