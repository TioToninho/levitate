import http from 'k6/http';
import { sleep, check } from 'k6';

export default function () {
    const res = http.get('http://localhost:8080/health');
    check(res, {
        'status is 200': (r) => r.status === 200,
    });
    sleep(1);
} 