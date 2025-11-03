import { useEffect, useState } from 'react';
import Pocketbase from 'pocketbase';
import { useAuth } from './auth';

const pb = new Pocketbase('');

function Home() {
  const [talks, setTalks] = useState<any[] | null>(null);

  useEffect(() => {
    pb.collection('talks')
      .getFullList({
        filter: 'conference.name="38c3"',
        expand: 'assignee',
        fields: 'id,title,state,assignee',
      })
      .then((talks) => {
        setTalks(talks);
      });
  }, []);

  return (
    <div>
      {talks && (
        <div>
          {talks.map((talk) => {
            return <p>{talk.title}</p>;
          })}
        </div>
      )}
    </div>
  );
}

function App() {
  const auth = useAuth();

  if (auth.loggedIn == null) {
    return <p>Loading...</p>;
  }

  if (auth.loggedIn === false) {
    return (
      <form
        onSubmit={(e) => {
          e.preventDefault();
          const data = new FormData(e.currentTarget);
          const username = data.get('username') as string;
          const password = data.get('password') as string;

          auth.login(username, password);
        }}
      >
        <input type="text" name="username" placeholder="username" />
        <input type="password" name="password" placeholder="Password" />
        <button type="submit">Login</button>
      </form>
    );
  }

  return <Home />;
}

export default App;
