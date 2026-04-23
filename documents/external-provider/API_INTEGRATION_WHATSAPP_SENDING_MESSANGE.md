# WhatsApp Server — Evolution API

Deployed on **safer** (`safer.easipath.com`) using Docker Compose.

---

## Stack

| Service                    | Image                          | Port     | Purpose                          |
|----------------------------|--------------------------------|----------|----------------------------------|
| `evolution-api`            | `atendai/evolution-api:latest` | `8080`   | WhatsApp REST API                |
| `evolution-postgres`       | `postgres:15`                  | internal | Persistent message/session DB    |
| `evolution-redis`          | `redis:7-alpine`               | internal | Cache / queue                    |
| `evolution-web-ui`         | `nginx:alpine`                 | `8081`   | Web management interface         |
| `evolution-webhook-bridge` | `whatsapp-kafka-bridge:latest` | `3100`   | Webhook → Kafka forwarder        |

---

## Access

| Resource     | URL                                   |
|--------------|---------------------------------------|
| Web UI       | http://safer.easipath.com:8081        |
| API          | http://safer.easipath.com:8080        |
| Kafka UI     | http://safer.easipath.com:9000        |
| API Docs     | https://doc.evolution-api.com         |

**API Key:** `mUtombo8544e4EGG25841serEEESSA`

---

## Active Instance

| Field          | Value                          |
|----------------|--------------------------------|
| Instance name  | `biacibenga`                   |
| WhatsApp number| `+27684011702`                 |
| Profile name   | `Bia`                          |
| Status         | `open` (connected)             |
| Proxy          | IPRoyal Residential (sticky 7d)|

---

## Files on Server

```
/apps/docker-compose-script/whatsapp-server/
├── docker-compose.yml
├── web-ui/
│   └── index.html               # Web management UI
└── webhook-bridge/
    ├── Dockerfile
    ├── package.json
    └── index.js                 # Webhook → Kafka bridge
```

---

## Managing the Stack

```bash
ssh safer
cd /apps/docker-compose-script/whatsapp-server
```

| Action                    | Command                                                   |
|---------------------------|-----------------------------------------------------------|
| Start all                 | `sudo docker compose up -d`                               |
| Stop all                  | `sudo docker compose down`                                |
| Restart API               | `sudo docker compose restart evolution-api`               |
| Restart bridge            | `sudo docker compose restart webhook-bridge`              |
| View API logs             | `sudo docker compose logs -f evolution-api`               |
| View bridge logs          | `sudo docker compose logs -f webhook-bridge`              |
| View all containers       | `sudo docker ps --filter name=evolution`                  |

---

## Important Configuration

### Baileys Version Fix
WhatsApp requires version `2.3000.1035194821`. This is patched in two places:

1. **`.env` inside the container** — `CONFIG_SESSION_PHONE_VERSION=2.3000.1035194821`
2. **Bundled JSON** — `/evolution/node_modules/baileys/lib/Defaults/baileys-version.json`

After any `docker pull` / image update, re-apply the patch:
```bash
ssh safer
sudo docker exec evolution-api sed -i \
  's/CONFIG_SESSION_PHONE_VERSION=.*/CONFIG_SESSION_PHONE_VERSION=2.3000.1035194821/' \
  /evolution/.env
sudo docker exec evolution-api sh -c \
  'echo "{\"version\":[2,3000,1035194821]}" > /evolution/node_modules/baileys/lib/Defaults/baileys-version.json'
sudo docker restart evolution-api
```

### Proxy (IPRoyal Residential)
Required because WhatsApp blocks datacenter IPs. Configured per instance:
```bash
curl -X POST http://safer.easipath.com:8080/proxy/set/biacibenga \
  -H "Content-Type: application/json" \
  -H "apikey: mUtombo8544e4EGG25841serEEESSA" \
  -d '{
    "enabled": true,
    "host": "geo.iproyal.com",
    "port": "12321",
    "protocol": "http",
    "username": "4Me1fyVLEhkrUCr2",
    "password": "VeSsySk7UDRPZg4J_session-rI7hGHgz_lifetime-168h"
  }'
```

---

## Web UI Features

Open **http://safer.easipath.com:8081**

- **Server Config** — set API URL and key, click Connect
- **Instances** — list, create, delete instances; see connection state
- **QR Code** — click QR button next to an instance to scan and link WhatsApp
- **Send Message** — send text from any connected instance
- **Webhook** — configure inbound webhook URL per instance
- **Activity Log** — live log of API activity

---

## Sending Messages via API

All requests require header: `apikey: mUtombo8544e4EGG25841serEEESSA`

