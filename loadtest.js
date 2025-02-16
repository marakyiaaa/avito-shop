import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Trend, Rate } from 'k6/metrics';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

let loginTime = new Trend('login_time');
let buyTime = new Trend('buy_time');
let userInfoTime = new Trend('user_info_time');
let failedRequests = new Rate('failed_requests');

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const USER_PASSWORD = __ENV.USER_PASSWORD || 'test_password';

function login(username, password) {
    let res = http.post(`${BASE_URL}/api/auth`, JSON.stringify({
        username: username,
        password: password
    }), {
        headers: { 'Content-Type': 'application/json' },
    });

    check(res, { 'auth success': (r) => r.status === 200 });
    loginTime.add(res.timings.duration);
    return res.json();
}

function buyItem(token, itemId) {
    let res = http.post(`${BASE_URL}/api/buy/${itemId}`, null, { headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' } });

    check(res, { 'buy success': (r) => r.status === 200 });
    buyTime.add(res.timings.duration);
}

function getUserInfo(token) {
    let res = http.get(`${BASE_URL}/api/info`, {
        headers: { 'Authorization': `Bearer ${token}` },
    });

    check(res, { 'user info success': (r) => r.status === 200 });
    userInfoTime.add(res.timings.duration);
}

export default function () {
    let username1 = `user_${uuidv4()}`;
    let username2 = `user_${uuidv4()}`;

    let token1, token2, userID1, userID2;

    group('Authentication', function () {
        let loginRes1 = login(username1, USER_PASSWORD);
        token1 = loginRes1.token;
        userID1 = loginRes1.userID;

        let loginRes2 = login(username2, USER_PASSWORD);
        token2 = loginRes2.token;
        userID2 = loginRes2.userID;

        if (!token1 || !token2 || !userID1 || !userID2) {
            failedRequests.add(1);
            return;
        }
    });

    group('Store', function () {
        buyItem(token1, 'cup');
    });
    group('User Info', function () {
        getUserInfo(token1);
    });

    sleep(1);
}

export let options = {
    stages: [
        { duration: '1m', target: 1000 },
    ],
    rps: 1000,
};