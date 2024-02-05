import http from 'k6/http';
import { sleep } from 'k6';
export const options = {
	scenarios: {
	  constant_request_rate: {
		executor: 'constant-arrival-rate',
		rate: 20000,
		timeUnit: '1s',
		duration: '30s',
		preAllocatedVUs: 50,
		maxVUs: 15000,
	  },
	},
};
export default function () {
	  http.get('http://192.168.0.101:8000/api/v1/ad');
	  sleep(1);
}