import { writable } from 'svelte/store';

export const listeners = writable();
export const loading = writable(false);
export const selectedListener = writable();
