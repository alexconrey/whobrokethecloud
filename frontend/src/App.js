import logo from './logo.svg';
import './App.css';
import Home from './Home';
import Footer from './Footer';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';

function App() {
  return (
    <Grid 
      container
      spacing={0}
      direction="column"
      alignItems="center"
      justifyContent="center"
      style={{ minHeight: '80vh' }}
    >
      <Grid item xs={3}>
        <Box
          component="img"
          sx={{
            height: 233,
            width: 350,
            maxHeight: { xs: 233, md: 167 },
            maxWidth: { xs: 350, md: 250 },
          }}          
          alt="logo"
          src={logo}
        />
      </Grid>
      <Grid item xs={3}>
        <Home />
      </Grid>
      <Grid item xs={3}>
        <Footer />
      </Grid>
    </Grid>
  );
}

export default App;
