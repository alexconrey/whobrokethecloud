import { render, screen } from '@testing-library/react';
import App from './App';

test('renders loading text', () => {
  render(<App />);
  const textElement = screen.getByText(/Loading.../i);
  expect(textElement).toBeInTheDocument();
});

test('renders fooder', () => {
  render(<App />);
  const copyrightElement = screen.getByText(/Copyright/i);
  expect(copyrightElement).toBeInTheDocument();
});
