import { useGetFullList } from '../pb';
import { Link } from '../components/Link';
import { useNavigate } from 'react-router';

function formatTime(secs: number) {
  const hours = Math.floor(secs / 3600);
  const minutes = Math.floor((secs % 3600) / 60);
  const seconds = secs % 60;
  return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
}

export function HomePage() {
  const navigate = useNavigate();

  const { data: talks } = useGetFullList('talks', {
    filter: 'conference.name="38c3"',
    expand: 'assignee',
    fields:
      'id,title,transcribee_state,assignee,expand.assignee,date,duration_secs,corrected_until_secs',
    sort: 'release_date',
  });

  return (
    <div className="mx-8">
      {talks && (
        <div className="w-full mb-8">
          <div className="sticky top-[52px] bg-main-background pt-2">
            <div className="flex text-sm font-bold py-2 bg-[#403c3b] rounded-t-xl border border-white/20 border-b-white/8 z-10">
              <div className="px-6 flex-1">Title</div>
              <div className="px-6 w-50">State</div>
              <div className="px-6 w-40">Assignee</div>
            </div>
          </div>
          <div className="*:border *:border-white/16 *:border-t-0 *:flex *:last:rounded-b-xl *:bg-white/5 *:hover:bg-white/10">
            {talks
              .sort((a, b) => +new Date(a.date) - +new Date(b.date))
              .map((talk) => {
                let state = 'todo';

                if (talk.transcribee_state === 'done' && talk.corrected_until_secs == 0) {
                  state = 'needs correction';
                } else if (
                  talk.transcribee_state === 'done' &&
                  talk.corrected_until_secs < talk.duration_secs
                ) {
                  state = `corrected until ${formatTime(talk.corrected_until_secs)}`;
                } else if (
                  talk.transcribee_state === 'done' &&
                  talk.corrected_until_secs === talk.duration_secs
                ) {
                  state = `done`;
                } else if (talk.transcribee_state === 'in_progress') {
                  state = `preparing`;
                } else {
                  state = `todo`;
                }

                return (
                  <div
                    key={talk.id}
                    onClick={() => {
                      navigate(`/talk/${talk.id}`);
                    }}
                  >
                    <div className="flex-1 py-3 px-6">
                      <Link to={`/talk/${talk.id}`} onClick={(e) => e.stopPropagation()}>
                        {talk.title}
                      </Link>
                    </div>
                    <div className="py-3 px-6 w-50">
                      <span className="bg-yellow-300 text-black py-0.5 px-1 text-sm font-semibold rounded">
                        {state}
                      </span>
                    </div>
                    <div className="py-3 px-6 w-40">{talk.expand?.assignee?.username || ''}</div>
                  </div>
                );
              })}
          </div>
        </div>
      )}
    </div>
  );
}
