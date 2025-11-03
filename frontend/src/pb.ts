import Pocketbase, { type RecordFullListOptions, type RecordOptions } from 'pocketbase';
import { useEffect, useState } from 'react';

export const pb = new Pocketbase('/');

export const useGetFullList = (collection: string, opts: RecordFullListOptions) => {
  const [data, setData] = useState<null | any[]>(null);

  useEffect(() => {
    pb.collection(collection).getFullList(opts).then(setData).catch(console.error);
  }, [collection, JSON.stringify(opts)]);

  return data;
};

export const useGetOne = (collection: string, id: string | undefined, opts?: RecordOptions) => {
  const [data, setData] = useState<null | any>(null);

  useEffect(() => {
    if (!id) return;
    pb.collection(collection).getOne(id, opts).then(setData).catch(console.error);
  }, [collection, id, JSON.stringify(opts)]);

  return data;
};
