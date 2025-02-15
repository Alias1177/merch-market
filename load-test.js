import http from 'k6/http';
import { sleep, check } from 'k6';

export let options = {
    scenarios: {
        constant_request_rate: {
            executor: 'constant-arrival-rate',
            rate: 1000,            
            timeUnit: '1s',          
            duration: '10s',          
            preAllocatedVUs: 1000,    
            maxVUs: 1000,           
        },
    },
    thresholds: {
       
        http_req_duration: ['p(95)<50'],
     
        http_req_failed: ['rate<0.0001'],
    },
};

export default function () {
//Вставте свой токен после регистрации
    let res = http.get('http://localhost:8080/api/info', {
        headers: {
            'Authorization': 'Bearer  <your_jwt_token>',
        },
    });

    check(res, {
        'status is 200': (r) => r.status === 200,
    });

  
    sleep(0.001);
}
