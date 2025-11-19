import { useNavigate } from 'react-router';
import { useAuth } from '../auth';
import { Button } from '../components/Button';

export const UserPage = () => {
  const auth = useAuth();
  const navigate = useNavigate();

  return (
    <div className="mx-8">
      <h1>Hello User</h1>
      <p>At some point here might be statistics and claimed tasks. But not today...</p>

      <br />
      <Button
        onClick={() => {
          auth.logout();
          navigate('/');
        }}
      >
        Logout
      </Button>
    </div>
  );
};
