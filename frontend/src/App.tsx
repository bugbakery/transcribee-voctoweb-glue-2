import { useAuth } from './auth';
import { HomePage } from './pages/HomePage';
import { Route, Switch } from 'wouter';
import { TalkPage } from './pages/TalkPage';
import { UserPage } from './pages/UserPage';
import { Link } from './components/Link';
import { Button } from './components/Button';

function App() {
  const auth = useAuth();

  return (
    <>
      <header className="flex justify-between items-center h-14 mb-4 border-b border-white/20">
        <h1 className="text-2xl font-bold pl-8">
          <Link href="/">C3Subtitles Baggage Claim</Link>
        </h1>
        {auth.loggedIn && (
          <Link
            className="flex w-10 h-10 rounded-full text-center text-5xl bg-[linear-gradient(#819fd9_60%,#4767e1_70%)]"
            href="/user"
          >
            🍿
          </Link>
        )}
      </header>
      <AppMain />
    </>
  );
}

function AppMain() {
  const auth = useAuth();

  if (auth.loggedIn == null) {
    return null;
  }

  if (auth.loggedIn === false) {
    return (
      <div className="flex justify-center items-center h-[calc(100vh-100px)]">
        <form
          className="flex flex-col w-90 grow-0 bg-white/5 border border-white/16 p-8 rounded-xl"
          onSubmit={(e) => {
            e.preventDefault();
            const data = new FormData(e.currentTarget);
            const username = data.get('username') as string;
            const password = data.get('password') as string;

            auth.login(username, password);
          }}
        >
          <div className="text-xl font-bold mb-4 text-center">Sign in</div>
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
          <Button className="w-full mt-2 py-2" type="submit">Login</Button>
        </form>
      </div>
    );
  }

  return (
    <main>
      <Switch>
        <Route path="/" component={HomePage} />
        <Route path="/user" component={UserPage} />
        <Route path="/talk/:id" component={TalkPage} />
      </Switch>
    </main>
  );
}

export default App;
