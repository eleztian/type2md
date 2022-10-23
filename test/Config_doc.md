# Config Doc
Config doc.

| key      | 类型      | 必填 | 默认值           | 描述          |
|----------|----------|-----|------------------|--------------|
|Pre|[Hook](#ext.Hook)|false|||
|Post|[Hook](#ext.Hook)|false|||
|servers.{string}.host|string|false|||
|servers.{string}.port|int|true||- `22`<br>- `65522`|

## ext.Hook
Hook hook config.

| key      | 类型      | 必填 | 默认值           | 描述          |
|----------|----------|-----|------------------|--------------|
|name|string|true|example|hook name.|
|commands.[].|string|false|||
|envs.{string}.|string|false|||
|mode|[Mode](#ext.Mode)|false|1|run mode.|

## ext.Mode
**Type:** int

Mode mode define.

| Value      | 描述          |
|----------|--------------|
|1|mode q.|
|2|mode a.|

---
GENERATED BY THE COMMAND [type2md](https://github.com/eleztian/type2md)
from github.com/eleztian/type2md/test.Config