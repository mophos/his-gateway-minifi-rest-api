# his-gateway-minifi-rest-api

`env.yaml` :

```yaml
data:
  path: /data
  triggerFile: /data/update.txt
  connections: /data/connections
cmd: /opt/minifi/bin/minifi.sh
```

- `data.path` ที่เก็บไฟล์ configure ของ flow.
- `data.triggerFile` ไฟล์ที่ใช้สำหรับ trigger ให้ minifi reload config.
- `data.connections` โฟลเดอร์สำหรับเก็บไฟล์ tables, table_manual
- `cmd` สำหรับอ้างอิง path ที่เก็บไฟล์ `minifi.sh`