### Send a text message
```bash
curl -X POST http://safer.easipath.com:8080/message/sendText/biacibenga \
  -H "Content-Type: application/json" \
  -H "apikey: mUtombo8544e4EGG25841serEEESSA" \
  -d '{
    "number": "27821234567",
    "text": "Hello from Evolution API!"
  }'
```

### Send an image
```bash
curl -X POST http://safer.easipath.com:8080/message/sendMedia/biacibenga \
  -H "Content-Type: application/json" \
  -H "apikey: mUtombo8544e4EGG25841serEEESSA" \
  -d '{
    "number": "27821234567",
    "mediatype": "image",
    "mimetype": "image/jpeg",
    "caption": "Check this out!",
    "media": "https://example.com/image.jpg"
  }'
```

### Send a document
```bash
curl -X POST http://safer.easipath.com:8080/message/sendMedia/biacibenga \
  -H "Content-Type: application/json" \
  -H "apikey: mUtombo8544e4EGG25841serEEESSA" \
  -d '{
    "number": "27821234567",
    "mediatype": "document",
    "mimetype": "application/pdf",
    "caption": "Invoice",
    "media": "https://example.com/invoice.pdf",
    "fileName": "invoice.pdf"
  }'
```

### Send to a group
```bash
# Use the group JID (get it from chats.upsert Kafka events)
curl -X POST http://safer.easipath.com:8080/message/sendText/biacibenga \
  -H "Content-Type: application/json" \
  -H "apikey: mUtombo8544e4EGG25841serEEESSA" \
  -d '{
    "number": "120363XXXXXXXXXX@g.us",
    "text": "Hello group!"
  }'
```

### From code (Node.js)
```js
const axios = require('axios');

const WA_URL = 'http://safer.easipath.com:8080';
const API_KEY = 'mUtombo8544e4EGG25841serEEESSA';
const INSTANCE = 'biacibenga';

async function sendMessage(number, text) {
  const res = await axios.post(
    `${WA_URL}/message/sendText/${INSTANCE}`,
    { number, text },
    { headers: { apikey: API_KEY } }
  );
  return res.data;
}

sendMessage('27821234567', 'Hello!').then(console.log);
```

### From code (Go)
```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

func SendWhatsApp(number, text string) error {
    body, _ := json.Marshal(map[string]string{"number": number, "text": text})
    req, _ := http.NewRequest("POST",
        "http://safer.easipath.com:8080/message/sendText/biacibenga",
        bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("apikey", "mUtombo8544e4EGG25841serEEESSA")
    resp, err := http.DefaultClient.Do(req)
    if err != nil { return err }
    defer resp.Body.Close()
    fmt.Println("Status:", resp.Status)
    return nil
}
```

---

## Kafka Integration

### Topic
| Topic               | Partitions | Replication | Purpose                     |
|---------------------|------------|-------------|-----------------------------|
| `WHATSAPP_MESSAGES` | 3          | 3           | All WhatsApp events         |

### Message format
Every Kafka message has this structure:

```json
{
  "event": "messages.upsert",
  "instance": "biacibenga",
  "data": { ... },
  "destination": "http://evolution-webhook-bridge:3000",
  "date_time": "2026-03-31T10:00:00.000Z",
  "sender": "27684011702@s.whatsapp.net",
  "server_url": "http://safer.easipath.com:8080"
}
```

### Kafka message headers (extracted by bridge)
| Header       | Description                          |
|--------------|--------------------------------------|
| `event`      | Event type (e.g. `messages.upsert`)  |
| `instance`   | Instance name (`biacibenga`)         |
| `from`       | Sender JID (`27821234567@s.whatsapp.net`) |
| `fromMe`     | `true` if you sent it, `false` if received |
| `text`       | Message text (first 200 chars)       |
| `status`     | Delivery status (on update events)   |
| `receivedAt` | ISO timestamp when bridge received it|

### Event types
| Event              | When it fires                                  |
|--------------------|------------------------------------------------|
| `messages.upsert`  | New message received or sent                   |
| `messages.update`  | Delivery/read receipt (DELIVERY_ACK, READ)     |
| `messages.delete`  | Message deleted                                |
| `send.message`     | Outgoing message dispatched                    |
| `chats.upsert`     | New chat created                               |
| `chats.update`     | Chat metadata changed                          |
| `contacts.upsert`  | New contact saved                              |
| `contacts.update`  | Contact info updated                           |
| `connection.update`| WhatsApp connected / disconnected              |
| `call`             | Incoming call                                  |

---

