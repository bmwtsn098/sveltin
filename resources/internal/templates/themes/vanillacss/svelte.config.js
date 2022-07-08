import { mdsvex } from 'mdsvex';
import mdsvexConfig from './mdsvex.config.js';

import preprocess from 'svelte-preprocess';
import adapter from '@sveltejs/adapter-static';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	extensions: ['.svelte', ...mdsvexConfig.extensions],
	// Learn more at https://github.com/sveltejs/svelte-preprocess
	preprocess: [
		mdsvex(mdsvexConfig),
		preprocess({
			postcss: false,
			preserve: ['ld+json'],
		}),
	],
	kit: {
		adapter: adapter({
			// default options are shown
			pages: 'build',
			assets: 'build',
			fallback: null,
		}),
		prerender: {
			default: true,
			entries: ['*'],
		},
	},
};

export default config;
