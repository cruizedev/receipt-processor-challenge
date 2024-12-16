import http from "k6/http";
import { check, group } from "k6";

const baseURL = "http://127.0.0.1:1378";

export let options = {
  stages: [
    {
      target: 35,
      duration: "2m",
    },
  ],
  thresholds: {
    http_req_duration: ["avg<10000", "p(100)<30000"],
    http_req_failed: ["rate<0.01"],
  },
};

export default function () {
  group("points", () => {
    let id = "";

    group("process", () => {
      let payload = `
{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}
      `;

      let res = http.post(`${baseURL}/receipts/process`, payload, {
        headers: {
          "Content-Type": "application/json",
        },
      });

      check(res, {
        success: (res) => res.status == 200,
      });

      id = res.json()["id"];
    });

    console.log(id);

    group("points", () => {
      let res = null;
      do {
        res = http.get(`${baseURL}/receipts/${id}/points`);
      } while (res.status == 202);

      check(res, {
        success: (res) => res.status == 200,
      });
    });
  });
}
