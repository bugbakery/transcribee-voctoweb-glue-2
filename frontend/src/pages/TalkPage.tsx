import { pb, useGetOne } from '../pb';
import { Button } from '../components/Button';
import { ScrollRestoration, useParams } from 'react-router';

export const TalkPage = () => {
  const { id } = useParams();
  const { data: talk, mutate } = useGetOne('talks', id!, { expand: 'conference,assignee' });

  return (
    <div className="mx-8">
      <div className="flex gap-10 items-start">
        <div>
          <h1 className="text-2xl font-bold mb-4">{talk.title}</h1>
          <h2 className="text-lg font-bold mb-4">{talk.subtitle}</h2>
          <div className="flex flex-row mb-4">
            <InfoFieldH label={talk.persons.length > 1 ? 'Speakers' : 'Speaker'}>
              {talk.persons.join(', ')}
            </InfoFieldH>
            <InfoFieldH label="Duration">{formatDuration(talk.duration_secs)}</InfoFieldH>
            <InfoFieldH label="Event">{talk.expand?.conference?.name}</InfoFieldH>
          </div>
          <div>{talk.description}</div>
        </div>
        <div className="min-w-80">
          <div className="flex justify-end gap-3 mb-4">
            {talk.assignee === pb.authStore.record?.id ? (
              <>
                <a
                  href={talk.transcribee_url + `&speaker_name_options=${talk.persons.join(',')}`}
                  target="_blank"
                  className="bg-white/80 border border-white text-black hover:bg-white text-sm font-semibold py-1 px-2 rounded-lg"
                >
                  Open Editor
                </a>
                <Button
                  type="button"
                  onClick={async () => {
                    await pb.collection('talks').update(talk.id, { assignee: null });
                    mutate();
                  }}
                >
                  Finish Work
                </Button>
              </>
            ) : (
              <>
                <Button
                  type="button"
                  onClick={async () => {
                    await pb
                      .collection('talks')
                      .update(talk.id, { assignee: pb.authStore.record?.id });
                    mutate();
                  }}
                >
                  Claim
                </Button>
              </>
            )}
          </div>
          <div className="p-8 border border-white/16 bg-white/5 rounded-2xl">
            <InfoField label="Assignee">{talk.expand?.assignee?.username || '-'}</InfoField>
            <InfoField label="Corrected until">
              {formatDuration(talk.corrected_until_secs)}
            </InfoField>
          </div>
        </div>
      </div>
      <ScrollRestoration />
    </div>
  );
};

function InfoFieldH({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="mr-4">
      <span className="font-bold">{label}:</span> <span>{children}</span>
    </div>
  );
}

function InfoField({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="flex flex-col not-last:mb-4">
      <span className="font-bold">{label}</span>
      <span>{children}</span>
    </div>
  );
}

function formatDuration(duration: number) {
  const hours = Math.floor(duration / 3600);
  const minutes = Math.floor((duration % 3600) / 60);
  const seconds = duration % 60;

  if (hours === 0) {
    return `${minutes}m ${seconds}s`;
  }

  return `${hours}h ${minutes}m ${seconds}s`;
}