## Consuming messages.upsert from Kafka

### CLI (test/debug)
```bash
ssh safer
sudo docker exec kafka0 kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic WHATSAPP_MESSAGES \
  --from-beginning
```

Filter only incoming messages:
```bash
sudo docker exec kafka0 kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic WHATSAPP_MESSAGES \
  --from-beginning | \
  python3 -c "
import sys, json
for line in sys.stdin:
    d = json.loads(line)
    if d.get('event') == 'messages.upsert':
        msgs = d['data'].get('messages', [d['data']])
        for m in msgs:
            if not m.get('key', {}).get('fromMe'):
                text = m.get('message', {}).get('conversation', '')
                print('FROM:', m['key']['remoteJid'], '| MSG:', text)
"
```

### Node.js consumer
```js
const { Kafka } = require('kafkajs');

const kafka = new Kafka({
  clientId: 'my-app',
  brokers: ['safer.easipath.com:9092']
});

const consumer = kafka.consumer({ groupId: 'my-app-group' });

async function run() {
  await consumer.connect();
  await consumer.subscribe({ topic: 'WHATSAPP_MESSAGES', fromBeginning: false });

  await consumer.run({
    eachMessage: async ({ message }) => {
      const payload = JSON.parse(message.value.toString());

      if (payload.event !== 'messages.upsert') return;

      const msgs = payload.data?.messages || [payload.data];
      for (const msg of msgs) {
        const key     = msg.key || {};
        const content = msg.message || {};
        const text    = content.conversation
                     || content.extendedTextMessage?.text
                     || '';

        if (key.fromMe) continue; // skip outgoing

        console.log({
          from:      key.remoteJid,
          text,
          timestamp: msg.messageTimestamp,
          instance:  payload.instance
        });

        // Reply to sender
        // await sendMessage(key.remoteJid.replace('@s.whatsapp.net',''), 'Got your message!');
      }
    }
  });
}

run().catch(console.error);
```

### Go consumer
```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/segmentio/kafka-go"
)

type WAPayload struct {
    Event    string                 `json:"event"`
    Instance string                 `json:"instance"`
    Data     map[string]interface{} `json:"data"`
}

func main() {
    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{"safer.easipath.com:9092"},
        Topic:   "WHATSAPP_MESSAGES",
        GroupID: "my-go-app",
    })
    defer r.Close()

    for {
        m, err := r.ReadMessage(context.Background())
        if err != nil { break }

        var payload WAPayload
        if err := json.Unmarshal(m.Value, &payload); err != nil { continue }

        if payload.Event != "messages.upsert" { continue }

        fmt.Printf("Event: %s | Instance: %s\n", payload.Event, payload.Instance)
        // Process payload.Data here
    }
}
```

---

## Other Useful API Calls

### Check instance status
```bash
curl http://safer.easipath.com:8080/instance/connectionState/biacibenga \
  -H "apikey: mUtombo8544e4EGG25841serEEESSA"
```

### List all instances
```bash
curl http://safer.easipath.com:8080/instance/fetchInstances \
  -H "apikey: mUtombo8544e4EGG25841serEEESSA"
```

### Get webhook config
```bash
curl http://safer.easipath.com:8080/webhook/find/biacibenga \
  -H "apikey: mUtombo8544e4EGG25841serEEESSA"
```

### Get contacts
```bash
curl http://safer.easipath.com:8080/contact/findContacts/biacibenga \
  -H "apikey: mUtombo8544e4EGG25841serEEESSA"
```

### Get chat messages
```bash
curl -X POST http://safer.easipath.com:8080/chat/findMessages/biacibenga \
  -H "Content-Type: application/json" \
  -H "apikey: mUtombo8544e4EGG25841serEEESSA" \
  -d '{"where": {"key": {"remoteJid": "27821234567@s.whatsapp.net"}}}'
```

---

## Database

- **Engine:** PostgreSQL 15
- **Database:** `evolution`  **User/Pass:** `evolution/evolution`
- **Volume:** `evolution_postgres_data` (persists across restarts)

```bash
ssh safer
sudo docker exec -it evolution-postgres psql -U evolution -d evolution
```

---

## Notes

- Sessions survive restarts — stored in PostgreSQL, no need to re-scan QR.
- After any image update (`docker pull`), re-apply the Baileys version patch (see above).
- Proxy renews every 7 days — update the password suffix `_lifetime-168h` session token if it expires.
- Kafka UI: http://safer.easipath.com:9000 → topic `WHATSAPP_MESSAGES`
- Bridge logs: `sudo docker logs -f evolution-webhook-bridge`
