import { useAuth } from './auth';
import { HomePage } from './pages/HomePage';
import { Route, Switch } from 'wouter';
import { TalkPage } from './pages/TalkPage';
import { UserPage } from './pages/UserPage';
import { Link } from './components/Link';

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
            className="flex w-10 h-10 rounded-full text-center text-4xl bg-[linear-gradient(#819fd9_60%,#4767e1_70%)]"
            href="/user"
          >
            🏄
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
      <div className="flex justify-center items-center h-full">
        <form
          className="flex-col w-50 grow-0 bg-gray-700 p-4 rounded"
          onSubmit={(e) => {
            e.preventDefault();
            const data = new FormData(e.currentTarget);
            const username = data.get('username') as string;
            const password = data.get('password') as string;

            auth.login(username, password);
          }}
        >
          <input
            type="text"
            name="username"
            placeholder="username"
            className="bg-white shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline mb-2"
          />
          <input
            type="password"
            name="password"
            placeholder="Password"
            className="bg-white shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline mb-2"
          />
          <button
            type="submit"
            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full"
          >
            Login
          </button>
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
