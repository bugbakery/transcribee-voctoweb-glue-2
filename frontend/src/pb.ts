import Pocketbase, { type RecordFullListOptions, type RecordOptions } from 'pocketbase';
import useSWR from 'swr';

export const pb = new Pocketbase('/');
pb.autoCancellation(false);

export const useGetFullList = (collection: string, opts: RecordFullListOptions) => {
  return useSWR(
    ['getFullList', collection, opts],
    () => pb.collection(collection).getFullList(opts),
    { suspense: true },
  );
};

export const useGetOne = (collection: string, id: string, opts?: RecordOptions) => {
  return useSWR(
    ['getOne', collection, id, opts],
    () => pb.collection(collection).getOne(id || '', opts),
    { suspense: true },
  );
};
