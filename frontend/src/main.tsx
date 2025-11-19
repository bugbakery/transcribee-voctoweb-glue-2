import { createRoot } from 'react-dom/client';
import App from './App.tsx';
import { AuthProvider } from './auth.tsx';
import { HomePage } from './pages/HomePage.tsx';
import { UserPage } from './pages/UserPage.tsx';
import { TalkPage } from './pages/TalkPage.tsx';
import { createBrowserRouter, RouterProvider } from 'react-router';

const Router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      {
        path: '/',
        element: <HomePage />,
      },
      {
        path: '/user',
        element: <UserPage />,
      },
      {
        path: '/talk/:id',
        element: <TalkPage />,
      },
    ],
  },
]);


createRoot(document.getElementById('root')!).render(
  <AuthProvider>
    <RouterProvider router={Router} />
  </AuthProvider>,
);
