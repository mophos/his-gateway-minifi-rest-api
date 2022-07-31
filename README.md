# his-gateway-minifi-rest-api

`env.yaml` :

```yaml
dataPath: /opt/minifi/data
settingFile: /opt/minifi/data/config/setting.yml
templatePath: /opt/minifi/data/template
outPath: /opt/minifi/conf
cmd: /opt/minifi/bin/minifi.sh
```

- `dataPath` ที่เก็บไฟล์ configure ของ flow.
- `settingFile` ที่อยู่ของไฟล์ที่เก็บค่า setting.yml
- `templatePath` โฟลเดอร์ที่เก็บไฟล์ template ของ Flow
- `outPath` ที่อยู่ของโฟลเดอร์ที่เก็บไฟล์ config.yml
- `cmd` สำหรับอ้างอิง path ที่เก็บไฟล์ `minifi.sh`
