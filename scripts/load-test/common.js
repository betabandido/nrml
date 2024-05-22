import http from "k6/http";

const baseUrl = "<ADD YOUR URL HERE>";
const urlByPrd = `${baseUrl}/api/v1/productByProductKey/tenant`;

export function getByPrd(locale, productKey) {
  return http.get(`${urlByPrd}/${locale}/${productKey}`);
}

export function loadProducts(csvFileName) {
  const lines = open(csvFileName).split("\n");
  return lines.map((line) => {
    const fields = line.split(",");
    const productKey = fields[0];
    const locale = fields[1].replace("\r", "");
    return {
      productKey,
      locale,
    };
  });
}
