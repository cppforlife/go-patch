## go-patch

More or less based on https://tools.ietf.org/html/rfc6902.

## Pointer (aka path)

More or less based on https://tools.ietf.org/html/rfc6901.

- Root
- Index (ex: `/0`, `/-1`)
- AfterLastIndex (ex: `/-`)
- MatchingIndex (ex: `/key=val`)
- Key: `/key`

## Operations

- Remove
- Replace

## Example

### Input

```yaml
releases:
- name: capi
  version: 0.1

instance_groups:
- name: cloud_controller
  instances: 0
  jobs:
  - name: cloud_controller
    release: capi

- name: uaa
  instances: 0
```

### Operations

```yaml
- type: replace
  path: /instance_groups/name=cloud_controller/instances
  value: 1

- type: replace
  path: /instance_groups/name=cloud_controller/jobs/name=cloud_controller/consumes?/db
  value:
    instances:
    - address: some-db.local
    properties:
      username: user
      password: pass

- type: replace
  path: /instance_groups/name=uaa/instances
  value: 1

- type: replace
  path: /instance_groups/-
  value:
    name: uaadb
    instances: 2
```

### Output

```yaml
releases:
- name: capi
  version: latest

instance_groups:
- name: cloud_controller
  instances: 1
  jobs:
  - name: cloud_controller
    release: capi
    consumes:
      db:
        instances:
        - address: some-db.local
        properties:
          username: user
          password: pass

- name: uaa
  instances: 1

- name: uaadb
  instances: 2
```
