import http.client
import json

conn = http.client.HTTPConnection("localhost", 1234)
payload = json.dumps({
  "longitude": 1234,
  "latitude": 3245345,
  "imei": "imei"
})
headers = {
  'Content-Type': 'application/json'
}
conn.request("POST", "/location", payload, headers)
res = conn.getresponse()
data = res.read()
print(data.decode("utf-8"))
