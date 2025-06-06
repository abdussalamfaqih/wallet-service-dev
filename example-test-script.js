import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    vus: 10,
    duration: '30s',
};

const baseUrl = 'http://localhost:8080/v1';

export default function () {
    // Generate account IDs >= 1
    const accountID1 = Math.floor(Math.random() * 100000) + 1;
    const accountID2 = accountID1 + 1;

    // 1. Create Account #1
    let res1 = http.post(`${baseUrl}/accounts`, JSON.stringify({
        account_id: accountID1,
        initial_balance: "200.23344"
    }), {
        headers: { 'Content-Type': 'application/json' }
    });

    check(res1, {
        'created account 1': (r) => r.status === 200,
    });

    // 2. Create Account #2
    let res2 = http.post(`${baseUrl}/accounts`, JSON.stringify({
        account_id: accountID2,
        initial_balance: "300.99999"
    }), {
        headers: { 'Content-Type': 'application/json' }
    });

    check(res2, {
        'created account 2': (r) => r.status === 200,
    });

    // 3. Get Account #1
    let res3 = http.get(`${baseUrl}/accounts/${accountID1}`);
    check(res3, {
        'get account 1': (r) => r.status === 200,
        'account 1 has balance': (r) => JSON.parse(r.body).balance !== undefined,
    });

    // 4. Create Transaction
    let res4 = http.post(`${baseUrl}/transactions`, JSON.stringify({
        source_account_id: accountID1,
        destination_account_id: accountID2,
        amount: "100.12345"
    }), {
        headers: { 'Content-Type': 'application/json' }
    });

    check(res4, {
        'transaction success': (r) => r.status === 200,
    });

    sleep(1);
}
