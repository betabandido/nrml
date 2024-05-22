import {BatchWriteItemCommand, DynamoDBClient} from "@aws-sdk/client-dynamodb"
import {parse} from "csv-parse/sync"
import {promises as fs} from "fs"

const tableName = "nr-memory-leak-investigation"

function buildPutRequests(records, start) {
    return records.slice(start, start + 25).map(record => {
        return {
            PutRequest: {
                Item: {
                    ProductKey: {S: `tenant#${record[1].toLowerCase()}#${record[0]}`},
                    ProductDetails: {
                        M: {
                            Key: {S: record[0]},
                            Version: {N: "1"},
                            Options: {M: {}}
                        }
                    }
                }
            }
        }
    })
}

async function main() {
    const content = await fs.readFile("../../scripts/load-test/products.csv")
    const records = parse(content, {bom: true})

    const client = new DynamoDBClient({region: "eu-west-1"})

    for (let idx = 0; ; idx += 25) {
        console.log(idx)
        const requests = buildPutRequests(records, idx)
        if (requests.length === 0) {
            break
        }

        await client.send(new BatchWriteItemCommand({
            RequestItems: {
                [tableName]: requests
            }
        }))
    }
}

main()
