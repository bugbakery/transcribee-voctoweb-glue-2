import { useGetFullList } from '../pb';
import { Link } from '../components/Link';
import { useNavigate } from 'react-router';
import type { RecordModel } from 'pocketbase';
import { cn } from '../cn';

export function HomePage() {
  const navigate = useNavigate();

  const { data: talks } = useGetFullList('talks', {
    filter: 'conference.name="39c3"',
    expand: 'assignee',
    fields:
      'id,title,transcribee_state,assignee,expand.assignee,date,duration_secs,corrected_until_secs',
    sort: 'release_date',
  });

  function sortOrder(state: string) {
    if (!state) return 0;
    if (state == 'todo') return 8;
    if (state == 'needs correction') return 3;
    if (state == 'done') return 7;
    if (state == 'preparing') return 6;
    if (state.startsWith('partially corrected')) return 2;
    return 0;
  }

  function stateColor(state: string) {
    if (!state) return 'bg-yellow-300';
    if (state == 'unclear') return 'bg-gray-300';
    if (state == 'preparing') return 'bg-gray-300';
    if (state == 'done') return 'bg-green-300';
    return 'bg-yellow-300';
  }

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
              .map((talk) => {
                let state = 'unclear';

                if (talk.transcribee_state === 'done' && talk.corrected_until_secs == 0) {
                  state = 'needs correction';
                } else if (
                  talk.transcribee_state === 'done' &&
                  talk.corrected_until_secs < talk.duration_secs
                ) {
                  state = `partially corrected`;
                } else if (
                  talk.transcribee_state === 'done' &&
                  talk.corrected_until_secs === talk.duration_secs
                ) {
                  state = `done`;
                } else if (talk.transcribee_state === 'in_progress') {
                  state = `preparing`;
                } else {
                  state = `preparing`;
                }
                return { ...talk, state } as RecordModel & { state: keyof typeof sortOrder };
              })
              .sort(
                (a, b) =>
                  sortOrder(a.state) * 100000000000000 +
                  +new Date(a.date) -
                  sortOrder(b.state) * 100000000000000 +
                  +new Date(b.date),
              )
              .map((talk) => {
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
                      <span
                        className={cn(
                          'text-black py-0.5 px-1 text-sm font-semibold rounded',
                          stateColor(talk.state),
                        )}
                      >
                        {talk.state}
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
