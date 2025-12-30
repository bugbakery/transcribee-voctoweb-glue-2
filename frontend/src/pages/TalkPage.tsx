import { pb, useGetOne } from '../pb';
import { Button } from '../components/Button';
import { useParams } from 'react-router';
import { Dialog } from '../components/Dialog';
import { useState } from 'react';

export const TalkPage = () => {
  const { id } = useParams();
  const { data: talk, mutate } = useGetOne('talks', id!, { expand: 'conference,assignee' });
  const [isPublishing, setIsPublishing] = useState(false);

  const [isCorrectionTimeDialogOpen, setIsCorrectionTimeDialogOpen] = useState(false);
  const [correctedUntilInput, setCorrectedUntilInput] = useState('');

  return (
    <>
      <Dialog
        title="Until when did you correct the transcript?"
        isOpen={isCorrectionTimeDialogOpen}
        onClose={() => setIsCorrectionTimeDialogOpen(false)}
      >
        <input
          placeholder="13:12"
          className="w-full border rounded p-2"
          value={correctedUntilInput}
          onChange={(e) => setCorrectedUntilInput(e.target.value)}
        />
        <div className="flex justify-between">
          <Button
            className="bg-white/50"
            onClick={async () => {
              if (!talk) return;
              await pb
                .collection('talks')
                .update(talk.id, { corrected_until_secs: talk.duration_secs, assignee: null });
              mutate();
              setIsCorrectionTimeDialogOpen(false);
            }}
          >
            Until the End
          </Button>
          <Button
            onClick={async () => {
              if (!talk) return;
              const parts = correctedUntilInput.split(':');
              let hours = '0';
              let mins = '0';
              let secs = '0';
              if (parts.length == 2) {
                hours = '0';
                mins = parts[0];
                secs = parts[1];
              } else {
                hours = parts[0];
                mins = parts[1];
                secs = parts[2];
              }

              const corrected_until_secs_user =
                parseInt(hours) * 60 * 60 + parseInt(mins) * 60 + parseInt(secs);
              const corrected_until_secs = Math.min(talk.duration_secs, corrected_until_secs_user);
              await pb
                .collection('talks')
                .update(talk.id, { corrected_until_secs, assignee: null });
              mutate();
              setIsCorrectionTimeDialogOpen(false);
            }}
          >
            Ok
          </Button>
        </div>
      </Dialog>

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
              {talk.assignee === pb.authStore.record?.id && (
                <>
                  {talk.transcribee_url && (
                    <a
                      href={
                        talk.transcribee_url + `&speaker_name_options=${talk.persons.join(',')}`
                      }
                      target="_blank"
                      className="bg-white/80 border border-white text-black hover:bg-white text-sm font-semibold py-1 px-2 rounded-lg"
                    >
                      Open Editor
                    </a>
                  )}
                  <Button
                    type="button"
                    onClick={() => {
                      setIsCorrectionTimeDialogOpen(true);
                    }}
                  >
                    Finish Work
                  </Button>
                </>
              )}

              {!talk.assignee && <Button
                type="button"
                onClick={async () => {
                  await pb
                    .collection('talks')
                    .update(talk.id, { assignee: pb.authStore.record?.id });
                  mutate();
                }}
              >
                Claim
              </Button>}

              <Button
                type="button"
                onClick={async () => {
                  try {
                    setIsPublishing(true);
                    await pb.send(`api/talks/${talk.id}/publish`, {
                      method: 'POST',
                    });
                    window.alert("Talk published!");
                  } catch (error) {
                    window.alert("Error: Publishing failed");
                  } finally{
                    setIsPublishing(false);
                  }
                }}
                disabled={isPublishing}
                title={isPublishing ? 'Publishing...' : undefined}
              >
                Publish
              </Button>
            </div>
            <div className="p-8 border border-white/16 bg-white/5 rounded-2xl">
              <InfoField label="Assignee">{talk.expand?.assignee?.username || '-'}</InfoField>
              <InfoField label="Corrected until">
                {formatDuration(talk.corrected_until_secs)}{' '}
                {talk.corrected_until_secs == talk.duration_secs ? '(done)' : ''}
              </InfoField>
              <InfoField label="Published at">
                {talk.published_at ? formatDateTime(talk.published_at) : '-'}
              </InfoField>
            </div>
          </div>
        </div>
      </div>
    </>
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

function formatDateTime(date: string) {
  const parsedDate = new Date(date);
  return (
    <div>
      {parsedDate.toLocaleDateString()} {parsedDate.toLocaleTimeString()}
    </div>
  );
}
