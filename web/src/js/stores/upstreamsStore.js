import { writable } from 'svelte/store';

export const upstreams = writable();
export const loading = writable(false);
