import { useParams } from 'wouter';
import { pb, useGetOne } from '../pb';
import { Button } from '../components/Button';

export const TalkPage = () => {
  const { id } = useParams();
  const { data: talk, refresh } = useGetOne('talks', id, { expand: 'conference,assignee' });

  if (!talk) {
    return null;
  }

  return (
    <div className="mx-8">
      <div>
        <h1 className="text-2xl font-bold mb-4">{talk.title}</h1>
        <h2 className="text-lg font-bold mb-4">{talk.subtitle}</h2>
        <div className="flex flex-row mb-4">
          <InfoFieldH label={talk.persons.length > 1 ? 'Speakers' : 'Speaker'}>
            {talk.persons.join(', ')}
          </InfoFieldH>
          <InfoFieldH label="Duration">{formatDuration(talk.duration_secs)}</InfoFieldH>
          <InfoFieldH label="Event">{talk.expand.conference.name}</InfoFieldH>
        </div>
        <div className="flex gap-10 items-start">
          <div>{talk.description}</div>

          <div className="min-w-80">
            <div className="flex justify-end gap-3 mb-4">
              {talk.assignee === pb.authStore.record?.id ? (
                <>
                  <Button type="button" onClick={async () => {}}>
                    Open Editor
                  </Button>
                  <Button
                    type="button"
                    onClick={async () => {
                      await pb.collection('talks').update(talk.id, { assignee: null });
                      refresh();
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
                      refresh();
                    }}
                  >
                    Claim
                  </Button>
                </>
              )}
            </div>
            <div className="p-8 border border-white/16 bg-white/5 rounded-2xl">
              <InfoField label="Assignee">{talk.expand.assignee?.username || '-'}</InfoField>
              <InfoField label="Corrected until">{formatDuration(talk.corrected_until_secs)}</InfoField>
            </div>
          </div>
        </div>
      </div>
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
