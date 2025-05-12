import axios from 'axios';

const fetch = axios.create({
	baseURL: '',
	timeout: 1000,
	headers: {
		'Accept': 'application/json',
	},
});

export default fetch;

