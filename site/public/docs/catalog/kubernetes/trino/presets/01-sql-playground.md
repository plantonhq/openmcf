---
title: "SQL playground"
description: "The two-line Trino: a coordinator and two workers with PASSWORD authentication on (module-generated admin credential in the `trino-auth` Secret — upstream's open, anyone-can-query default never..."
type: "preset"
rank: "01"
presetSlug: "01-sql-playground"
componentSlug: "trino"
componentTitle: "Trino"
provider: "kubernetes"
icon: "package"
order: 1
---

# SQL playground

The two-line Trino: a coordinator and two workers with PASSWORD
authentication on (module-generated admin credential in the
`trino-auth` Secret — upstream's open, anyone-can-query default never
ships) and the in-image `tpch`/`tpcds` sample catalogs enabled, so a
fresh install answers `SELECT count(*) FROM tpch.tiny.nation`
immediately — no data source required.

Connect any SQL client at the exported coordinator endpoint (port
8080) with the admin username (`trino`) and the generated password:

```sql
SELECT nationkey, name FROM tpch.tiny.nation ORDER BY name LIMIT 5;
```

Grow it by declaring `catalogs.postgres` / `catalogs.mysql` entries —
each composed database becomes one more queryable catalog you can JOIN
across.
