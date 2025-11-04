import Pocketbase, { type RecordFullListOptions, type RecordOptions } from 'pocketbase';
import { useCallback, useEffect, useState } from 'react';

export const pb = new Pocketbase('/');

const cache = new Map<string, any>();
const fromCache = (key: any) => cache.get(JSON.stringify(key)) || null;
const toCache = (key: any, value: any) => cache.set(JSON.stringify(key), value);

export const useGetFullList = (collection: string, opts: RecordFullListOptions) => {
  const [data, setData] = useState<null | any[]>(fromCache(['getFullList', collection, opts]));

  useEffect(() => {
    pb.collection(collection)
      .getFullList(opts)
      .then((data) => {
        toCache(['getFullList', collection, opts], data);
        setData(data);
      })
      .catch(console.error);
  }, [collection, JSON.stringify(opts)]);

  return data;
};

export const useGetOne = (collection: string, id: string | undefined, opts?: RecordOptions) => {
  const [data, setData] = useState<null | any>(fromCache(['getOne', collection, id, opts]));

  const load = useCallback(() => {
    if (!id) return;
    pb.collection(collection).getOne(id, opts).then((data) => {
      toCache(['getOne', collection, id, opts], data);
      setData(data);
    }).catch(console.error);
  }, [id, collection, JSON.stringify(opts)]);

  useEffect(() => {
    load();
  }, [load]);

  return { data, refresh: load };
};
