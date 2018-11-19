# am-silenced

A small go program to schedule the Alert Manager silences on daily basis using the cron syntax.

The working directory from where you run the program should have process.yaml and config.yaml . See the files included in the repository for configuration syntax. See below e.g config.yaml file 

```yaml
---
# Cron syntax used: "seconds minutes hours day-of-month month day-of-week"
- silencename: "core-app"
  starttime: "0 0 0 * * *"
  duration: "6h"
  matchers:
  - name: "alertname"
    value: "backend"
    isRegex: false
  - name: "environment"
    value: "uat"
    isRegex: false
```

The matchers are used in same way as we would use them while creating silence from the Web UI.
