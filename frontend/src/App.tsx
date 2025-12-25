import { useAuth } from './auth';

import { Link } from './components/Link';
import { Button } from './components/Button';
import { Outlet, ScrollRestoration } from 'react-router';
import { Suspense, useState } from 'react';

function App() {
  const auth = useAuth();

  return (
    <>
      <header className="flex justify-between items-center h-14 mb-4 border-b border-t border-white/20 border-b-white/8 bg-[#1d1817] sticky top-0 z-100">
        <h1 className="text-2xl font-bold pl-8">
          <Link to="/">C3Subtitles Baggage Claim</Link>
        </h1>
        {auth.loggedIn && (
          <Link
            className="relative overflow-hidden w-9 h-9 mr-4 rounded-full bg-[linear-gradient(#819fd9_60%,#4767e1_70%)]"
            to="/user"
          >
            <div className="absolute text-[35px] -top-1 w-full text-center">🍿</div>
          </Link>
        )}
      </header>
      <AppMain />
      <ScrollRestoration />
    </>
  );
}

function AppMain() {
  const auth = useAuth();
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  if (auth.loggedIn == null) {
    return null;
  }

  if (auth.loggedIn === false) {
    return (
      <div className="flex justify-center items-center h-[calc(100vh-100px)]">
        <form
          className="flex flex-col w-90 grow-0 bg-white/5 border border-white/16 p-8 rounded-xl"
          onSubmit={async (e) => {
            e.preventDefault();
            const data = new FormData(e.currentTarget);
            const username = data.get('username') as string;
            const password = data.get('password') as string;

            try {
              await auth.login(username, password);
            } catch (error: any) {
              setErrorMessage(error.message);
            }
          }}
        >
          <div className="text-xl font-bold mb-4 text-center">Sign in</div>
          {errorMessage != null && (
            <div
              role="alert"
              className="border border-red-700/50 bg-red-700/20 color-white py-3 px-4 rounded-lg mb-6 text-sm"
            >
              {errorMessage}
            </div>
          )}
          <label className="text-sm mb-1 font-semibold">Username</label>
          <input
            type="text"
            name="username"
            className="bg-white/3 border border-white/20 appearance-none rounded-md w-full py-2 px-3 text-white leading-tight focus:outline-1 focus:shadow-outline mb-4"
          />
          <label className="text-sm mb-1 font-semibold">Password</label>
          <input
            type="password"
            name="password"
            className="bg-white/3 border border-white/20 appearance-none rounded-md w-full py-2 px-3 text-white leading-tight focus:outline-1 focus:shadow-outline mb-4"
          />
          <Button className="w-full mt-2 py-2" type="submit">
            Login
          </Button>
        </form>
      </div>
    );
  }

  return (
    <main>
      <Suspense fallback={<div className="px-8">Loading...</div>}>
        <Outlet />
      </Suspense>
    </main>
  );
}

export default App;
