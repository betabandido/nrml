import { randomIntBetween } from "https://jslib.k6.io/k6-utils/1.0.0/index.js";
import { check, randomSeed } from "k6";
import { SharedArray } from "k6/data";
import { loadProducts, getByPrd } from "./common.js";

const productSources = new SharedArray("products", () => {
  return loadProducts("./products.csv");
});

export function setup() {
  randomSeed(new Date().getMilliseconds());
  console.log(`Using ${productSources.length} product sources`);
}

export const options = {
  scenarios: {
    constant: {
      executor: 'constant-arrival-rate',
      duration: '2s',
      rate: 100,
      timeUnit: '1s',
      preAllocatedVUs: 10,
      maxVUs: 100,
    }
  },
  summaryTrendStats: ["med", "p(75)", "p(90)", "p(99)", "p(99.9)"],
  thresholds: {
    http_req_failed: ["rate<0.05"],
  }
};

export default function () {
  const source = productSources[randomIntBetween(0, productSources.length - 1)];

  const response = getByPrd(source.locale, source.productKey)

  check(response, {
    "status was 200": (r) => r.status === 200,
    "status was not 400": (r) => r.status !== 400,
    "status was not 404": (r) => r.status !== 404,
    "status was not 500": (r) => r.status !== 500,
  });

  console.log(response.status)

  if (response.status >= 300 && response.status !== 404) {
    console.log(
      JSON.stringify({
        source: JSON.stringify(source),
        url: response.request.url,
        error: response.body,
      })
    );
  }
}
